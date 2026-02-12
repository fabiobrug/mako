package ai

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	preferencesFile = ".mako/preferences.json"
)

// CommandPreference stores learned preferences for a command
type CommandPreference struct {
	BaseCommand     string            `json:"base_command"`     // e.g., "ls", "git commit"
	CommonFlags     map[string]int    `json:"common_flags"`     // e.g., {"-lah": 10, "-la": 2}
	PreferredFlag   string            `json:"preferred_flag"`   // Most commonly used flag combo
	UsageCount      int               `json:"usage_count"`      // Total times this command was used
	LastUsed        string            `json:"last_used"`        // Timestamp of last use
}

// PersonalizationStore manages user preferences
type PersonalizationStore struct {
	Preferences map[string]*CommandPreference `json:"preferences"`
}

// LoadPreferences loads the personalization store from disk
func LoadPreferences() (*PersonalizationStore, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	prefPath := filepath.Join(homeDir, preferencesFile)
	
	data, err := os.ReadFile(prefPath)
	if err != nil {
		if os.IsNotExist(err) {
			// No preferences exist, start fresh
			return &PersonalizationStore{
				Preferences: make(map[string]*CommandPreference),
			}, nil
		}
		return nil, fmt.Errorf("failed to read preferences: %w", err)
	}

	var store PersonalizationStore
	if err := json.Unmarshal(data, &store); err != nil {
		// Corrupted file, start fresh
		return &PersonalizationStore{
			Preferences: make(map[string]*CommandPreference),
		}, nil
	}

	if store.Preferences == nil {
		store.Preferences = make(map[string]*CommandPreference)
	}

	return &store, nil
}

// Save saves the personalization store to disk
func (p *PersonalizationStore) Save() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	prefPath := filepath.Join(homeDir, preferencesFile)
	
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(prefPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal preferences: %w", err)
	}

	if err := os.WriteFile(prefPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write preferences: %w", err)
	}

	return nil
}

// LearnFromCommand learns user preferences from an executed command
func (p *PersonalizationStore) LearnFromCommand(command string) {
	// Extract base command and flags
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return
	}

	baseCmd := parts[0]
	
	// Handle git subcommands as separate commands
	if baseCmd == "git" && len(parts) > 1 {
		baseCmd = "git " + parts[1]
		parts = parts[1:]
	}

	// Extract flags (starts with -)
	var flags []string
	for i := 1; i < len(parts); i++ {
		if strings.HasPrefix(parts[i], "-") {
			flags = append(flags, parts[i])
		}
	}

	// Only learn if there are flags
	if len(flags) == 0 {
		return
	}

	// Get or create preference entry
	pref, exists := p.Preferences[baseCmd]
	if !exists {
		pref = &CommandPreference{
			BaseCommand: baseCmd,
			CommonFlags: make(map[string]int),
		}
		p.Preferences[baseCmd] = pref
	}

	// Update usage count
	pref.UsageCount++
	pref.LastUsed = fmt.Sprintf("%d", os.Getpid()) // Simple timestamp

	// Combine flags into a single string
	flagCombo := strings.Join(flags, " ")
	pref.CommonFlags[flagCombo]++

	// Update preferred flag
	maxCount := 0
	for flag, count := range pref.CommonFlags {
		if count > maxCount {
			maxCount = count
			pref.PreferredFlag = flag
		}
	}
}

// GetPreferenceHints returns AI-friendly hints about user preferences
func (p *PersonalizationStore) GetPreferenceHints() string {
	if len(p.Preferences) == 0 {
		return ""
	}

	var hints strings.Builder
	hints.WriteString("\nLEARNED USER PREFERENCES:\n")

	// Show top preferences (commands used more than twice)
	for cmd, pref := range p.Preferences {
		if pref.UsageCount >= 3 && pref.PreferredFlag != "" {
			hints.WriteString(fmt.Sprintf("- User typically uses '%s %s' (used %d times)\n", 
				cmd, pref.PreferredFlag, pref.CommonFlags[pref.PreferredFlag]))
		}
	}

	return hints.String()
}

// GetPreferenceForCommand returns the preferred flags for a specific command
func (p *PersonalizationStore) GetPreferenceForCommand(baseCommand string) string {
	pref, exists := p.Preferences[baseCommand]
	if !exists || pref.PreferredFlag == "" {
		return ""
	}

	// Only suggest if used at least 3 times
	if pref.UsageCount < 3 {
		return ""
	}

	return pref.PreferredFlag
}

// GetTopCommands returns the most frequently used commands
func (p *PersonalizationStore) GetTopCommands(limit int) []CommandPreference {
	// Convert map to slice
	prefs := make([]CommandPreference, 0, len(p.Preferences))
	for _, pref := range p.Preferences {
		prefs = append(prefs, *pref)
	}

	// Simple bubble sort by usage count (good enough for small datasets)
	for i := 0; i < len(prefs)-1; i++ {
		for j := 0; j < len(prefs)-i-1; j++ {
			if prefs[j].UsageCount < prefs[j+1].UsageCount {
				prefs[j], prefs[j+1] = prefs[j+1], prefs[j]
			}
		}
	}

	// Return top N
	if limit > len(prefs) {
		limit = len(prefs)
	}

	return prefs[:limit]
}
