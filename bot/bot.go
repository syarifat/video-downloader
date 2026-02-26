package bot

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"go-video-bot/config"
	"go-video-bot/downloader"

	tele "gopkg.in/telebot.v3"
)

// Bot wraps the Telegram bot with video download capabilities
type Bot struct {
	bot    *tele.Bot
	cfg    *config.Config
	active sync.Map // Track active downloads per user
}

// CheckYTDLP verifies that yt-dlp is installed and accessible
func CheckYTDLP() error {
	cmd := exec.Command("yt-dlp", "--version")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("yt-dlp is not installed or not in PATH: %w", err)
	}
	log.Printf("📦 yt-dlp version: %s", strings.TrimSpace(string(output)))
	return nil
}

// New creates a new Bot instance
func New(cfg *config.Config) (*Bot, error) {
	pref := tele.Settings{
		Token:  cfg.TelegramToken,
		Poller: &tele.LongPoller{Timeout: 30 * time.Second},
	}

	teleBot, err := tele.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram bot: %w", err)
	}

	b := &Bot{
		bot: teleBot,
		cfg: cfg,
	}

	b.registerHandlers()

	return b, nil
}

// Start starts the bot polling
func (b *Bot) Start() {
	b.bot.Start()
}

// Stop stops the bot
func (b *Bot) Stop() {
	b.bot.Stop()
}

// registerHandlers sets up all the message handlers
func (b *Bot) registerHandlers() {
	b.bot.Handle("/start", b.handleStart)
	b.bot.Handle("/help", b.handleHelp)
	b.bot.Handle(tele.OnText, b.handleText)
}

// handleStart handles the /start command
func (b *Bot) handleStart(c tele.Context) error {
	msg := `🎬 *Video Download Bot*

Selamat datang! Saya bisa mendownload video dari berbagai platform.

*Cara penggunaan:*
Kirim link video dan saya akan mendownload lalu mengirimkannya ke Anda.

*Platform yang didukung:*
🔴 YouTube
🟢 Doodstream / Dood
🔵 Videy
🟡 Dailymotion
🟣 Vimeo
⚫ TikTok
🟠 Instagram
📌 Twitter/X
📎 Facebook
🎵 SoundCloud (audio)
Dan masih banyak lagi!

Ketik /help untuk bantuan lebih lanjut.`

	return c.Send(msg, &tele.SendOptions{ParseMode: tele.ModeMarkdown})
}

// handleHelp handles the /help command
func (b *Bot) handleHelp(c tele.Context) error {
	msg := `📖 *Panduan Penggunaan*

1️⃣ Copy link video dari browser/aplikasi
2️⃣ Paste dan kirim ke bot ini
3️⃣ Tunggu proses download selesai
4️⃣ Video akan dikirim ke chat Anda

⚠️ *Batasan:*
• Ukuran file maksimal: %d MB
• Durasi video tidak dibatasi (tapi file harus < %d MB)
• Satu download per waktu

💡 *Tips:*
• Pastikan link valid dan bisa diakses
• Video pendek akan lebih cepat diproses
• Jika gagal, coba lagi beberapa saat kemudian

🔧 *Perintah:*
/start - Mulai bot
/help - Tampilkan bantuan ini`

	return c.Send(fmt.Sprintf(msg, b.cfg.MaxFileSizeMB, b.cfg.MaxFileSizeMB), &tele.SendOptions{ParseMode: tele.ModeMarkdown})
}

// handleText handles incoming text messages (checks for URLs)
func (b *Bot) handleText(c tele.Context) error {
	text := strings.TrimSpace(c.Text())

	// Check if user is allowed
	if !b.isAllowed(c.Sender().ID) {
		return c.Send("⛔ Maaf, Anda tidak memiliki akses ke bot ini.")
	}

	// Extract URL from message
	videoURL := extractURL(text)
	if videoURL == "" {
		return c.Send("❌ Tidak ada link yang valid ditemukan.\n\nKirim link video untuk mendownload. Contoh:\nhttps://www.youtube.com/watch?v=xxx")
	}

	// Check if user already has an active download
	userID := c.Sender().ID
	if _, loaded := b.active.LoadOrStore(userID, true); loaded {
		return c.Send("⏳ Anda masih memiliki download yang sedang berjalan. Mohon tunggu sampai selesai.")
	}
	defer b.active.Delete(userID)

	return b.processDownload(c, videoURL)
}

// processDownload handles the full download workflow
func (b *Bot) processDownload(c tele.Context, videoURL string) error {
	// Send initial status
	statusMsg, err := b.bot.Send(c.Recipient(), "🔍 Mengambil informasi video...")
	if err != nil {
		return err
	}

	// Create a context with timeout (10 minutes max)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Get video info first
	info, err := downloader.GetVideoInfo(ctx, videoURL)
	if err != nil {
		b.bot.Edit(statusMsg, fmt.Sprintf("❌ Gagal mengambil info video:\n%s", truncateError(err.Error())))
		return nil
	}

	// Update status with video info
	infoText := fmt.Sprintf("📹 *%s*\n", escapeMarkdown(info.Title))
	if info.Duration > 0 {
		infoText += fmt.Sprintf("⏱ Durasi: %s\n", downloader.FormatDuration(info.Duration))
	}
	if info.Extractor != "" {
		infoText += fmt.Sprintf("🌐 Sumber: %s\n", info.Extractor)
	}
	infoText += "\n⬇️ Mendownload video..."

	b.bot.Edit(statusMsg, infoText, &tele.SendOptions{ParseMode: tele.ModeMarkdown})

	// Download the video
	result, err := downloader.Download(ctx, videoURL, b.cfg.DownloadDir, b.cfg.MaxFileSizeMB)
	if err != nil {
		b.bot.Edit(statusMsg, fmt.Sprintf("❌ Gagal mendownload video:\n%s", truncateError(err.Error())))
		return nil
	}
	defer downloader.Cleanup(result.FilePath)

	// Update status
	b.bot.Edit(statusMsg, "📤 Mengirim video ke Telegram...")

	// Send the video
	file, err := os.Open(result.FilePath)
	if err != nil {
		b.bot.Edit(statusMsg, "❌ Gagal membuka file video.")
		return nil
	}
	defer file.Close()

	video := &tele.Video{
		File:    tele.FromReader(file),
		Caption: fmt.Sprintf("🎬 %s\n📦 %s", info.Title, downloader.FormatFileSize(result.FileSize)),
	}

	// Set filename
	video.FileName = sanitizeFilename(info.Title) + ".mp4"

	err = c.Send(video)
	if err != nil {
		// If sending as video fails, try as document
		file.Seek(0, 0)
		doc := &tele.Document{
			File:     tele.FromReader(file),
			FileName: sanitizeFilename(info.Title) + ".mp4",
			Caption:  fmt.Sprintf("🎬 %s\n📦 %s", info.Title, downloader.FormatFileSize(result.FileSize)),
		}
		err = c.Send(doc)
		if err != nil {
			b.bot.Edit(statusMsg, fmt.Sprintf("❌ Gagal mengirim video. File mungkin terlalu besar (%s).", downloader.FormatFileSize(result.FileSize)))
			return nil
		}
	}

	// Delete status message after successful send
	b.bot.Delete(statusMsg)
	return nil
}

// isAllowed checks if a user ID is in the allowed list
func (b *Bot) isAllowed(userID int64) bool {
	if len(b.cfg.AllowedUsers) == 0 {
		return true // No restrictions
	}
	for _, id := range b.cfg.AllowedUsers {
		if id == userID {
			return true
		}
	}
	return false
}

// extractURL extracts the first URL from a text message
func extractURL(text string) string {
	words := strings.Fields(text)
	for _, word := range words {
		// Clean up the word
		word = strings.Trim(word, "<>()[]")

		// Check if it looks like a URL
		if strings.HasPrefix(word, "http://") || strings.HasPrefix(word, "https://") {
			if u, err := url.Parse(word); err == nil && u.Host != "" {
				return word
			}
		}

		// Try adding https:// prefix for common patterns
		if strings.Contains(word, ".") && !strings.Contains(word, " ") {
			withScheme := "https://" + word
			if u, err := url.Parse(withScheme); err == nil && u.Host != "" && strings.Contains(u.Host, ".") {
				return withScheme
			}
		}
	}
	return ""
}

// escapeMarkdown escapes special Markdown characters
func escapeMarkdown(s string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"`", "\\`",
	)
	return replacer.Replace(s)
}

// sanitizeFilename removes invalid characters from filename
func sanitizeFilename(name string) string {
	// Replace invalid characters
	invalid := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	result := name
	for _, ch := range invalid {
		result = strings.ReplaceAll(result, ch, "_")
	}
	// Truncate to reasonable length
	if len(result) > 60 {
		result = result[:60]
	}
	return strings.TrimSpace(result)
}

// truncateError truncates long error messages
func truncateError(msg string) string {
	if len(msg) > 300 {
		return msg[:300] + "..."
	}
	return msg
}
