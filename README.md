# 🎬 Go Video Bot

Bot Telegram untuk mendownload video dari berbagai platform dan mengirimkannya langsung ke chat Telegram.

## ✨ Fitur

- 📥 Download video dari berbagai platform (YouTube, Doodstream, Videy, Dailymotion, Vimeo, TikTok, Instagram, Twitter/X, Facebook, dll)
- 🤖 Integrasi penuh dengan Telegram Bot
- 📦 Otomatis konversi ke format MP4
- 🔒 Pembatasan akses user (opsional)
- ⚡ Satu download per user untuk menghindari spam
- 🧹 Auto cleanup file setelah dikirim

## 📋 Prasyarat

1. **Go** 1.21+ terinstal
2. **yt-dlp** terinstal dan tersedia di PATH
   - Windows: `winget install yt-dlp` atau download dari [GitHub](https://github.com/yt-dlp/yt-dlp/releases)
   - Linux: `pip install yt-dlp` atau `sudo apt install yt-dlp`
   - macOS: `brew install yt-dlp`
3. **ffmpeg** terinstal (untuk konversi format)
   - Windows: `winget install ffmpeg` atau download dari [FFmpeg.org](https://ffmpeg.org/download.html)
   - Linux: `sudo apt install ffmpeg`
   - macOS: `brew install ffmpeg`
4. **Telegram Bot Token** dari [@BotFather](https://t.me/BotFather)

## 🚀 Instalasi

### 1. Clone repository

```bash
git clone <repository-url>
cd go-video-bot
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Konfigurasi

Copy file `.env.example` menjadi `.env` dan isi konfigurasi:

```bash
cp .env.example .env
```

Edit file `.env`:

```env
TELEGRAM_BOT_TOKEN=your-telegram-bot-token-here
DOWNLOAD_DIR=./downloads
MAX_FILE_SIZE_MB=50
ALLOWED_USERS=
```

### 4. Set environment variable dan jalankan

**Windows (CMD):**

```cmd
set TELEGRAM_BOT_TOKEN=your-token-here
go run .
```

**Windows (PowerShell):**

```powershell
$env:TELEGRAM_BOT_TOKEN="your-token-here"
go run .
```

**Linux/macOS:**

```bash
export TELEGRAM_BOT_TOKEN=your-token-here
go run .
```

### 5. Build (opsional)

```bash
go build -o video-bot.exe .
```

## ⚙️ Konfigurasi

| Variable             | Deskripsi                                    | Default         | Required |
| -------------------- | -------------------------------------------- | --------------- | -------- |
| `TELEGRAM_BOT_TOKEN` | Token bot dari @BotFather                    | -               | ✅       |
| `DOWNLOAD_DIR`       | Direktori untuk menyimpan download sementara | System temp dir | ❌       |
| `MAX_FILE_SIZE_MB`   | Ukuran file maksimal (MB)                    | 50              | ❌       |
| `ALLOWED_USERS`      | User ID yang diizinkan (comma-separated)     | Semua diizinkan | ❌       |

## 📱 Cara Penggunaan

1. Buka bot Telegram Anda
2. Ketik `/start` untuk memulai
3. Kirim/paste link video
4. Tunggu bot mendownload dan mengirim video

## 🌐 Platform yang Didukung

Bot ini menggunakan `yt-dlp` sebagai backend, sehingga mendukung [1000+ situs](https://github.com/yt-dlp/yt-dlp/blob/master/supportedsites.md), termasuk:

- YouTube
- Doodstream / Dood
- Videy
- Dailymotion
- Vimeo
- TikTok
- Instagram (Reels, Posts)
- Twitter/X
- Facebook
- Reddit
- Twitch
- SoundCloud
- Dan masih banyak lagi...

## 📁 Struktur Project

```
go-video-bot/
├── main.go              # Entry point
├── bot/
│   └── bot.go           # Telegram bot handler
├── config/
│   └── config.go        # Configuration management
├── downloader/
│   └── downloader.go    # Video download logic (yt-dlp wrapper)
├── go.mod
├── go.sum
├── .env.example
├── .gitignore
└── README.md
```

## 🚀 Deploy ke Hosting Gratis

### Opsi 1: Railway.app (Rekomendasi ⭐)

Railway memberikan **$5 credit gratis per bulan** (~500 jam runtime), cukup untuk bot Telegram 24/7.

**Langkah-langkah:**

1. Buat akun di [railway.app](https://railway.app)
2. Push project ini ke GitHub repository
3. Di Railway dashboard, klik **"New Project"** → **"Deploy from GitHub Repo"**
4. Pilih repository `go-video-bot`
5. Railway akan otomatis mendeteksi `Dockerfile`
6. Tambahkan environment variable:
   - Klik service → **Variables** → **New Variable**
   - Tambahkan `TELEGRAM_BOT_TOKEN` = `token-dari-botfather`
7. Klik **Deploy** — selesai! 🎉

**Via Railway CLI:**

```bash
# Install Railway CLI
npm install -g @railway/cli

# Login
railway login

# Link ke project
railway link

# Set token
railway variables set TELEGRAM_BOT_TOKEN=your-token-here

# Deploy
railway up
```

---

### Opsi 2: Fly.io (Alternatif)

Fly.io memberikan **3 VM shared gratis**, cocok untuk bot yang selalu aktif.

**Langkah-langkah:**

1. Install Fly CLI: [fly.io/docs/hands-on/install-flyctl](https://fly.io/docs/hands-on/install-flyctl/)
2. Login: `fly auth login`
3. Deploy:

```bash
# Launch app (pertama kali)
fly launch --name go-video-bot --region sin --no-deploy

# Set secret token
fly secrets set TELEGRAM_BOT_TOKEN=your-token-here

# Deploy
fly deploy
```

---

### Opsi 3: Render.com

1. Buat akun di [render.com](https://render.com)
2. New → **Background Worker**
3. Connect GitHub repo
4. Runtime: **Docker**
5. Tambahkan `TELEGRAM_BOT_TOKEN` di Environment
6. Deploy

> ⚠️ Free tier Render akan spin down setelah idle, tapi akan restart otomatis saat ada request.

## ⚠️ Batasan

- **Ukuran file**: Telegram Bot API membatasi upload file hingga 50MB
- **Format**: Video akan dikonversi ke MP4 secara otomatis
- **Concurrent**: Satu user hanya bisa download satu video pada satu waktu

## 📝 Lisensi

MIT License
