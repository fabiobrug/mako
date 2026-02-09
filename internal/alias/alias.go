package alias

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type AliasStore struct {
	Aliases map[string]string `json:"aliases"`
	path    string
}

// NewAliasStore creates or loads the alias store from ~/.mako/aliases.json
func NewAliasStore() (*AliasStore, error) {
	homeDir := os.Getenv("HOME")
	makoDir := filepath.Join(homeDir, ".mako")
	aliasPath := filepath.Join(makoDir, "aliases.json")

	// Ensure .mako directory exists
	if err := os.MkdirAll(makoDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create .mako directory: %w", err)
	}

	store := &AliasStore{
		Aliases: make(map[string]string),
		path:    aliasPath,
	}

	// Load existing aliases if file exists
	if _, err := os.Stat(aliasPath); err == nil {
		data, err := os.ReadFile(aliasPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read aliases file: %w", err)
		}

		if err := json.Unmarshal(data, store); err != nil {
			return nil, fmt.Errorf("failed to parse aliases file: %w", err)
		}
	}

	return store, nil
}

// Save writes the aliases to disk
func (s *AliasStore) Save() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal aliases: %w", err)
	}

	if err := os.WriteFile(s.path, data, 0644); err != nil {
		return fmt.Errorf("failed to write aliases file: %w", err)
	}

	return nil
}

// Set adds or updates an alias
func (s *AliasStore) Set(name, command string) error {
	if name == "" {
		return fmt.Errorf("alias name cannot be empty")
	}
	if command == "" {
		return fmt.Errorf("command cannot be empty")
	}

	s.Aliases[name] = command
	return s.Save()
}

// Get retrieves an alias by name
func (s *AliasStore) Get(name string) (string, bool) {
	cmd, ok := s.Aliases[name]
	return cmd, ok
}

// Delete removes an alias
func (s *AliasStore) Delete(name string) error {
	if _, ok := s.Aliases[name]; !ok {
		return fmt.Errorf("alias '%s' not found", name)
	}

	delete(s.Aliases, name)
	return s.Save()
}

// List returns all aliases
func (s *AliasStore) List() map[string]string {
	return s.Aliases
}
