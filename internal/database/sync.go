package database

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// BashHistoryEntry represents a single bash history entry
type BashHistoryEntry struct {
	Command   string
	Timestamp time.Time
}

// SyncBashHistory synchronizes new bash history commands since last sync
func (db *DB) SyncBashHistory(historyPath string, limit int) (int, error) {
	// Get last sync time
	lastSync, err := db.GetLastSyncTime()
	if err != nil {
		return 0, fmt.Errorf("failed to get last sync time: %w", err)
	}

	// Parse bash history
	entries, err := parseBashHistory(historyPath, lastSync)
	if err != nil {
		return 0, fmt.Errorf("failed to parse bash history: %w", err)
	}

	// Limit number of entries
	if len(entries) > limit {
		entries = entries[len(entries)-limit:]
	}

	if len(entries) == 0 {
		return 0, nil
	}

	// Convert to commands
	commands := make([]Command, 0, len(entries))
	for _, entry := range entries {
		// Skip mako commands to avoid recursive storage
		if strings.HasPrefix(entry.Command, "mako ") {
			continue
		}

		commands = append(commands, Command{
			Command:         entry.Command,
			Timestamp:       entry.Timestamp,
			ExitCode:        0, // Unknown from history file
			Duration:        0, // Unknown from history file
			WorkingDir:      "", // Unknown from history file
			EmbeddingStatus: "pending",
		})
	}

	// Bulk insert
	if err := db.BulkInsertCommands(commands); err != nil {
		return 0, fmt.Errorf("failed to bulk insert commands: %w", err)
	}

	// Update last sync time
	if len(entries) > 0 {
		latestTime := entries[len(entries)-1].Timestamp
		if err := db.SetLastSyncTime(latestTime); err != nil {
			return 0, fmt.Errorf("failed to update sync time: %w", err)
		}
	}

	return len(commands), nil
}

// parseBashHistory parses bash history file
// Supports both timestamped (#1234567890) and non-timestamped formats
func parseBashHistory(historyPath string, since time.Time) ([]BashHistoryEntry, error) {
	file, err := os.Open(historyPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []BashHistoryEntry
	scanner := bufio.NewScanner(file)
	
	var currentTimestamp time.Time
	useFileModTime := false

	// Check if file has timestamps by reading first few lines
	firstLines := make([]string, 0, 10)
	for i := 0; i < 10 && scanner.Scan(); i++ {
		firstLines = append(firstLines, scanner.Text())
	}

	// Check if any line starts with #
	hasTimestamps := false
	for _, line := range firstLines {
		if strings.HasPrefix(line, "#") && len(line) > 1 {
			hasTimestamps = true
			break
		}
	}

	// Reset file
	file.Seek(0, 0)
	scanner = bufio.NewScanner(file)

	if !hasTimestamps {
		// Use file modification time as fallback
		info, err := file.Stat()
		if err == nil {
			currentTimestamp = info.ModTime()
			useFileModTime = true
		} else {
			currentTimestamp = time.Now()
		}
	}

	for scanner.Scan() {
		line := scanner.Text()

		// Check for timestamp line
		if strings.HasPrefix(line, "#") && len(line) > 1 {
			// Parse timestamp: #1234567890
			timestampStr := strings.TrimPrefix(line, "#")
			timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
			if err == nil {
				currentTimestamp = time.Unix(timestamp, 0)
			}
			continue
		}

		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		// This is a command line
		if useFileModTime {
			// All commands get same timestamp (file mod time)
			if currentTimestamp.After(since) {
				entries = append(entries, BashHistoryEntry{
					Command:   line,
					Timestamp: currentTimestamp,
				})
			}
		} else {
			// Use parsed timestamp
			if currentTimestamp.IsZero() {
				currentTimestamp = time.Now()
			}

			if currentTimestamp.After(since) {
				entries = append(entries, BashHistoryEntry{
					Command:   line,
					Timestamp: currentTimestamp,
				})
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

// GetDefaultHistoryPath returns the default bash history path
func GetDefaultHistoryPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	// Try .bash_history first, then .zsh_history
	bashHistory := filepath.Join(home, ".bash_history")
	if _, err := os.Stat(bashHistory); err == nil {
		return bashHistory
	}

	zshHistory := filepath.Join(home, ".zsh_history")
	if _, err := os.Stat(zshHistory); err == nil {
		return zshHistory
	}

	return bashHistory // Return default even if doesn't exist
}
