package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/fabiobrug/mako.git/internal/testutil"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	
	if config.Version != "1.0" {
		t.Errorf("Expected version 1.0, got %s", config.Version)
	}
	if config.LLMProvider != "gemini" {
		t.Errorf("Expected default provider gemini, got %s", config.LLMProvider)
	}
	if config.Theme != "ocean" {
		t.Errorf("Expected theme ocean, got %s", config.Theme)
	}
	if config.CacheSize != 10000 {
		t.Errorf("Expected cache size 10000, got %d", config.CacheSize)
	}
	if config.Telemetry != false {
		t.Error("Expected telemetry to be false by default")
	}
	if config.AutoUpdate != true {
		t.Error("Expected auto_update to be true by default")
	}
	if config.HistoryLimit != 100000 {
		t.Errorf("Expected history limit 100000, got %d", config.HistoryLimit)
	}
	if config.SafetyLevel != "medium" {
		t.Errorf("Expected safety level medium, got %s", config.SafetyLevel)
	}
}

func TestLoadConfigNotExists(t *testing.T) {
	tmpHome := testutil.MockHomeDir(t)
	
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}
	
	// Should return default config when file doesn't exist
	if config.LLMProvider != "gemini" {
		t.Errorf("Expected default provider, got %s", config.LLMProvider)
	}
	
	// Config file should not be created by LoadConfig
	configPath := filepath.Join(tmpHome, ".mako", "config.json")
	if _, err := os.Stat(configPath); err == nil {
		t.Error("Config file should not be created by LoadConfig")
	}
}

func TestLoadConfigWithEnvVar(t *testing.T) {
	tmpHome := testutil.MockHomeDir(t)
	testutil.SetEnv(t, "GEMINI_API_KEY", "test-api-key-123")
	
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}
	
	if config.APIKey != "test-api-key-123" {
		t.Errorf("Expected API key from env, got %s", config.APIKey)
	}
	
	_ = tmpHome // Use variable
}

func TestLoadConfigFromFile(t *testing.T) {
	tmpHome := testutil.MockHomeDir(t)
	makoDir := filepath.Join(tmpHome, ".mako")
	os.MkdirAll(makoDir, 0755)
	
	// Create config file
	configData := map[string]interface{}{
		"version":      "1.0",
		"llm_provider": "openai",
		"llm_model":    "gpt-4",
		"theme":        "dark",
		"cache_size":   5000,
		"telemetry":    true,
	}
	data, _ := json.Marshal(configData)
	configPath := filepath.Join(makoDir, "config.json")
	os.WriteFile(configPath, data, 0644)
	
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}
	
	if config.LLMProvider != "openai" {
		t.Errorf("Expected provider openai, got %s", config.LLMProvider)
	}
	if config.LLMModel != "gpt-4" {
		t.Errorf("Expected model gpt-4, got %s", config.LLMModel)
	}
	if config.Theme != "dark" {
		t.Errorf("Expected theme dark, got %s", config.Theme)
	}
	if config.CacheSize != 5000 {
		t.Errorf("Expected cache size 5000, got %d", config.CacheSize)
	}
	if config.Telemetry != true {
		t.Error("Expected telemetry to be true")
	}
}

func TestLoadConfigMergesDefaults(t *testing.T) {
	tmpHome := testutil.MockHomeDir(t)
	makoDir := filepath.Join(tmpHome, ".mako")
	os.MkdirAll(makoDir, 0755)
	
	// Create partial config (missing some fields)
	configData := map[string]interface{}{
		"llm_provider": "anthropic",
		"theme":        "light",
	}
	data, _ := json.Marshal(configData)
	configPath := filepath.Join(makoDir, "config.json")
	os.WriteFile(configPath, data, 0644)
	
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}
	
	// Should have user values
	if config.LLMProvider != "anthropic" {
		t.Errorf("Expected provider anthropic, got %s", config.LLMProvider)
	}
	if config.Theme != "light" {
		t.Errorf("Expected theme light, got %s", config.Theme)
	}
	
	// Should have default values for missing fields
	if config.CacheSize != 10000 {
		t.Errorf("Expected default cache size 10000, got %d", config.CacheSize)
	}
	if config.AutoUpdate != true {
		t.Error("Expected default auto_update to be true")
	}
}

func TestSaveConfig(t *testing.T) {
	tmpHome := testutil.MockHomeDir(t)
	makoDir := filepath.Join(tmpHome, ".mako")
	
	config := &Config{
		Version:      "1.0",
		LLMProvider:  "openai",
		LLMModel:     "gpt-4-turbo",
		Theme:        "dark",
		CacheSize:    20000,
		Telemetry:    false,
		AutoUpdate:   true,
		HistoryLimit: 50000,
		SafetyLevel:  "high",
	}
	
	err := config.Save()
	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}
	
	// Verify file was created
	configPath := filepath.Join(makoDir, "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}
	
	// Load and verify content
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}
	
	var loaded Config
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("Failed to parse config file: %v", err)
	}
	
	if loaded.LLMProvider != "openai" {
		t.Errorf("Expected provider openai, got %s", loaded.LLMProvider)
	}
	if loaded.CacheSize != 20000 {
		t.Errorf("Expected cache size 20000, got %d", loaded.CacheSize)
	}
	if loaded.SafetyLevel != "high" {
		t.Errorf("Expected safety level high, got %s", loaded.SafetyLevel)
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	tmpHome := testutil.MockHomeDir(t)
	_ = tmpHome
	
	// Create config
	config := DefaultConfig()
	config.LLMProvider = "anthropic"
	config.Theme = "dark"
	config.CacheSize = 5000
	
	// Save
	err := config.Save()
	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}
	
	// Load and verify
	loaded, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}
	
	if loaded.LLMProvider != "anthropic" {
		t.Errorf("Expected anthropic, got %s", loaded.LLMProvider)
	}
	if loaded.Theme != "dark" {
		t.Errorf("Expected dark, got %s", loaded.Theme)
	}
	if loaded.CacheSize != 5000 {
		t.Errorf("Expected 5000, got %d", loaded.CacheSize)
	}
}

func TestConfigFields(t *testing.T) {
	config := DefaultConfig()
	
	// Test field access
	if config.LLMProvider != "gemini" {
		t.Errorf("Expected gemini, got %s", config.LLMProvider)
	}
	
	// Modify and save
	config.LLMProvider = "openai"
	config.Theme = "dark"
	
	if config.LLMProvider != "openai" {
		t.Error("Field modification failed")
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name:    "Valid config",
			config:  DefaultConfig(),
			wantErr: false,
		},
		{
			name: "Invalid cache size",
			config: &Config{
				CacheSize: -100,
			},
			wantErr: false, // Currently no validation
		},
		{
			name: "Empty provider",
			config: &Config{
				LLMProvider: "",
			},
			wantErr: false, // Currently no validation
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpHome := testutil.MockHomeDir(t)
			_ = tmpHome
			
			err := tt.config.Save()
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func BenchmarkLoadConfig(b *testing.B) {
	tmpHome, _ := os.MkdirTemp("", "mako-bench-*")
	defer os.RemoveAll(tmpHome)
	os.Setenv("HOME", tmpHome)
	defer os.Unsetenv("HOME")
	
	config := DefaultConfig()
	config.Save()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LoadConfig()
	}
}

func BenchmarkSaveConfig(b *testing.B) {
	tmpHome, _ := os.MkdirTemp("", "mako-bench-*")
	defer os.RemoveAll(tmpHome)
	os.Setenv("HOME", tmpHome)
	defer os.Unsetenv("HOME")
	
	config := DefaultConfig()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config.Save()
	}
}
