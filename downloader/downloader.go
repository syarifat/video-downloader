package downloader

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// VideoInfo contains metadata about the video
type VideoInfo struct {
	Title     string  `json:"title"`
	Duration  float64 `json:"duration"`
	URL       string  `json:"webpage_url"`
	Thumbnail string  `json:"thumbnail"`
	Extractor string  `json:"extractor"`
	FileSize  int64   `json:"filesize_approx"`
	Format    string  `json:"format"`
}

// Result contains the download result
type Result struct {
	FilePath  string
	FileName  string
	Title     string
	Duration  float64
	FileSize  int64
	Extractor string
}

// GetVideoInfo fetches video metadata without downloading
func GetVideoInfo(ctx context.Context, url string) (*VideoInfo, error) {
	cmd := exec.CommandContext(ctx, "yt-dlp",
		"--dump-json",
		"--no-playlist",
		"--no-warnings",
		url,
	)

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("yt-dlp error: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("failed to get video info: %w", err)
	}

	var info VideoInfo
	if err := json.Unmarshal(output, &info); err != nil {
		return nil, fmt.Errorf("failed to parse video info: %w", err)
	}

	return &info, nil
}

// Download downloads a video from the given URL
func Download(ctx context.Context, url string, downloadDir string, maxSizeMB int) (*Result, error) {
	// Generate unique filename
	timestamp := time.Now().UnixNano()
	outputTemplate := filepath.Join(downloadDir, fmt.Sprintf("%d_%%(title).50s.%%(ext)s", timestamp))

	// Build yt-dlp command with optimized settings
	args := []string{
		// Output
		"-o", outputTemplate,
		// Format selection: prefer mp4 under size limit
		"-f", fmt.Sprintf("best[filesize<%dM]/bestvideo[filesize<%dM]+bestaudio/best", maxSizeMB, maxSizeMB),
		// Merge to mp4
		"--merge-output-format", "mp4",
		// Re-encode to mp4 if needed
		"--recode-video", "mp4",
		// No playlist
		"--no-playlist",
		// No warnings
		"--no-warnings",
		// Limit download speed to prevent issues
		"--no-part",
		// Embed thumbnail if possible
		"--no-embed-thumbnail",
		// Restrict filenames
		"--restrict-filenames",
		// Add URL
		url,
	}

	cmd := exec.CommandContext(ctx, "yt-dlp", args...)
	cmd.Stderr = os.Stderr

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("download failed: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("download failed: %w", err)
	}

	// Find the downloaded file
	filePath, err := findDownloadedFile(downloadDir, timestamp, string(output))
	if err != nil {
		return nil, fmt.Errorf("could not find downloaded file: %w", err)
	}

	// Get file info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not stat downloaded file: %w", err)
	}

	// Check file size
	fileSizeMB := fileInfo.Size() / (1024 * 1024)
	if fileSizeMB > int64(maxSizeMB) {
		os.Remove(filePath)
		return nil, fmt.Errorf("file too large (%dMB > %dMB limit). Try a shorter video", fileSizeMB, maxSizeMB)
	}

	return &Result{
		FilePath: filePath,
		FileName: fileInfo.Name(),
		FileSize: fileInfo.Size(),
	}, nil
}

// findDownloadedFile locates the downloaded file in the directory
func findDownloadedFile(dir string, timestamp int64, output string) (string, error) {
	prefix := fmt.Sprintf("%d_", timestamp)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), prefix) {
			return filepath.Join(dir, entry.Name()), nil
		}
	}

	// Fallback: try to extract from yt-dlp output
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "[Merger]") || strings.Contains(line, "Destination:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				candidate := strings.TrimSpace(parts[1])
				if _, err := os.Stat(candidate); err == nil {
					return candidate, nil
				}
			}
		}
	}

	return "", fmt.Errorf("no file found with prefix %s in %s", prefix, dir)
}

// Cleanup removes a downloaded file
func Cleanup(filePath string) {
	if filePath != "" {
		os.Remove(filePath)
	}
}

// FormatDuration formats seconds to HH:MM:SS
func FormatDuration(seconds float64) string {
	d := time.Duration(seconds) * time.Second
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60

	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}

// FormatFileSize formats bytes to human-readable size
func FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
