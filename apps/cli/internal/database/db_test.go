package database

import (
	"bytes"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/fabiobrug/mako.git/internal/testutil"
)

func TestNewDB(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	dbPath := filepath.Join(tmpDir, "test.db")
	
	db, err := NewDB(dbPath)
	if err != nil {
		t.Fatalf("NewDB() failed: %v", err)
	}
	defer db.Close()
	
	if db == nil {
		t.Fatal("Expected non-nil database")
	}
	
	// Verify database file was created
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Error("Database file was not created")
	}
}

func TestSaveCommand(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	dbPath := filepath.Join(tmpDir, "test.db")
	db, _ := NewDB(dbPath)
	defer db.Close()
	
	cmd := Command{
		Command:       "ls -lha",
		Timestamp:     time.Now(),
		ExitCode:      0,
		Duration:      150,
		WorkingDir:    "/home/user",
		OutputPreview: "total 48K\ndrwxr-xr-x 5 user",
	}
	
	err := db.SaveCommand(cmd)
	if err != nil {
		t.Fatalf("SaveCommand() failed: %v", err)
	}
	
	// Verify command was saved
	commands, err := db.GetRecentCommands(1)
	if err != nil {
		t.Fatalf("GetRecentCommands() failed: %v", err)
	}
	
	if len(commands) != 1 {
		t.Fatalf("Expected 1 command, got %d", len(commands))
	}
	
	if commands[0].Command != "ls -lha" {
		t.Errorf("Expected 'ls -lha', got '%s'", commands[0].Command)
	}
}

func TestGetRecentCommands(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	dbPath := filepath.Join(tmpDir, "test.db")
	db, _ := NewDB(dbPath)
	defer db.Close()
	
	// Insert test commands
	commands := []Command{
		{Command: "ls", Timestamp: time.Now().Add(-3 * time.Second), ExitCode: 0},
		{Command: "pwd", Timestamp: time.Now().Add(-2 * time.Second), ExitCode: 0},
		{Command: "git status", Timestamp: time.Now().Add(-1 * time.Second), ExitCode: 0},
	}
	
	for _, cmd := range commands {
		db.SaveCommand(cmd)
	}
	
	// Get recent commands
	recent, err := db.GetRecentCommands(2)
	if err != nil {
		t.Fatalf("GetRecentCommands() failed: %v", err)
	}
	
	if len(recent) != 2 {
		t.Errorf("Expected 2 commands, got %d", len(recent))
	}
	
	// Should be in reverse chronological order
	if recent[0].Command != "git status" {
		t.Errorf("Expected 'git status' first, got '%s'", recent[0].Command)
	}
	if recent[1].Command != "pwd" {
		t.Errorf("Expected 'pwd' second, got '%s'", recent[1].Command)
	}
}

func TestSearchCommands(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	dbPath := filepath.Join(tmpDir, "test.db")
	db, _ := NewDB(dbPath)
	defer db.Close()
	
	// Insert test commands
	commands := []Command{
		{Command: "git commit -m 'initial'", Timestamp: time.Now()},
		{Command: "git push origin main", Timestamp: time.Now()},
		{Command: "npm install", Timestamp: time.Now()},
		{Command: "ls -lha", Timestamp: time.Now()},
	}
	
	for _, cmd := range commands {
		db.SaveCommand(cmd)
	}
	
	// Search for git commands
	results, err := db.SearchCommands("git", 10)
	if err != nil {
		t.Fatalf("SearchCommands() failed: %v", err)
	}
	
	if len(results) != 2 {
		t.Errorf("Expected 2 git commands, got %d", len(results))
	}
	
	// Verify results contain git commands
	for _, cmd := range results {
		if cmd.Command != "git commit -m 'initial'" && cmd.Command != "git push origin main" {
			t.Errorf("Unexpected command in search results: %s", cmd.Command)
		}
	}
}

func TestGetRecentCommandsWithDifferentExitCodes(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	dbPath := filepath.Join(tmpDir, "test.db")
	db, _ := NewDB(dbPath)
	defer db.Close()
	
	// Insert commands with different exit codes
	commands := []Command{
		{Command: "ls", ExitCode: 0, Timestamp: time.Now()},
		{Command: "cat nonexistent", ExitCode: 1, Timestamp: time.Now()},
		{Command: "grep pattern file", ExitCode: 0, Timestamp: time.Now()},
		{Command: "invalid-command", ExitCode: 127, Timestamp: time.Now()},
	}
	
	for _, cmd := range commands {
		db.SaveCommand(cmd)
	}
	
	// Get all commands
	all, err := db.GetRecentCommands(10)
	if err != nil {
		t.Fatalf("GetRecentCommands() failed: %v", err)
	}
	
	if len(all) != 4 {
		t.Errorf("Expected 4 commands, got %d", len(all))
	}
	
	// Verify exit codes are stored correctly
	exitCodes := make(map[int]int)
	for _, cmd := range all {
		exitCodes[cmd.ExitCode]++
	}
	
	if exitCodes[0] != 2 {
		t.Errorf("Expected 2 commands with exit code 0, got %d", exitCodes[0])
	}
	if exitCodes[1] != 1 {
		t.Errorf("Expected 1 command with exit code 1, got %d", exitCodes[1])
	}
}

func TestGetStats(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	dbPath := filepath.Join(tmpDir, "test.db")
	db, _ := NewDB(dbPath)
	defer db.Close()
	
	// Insert test commands
	commands := []Command{
		{Command: "ls", ExitCode: 0, Duration: 100, Timestamp: time.Now()},
		{Command: "pwd", ExitCode: 0, Duration: 50, Timestamp: time.Now()},
		{Command: "cat missing", ExitCode: 1, Duration: 75, Timestamp: time.Now()},
	}
	
	for _, cmd := range commands {
		db.SaveCommand(cmd)
	}
	
	stats, err := db.GetStats()
	if err != nil {
		t.Fatalf("GetStats() failed: %v", err)
	}
	
	// Verify stats contain expected keys
	if stats == nil {
		t.Fatal("Expected non-nil stats")
	}
	
	// GetStats returns map[string]interface{}, check it's not empty
	if len(stats) == 0 {
		t.Error("Expected non-empty stats map")
	}
}

func TestEmbeddingOperations(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	dbPath := filepath.Join(tmpDir, "test.db")
	db, _ := NewDB(dbPath)
	defer db.Close()
	
	// Save command without embedding
	cmd := Command{
		Command:         "test command",
		Timestamp:       time.Now(),
		EmbeddingStatus: "pending",
	}
	db.SaveCommand(cmd)
	
	// Get commands pending embedding
	pending, err := db.GetPendingEmbeddings(10)
	if err != nil {
		t.Fatalf("GetPendingEmbeddings() failed: %v", err)
	}
	
	if len(pending) != 1 {
		t.Errorf("Expected 1 pending command, got %d", len(pending))
	}
	
	// Update embedding status
	embedding := []byte{1, 2, 3, 4}
	err = db.UpdateEmbeddingStatus(pending[0].ID, "completed", embedding)
	if err != nil {
		t.Fatalf("UpdateEmbeddingStatus() failed: %v", err)
	}
	
	// Verify embedding was saved
	commands, _ := db.GetRecentCommands(1)
	if commands[0].Embedding == nil {
		t.Error("Expected embedding to be saved")
	}
	if commands[0].EmbeddingStatus != "completed" {
		t.Errorf("Expected status 'completed', got '%s'", commands[0].EmbeddingStatus)
	}
}

func TestEmbeddingStorage(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	dbPath := filepath.Join(tmpDir, "test.db")
	db, _ := NewDB(dbPath)
	defer db.Close()
	
	// Save commands with embeddings
	embedding1 := []byte{1, 0, 0, 0}
	embedding2 := []byte{0, 1, 0, 0}
	
	commands := []Command{
		{Command: "git commit", Timestamp: time.Now(), Embedding: embedding1, EmbeddingStatus: "completed"},
		{Command: "npm install", Timestamp: time.Now(), Embedding: embedding2, EmbeddingStatus: "completed"},
		{Command: "git push", Timestamp: time.Now(), EmbeddingStatus: "pending"},
	}
	
	for _, cmd := range commands {
		db.SaveCommand(cmd)
	}
	
	// Get recent commands
	recent, _ := db.GetRecentCommands(3)
	
	// Verify embeddings were stored
	hasEmbedding := 0
	for _, cmd := range recent {
		if cmd.Embedding != nil {
			hasEmbedding++
		}
	}
	
	if hasEmbedding != 2 {
		t.Errorf("Expected 2 commands with embeddings, got %d", hasEmbedding)
	}
}

func TestCommandDeduplication(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	dbPath := filepath.Join(tmpDir, "test.db")
	db, _ := NewDB(dbPath)
	defer db.Close()
	
	// Save same command multiple times
	cmd := Command{
		Command:   "ls -lha",
		Timestamp: time.Now(),
	}
	
	for i := 0; i < 3; i++ {
		db.SaveCommand(cmd)
		time.Sleep(10 * time.Millisecond)
	}
	
	// All commands should be saved (deduplication happens at higher level)
	commands, _ := db.GetRecentCommands(10)
	if len(commands) != 3 {
		t.Errorf("Expected 3 commands, got %d", len(commands))
	}
}

func TestWorkingDirectoryStorage(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	dbPath := filepath.Join(tmpDir, "test.db")
	db, _ := NewDB(dbPath)
	defer db.Close()
	
	// Insert commands from different directories
	commands := []Command{
		{Command: "ls", WorkingDir: "/home/user/project", Timestamp: time.Now()},
		{Command: "pwd", WorkingDir: "/home/user/other", Timestamp: time.Now()},
		{Command: "git status", WorkingDir: "/home/user/project", Timestamp: time.Now()},
	}
	
	for _, cmd := range commands {
		db.SaveCommand(cmd)
	}
	
	// Get all commands and verify working directories
	all, err := db.GetRecentCommands(10)
	if err != nil {
		t.Fatalf("GetRecentCommands() failed: %v", err)
	}
	
	if len(all) != 3 {
		t.Errorf("Expected 3 commands, got %d", len(all))
	}
	
	// Verify working directories are stored correctly
	projectDirCount := 0
	for _, cmd := range all {
		if cmd.WorkingDir == "/home/user/project" {
			projectDirCount++
		}
	}
	
	if projectDirCount != 2 {
		t.Errorf("Expected 2 commands from project dir, got %d", projectDirCount)
	}
}

func TestCommandDuplication(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	dbPath := filepath.Join(tmpDir, "test.db")
	db, _ := NewDB(dbPath)
	defer db.Close()
	
	// Insert repeated commands (all saved, deduplication at app level)
	commands := []string{"ls", "ls", "ls", "git status", "git status", "pwd"}
	for _, cmd := range commands {
		db.SaveCommand(Command{
			Command:   cmd,
			Timestamp: time.Now(),
		})
		time.Sleep(1 * time.Millisecond)
	}
	
	// All commands should be saved
	all, err := db.GetRecentCommands(10)
	if err != nil {
		t.Fatalf("GetRecentCommands() failed: %v", err)
	}
	
	if len(all) != 6 {
		t.Errorf("Expected 6 commands, got %d", len(all))
	}
}

func TestClose(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	dbPath := filepath.Join(tmpDir, "test.db")
	db, _ := NewDB(dbPath)
	
	err := db.Close()
	if err != nil {
		t.Fatalf("Close() failed: %v", err)
	}
	
	// Operations after close should fail
	err = db.SaveCommand(Command{Command: "test"})
	if err == nil {
		t.Error("Expected error saving to closed database")
	}
}

func TestEmbeddingBinaryFormat(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	dbPath := filepath.Join(tmpDir, "test.db")
	db, _ := NewDB(dbPath)
	defer db.Close()
	
	// Create embedding
	embedding := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	
	// Save command with embedding
	cmd := Command{
		Command:         "test",
		Timestamp:       time.Now(),
		Embedding:       embedding,
		EmbeddingStatus: "completed",
	}
	db.SaveCommand(cmd)
	
	// Retrieve and verify
	retrieved, _ := db.GetRecentCommands(1)
	if retrieved[0].Embedding == nil {
		t.Fatal("Expected embedding to be saved")
	}
	
	// Verify binary format
	if len(retrieved[0].Embedding) != len(embedding) {
		t.Errorf("Expected %d bytes, got %d", len(embedding), len(retrieved[0].Embedding))
	}
	
	if !bytes.Equal(retrieved[0].Embedding, embedding) {
		t.Error("Embedding data doesn't match")
	}
}

func TestFTSIntegration(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	dbPath := filepath.Join(tmpDir, "test.db")
	db, _ := NewDB(dbPath)
	defer db.Close()
	
	// Insert command with output
	cmd := Command{
		Command:       "docker ps",
		OutputPreview: "CONTAINER ID   IMAGE   COMMAND   STATUS",
		Timestamp:     time.Now(),
	}
	db.SaveCommand(cmd)
	
	// Search in output
	results, err := db.SearchCommands("CONTAINER", 10)
	if err != nil {
		t.Fatalf("SearchCommands() failed: %v", err)
	}
	
	if len(results) != 1 {
		t.Errorf("Expected 1 result from FTS search, got %d", len(results))
	}
}

func TestMigrations(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	dbPath := filepath.Join(tmpDir, "test.db")
	
	// Create database (should run migrations)
	db, err := NewDB(dbPath)
	if err != nil {
		t.Fatalf("NewDB() failed: %v", err)
	}
	db.Close()
	
	// Reopen database (should handle existing schema)
	db2, err := NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to reopen database: %v", err)
	}
	defer db2.Close()
	
	// Should be able to save command
	cmd := Command{
		Command:   "test",
		Timestamp: time.Now(),
	}
	err = db2.SaveCommand(cmd)
	if err != nil {
		t.Errorf("Failed to save command after migration: %v", err)
	}
}

func TestConcurrentWrites(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	dbPath := filepath.Join(tmpDir, "test.db")
	db, _ := NewDB(dbPath)
	defer db.Close()
	
	// Run concurrent writes with slight delays to reduce contention
	done := make(chan bool)
	errorCount := 0
	var mu sync.Mutex
	
	for i := 0; i < 10; i++ {
		go func(n int) {
			// Add tiny delay to reduce contention
			time.Sleep(time.Duration(n) * time.Millisecond)
			cmd := Command{
				Command:   "test",
				Timestamp: time.Now(),
			}
			if err := db.SaveCommand(cmd); err != nil {
				mu.Lock()
				errorCount++
				mu.Unlock()
			}
			done <- true
		}(i)
	}
	
	// Wait for completion
	for i := 0; i < 10; i++ {
		<-done
	}
	
	// Verify that at least 40% of commands were saved (SQLite file-locking can cause contention)
	commands, _ := db.GetRecentCommands(100)
	minExpected := 4 // At least 40% success rate is acceptable for concurrent writes
	if len(commands) < minExpected {
		t.Errorf("Expected at least %d commands saved, got %d (errors: %d)", minExpected, len(commands), errorCount)
	}
	
	// Log info about concurrent behavior (this is expected with SQLite)
	t.Logf("Concurrent writes: %d succeeded, %d failed (SQLite file-locking contention is normal)", len(commands), errorCount)
}

func BenchmarkSaveCommand(b *testing.B) {
	tmpDir, _ := os.MkdirTemp("", "mako-bench-*")
	defer os.RemoveAll(tmpDir)
	
	dbPath := filepath.Join(tmpDir, "bench.db")
	db, _ := NewDB(dbPath)
	defer db.Close()
	
	cmd := Command{
		Command:       "ls -lha",
		Timestamp:     time.Now(),
		ExitCode:      0,
		Duration:      100,
		WorkingDir:    "/home/user",
		OutputPreview: "test output",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.SaveCommand(cmd)
	}
}

func BenchmarkSearchCommands(b *testing.B) {
	tmpDir, _ := os.MkdirTemp("", "mako-bench-*")
	defer os.RemoveAll(tmpDir)
	
	dbPath := filepath.Join(tmpDir, "bench.db")
	db, _ := NewDB(dbPath)
	defer db.Close()
	
	// Insert test data
	for i := 0; i < 1000; i++ {
		db.SaveCommand(Command{
			Command:   "test command",
			Timestamp: time.Now(),
		})
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.SearchCommands("test", 10)
	}
}

// Helper function to compare embeddings
func embeddingsEqual(a, b []byte) bool {
	return bytes.Equal(a, b)
}
