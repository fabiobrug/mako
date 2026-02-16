package shell

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/fabiobrug/mako.git/internal/database"
	"github.com/fabiobrug/mako.git/internal/export"
)

func handleExport(args []string, db *database.DB) (string, error) {
	if db == nil {
		return "\r\n✗ Database not available\r\n\r\n", nil
	}
	
	if len(args) == 0 {
		return "Usage: mako export [--last N] [--dir /path] > output.json\r\n", nil
	}
	
	opts := export.ExportOptions{}
	
	// Parse arguments
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--last":
			if i+1 < len(args) {
				fmt.Sscanf(args[i+1], "%d", &opts.Last)
				i++
			}
		case "--dir":
			if i+1 < len(args) {
				opts.WorkingDir = args[i+1]
				i++
			}
		case "--success":
			opts.SuccessOnly = true
		case "--failed":
			opts.FailedOnly = true
		}
	}
	
	// Default to last 1000 if nothing specified
	if opts.Last == 0 && opts.WorkingDir == "" && !opts.SuccessOnly && !opts.FailedOnly {
		opts.Last = 1000
	}
	
	// Create exporter
	exporter := export.NewExporter(db)
	
	// Export to stdout
	var buf bytes.Buffer
	if err := exporter.Export(&buf, opts); err != nil {
		return "", fmt.Errorf("export failed: %w", err)
	}
	
	// Convert to proper line endings
	output := strings.ReplaceAll(buf.String(), "\n", "\r\n")
	return output, nil
}

func handleImport(args []string, db *database.DB) (string, error) {
	if db == nil {
		return "\r\n✗ Database not available\r\n\r\n", nil
	}
	
	if len(args) == 0 {
		return "Usage: mako import [--merge|--skip|--overwrite] <file.json>\r\n", nil
	}
	
	opts := export.ImportOptions{
		ConflictStrategy: export.ConflictSkip, // Default
	}
	
	var filename string
	
	// Parse arguments
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--merge":
			opts.ConflictStrategy = export.ConflictMerge
		case "--skip":
			opts.ConflictStrategy = export.ConflictSkip
		case "--overwrite":
			opts.ConflictStrategy = export.ConflictOverwrite
		case "--dry-run":
			opts.DryRun = true
		default:
			filename = args[i]
		}
	}
	
	if filename == "" {
		return "Error: No file specified\r\n", nil
	}
	
	// Open file
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	
	// Create importer
	importer := export.NewImporter(db)
	
	// Import
	result, err := importer.Import(file, opts)
	if err != nil {
		return "", fmt.Errorf("import failed: %w", err)
	}
	
	// Format result
	output := fmt.Sprintf("\r\nImport complete:\r\n")
	output += fmt.Sprintf("  Total: %d\r\n", result.TotalCommands)
	output += fmt.Sprintf("  Imported: %d\r\n", result.ImportedNew)
	output += fmt.Sprintf("  Skipped: %d\r\n", result.Skipped)
	output += fmt.Sprintf("  Updated: %d\r\n", result.Updated)
	
	if len(result.Errors) > 0 {
		output += fmt.Sprintf("  Errors: %d\r\n", len(result.Errors))
		for i, errMsg := range result.Errors {
			if i < 5 { // Show first 5 errors
				output += fmt.Sprintf("    - %s\r\n", errMsg)
			}
		}
		if len(result.Errors) > 5 {
			output += fmt.Sprintf("    ... and %d more\r\n", len(result.Errors)-5)
		}
	}
	
	return output, nil
}

func handleSync(db *database.DB) (string, error) {
	if db == nil {
		return "\r\n✗ Database not available\r\n\r\n", nil
	}
	
	// Get default history path
	historyPath := database.GetDefaultHistoryPath()
	
	if historyPath == "" {
		return "Error: Could not find bash history file\r\n", nil
	}
	
	// Sync with limit of 100 new commands
	count, err := db.SyncBashHistory(historyPath, 100)
	if err != nil {
		return "", fmt.Errorf("sync failed: %w", err)
	}
	
	output := fmt.Sprintf("Synced %d new commands from %s\r\n", count, historyPath)
	return output, nil
}
