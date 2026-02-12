package export

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/fabiobrug/mako.git/internal/database"
)

// ExportOptions configures what to export
type ExportOptions struct {
	Last         int       // Last N commands
	Semantic     string    // Semantic search query
	DateFrom     time.Time // Start date
	DateTo       time.Time // End date
	WorkingDir   string    // Filter by directory
	SuccessOnly  bool      // Only successful commands
	FailedOnly   bool      // Only failed commands
}

// Exporter handles command history export
type Exporter struct {
	db *database.DB
}

// NewExporter creates a new exporter
func NewExporter(db *database.DB) *Exporter {
	return &Exporter{db: db}
}

// Export writes commands to the writer in JSON format
func (e *Exporter) Export(w io.Writer, opts ExportOptions) error {
	commands, err := e.fetchCommands(opts)
	if err != nil {
		return fmt.Errorf("failed to fetch commands: %w", err)
	}

	// Convert to export format
	exportData := ExportFormat{
		Version:    CurrentVersion,
		ExportedAt: time.Now(),
		Commands:   make([]ExportedCommand, 0, len(commands)),
	}

	for _, cmd := range commands {
		exportData.Commands = append(exportData.Commands, ExportedCommand{
			Command:       cmd.Command,
			Timestamp:     cmd.Timestamp,
			ExitCode:      cmd.ExitCode,
			DurationMS:    cmd.Duration,
			WorkingDir:    cmd.WorkingDir,
			OutputPreview: cmd.OutputPreview,
		})
	}

	// Encode to JSON with pretty printing
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	
	if err := encoder.Encode(exportData); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

// fetchCommands retrieves commands based on options
func (e *Exporter) fetchCommands(opts ExportOptions) ([]database.Command, error) {
	// Handle semantic search
	if opts.Semantic != "" {
		// This would require embedding generation and semantic search
		// For now, fall back to FTS search
		return e.db.SearchCommands(opts.Semantic, 10000)
	}

	// Handle last N commands
	if opts.Last > 0 {
		return e.db.GetRecentCommands(opts.Last)
	}

	// Handle directory filter
	if opts.WorkingDir != "" {
		return e.db.GetCommandsByDirectory(opts.WorkingDir, 10000)
	}

	// Handle success/failed filter
	if opts.SuccessOnly {
		return e.db.GetCommandsByExitCode(true, 10000)
	}
	if opts.FailedOnly {
		return e.db.GetCommandsByExitCode(false, 10000)
	}

	// Handle date range
	if !opts.DateFrom.IsZero() || !opts.DateTo.IsZero() {
		return e.fetchByDateRange(opts.DateFrom, opts.DateTo)
	}

	// Default: all commands (limited to 10k for safety)
	return e.db.GetRecentCommands(10000)
}

// fetchByDateRange retrieves commands within a date range
func (e *Exporter) fetchByDateRange(from, to time.Time) ([]database.Command, error) {
	query := `
		SELECT id, command, timestamp, exit_code, duration_ms, working_dir, output_preview
		FROM commands
		WHERE 1=1
	`
	args := []interface{}{}

	if !from.IsZero() {
		query += " AND timestamp >= ?"
		args = append(args, from)
	}

	if !to.IsZero() {
		query += " AND timestamp <= ?"
		args = append(args, to)
	}

	query += " ORDER BY timestamp DESC LIMIT 10000"

	rows, err := e.db.GetConn().Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commands []database.Command
	for rows.Next() {
		var cmd database.Command
		err := rows.Scan(
			&cmd.ID,
			&cmd.Command,
			&cmd.Timestamp,
			&cmd.ExitCode,
			&cmd.Duration,
			&cmd.WorkingDir,
			&cmd.OutputPreview,
		)
		if err != nil {
			continue
		}
		commands = append(commands, cmd)
	}

	return commands, nil
}
