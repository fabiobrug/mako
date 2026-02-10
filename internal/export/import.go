package export

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/fabiobrug/mako.git/internal/database"
)

// ConflictStrategy defines how to handle duplicate commands
type ConflictStrategy string

const (
	ConflictSkip      ConflictStrategy = "skip"      // Skip duplicates
	ConflictMerge     ConflictStrategy = "merge"     // Update timestamp
	ConflictOverwrite ConflictStrategy = "overwrite" // Replace existing
)

// ImportOptions configures import behavior
type ImportOptions struct {
	ConflictStrategy ConflictStrategy
	DryRun           bool // Don't actually import, just validate
}

// Importer handles command history import
type Importer struct {
	db *database.DB
}

// NewImporter creates a new importer
func NewImporter(db *database.DB) *Importer {
	return &Importer{db: db}
}

// ImportResult contains import statistics
type ImportResult struct {
	TotalCommands   int
	ImportedNew     int
	Skipped         int
	Updated         int
	Errors          []string
}

// Import reads commands from JSON and imports them
func (i *Importer) Import(r io.Reader, opts ImportOptions) (*ImportResult, error) {
	// Decode JSON
	var exportData ExportFormat
	decoder := json.NewDecoder(r)
	
	if err := decoder.Decode(&exportData); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	// Validate version
	if exportData.Version != CurrentVersion {
		return nil, fmt.Errorf("unsupported export version: %s (expected %s)", exportData.Version, CurrentVersion)
	}

	result := &ImportResult{
		TotalCommands: len(exportData.Commands),
		Errors:        make([]string, 0),
	}

	// Process each command
	for idx, exportCmd := range exportData.Commands {
		// Validate command
		if err := i.validateCommand(exportCmd); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("command %d: %v", idx, err))
			continue
		}

		if opts.DryRun {
			result.ImportedNew++
			continue
		}

		// Import command
		if err := i.importCommand(exportCmd, opts.ConflictStrategy, result); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("command %d: %v", idx, err))
		}
	}

	return result, nil
}

// validateCommand checks if a command is valid for import
func (i *Importer) validateCommand(cmd ExportedCommand) error {
	if cmd.Command == "" {
		return fmt.Errorf("empty command")
	}

	if cmd.Timestamp.IsZero() {
		return fmt.Errorf("missing timestamp")
	}

	// Check for suspicious commands (optional safety)
	// Could add checks for dangerous commands like "rm -rf /"

	return nil
}

// importCommand imports a single command
func (i *Importer) importCommand(exportCmd ExportedCommand, strategy ConflictStrategy, result *ImportResult) error {
	// Convert to database command
	cmd := database.Command{
		Command:       exportCmd.Command,
		Timestamp:     exportCmd.Timestamp,
		ExitCode:      exportCmd.ExitCode,
		Duration:      exportCmd.DurationMS,
		WorkingDir:    exportCmd.WorkingDir,
		OutputPreview: exportCmd.OutputPreview,
	}

	switch strategy {
	case ConflictSkip:
		// Try to save with deduplication
		isNew, _, err := i.db.SaveCommandDeduplicated(cmd)
		if err != nil {
			return err
		}
		if isNew {
			result.ImportedNew++
		} else {
			result.Skipped++
		}

	case ConflictMerge:
		// Always use deduplicated save (updates last_used)
		isNew, _, err := i.db.SaveCommandDeduplicated(cmd)
		if err != nil {
			return err
		}
		if isNew {
			result.ImportedNew++
		} else {
			result.Updated++
		}

	case ConflictOverwrite:
		// Delete existing and insert new
		hash := hashCommand(cmd.Command)
		i.db.GetConn().Exec("DELETE FROM commands WHERE command_hash = ?", hash)
		
		_, err := i.db.SaveCommandAsync(cmd)
		if err != nil {
			return err
		}
		result.ImportedNew++

	default:
		return fmt.Errorf("unknown conflict strategy: %s", strategy)
	}

	return nil
}

// hashCommand generates a hash for the command (helper)
func hashCommand(command string) string {
	// Simple implementation - would be better to share this with database/db.go
	// For now, return the command itself as placeholder
	// The actual hashing is done in SaveCommandDeduplicated
	return command
}
