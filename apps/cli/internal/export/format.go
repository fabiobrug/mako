package export

import (
	"time"
)

// ExportFormat defines the JSON export structure
type ExportFormat struct {
	Version    string           `json:"version"`
	ExportedAt time.Time        `json:"exported_at"`
	Commands   []ExportedCommand `json:"commands"`
}

// ExportedCommand represents a command in the export
type ExportedCommand struct {
	Command      string    `json:"command"`
	Timestamp    time.Time `json:"timestamp"`
	ExitCode     int       `json:"exit_code"`
	DurationMS   int64     `json:"duration_ms"`
	WorkingDir   string    `json:"working_dir"`
	OutputPreview string   `json:"output_preview,omitempty"`
}

const CurrentVersion = "1.0"
