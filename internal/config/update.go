package config

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	GitHubAPIURL = "https://api.github.com/repos/fabiobrug/mako/releases/latest"
	CurrentVersion = "1.0.0"
)

// UpdateInfo contains information about available updates
type UpdateInfo struct {
	Available      bool
	CurrentVersion string
	LatestVersion  string
	ReleaseNotes   string
	DownloadURL    string
}

// GitHubRelease represents a GitHub release response
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Body    string `json:"body"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// CheckForUpdates checks if a new version is available
func CheckForUpdates() (*UpdateInfo, error) {
	resp, err := http.Get(GitHubAPIURL)
	if err != nil {
		return nil, fmt.Errorf("failed to check for updates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to check for updates: HTTP %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to parse release info: %w", err)
	}

	info := &UpdateInfo{
		CurrentVersion: CurrentVersion,
		LatestVersion:  strings.TrimPrefix(release.TagName, "v"),
		ReleaseNotes:   release.Body,
	}

	// Compare versions (simple string comparison works for semantic versioning)
	info.Available = info.LatestVersion > info.CurrentVersion

	// Find the appropriate binary for this platform
	if info.Available {
		binaryName := fmt.Sprintf("mako-%s-%s", runtime.GOOS, runtime.GOARCH)
		for _, asset := range release.Assets {
			if asset.Name == binaryName {
				info.DownloadURL = asset.BrowserDownloadURL
				break
			}
		}
	}

	return info, nil
}

// InstallUpdate downloads and installs the latest version
func InstallUpdate(info *UpdateInfo) error {
	if !info.Available {
		return fmt.Errorf("no update available")
	}

	if info.DownloadURL == "" {
		return fmt.Errorf("no download URL for this platform")
	}

	cyan := "\033[38;2;0;209;255m"
	lightBlue := "\033[38;2;93;173;226m"
	reset := "\033[0m"

	fmt.Printf("%sDownloading Mako v%s...%s\n", lightBlue, info.LatestVersion, reset)

	// Download the new binary
	resp, err := http.Get(info.DownloadURL)
	if err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}
	defer resp.Body.Close()

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "mako-update-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write downloaded content
	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		return fmt.Errorf("failed to save update: %w", err)
	}
	tmpFile.Close()

	// Make it executable
	if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
		return fmt.Errorf("failed to make binary executable: %w", err)
	}

	fmt.Printf("%sInstalling...%s\n", lightBlue, reset)

	// Get current binary path
	currentBinary, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current binary path: %w", err)
	}

	// Resolve symlinks
	currentBinary, err = filepath.EvalSymlinks(currentBinary)
	if err != nil {
		return fmt.Errorf("failed to resolve binary path: %w", err)
	}

	// Create backup of current binary
	backupPath := currentBinary + ".backup"
	if err := copyFile(currentBinary, backupPath); err != nil {
		return fmt.Errorf("failed to backup current binary: %w", err)
	}

	// Replace current binary with new one
	if err := copyFile(tmpFile.Name(), currentBinary); err != nil {
		// Restore backup on failure
		copyFile(backupPath, currentBinary)
		os.Remove(backupPath)
		return fmt.Errorf("failed to install update: %w", err)
	}

	// Remove backup
	os.Remove(backupPath)

	fmt.Printf("%s✓ Updated to v%s successfully!%s\n", cyan, info.LatestVersion, reset)
	fmt.Printf("%sRestart Mako to use the new version%s\n\n", lightBlue, reset)

	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	// Copy permissions
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, sourceInfo.Mode())
}

// CheckUpdateOnStartup checks for updates and notifies the user
func CheckUpdateOnStartup(config *Config) {
	if !config.AutoUpdate {
		return
	}

	// Check for updates in background
	go func() {
		info, err := CheckForUpdates()
		if err != nil {
			return // Silent fail for startup check
		}

		if info.Available {
			lightBlue := "\033[38;2;93;173;226m"
			cyan := "\033[38;2;0;209;255m"
			reset := "\033[0m"

			fmt.Printf("%sℹ  New version available: %sv%s%s (you have v%s)\n",
				lightBlue, cyan, info.LatestVersion, reset, info.CurrentVersion)
			fmt.Printf("%s   Run 'mako update install' to update%s\n\n",
				lightBlue, reset)
		}
	}()
}
