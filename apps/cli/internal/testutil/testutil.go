package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

// TempDir creates a temporary directory for tests
func TempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "mako-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() {
		os.RemoveAll(dir)
	})
	return dir
}

// TempFile creates a temporary file with content
func TempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}

// SetEnv sets an environment variable for the duration of the test
func SetEnv(t *testing.T, key, value string) {
	t.Helper()
	old := os.Getenv(key)
	os.Setenv(key, value)
	t.Cleanup(func() {
		if old != "" {
			os.Setenv(key, old)
		} else {
			os.Unsetenv(key)
		}
	})
}

// MockHomeDir sets HOME environment variable to a temp directory
func MockHomeDir(t *testing.T) string {
	t.Helper()
	tmpHome := TempDir(t)
	SetEnv(t, "HOME", tmpHome)
	return tmpHome
}
