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

## 🚀 Deploy ke Render.com (Gratis)

Render menyediakan **Background Worker gratis**, cocok untuk bot Telegram.

### Langkah-langkah:

1. **Buat akun** di [render.com](https://render.com)
2. **Push project** ini ke GitHub repository
3. Di Render dashboard, klik **"New +"** → **"Background Worker"**
4. **Connect** GitHub repo `go-video-bot`
5. Isi konfigurasi:
   - **Name**: `go-video-bot`
   - **Region**: `Singapore (Southeast Asia)`
   - **Runtime**: `Docker`
   - **Instance Type**: `Free`
6. Tambahkan **Environment Variable**:
   - Key: `TELEGRAM_BOT_TOKEN`
   - Value: `token-dari-botfather`
7. Klik **"Create Background Worker"** — selesai! 🎉

### Atau deploy via Blueprint:

1. Di Render dashboard, klik **"New +"** → **"Blueprint"**
2. Connect repo yang berisi file `render.yaml`
3. Render akan otomatis membaca konfigurasi
4. Isi value `TELEGRAM_BOT_TOKEN`
5. Deploy

## ⚠️ Batasan

- **Ukuran file**: Telegram Bot API membatasi upload file hingga 50MB
- **Format**: Video akan dikonversi ke MP4 secara otomatis
- **Concurrent**: Satu user hanya bisa download satu video pada satu waktu

## 📝 Lisensi

MIT License
