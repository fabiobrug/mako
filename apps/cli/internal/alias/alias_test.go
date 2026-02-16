package alias

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fabiobrug/mako.git/internal/testutil"
)

func TestNewAliasStore(t *testing.T) {
	tmpHome := testutil.MockHomeDir(t)
	
	store, err := NewAliasStore()
	if err != nil {
		t.Fatalf("NewAliasStore() failed: %v", err)
	}
	
	if store == nil {
		t.Fatal("Expected non-nil store")
	}
	if store.Aliases == nil {
		t.Fatal("Expected initialized aliases map")
	}
	if len(store.Aliases) != 0 {
		t.Errorf("Expected empty aliases, got %d", len(store.Aliases))
	}
	
	// Verify .mako directory was created
	makoDir := filepath.Join(tmpHome, ".mako")
	if _, err := os.Stat(makoDir); os.IsNotExist(err) {
		t.Error(".mako directory was not created")
	}
}

func TestAliasSetAndGet(t *testing.T) {
	testutil.MockHomeDir(t)
	
	store, err := NewAliasStore()
	if err != nil {
		t.Fatalf("NewAliasStore() failed: %v", err)
	}
	
	// Set alias
	err = store.Set("ll", "ls -lha", []string{"list", "files"})
	if err != nil {
		t.Fatalf("Set() failed: %v", err)
	}
	
	// Get alias
	cmd, ok := store.Get("ll")
	if !ok {
		t.Error("Expected alias to exist")
	}
	if cmd != "ls -lha" {
		t.Errorf("Expected 'ls -lha', got '%s'", cmd)
	}
	
	// Get full info
	info, ok := store.GetInfo("ll")
	if !ok {
		t.Error("Expected alias info to exist")
	}
	if info.Command != "ls -lha" {
		t.Errorf("Expected 'ls -lha', got '%s'", info.Command)
	}
	if len(info.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(info.Tags))
	}
}

func TestAliasSetValidation(t *testing.T) {
	testutil.MockHomeDir(t)
	
	store, _ := NewAliasStore()
	
	tests := []struct {
		name    string
		alias   string
		command string
		wantErr bool
	}{
		{"Valid alias", "ll", "ls -lha", false},
		{"Empty name", "", "ls -lha", true},
		{"Empty command", "ll", "", true},
		{"Valid with args", "gs", "git status", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.Set(tt.alias, tt.command, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("Set(%s, %s) error = %v, wantErr %v", 
					tt.alias, tt.command, err, tt.wantErr)
			}
		})
	}
}

func TestAliasDelete(t *testing.T) {
	testutil.MockHomeDir(t)
	
	store, _ := NewAliasStore()
	store.Set("ll", "ls -lha", nil)
	
	// Delete existing alias
	err := store.Delete("ll")
	if err != nil {
		t.Fatalf("Delete() failed: %v", err)
	}
	
	// Verify it's gone
	_, ok := store.Get("ll")
	if ok {
		t.Error("Alias should have been deleted")
	}
	
	// Delete non-existent alias
	err = store.Delete("nonexistent")
	if err == nil {
		t.Error("Expected error deleting non-existent alias")
	}
}

func TestAliasList(t *testing.T) {
	testutil.MockHomeDir(t)
	
	store, _ := NewAliasStore()
	store.Set("ll", "ls -lha", []string{"list"})
	store.Set("gs", "git status", []string{"git"})
	store.Set("gp", "git push", []string{"git"})
	
	aliases := store.List()
	if len(aliases) != 3 {
		t.Errorf("Expected 3 aliases, got %d", len(aliases))
	}
	
	if _, ok := aliases["ll"]; !ok {
		t.Error("Expected 'll' alias")
	}
	if _, ok := aliases["gs"]; !ok {
		t.Error("Expected 'gs' alias")
	}
}

func TestAliasListByTag(t *testing.T) {
	testutil.MockHomeDir(t)
	
	store, _ := NewAliasStore()
	store.Set("ll", "ls -lha", []string{"list", "files"})
	store.Set("gs", "git status", []string{"git"})
	store.Set("gp", "git push", []string{"git", "remote"})
	store.Set("du", "du -sh", []string{"disk", "files"})
	
	// Filter by git tag
	gitAliases := store.ListByTag("git")
	if len(gitAliases) != 2 {
		t.Errorf("Expected 2 git aliases, got %d", len(gitAliases))
	}
	
	// Filter by files tag
	fileAliases := store.ListByTag("files")
	if len(fileAliases) != 2 {
		t.Errorf("Expected 2 file aliases, got %d", len(fileAliases))
	}
	
	// Non-existent tag
	noneAliases := store.ListByTag("nonexistent")
	if len(noneAliases) != 0 {
		t.Errorf("Expected 0 aliases for non-existent tag, got %d", len(noneAliases))
	}
}

func TestGetAllTags(t *testing.T) {
	testutil.MockHomeDir(t)
	
	store, _ := NewAliasStore()
	store.Set("ll", "ls -lha", []string{"list", "files"})
	store.Set("gs", "git status", []string{"git"})
	store.Set("gp", "git push", []string{"git", "remote"})
	
	tags := store.GetAllTags()
	expectedTags := map[string]bool{
		"list":   true,
		"files":  true,
		"git":    true,
		"remote": true,
	}
	
	if len(tags) != len(expectedTags) {
		t.Errorf("Expected %d unique tags, got %d", len(expectedTags), len(tags))
	}
	
	for _, tag := range tags {
		if !expectedTags[tag] {
			t.Errorf("Unexpected tag: %s", tag)
		}
	}
}

func TestExpandParameters(t *testing.T) {
	tests := []struct {
		command string
		args    []string
		want    string
	}{
		{
			"echo $1",
			[]string{"hello"},
			"echo hello",
		},
		{
			"echo $1 $2",
			[]string{"hello", "world"},
			"echo hello world",
		},
		{
			"git commit -m '$1'",
			[]string{"Initial commit"},
			"git commit -m 'Initial commit'",
		},
		{
			"echo $@",
			[]string{"hello", "world", "test"},
			"echo hello world test",
		},
		{
			"echo Count: $#",
			[]string{"a", "b", "c"},
			"echo Count: 3",
		},
		{
			"cp $1 $2",
			[]string{"source.txt", "dest.txt"},
			"cp source.txt dest.txt",
		},
		{
			"docker run --name $1 -p $2:80 nginx",
			[]string{"myapp", "8080"},
			"docker run --name myapp -p 8080:80 nginx",
		},
		{
			"no parameters here",
			[]string{"arg1", "arg2"},
			"no parameters here",
		},
		{
			"echo $1 $2 $3",
			[]string{"one"},
			"echo one $2 $3", // Unreplaced params stay
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.command, func(t *testing.T) {
			got := ExpandParameters(tt.command, tt.args)
			if got != tt.want {
				t.Errorf("ExpandParameters(%q, %v) = %q, want %q",
					tt.command, tt.args, got, tt.want)
			}
		})
	}
}

func TestAliasPersistence(t *testing.T) {
	tmpHome := testutil.MockHomeDir(t)
	
	// Create and save aliases
	store1, _ := NewAliasStore()
	store1.Set("ll", "ls -lha", []string{"list"})
	store1.Set("gs", "git status", []string{"git"})
	
	// Load in new store
	store2, err := NewAliasStore()
	if err != nil {
		t.Fatalf("Failed to load store: %v", err)
	}
	
	// Verify aliases persisted
	if len(store2.Aliases) != 2 {
		t.Errorf("Expected 2 persisted aliases, got %d", len(store2.Aliases))
	}
	
	cmd, ok := store2.Get("ll")
	if !ok || cmd != "ls -lha" {
		t.Error("Alias 'll' not persisted correctly")
	}
	
	_ = tmpHome
}

func TestAliasExportImport(t *testing.T) {
	tmpHome := testutil.MockHomeDir(t)
	tmpDir := testutil.TempDir(t)
	
	// Create aliases
	store, _ := NewAliasStore()
	store.Set("ll", "ls -lha", []string{"list"})
	store.Set("gs", "git status", []string{"git"})
	
	// Export
	exportPath := filepath.Join(tmpDir, "aliases-export.json")
	err := store.ExportToFile(exportPath)
	if err != nil {
		t.Fatalf("ExportToFile() failed: %v", err)
	}
	
	// Verify file exists
	if _, err := os.Stat(exportPath); os.IsNotExist(err) {
		t.Error("Export file was not created")
	}
	
	// Create new store and import
	store2, _ := NewAliasStore()
	err = store2.ImportFromFile(exportPath)
	if err != nil {
		t.Fatalf("ImportFromFile() failed: %v", err)
	}
	
	// Verify imported aliases
	if len(store2.Aliases) != 2 {
		t.Errorf("Expected 2 imported aliases, got %d", len(store2.Aliases))
	}
	
	cmd, ok := store2.Get("ll")
	if !ok || cmd != "ls -lha" {
		t.Error("Alias 'll' not imported correctly")
	}
	
	_ = tmpHome
}

func TestAliasImportFromReader(t *testing.T) {
	testutil.MockHomeDir(t)
	
	store, _ := NewAliasStore()
	
	// Create JSON input
	jsonData := `{
		"ll": {"command": "ls -lha", "tags": ["list"]},
		"gs": {"command": "git status", "tags": ["git"]}
	}`
	
	reader := strings.NewReader(jsonData)
	err := store.ImportFromReader(reader)
	if err != nil {
		t.Fatalf("ImportFromReader() failed: %v", err)
	}
	
	// Verify imported aliases
	if len(store.Aliases) != 2 {
		t.Errorf("Expected 2 imported aliases, got %d", len(store.Aliases))
	}
	
	info, ok := store.GetInfo("ll")
	if !ok {
		t.Error("Expected 'll' alias")
	}
	if info.Command != "ls -lha" {
		t.Errorf("Expected 'ls -lha', got '%s'", info.Command)
	}
	if len(info.Tags) != 1 || info.Tags[0] != "list" {
		t.Errorf("Expected tags ['list'], got %v", info.Tags)
	}
}

func TestAliasBackwardCompatibility(t *testing.T) {
	tmpHome := testutil.MockHomeDir(t)
	makoDir := filepath.Join(tmpHome, ".mako")
	os.MkdirAll(makoDir, 0755)
	
	// Write old format (without tags)
	oldFormat := `{
		"aliases": {
			"ll": "ls -lha",
			"gs": "git status"
		}
	}`
	
	aliasPath := filepath.Join(makoDir, "aliases.json")
	os.WriteFile(aliasPath, []byte(oldFormat), 0644)
	
	// Load with new code
	store, err := NewAliasStore()
	if err != nil {
		t.Fatalf("NewAliasStore() failed: %v", err)
	}
	
	// Verify aliases loaded
	if len(store.Aliases) != 2 {
		t.Errorf("Expected 2 aliases from old format, got %d", len(store.Aliases))
	}
	
	info, ok := store.GetInfo("ll")
	if !ok {
		t.Error("Expected 'll' alias")
	}
	if info.Command != "ls -lha" {
		t.Errorf("Expected 'ls -lha', got '%s'", info.Command)
	}
	// Tags should be empty for old format
	if len(info.Tags) != 0 {
		t.Errorf("Expected empty tags for old format, got %v", info.Tags)
	}
}

func TestAliasUpdate(t *testing.T) {
	testutil.MockHomeDir(t)
	
	store, _ := NewAliasStore()
	
	// Set initial alias
	store.Set("ll", "ls -l", []string{"list"})
	
	// Update alias
	store.Set("ll", "ls -lha", []string{"list", "all"})
	
	// Verify update
	info, ok := store.GetInfo("ll")
	if !ok {
		t.Error("Expected alias to exist")
	}
	if info.Command != "ls -lha" {
		t.Errorf("Expected updated command 'ls -lha', got '%s'", info.Command)
	}
	if len(info.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(info.Tags))
	}
}

func BenchmarkExpandParameters(b *testing.B) {
	command := "docker run --name $1 -p $2:80 -v $3:/app -e ENV=$4 nginx"
	args := []string{"myapp", "8080", "/data", "production"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ExpandParameters(command, args)
	}
}

func BenchmarkAliasGet(b *testing.B) {
	tmpHome, _ := os.MkdirTemp("", "mako-bench-*")
	defer os.RemoveAll(tmpHome)
	os.Setenv("HOME", tmpHome)
	defer os.Unsetenv("HOME")
	
	store, _ := NewAliasStore()
	store.Set("ll", "ls -lha", nil)
	store.Set("gs", "git status", nil)
	store.Set("gp", "git push", nil)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.Get("gs")
	}
}
