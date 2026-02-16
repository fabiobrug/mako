package shell

import (
	"fmt"
	"os"
	"path/filepath"
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
		return fmt.Sprintf("Usage: mako config <list|get|set|reset|providers|switch>\r\n"), nil
	}

	switch args[0] {
	case "providers":
		return handleProvidersList(), nil
	
	case "switch":
		if len(args) < 2 {
			return "Usage: mako config switch <provider>\r\n", nil
		}
		return handleProviderSwitch(args[1], cfg)
	
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

// Provider colors for display
const (
	colorGemini     = "\033[38;2;66;133;244m"  // Google Blue
	colorClaude     = "\033[38;2;255;138;76m"  // Claude Orange
	colorOpenAI     = "\033[38;2;16;163;127m"  // OpenAI Green
	colorDeepSeek   = "\033[38;2;138;43;226m"  // Purple
	colorOllama     = "\033[38;2;255;215;0m"   // Gold/Yellow
	colorOpenRouter = "\033[38;2;255;20;147m"  // Deep Pink/Magenta
	colorReset      = "\033[0m"
	colorDimBlue    = "\033[38;2;120;150;180m"
	colorLightBlue  = "\033[38;2;93;173;226m"
	colorCyan       = "\033[38;2;0;209;255m"
	colorGreen      = "\033[38;2;46;204;113m"
)

// providerColorMap maps provider names to their display colors
var providerColorMap = map[string]string{
	"gemini":     colorGemini,
	"anthropic":  colorClaude,
	"openai":     colorOpenAI,
	"deepseek":   colorDeepSeek,
	"ollama":     colorOllama,
	"openrouter": colorOpenRouter,
}

// handleProvidersList displays all configured providers
func handleProvidersList() string {
	output := fmt.Sprintf("\r\n%sConfigured AI Providers%s\r\n", colorCyan, colorReset)
	output += fmt.Sprintf("%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\r\n\r\n", colorDimBlue, colorReset)
	
	// Load current config to see active provider
	cfg, err := config.LoadConfig()
	activeProvider := "gemini" // default
	if err == nil {
		activeProvider = cfg.LLMProvider
	}
	
	// Read .env file to find configured providers
	makoDir := config.GetMakoDir()
	envPath := filepath.Join(makoDir, ".env")
	
	configuredProviders := make(map[string]bool)
	
	if data, err := os.ReadFile(envPath); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.Contains(line, "GEMINI_API_KEY=") && !strings.HasSuffix(line, "=") {
				configuredProviders["gemini"] = true
			}
			if strings.Contains(line, "ANTHROPIC_API_KEY=") && !strings.HasSuffix(line, "=") {
				configuredProviders["anthropic"] = true
			}
			if strings.Contains(line, "OPENAI_API_KEY=") && !strings.HasSuffix(line, "=") {
				configuredProviders["openai"] = true
			}
			if strings.Contains(line, "DEEPSEEK_API_KEY=") && !strings.HasSuffix(line, "=") {
				configuredProviders["deepseek"] = true
			}
			if strings.Contains(line, "OPENROUTER_API_KEY=") && !strings.HasSuffix(line, "=") {
				configuredProviders["openrouter"] = true
			}
		}
	}
	
	// Also check environment variables
	if os.Getenv("GEMINI_API_KEY") != "" {
		configuredProviders["gemini"] = true
	}
	if os.Getenv("ANTHROPIC_API_KEY") != "" {
		configuredProviders["anthropic"] = true
	}
	if os.Getenv("OPENAI_API_KEY") != "" {
		configuredProviders["openai"] = true
	}
	if os.Getenv("DEEPSEEK_API_KEY") != "" {
		configuredProviders["deepseek"] = true
	}
	if os.Getenv("OPENROUTER_API_KEY") != "" {
		configuredProviders["openrouter"] = true
	}
	
	// Ollama is always available (local)
	configuredProviders["ollama"] = true
	
	providers := []struct {
		name        string
		description string
		isLocal     bool
	}{
		{"gemini", "Google's Gemini", false},
		{"anthropic", "Anthropic Claude", false},
		{"openai", "OpenAI GPT Models", false},
		{"deepseek", "DeepSeek", false},
		{"openrouter", "OpenRouter", false},
		{"ollama", "Ollama (Local)", true},
	}
	
	for _, p := range providers {
		color := providerColorMap[p.name]
		if color == "" {
			color = colorReset
		}
		
		status := " "
		statusColor := colorDimBlue
		
		if p.name == activeProvider {
			status = "●"
			statusColor = colorGreen
		} else if configuredProviders[p.name] {
			status = "○"
			statusColor = colorDimBlue
		} else {
			status = "✕"
			statusColor = colorDimBlue
		}
		
		output += fmt.Sprintf("  %s%s%s %s%-12s%s %s%s\r\n", 
			statusColor, status, colorReset,
			color, p.name, colorReset,
			colorDimBlue, p.description)
	}
	
	output += fmt.Sprintf("\r\n%s● Active  ○ Configured  ✕ Not configured%s\r\n\r\n", colorDimBlue, colorReset)
	output += fmt.Sprintf("%sSwitch provider:%s %smako config switch <provider>%s\r\n", colorDimBlue, colorReset, colorCyan, colorReset)
	output += fmt.Sprintf("%sSetup wizard:%s   %smako setup%s\r\n\r\n", colorDimBlue, colorReset, colorCyan, colorReset)
	
	return output
}

// handleProviderSwitch switches the active AI provider
func handleProviderSwitch(provider string, cfg *config.Config) (string, error) {
	provider = strings.ToLower(strings.TrimSpace(provider))
	
	// Validate provider
	validProviders := map[string]string{
		"gemini":     "gemini",
		"anthropic":  "anthropic",
		"claude":     "anthropic", // alias
		"openai":     "openai",
		"gpt":        "openai", // alias
		"deepseek":   "deepseek",
		"openrouter": "openrouter",
		"ollama":     "ollama",
	}
	
	normalizedProvider, ok := validProviders[provider]
	if !ok {
		return fmt.Sprintf("%sError:%s Unknown provider '%s'\r\n\r\nValid providers: gemini, anthropic, openai, deepseek, openrouter, ollama\r\n", 
			colorCyan, colorReset, provider), nil
	}
	
	// Check if provider has API key (except ollama)
	if normalizedProvider != "ollama" {
		hasKey := false
		
		// Check .env file
		makoDir := config.GetMakoDir()
		envPath := filepath.Join(makoDir, ".env")
		
		envVarName := ""
		switch normalizedProvider {
		case "gemini":
			envVarName = "GEMINI_API_KEY"
		case "anthropic":
			envVarName = "ANTHROPIC_API_KEY"
		case "openai":
			envVarName = "OPENAI_API_KEY"
		case "deepseek":
			envVarName = "DEEPSEEK_API_KEY"
		case "openrouter":
			envVarName = "OPENROUTER_API_KEY"
		}
		
		// Check environment
		if os.Getenv(envVarName) != "" {
			hasKey = true
		}
		
		// Check .env file
		if !hasKey {
			if data, err := os.ReadFile(envPath); err == nil {
				if strings.Contains(string(data), envVarName+"=") && 
				   !strings.Contains(string(data), envVarName+"=\n") &&
				   !strings.Contains(string(data), envVarName+"= ") {
					hasKey = true
				}
			}
		}
		
		if !hasKey {
			return fmt.Sprintf("%sWarning:%s Provider '%s%s%s' is not configured.\r\n\r\n%sRun %smako setup%s to configure it.%s\r\n", 
				colorCyan, colorReset, 
				providerColorMap[normalizedProvider], normalizedProvider, colorReset,
				colorDimBlue, colorCyan, colorDimBlue, colorReset), nil
		}
	}
	
	// Update config
	cfg.LLMProvider = normalizedProvider
	if err := cfg.Save(); err != nil {
		return "", fmt.Errorf("failed to save config: %w", err)
	}
	
	color := providerColorMap[normalizedProvider]
	if color == "" {
		color = colorReset
	}
	
	return fmt.Sprintf("%s✓ Switched to %s%s%s\r\n", colorGreen, color, normalizedProvider, colorReset), nil
}
