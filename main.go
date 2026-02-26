package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-video-bot/bot"
	"go-video-bot/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	// Check if yt-dlp is available
	if err := bot.CheckYTDLP(); err != nil {
		log.Fatalf("❌ yt-dlp not found: %v\nPlease install yt-dlp: https://github.com/yt-dlp/yt-dlp#installation", err)
	}

	fmt.Println("🤖 Starting Video Download Bot...")
	fmt.Printf("📁 Download directory: %s\n", cfg.DownloadDir)
	fmt.Printf("📏 Max file size: %d MB\n", cfg.MaxFileSizeMB)

	// Create and start bot
	b, err := bot.New(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to create bot: %v", err)
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		fmt.Println("\n🛑 Shutting down bot...")
		b.Stop()
	}()

	fmt.Println("✅ Bot is running! Send a video link to download.")
	b.Start()
}
