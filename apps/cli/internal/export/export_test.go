package export

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fabiobrug/mako.git/internal/database"
	"github.com/fabiobrug/mako.git/internal/testutil"
)

func setupTestDB(t *testing.T) *database.DB {
	t.Helper()
	tmpDir := testutil.TempDir(t)
	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := database.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	return db
}

func TestExportBasic(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	// Insert test commands
	commands := []database.Command{
		{
			Command:       "ls -lha",
			Timestamp:     time.Now(),
			ExitCode:      0,
			Duration:      100,
			WorkingDir:    "/home/user",
			OutputPreview: "total 48K",
		},
		{
			Command:    "git status",
			Timestamp:  time.Now(),
			ExitCode:   0,
			WorkingDir: "/home/user/project",
		},
	}
	
	for _, cmd := range commands {
		db.SaveCommand(cmd)
	}
	
	// Export
	exporter := NewExporter(db)
	var buf bytes.Buffer
	err := exporter.Export(&buf, ExportOptions{Last: 10})
	if err != nil {
		t.Fatalf("Export() failed: %v", err)
	}
	
	// Verify output
	if buf.Len() == 0 {
		t.Error("Expected non-empty export output")
	}
	
	// Parse JSON
	var exportData ExportFormat
	err = json.Unmarshal(buf.Bytes(), &exportData)
	if err != nil {
		t.Fatalf("Failed to parse export JSON: %v", err)
	}
	
	if exportData.Version != CurrentVersion {
		t.Errorf("Expected version %s, got %s", CurrentVersion, exportData.Version)
	}
	
	if len(exportData.Commands) != 2 {
		t.Errorf("Expected 2 commands, got %d", len(exportData.Commands))
	}
}

func TestExportWithFilters(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	// Insert commands with different exit codes
	commands := []database.Command{
		{Command: "ls", ExitCode: 0, Timestamp: time.Now()},
		{Command: "cat missing", ExitCode: 1, Timestamp: time.Now()},
		{Command: "pwd", ExitCode: 0, Timestamp: time.Now()},
	}
	
	for _, cmd := range commands {
		db.SaveCommand(cmd)
	}
	
	// Export only successful
	exporter := NewExporter(db)
	var buf bytes.Buffer
	err := exporter.Export(&buf, ExportOptions{
		SuccessOnly: true,
		Last:        10,
	})
	if err != nil {
		t.Fatalf("Export() failed: %v", err)
	}
	
	var exportData ExportFormat
	json.Unmarshal(buf.Bytes(), &exportData)
	
	if len(exportData.Commands) != 2 {
		t.Errorf("Expected 2 successful commands, got %d", len(exportData.Commands))
	}
	
	// Verify all have exit code 0
	for _, cmd := range exportData.Commands {
		if cmd.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", cmd.ExitCode)
		}
	}
}

func TestExportLast(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	// Insert many commands
	for i := 0; i < 10; i++ {
		db.SaveCommand(database.Command{
			Command:   "test",
			Timestamp: time.Now(),
		})
		time.Sleep(1 * time.Millisecond)
	}
	
	// Export last 3
	exporter := NewExporter(db)
	var buf bytes.Buffer
	err := exporter.Export(&buf, ExportOptions{Last: 3})
	if err != nil {
		t.Fatalf("Export() failed: %v", err)
	}
	
	var exportData ExportFormat
	json.Unmarshal(buf.Bytes(), &exportData)
	
	if len(exportData.Commands) != 3 {
		t.Errorf("Expected 3 commands, got %d", len(exportData.Commands))
	}
}

func TestExportByWorkingDir(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	// Insert commands from different directories
	commands := []database.Command{
		{Command: "ls", WorkingDir: "/home/user/project1", Timestamp: time.Now()},
		{Command: "pwd", WorkingDir: "/home/user/project2", Timestamp: time.Now()},
		{Command: "git status", WorkingDir: "/home/user/project1", Timestamp: time.Now()},
	}
	
	for _, cmd := range commands {
		db.SaveCommand(cmd)
	}
	
	// Export from specific directory
	exporter := NewExporter(db)
	var buf bytes.Buffer
	err := exporter.Export(&buf, ExportOptions{
		WorkingDir: "/home/user/project1",
		Last:       10,
	})
	if err != nil {
		t.Fatalf("Export() failed: %v", err)
	}
	
	var exportData ExportFormat
	json.Unmarshal(buf.Bytes(), &exportData)
	
	if len(exportData.Commands) != 2 {
		t.Errorf("Expected 2 commands from project1, got %d", len(exportData.Commands))
	}
}

func TestImportBasic(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	// Create import data
	importData := ExportFormat{
		Version:    "1.0",
		ExportedAt: time.Now(),
		Commands: []ExportedCommand{
			{
				Command:    "ls -lha",
				Timestamp:  time.Now(),
				ExitCode:   0,
				DurationMS: 100,
				WorkingDir: "/home/user",
			},
			{
				Command:    "git status",
				Timestamp:  time.Now(),
				ExitCode:   0,
				WorkingDir: "/home/user/project",
			},
		},
	}
	
	// Convert to JSON
	jsonData, _ := json.Marshal(importData)
	reader := bytes.NewReader(jsonData)
	
	// Import
	importer := NewImporter(db)
	stats, err := importer.Import(reader, ImportOptions{ConflictStrategy: ConflictSkip})
	if err != nil {
		t.Fatalf("Import() failed: %v", err)
	}
	
	if stats.ImportedNew != 2 {
		t.Errorf("Expected 2 imported commands, got %d", stats.ImportedNew)
	}
	
	// Verify commands were saved
	commands, _ := db.GetRecentCommands(10)
	if len(commands) != 2 {
		t.Errorf("Expected 2 commands in database, got %d", len(commands))
	}
}

func TestImportSkipDuplicates(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	// Save existing command
	existing := database.Command{
		Command:   "ls -lha",
		Timestamp: time.Now().Add(-1 * time.Hour),
	}
	db.SaveCommand(existing)
	
	// Try to import same command
	importData := ExportFormat{
		Version: "1.0",
		Commands: []ExportedCommand{
			{
				Command:   "ls -lha",
				Timestamp: time.Now().Add(-1 * time.Hour),
			},
			{
				Command:   "pwd",
				Timestamp: time.Now(),
			},
		},
	}
	
	jsonData, _ := json.Marshal(importData)
	reader := bytes.NewReader(jsonData)
	
	importer := NewImporter(db)
	stats, err := importer.Import(reader, ImportOptions{ConflictStrategy: ConflictSkip})
	if err != nil {
		t.Fatalf("Import() failed: %v", err)
	}
	
	// Should skip duplicate, import only new command
	if stats.Skipped != 1 {
		t.Errorf("Expected 1 skipped command, got %d", stats.Skipped)
	}
	if stats.ImportedNew != 1 {
		t.Errorf("Expected 1 imported command, got %d", stats.ImportedNew)
	}
}

func TestExportImportRoundtrip(t *testing.T) {
	db1 := setupTestDB(t)
	defer db1.Close()
	
	// Insert test commands
	originalCommands := []database.Command{
		{
			Command:       "ls -lha",
			Timestamp:     time.Now(),
			ExitCode:      0,
			Duration:      100,
			WorkingDir:    "/home/user",
			OutputPreview: "test output",
		},
		{
			Command:    "git status",
			Timestamp:  time.Now(),
			ExitCode:   0,
			WorkingDir: "/home/user/project",
		},
	}
	
	for _, cmd := range originalCommands {
		db1.SaveCommand(cmd)
	}
	
	// Export
	exporter := NewExporter(db1)
	var buf bytes.Buffer
	err := exporter.Export(&buf, ExportOptions{Last: 10})
	if err != nil {
		t.Fatalf("Export() failed: %v", err)
	}
	
	// Import into new database
	db2 := setupTestDB(t)
	defer db2.Close()
	
	importer := NewImporter(db2)
	reader := bytes.NewReader(buf.Bytes())
	stats, err := importer.Import(reader, ImportOptions{ConflictStrategy: ConflictSkip})
	if err != nil {
		t.Fatalf("Import() failed: %v", err)
	}
	
	if stats.ImportedNew != 2 {
		t.Errorf("Expected 2 imported commands, got %d", stats.ImportedNew)
	}
	
	// Verify commands match
	commands, _ := db2.GetRecentCommands(10)
	if len(commands) != 2 {
		t.Fatalf("Expected 2 commands, got %d", len(commands))
	}
	
	// Check command content (order may differ)
	commandMap := make(map[string]bool)
	for _, cmd := range commands {
		commandMap[cmd.Command] = true
	}
	
	if !commandMap["ls -lha"] || !commandMap["git status"] {
		t.Error("Imported commands don't match original")
	}
}

func TestExportToFile(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	// Insert test command
	db.SaveCommand(database.Command{
		Command:   "test command",
		Timestamp: time.Now(),
	})
	
	// Export to file
	tmpDir := testutil.TempDir(t)
	exportPath := filepath.Join(tmpDir, "export.json")
	
	exporter := NewExporter(db)
	file, err := os.Create(exportPath)
	if err != nil {
		t.Fatalf("Failed to create export file: %v", err)
	}
	
	err = exporter.Export(file, ExportOptions{Last: 10})
	file.Close()
	if err != nil {
		t.Fatalf("Export() failed: %v", err)
	}
	
	// Verify file exists and is valid JSON
	data, err := os.ReadFile(exportPath)
	if err != nil {
		t.Fatalf("Failed to read export file: %v", err)
	}
	
	var exportData ExportFormat
	err = json.Unmarshal(data, &exportData)
	if err != nil {
		t.Fatalf("Export file contains invalid JSON: %v", err)
	}
	
	if len(exportData.Commands) != 1 {
		t.Errorf("Expected 1 command in export, got %d", len(exportData.Commands))
	}
}

func TestInvalidImportData(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	// Try to import invalid JSON
	invalidJSON := bytes.NewReader([]byte("{ invalid json }"))
	
	importer := NewImporter(db)
	_, err := importer.Import(invalidJSON, ImportOptions{ConflictStrategy: ConflictSkip})
	if err == nil {
		t.Error("Expected error importing invalid JSON")
	}
}

func TestImportStats(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	// Save existing command
	db.SaveCommand(database.Command{
		Command:   "existing",
		Timestamp: time.Now().Add(-1 * time.Hour),
	})
	
	// Import mix of new and duplicate
	importData := ExportFormat{
		Version: "1.0",
		Commands: []ExportedCommand{
			{Command: "existing", Timestamp: time.Now().Add(-1 * time.Hour)},
			{Command: "new1", Timestamp: time.Now()},
			{Command: "new2", Timestamp: time.Now()},
		},
	}
	
	jsonData, _ := json.Marshal(importData)
	reader := bytes.NewReader(jsonData)
	
	importer := NewImporter(db)
	stats, err := importer.Import(reader, ImportOptions{ConflictStrategy: ConflictSkip})
	if err != nil {
		t.Fatalf("Import() failed: %v", err)
	}
	
	if stats.TotalCommands != 3 {
		t.Errorf("Expected total 3, got %d", stats.TotalCommands)
	}
	if stats.ImportedNew != 2 {
		t.Errorf("Expected 2 imported, got %d", stats.ImportedNew)
	}
	if stats.Skipped != 1 {
		t.Errorf("Expected 1 skipped, got %d", stats.Skipped)
	}
	if len(stats.Errors) != 0 {
		t.Errorf("Expected 0 failed, got %d", len(stats.Errors))
	}
}

func BenchmarkExport(b *testing.B) {
	tmpDir, _ := os.MkdirTemp("", "mako-bench-*")
	defer os.RemoveAll(tmpDir)
	
	dbPath := filepath.Join(tmpDir, "bench.db")
	db, _ := database.NewDB(dbPath)
	defer db.Close()
	
	// Insert test data
	for i := 0; i < 100; i++ {
		db.SaveCommand(database.Command{
			Command:   "test command",
			Timestamp: time.Now(),
		})
	}
	
	exporter := NewExporter(db)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		exporter.Export(&buf, ExportOptions{Last: 100})
	}
}

func BenchmarkImport(b *testing.B) {
	tmpDir, _ := os.MkdirTemp("", "mako-bench-*")
	defer os.RemoveAll(tmpDir)
	
	dbPath := filepath.Join(tmpDir, "bench.db")
	db, _ := database.NewDB(dbPath)
	defer db.Close()
	
	// Create import data
	importData := ExportFormat{
		Version:  "1.0",
		Commands: make([]ExportedCommand, 100),
	}
	
	for i := 0; i < 100; i++ {
		importData.Commands[i] = ExportedCommand{
			Command:   "test command",
			Timestamp: time.Now(),
		}
	}
	
	jsonData, _ := json.Marshal(importData)
	
	importer := NewImporter(db)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(jsonData)
		importer.Import(reader, ImportOptions{ConflictStrategy: ConflictSkip})
	}
}
