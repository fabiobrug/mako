package alias

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type AliasInfo struct {
	Command string   `json:"command"`
	Tags    []string `json:"tags,omitempty"`
}

type AliasStore struct {
	Aliases map[string]AliasInfo `json:"aliases"`
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
		Aliases: make(map[string]AliasInfo),
		path:    aliasPath,
	}

	// Load existing aliases if file exists
	if _, err := os.Stat(aliasPath); err == nil {
		data, err := os.ReadFile(aliasPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read aliases file: %w", err)
		}

		// Try new format first
		var newFormat struct {
			Aliases map[string]AliasInfo `json:"aliases"`
		}
		if err := json.Unmarshal(data, &newFormat); err == nil && len(newFormat.Aliases) > 0 {
			store.Aliases = newFormat.Aliases
		} else {
			// Try old format (backward compatibility)
			var oldFormat struct {
				Aliases map[string]string `json:"aliases"`
			}
			if err := json.Unmarshal(data, &oldFormat); err == nil {
				// Convert old format to new
				for name, command := range oldFormat.Aliases {
					store.Aliases[name] = AliasInfo{Command: command, Tags: []string{}}
				}
			} else {
				return nil, fmt.Errorf("failed to parse aliases file: %w", err)
			}
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
func (s *AliasStore) Set(name, command string, tags []string) error {
	if name == "" {
		return fmt.Errorf("alias name cannot be empty")
	}
	if command == "" {
		return fmt.Errorf("command cannot be empty")
	}

	s.Aliases[name] = AliasInfo{Command: command, Tags: tags}
	return s.Save()
}

// Get retrieves an alias command by name
func (s *AliasStore) Get(name string) (string, bool) {
	info, ok := s.Aliases[name]
	if !ok {
		return "", false
	}
	return info.Command, true
}

// GetInfo retrieves full alias info by name
func (s *AliasStore) GetInfo(name string) (AliasInfo, bool) {
	info, ok := s.Aliases[name]
	return info, ok
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
func (s *AliasStore) List() map[string]AliasInfo {
	return s.Aliases
}

// ListByTag returns aliases filtered by tag
func (s *AliasStore) ListByTag(tag string) map[string]AliasInfo {
	filtered := make(map[string]AliasInfo)
	for name, info := range s.Aliases {
		for _, t := range info.Tags {
			if t == tag {
				filtered[name] = info
				break
			}
		}
	}
	return filtered
}

// GetAllTags returns all unique tags used
func (s *AliasStore) GetAllTags() []string {
	tagSet := make(map[string]bool)
	for _, info := range s.Aliases {
		for _, tag := range info.Tags {
			tagSet[tag] = true
		}
	}
	
	tags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		tags = append(tags, tag)
	}
	return tags
}

// ExpandParameters replaces $1, $2, ... $n with actual arguments
func ExpandParameters(command string, args []string) string {
	result := command
	
	// Replace numbered parameters
	for i, arg := range args {
		placeholder := fmt.Sprintf("$%d", i+1)
		result = strings.ReplaceAll(result, placeholder, arg)
	}
	
	// Replace $@ with all arguments
	result = strings.ReplaceAll(result, "$@", strings.Join(args, " "))
	
	// Replace $# with argument count
	result = strings.ReplaceAll(result, "$#", strconv.Itoa(len(args)))
	
	return result
}

// ExportToFile exports aliases to a specified file path
func (s *AliasStore) ExportToFile(filepath string) error {
	data, err := json.MarshalIndent(s.Aliases, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal aliases: %w", err)
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write export file: %w", err)
	}

	return nil
}

// ImportFromFile imports aliases from a specified file path
func (s *AliasStore) ImportFromFile(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read import file: %w", err)
	}

	// Try new format first
	var imported map[string]AliasInfo
	if err := json.Unmarshal(data, &imported); err == nil {
		// New format with tags
		for name, info := range imported {
			s.Aliases[name] = info
		}
	} else {
		// Try old format (backward compatibility)
		var oldImported map[string]string
		if err := json.Unmarshal(data, &oldImported); err != nil {
			return fmt.Errorf("failed to parse import file: %w", err)
		}
		for name, command := range oldImported {
			s.Aliases[name] = AliasInfo{Command: command, Tags: []string{}}
		}
	}

	return s.Save()
}

// ImportFromReader imports aliases from a reader (for stdin support)
func (s *AliasStore) ImportFromReader(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	// Try new format first
	var imported map[string]AliasInfo
	if err := json.Unmarshal(data, &imported); err == nil {
		// New format with tags
		for name, info := range imported {
			s.Aliases[name] = info
		}
	} else {
		// Try old format
		var oldImported map[string]string
		if err := json.Unmarshal(data, &oldImported); err != nil {
			return fmt.Errorf("failed to parse input: %w", err)
		}
		for name, command := range oldImported {
			s.Aliases[name] = AliasInfo{Command: command, Tags: []string{}}
		}
	}

	return s.Save()
}
