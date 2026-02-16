package shell

import (
	"fmt"
	"os"
	"strings"

	"github.com/fabiobrug/mako.git/internal/config"
)

// handleConfig handles the 'mako config' command
func handleConfig(args []string) (string, error) {
	cyan := "\033[38;2;0;209;255m"
	lightBlue := "\033[38;2;93;173;226m"
	dimBlue := "\033[38;2;120;150;180m"
	reset := "\033[0m"

	cfg, err := config.LoadConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	if len(args) == 0 {
		return fmt.Sprintf("Usage: mako config <list|get|set|reset>\r\n"), nil
	}

	switch args[0] {
	case "list":
		output := fmt.Sprintf("\r\n%sMako Configuration%s\r\n", cyan, reset)
		output += fmt.Sprintf("%s━━━━━━━━━━━━━━━━━━━━━━%s\r\n\r\n", dimBlue, reset)
		
		// File-based configuration
		output += fmt.Sprintf("%sFile Config (~/.mako/config.json):%s\r\n", dimBlue, reset)
		settings := cfg.List()
		for key, value := range settings {
			// Hide API key (show only first 10 chars)
			if key == "api_key" {
				if str, ok := value.(string); ok && len(str) > 10 {
					value = str[:10] + "..."
				} else if str, ok := value.(string); ok && str == "" {
					value = "(not set)"
				}
			}
			output += fmt.Sprintf("  %s%-20s%s %v\r\n", lightBlue, key, reset, value)
		}
		
		// Environment variables for LLM provider
		output += fmt.Sprintf("\r\n%sLLM Provider (from environment):%s\r\n", dimBlue, reset)
		llmProvider := os.Getenv("LLM_PROVIDER")
		llmModel := os.Getenv("LLM_MODEL")
		llmAPIKey := os.Getenv("LLM_API_KEY")
		llmBaseURL := os.Getenv("LLM_API_BASE")
		
		if llmProvider == "" {
			llmProvider = "(using config file or default: gemini)"
		}
		if llmModel == "" {
			llmModel = "(provider default)"
		}
		if llmAPIKey == "" {
			llmAPIKey = "(not set)"
		} else if len(llmAPIKey) > 10 {
			llmAPIKey = llmAPIKey[:10] + "..."
		}
		if llmBaseURL == "" {
			llmBaseURL = "(provider default)"
		}
		
		output += fmt.Sprintf("  %s%-20s%s %v\r\n", lightBlue, "provider", reset, llmProvider)
		output += fmt.Sprintf("  %s%-20s%s %v\r\n", lightBlue, "model", reset, llmModel)
		output += fmt.Sprintf("  %s%-20s%s %v\r\n", lightBlue, "api_key", reset, llmAPIKey)
		output += fmt.Sprintf("  %s%-20s%s %v\r\n", lightBlue, "base_url", reset, llmBaseURL)
		
		// Environment variables for Embedding provider
		output += fmt.Sprintf("\r\n%sEmbedding Provider (from environment):%s\r\n", dimBlue, reset)
		embProvider := os.Getenv("EMBEDDING_PROVIDER")
		embModel := os.Getenv("EMBEDDING_MODEL")
		embAPIKey := os.Getenv("EMBEDDING_API_KEY")
		embBaseURL := os.Getenv("EMBEDDING_API_BASE")
		
		if embProvider == "" {
			embProvider = "(using LLM provider)"
		}
		if embModel == "" {
			embModel = "(provider default)"
		}
		if embAPIKey == "" {
			embAPIKey = "(using LLM API key)"
		} else if len(embAPIKey) > 10 {
			embAPIKey = embAPIKey[:10] + "..."
		}
		if embBaseURL == "" {
			embBaseURL = "(provider default)"
		}
		
		output += fmt.Sprintf("  %s%-20s%s %v\r\n", lightBlue, "provider", reset, embProvider)
		output += fmt.Sprintf("  %s%-20s%s %v\r\n", lightBlue, "model", reset, embModel)
		output += fmt.Sprintf("  %s%-20s%s %v\r\n", lightBlue, "api_key", reset, embAPIKey)
		output += fmt.Sprintf("  %s%-20s%s %v\r\n", lightBlue, "base_url", reset, embBaseURL)
		
		output += fmt.Sprintf("\r\n%sTip:%s Use 'mako health' to validate your configuration\r\n", dimBlue, reset)
		output += "\r\n"
		return output, nil

	case "get":
		if len(args) < 2 {
			return "Usage: mako config get <key>\r\n", nil
		}
		value, err := cfg.Get(args[1])
		if err != nil {
			return fmt.Sprintf("Error: %v\r\n", err), nil
		}
		
		// Hide API key
		if args[1] == "api_key" {
			if str, ok := value.(string); ok && len(str) > 10 {
				value = str[:10] + "..."
			}
		}
		
		return fmt.Sprintf("%s: %v\r\n", args[1], value), nil

	case "set":
		if len(args) < 3 {
			return "Usage: mako config set <key> <value>\r\n", nil
		}
		
		// Convert string value to appropriate type based on key
		key := args[1]
		valueStr := strings.Join(args[2:], " ")
		
		var value interface{} = valueStr
		
		// Type conversions for known keys
		switch key {
		case "cache_size", "history_limit", "embedding_batch_size":
			var intVal int
			if _, err := fmt.Sscanf(valueStr, "%d", &intVal); err != nil {
				return fmt.Sprintf("Error: %s must be an integer\r\n", key), nil
			}
			value = intVal
		case "telemetry", "auto_update":
			value = valueStr == "true" || valueStr == "1" || valueStr == "yes"
		}
		
		if err := cfg.Set(key, value); err != nil {
			return fmt.Sprintf("Error: %v\r\n", err), nil
		}
		
		if err := cfg.Save(); err != nil {
			return "", fmt.Errorf("failed to save config: %w", err)
		}
		
		return fmt.Sprintf("%s✓ Set %s%s\r\n", lightBlue, key, reset), nil

	case "reset":
		if err := cfg.Reset(); err != nil {
			return "", fmt.Errorf("failed to reset config: %w", err)
		}
		return fmt.Sprintf("%s✓ Configuration reset to defaults%s\r\n", lightBlue, reset), nil

	default:
		return fmt.Sprintf("Unknown config command: %s\r\n", args[0]), nil
	}
}

// handleUpdate handles the 'mako update' command
func handleUpdate(args []string) (string, error) {
	cyan := "\033[38;2;0;209;255m"
	lightBlue := "\033[38;2;93;173;226m"
	dimBlue := "\033[38;2;120;150;180m"
	reset := "\033[0m"

	if len(args) == 0 {
		return "Usage: mako update <check|install>\r\n", nil
	}

	switch args[0] {
	case "check":
		output := fmt.Sprintf("\r\n%sChecking for updates...%s\r\n\r\n", lightBlue, reset)
		
		info, err := config.CheckForUpdates()
		if err != nil {
			return "", fmt.Errorf("failed to check for updates: %w", err)
		}

		if info.Available {
			output += fmt.Sprintf("%sNew version available:%s %sv%s%s (you have v%s)\r\n\r\n",
				lightBlue, reset, cyan, info.LatestVersion, reset, info.CurrentVersion)
			
			if info.ReleaseNotes != "" {
				output += fmt.Sprintf("%sChanges:%s\r\n", dimBlue, reset)
				// Show first 5 lines of release notes
				lines := strings.Split(info.ReleaseNotes, "\n")
				for i, line := range lines {
					if i >= 5 {
						output += "  ...\r\n"
						break
					}
					output += fmt.Sprintf("  %s\r\n", line)
				}
				output += "\r\n"
			}
			
			output += fmt.Sprintf("%sRun:%s %smako update install%s\r\n\r\n", lightBlue, reset, cyan, reset)
		} else {
			output += fmt.Sprintf("%s✓ You're running the latest version (v%s)%s\r\n\r\n",
				cyan, info.CurrentVersion, reset)
		}
		
		return output, nil

	case "install":
		info, err := config.CheckForUpdates()
		if err != nil {
			return "", fmt.Errorf("failed to check for updates: %w", err)
		}

		if !info.Available {
			return fmt.Sprintf("\r\n%s✓ You're already running the latest version%s\r\n\r\n",
				cyan, reset), nil
		}

		if err := config.InstallUpdate(info); err != nil {
			return "", fmt.Errorf("failed to install update: %w", err)
		}

		return "", nil // Success message printed by InstallUpdate

	default:
		return fmt.Sprintf("Unknown update command: %s\r\n", args[0]), nil
	}
}
