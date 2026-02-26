package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// Config holds the bot configuration
type Config struct {
	// Telegram bot token from @BotFather
	TelegramToken string

	// Directory to store downloaded videos temporarily
	DownloadDir string

	// Maximum file size in MB (Telegram limit is 50MB for bots)
	MaxFileSizeMB int

	// Allowed user IDs (empty = allow all)
	AllowedUsers []int64
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN environment variable is required")
	}

	downloadDir := os.Getenv("DOWNLOAD_DIR")
	if downloadDir == "" {
		downloadDir = filepath.Join(os.TempDir(), "video-bot-downloads")
	}

	// Create download directory if it doesn't exist
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create download directory: %w", err)
	}

	maxSize := 50 // Default 50MB (Telegram bot limit)
	if s := os.Getenv("MAX_FILE_SIZE_MB"); s != "" {
		parsed, err := strconv.Atoi(s)
		if err == nil && parsed > 0 {
			maxSize = parsed
		}
	}

	cfg := &Config{
		TelegramToken: token,
		DownloadDir:   downloadDir,
		MaxFileSizeMB: maxSize,
	}

	// Parse allowed users
	if users := os.Getenv("ALLOWED_USERS"); users != "" {
		cfg.AllowedUsers = parseUserIDs(users)
	}

	return cfg, nil
}

// parseUserIDs parses comma-separated user IDs
func parseUserIDs(s string) []int64 {
	var ids []int64
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			part := s[start:i]
			// Trim spaces
			for len(part) > 0 && part[0] == ' ' {
				part = part[1:]
			}
			for len(part) > 0 && part[len(part)-1] == ' ' {
				part = part[:len(part)-1]
			}
			if id, err := strconv.ParseInt(part, 10, 64); err == nil {
				ids = append(ids, id)
			}
			start = i + 1
		}
	}
	return ids
}
