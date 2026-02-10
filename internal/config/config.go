package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
)

// Config represents Mako's configuration
type Config struct {
	Version            string `json:"version"`
	APIKey             string `json:"api_key,omitempty"`
	Theme              string `json:"theme"`
	CacheSize          int    `json:"cache_size"`
	Telemetry          bool   `json:"telemetry"`
	AutoUpdate         bool   `json:"auto_update"`
	HistoryLimit       int    `json:"history_limit"`
	SafetyLevel        string `json:"safety_level"`
	EmbeddingBatchSize int    `json:"embedding_batch_size"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Version:            "1.0",
		Theme:              "ocean",
		CacheSize:          10000,
		Telemetry:          false,
		AutoUpdate:         true,
		HistoryLimit:       100000,
		SafetyLevel:        "medium",
		EmbeddingBatchSize: 10,
	}
}

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	home := os.Getenv("HOME")
	return filepath.Join(home, ".mako", "config.json")
}

// GetMakoDir returns the Mako directory path
func GetMakoDir() string {
	home := os.Getenv("HOME")
	return filepath.Join(home, ".mako")
}

// LoadConfig loads configuration from file, creating defaults if not exists
func LoadConfig() (*Config, error) {
	configPath := GetConfigPath()

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config
		config := DefaultConfig()
		
		// Try to load API key from .env if it exists
		envPath := filepath.Join(GetMakoDir(), ".env")
		if data, err := os.ReadFile(envPath); err == nil {
			// Simple .env parsing for GEMINI_API_KEY
			lines := string(data)
			for _, line := range []string{lines} {
				if len(line) > 16 && line[:16] == "GEMINI_API_KEY=" {
					config.APIKey = line[16:]
					break
				}
			}
		}
		
		// Also check environment variable
		if envKey := os.Getenv("GEMINI_API_KEY"); envKey != "" && config.APIKey == "" {
			config.APIKey = envKey
		}
		
		return config, nil
	}

	// Read existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Merge with defaults for missing fields
	defaults := DefaultConfig()
	mergeDefaults(&config, defaults)

	return &config, nil
}

// mergeDefaults fills in missing fields from defaults
func mergeDefaults(config, defaults *Config) {
	configVal := reflect.ValueOf(config).Elem()
	defaultsVal := reflect.ValueOf(defaults).Elem()

	for i := 0; i < configVal.NumField(); i++ {
		field := configVal.Field(i)
		defaultField := defaultsVal.Field(i)

		// Skip if field is not zero value
		if !field.IsZero() {
			continue
		}

		// Set to default value
		if field.CanSet() {
			field.Set(defaultField)
		}
	}
}

// Save writes the configuration to disk
func (c *Config) Save() error {
	// Ensure .mako directory exists
	makoDir := GetMakoDir()
	if err := os.MkdirAll(makoDir, 0755); err != nil {
		return fmt.Errorf("failed to create .mako directory: %w", err)
	}

	// Marshal to JSON with pretty print
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	configPath := GetConfigPath()
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// Get retrieves a configuration value by key
func (c *Config) Get(key string) (interface{}, error) {
	configVal := reflect.ValueOf(c).Elem()
	configType := configVal.Type()

	for i := 0; i < configVal.NumField(); i++ {
		field := configType.Field(i)
		jsonTag := field.Tag.Get("json")
		
		// Remove ,omitempty suffix if present
		if len(jsonTag) > 10 && jsonTag[len(jsonTag)-10:] == ",omitempty" {
			jsonTag = jsonTag[:len(jsonTag)-10]
		}

		if jsonTag == key {
			return configVal.Field(i).Interface(), nil
		}
	}

	return nil, fmt.Errorf("unknown config key: %s", key)
}

// Set updates a configuration value by key
func (c *Config) Set(key string, value interface{}) error {
	configVal := reflect.ValueOf(c).Elem()
	configType := configVal.Type()

	for i := 0; i < configVal.NumField(); i++ {
		field := configType.Field(i)
		jsonTag := field.Tag.Get("json")
		
		// Remove ,omitempty suffix if present
		if len(jsonTag) > 10 && jsonTag[len(jsonTag)-10:] == ",omitempty" {
			jsonTag = jsonTag[:len(jsonTag)-10]
		}

		if jsonTag == key {
			fieldVal := configVal.Field(i)
			if !fieldVal.CanSet() {
				return fmt.Errorf("cannot set field: %s", key)
			}

			// Convert value to appropriate type
			val := reflect.ValueOf(value)
			if val.Type().ConvertibleTo(fieldVal.Type()) {
				fieldVal.Set(val.Convert(fieldVal.Type()))
				return nil
			}

			return fmt.Errorf("invalid type for %s: expected %s, got %s", 
				key, fieldVal.Type(), val.Type())
		}
	}

	return fmt.Errorf("unknown config key: %s", key)
}

// Reset restores default configuration
func (c *Config) Reset() error {
	defaults := DefaultConfig()
	*c = *defaults
	return c.Save()
}

// List returns all configuration as a map
func (c *Config) List() map[string]interface{} {
	result := make(map[string]interface{})
	configVal := reflect.ValueOf(c).Elem()
	configType := configVal.Type()

	for i := 0; i < configVal.NumField(); i++ {
		field := configType.Field(i)
		jsonTag := field.Tag.Get("json")
		
		// Remove ,omitempty suffix if present
		if len(jsonTag) > 10 && jsonTag[len(jsonTag)-10:] == ",omitempty" {
			jsonTag = jsonTag[:len(jsonTag)-10]
		}

		if jsonTag != "" && jsonTag != "-" {
			result[jsonTag] = configVal.Field(i).Interface()
		}
	}

	return result
}
