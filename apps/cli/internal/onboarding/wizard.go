package onboarding

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fabiobrug/mako.git/internal/config"
)

// Provider colors for consistent theming
const (
	ColorGemini     = "\033[38;2;66;133;244m"  // Google Blue
	ColorClaude     = "\033[38;2;255;138;76m"  // Claude Orange
	ColorOpenAI     = "\033[38;2;16;163;127m"  // OpenAI Green
	ColorDeepSeek   = "\033[38;2;138;43;226m"  // Purple
	ColorOllama     = "\033[38;2;255;215;0m"   // Gold/Yellow
	ColorOpenRouter = "\033[38;2;255;20;147m"  // Deep Pink/Magenta
	ColorCyan       = "\033[38;2;0;209;255m"   // Mako Cyan
	ColorLightBlue  = "\033[38;2;93;173;226m"  // Mako Light Blue
	ColorDimBlue    = "\033[38;2;120;150;180m" // Mako Dim Blue
	ColorGreen      = "\033[38;2;46;204;113m"  // Success Green
	ColorYellow     = "\033[38;2;241;196;15m"  // Warning Yellow
	ColorRed        = "\033[38;2;231;76;60m"   // Error Red
	ColorReset      = "\033[0m"
	ColorBold       = "\033[1m"
	ColorDim        = "\033[2m"
)

// ProviderInfo holds display information for each provider
type ProviderInfo struct {
	Name        string
	Color       string
	Description string
	APIKeyName  string
	GetAPIURL   string
	IsLocal     bool
}

var providers = []ProviderInfo{
	{
		Name:        "gemini",
		Color:       ColorGemini,
		Description: "Google's Gemini (Fast & Free tier available)",
		APIKeyName:  "GEMINI_API_KEY",
		GetAPIURL:   "https://aistudio.google.com/app/apikey",
		IsLocal:     false,
	},
	{
		Name:        "anthropic",
		Color:       ColorClaude,
		Description: "Anthropic Claude (High quality reasoning)",
		APIKeyName:  "ANTHROPIC_API_KEY",
		GetAPIURL:   "https://console.anthropic.com/",
		IsLocal:     false,
	},
	{
		Name:        "openai",
		Color:       ColorOpenAI,
		Description: "OpenAI GPT Models (Industry standard)",
		APIKeyName:  "OPENAI_API_KEY",
		GetAPIURL:   "https://platform.openai.com/api-keys",
		IsLocal:     false,
	},
	{
		Name:        "deepseek",
		Color:       ColorDeepSeek,
		Description: "DeepSeek (Cost-effective alternative)",
		APIKeyName:  "DEEPSEEK_API_KEY",
		GetAPIURL:   "https://platform.deepseek.com/",
		IsLocal:     false,
	},
	{
		Name:        "openrouter",
		Color:       ColorOpenRouter,
		Description: "OpenRouter (Access to multiple models)",
		APIKeyName:  "OPENROUTER_API_KEY",
		GetAPIURL:   "https://openrouter.ai/keys",
		IsLocal:     false,
	},
	{
		Name:        "ollama",
		Color:       ColorOllama,
		Description: "Ollama (Local, private, free)",
		APIKeyName:  "",
		GetAPIURL:   "https://ollama.ai/",
		IsLocal:     true,
	},
}

// RunWizard runs the interactive onboarding wizard
func RunWizard() error {
	reader := bufio.NewReader(os.Stdin)

	// Clear screen and show welcome
	clearScreen()
	showWelcomeBanner()

	// Allow skipping
	fmt.Printf("\n%sPress Enter to continue, or type 'skip' to use defaults: %s", ColorDimBlue, ColorReset)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	
	if response == "skip" || response == "s" {
		fmt.Printf("%s\n✓ Skipping setup wizard. Using default configuration.%s\n", ColorYellow, ColorReset)
		return createDefaultConfig()
	}

	clearScreen()

	// Step 1: Provider Selection
	selectedProviders, err := selectProviders(reader)
	if err != nil {
		return err
	}

	if len(selectedProviders) == 0 {
		fmt.Printf("%s\n⚠  No providers selected. Using default (Gemini).%s\n", ColorYellow, ColorReset)
		return createDefaultConfig()
	}

	// Step 2: Configure API Keys
	apiKeys := make(map[string]string)
	for _, provider := range selectedProviders {
		if provider.IsLocal {
			continue // Skip API key setup for local providers
		}
		
		key, err := configureAPIKey(reader, provider)
		if err != nil {
			return err
		}
		if key != "" {
			apiKeys[provider.Name] = key
		}
	}

	// Step 3: Select Default Provider
	defaultProvider := selectDefaultProvider(reader, selectedProviders)

	// Step 4: Additional Settings
	clearScreen()
	showSectionHeader("Additional Settings")
	
	safetyLevel := configureSafety(reader)
	autoUpdate := configureAutoUpdate(reader)

	// Step 5: Save Configuration
	if err := saveConfiguration(defaultProvider, apiKeys, safetyLevel, autoUpdate); err != nil {
		return err
	}

	// Step 6: Show completion message
	showCompletionMessage(defaultProvider)

	return nil
}

func selectProviders(reader *bufio.Reader) ([]ProviderInfo, error) {
	clearScreen()
	showSectionHeader("Select AI Providers")
	
	fmt.Printf("%sYou can configure multiple providers and switch between them later.%s\n\n", ColorDimBlue, ColorReset)
	
	// Display providers
	for i, p := range providers {
		icon := "☁"
		if p.IsLocal {
			icon = "🏠"
		}
		fmt.Printf("  %s%d%s. %s%s%s %s %s\n", 
			ColorCyan, i+1, ColorReset,
			p.Color, p.Name, ColorReset,
			icon, p.Description)
	}
	
	fmt.Printf("\n%sEnter provider numbers separated by commas (e.g., 1,2,3)%s\n", ColorDimBlue, ColorReset)
	fmt.Printf("%sOr press Enter for Gemini (recommended): %s", ColorLightBlue, ColorReset)
	
	response, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	
	response = strings.TrimSpace(response)
	if response == "" {
		// Default to Gemini
		return []ProviderInfo{providers[0]}, nil
	}
	
	// Parse selections
	selections := strings.Split(response, ",")
	var selected []ProviderInfo
	seen := make(map[string]bool)
	
	for _, sel := range selections {
		sel = strings.TrimSpace(sel)
		var idx int
		if _, err := fmt.Sscanf(sel, "%d", &idx); err == nil && idx >= 1 && idx <= len(providers) {
			p := providers[idx-1]
			if !seen[p.Name] {
				selected = append(selected, p)
				seen[p.Name] = true
			}
		}
	}
	
	return selected, nil
}

func configureAPIKey(reader *bufio.Reader, provider ProviderInfo) (string, error) {
	clearScreen()
	showSectionHeader(fmt.Sprintf("Configure %s", provider.Name))
	
	fmt.Printf("%s%s%s\n\n", provider.Color, provider.Description, ColorReset)
	fmt.Printf("%sGet your API key: %s%s%s\n\n", ColorDimBlue, ColorCyan, provider.GetAPIURL, ColorReset)
	
	// Check if key already exists in environment
	envKey := os.Getenv(provider.APIKeyName)
	if envKey != "" {
		fmt.Printf("%s✓ Found %s in environment%s\n", ColorGreen, provider.APIKeyName, ColorReset)
		fmt.Printf("%sUse this key? (Y/n): %s", ColorLightBlue, ColorReset)
		
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		
		if response == "" || response == "y" || response == "yes" {
			return envKey, nil
		}
	}
	
	fmt.Printf("%sEnter your %s API key (or press Enter to skip): %s", ColorLightBlue, provider.Name, ColorReset)
	
	key, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	
	key = strings.TrimSpace(key)
	
	if key == "" {
		fmt.Printf("%s⚠  Skipping %s configuration%s\n", ColorYellow, provider.Name, ColorReset)
		fmt.Printf("%sPress Enter to continue...%s", ColorDimBlue, ColorReset)
		reader.ReadString('\n')
		return "", nil
	}
	
	// Validate key format (basic check)
	if len(key) < 10 {
		fmt.Printf("%s⚠  API key seems too short. Continue anyway? (y/N): %s", ColorYellow, ColorReset)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		
		if response != "y" && response != "yes" {
			return "", nil
		}
	}
	
	fmt.Printf("%s✓ API key configured%s\n", ColorGreen, ColorReset)
	fmt.Printf("%sPress Enter to continue...%s", ColorDimBlue, ColorReset)
	reader.ReadString('\n')
	
	return key, nil
}

func selectDefaultProvider(reader *bufio.Reader, providers []ProviderInfo) string {
	if len(providers) == 1 {
		return providers[0].Name
	}
	
	clearScreen()
	showSectionHeader("Select Default Provider")
	
	fmt.Printf("%sWhich provider would you like to use by default?%s\n\n", ColorDimBlue, ColorReset)
	
	for i, p := range providers {
		fmt.Printf("  %s%d%s. %s%s%s\n", ColorCyan, i+1, ColorReset, p.Color, p.Name, ColorReset)
	}
	
	fmt.Printf("\n%sEnter number (default: 1): %s", ColorLightBlue, ColorReset)
	
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(response)
	
	if response == "" {
		return providers[0].Name
	}
	
	var idx int
	if _, err := fmt.Sscanf(response, "%d", &idx); err == nil && idx >= 1 && idx <= len(providers) {
		return providers[idx-1].Name
	}
	
	return providers[0].Name
}

func configureSafety(reader *bufio.Reader) string {
	fmt.Printf("\n%sSafety Level (confirm dangerous commands):%s\n", ColorDimBlue, ColorReset)
	fmt.Printf("  1. %sLow%s    - No confirmations\n", ColorRed, ColorReset)
	fmt.Printf("  2. %sMedium%s - Confirm destructive commands (recommended)\n", ColorYellow, ColorReset)
	fmt.Printf("  3. %sHigh%s   - Confirm all commands\n", ColorGreen, ColorReset)
	fmt.Printf("\n%sEnter number (default: 2): %s", ColorLightBlue, ColorReset)
	
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(response)
	
	switch response {
	case "1":
		return "low"
	case "3":
		return "high"
	default:
		return "medium"
	}
}

func configureAutoUpdate(reader *bufio.Reader) bool {
	fmt.Printf("\n%sEnable automatic updates? (Y/n): %s", ColorLightBlue, ColorReset)
	
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	
	return response == "" || response == "y" || response == "yes"
}

func saveConfiguration(defaultProvider string, apiKeys map[string]string, safetyLevel string, autoUpdate bool) error {
	// Create config
	cfg := config.DefaultConfig()
	cfg.LLMProvider = defaultProvider
	cfg.SafetyLevel = safetyLevel
	cfg.AutoUpdate = autoUpdate
	
	// For backward compatibility, set APIKey field if using single provider
	if key, ok := apiKeys[defaultProvider]; ok {
		cfg.APIKey = key
	}
	
	// Save main config
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}
	
	// Save API keys to .env file for multi-provider support
	if len(apiKeys) > 0 {
		if err := saveAPIKeys(apiKeys); err != nil {
			return fmt.Errorf("failed to save API keys: %w", err)
		}
	}
	
	return nil
}

func saveAPIKeys(apiKeys map[string]string) error {
	makoDir := config.GetMakoDir()
	envPath := fmt.Sprintf("%s/.env", makoDir)
	
	// Read existing .env if it exists
	existing := make(map[string]string)
	if data, err := os.ReadFile(envPath); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				existing[parts[0]] = parts[1]
			}
		}
	}
	
	// Update with new keys
	envVarMap := map[string]string{
		"gemini":     "GEMINI_API_KEY",
		"anthropic":  "ANTHROPIC_API_KEY",
		"openai":     "OPENAI_API_KEY",
		"deepseek":   "DEEPSEEK_API_KEY",
		"openrouter": "OPENROUTER_API_KEY",curl -sSL https://raw.githubusercontent.com/fabiobrug/mako/dev/scripts/install.sh | bash
	}
	
	for provider, key := range apiKeys {
		if envVar, ok := envVarMap[provider]; ok {
			existing[envVar] = key
		}
	}
	
	// Write back to file
	var lines []string
	lines = append(lines, "# Mako AI Provider API Keys")
	lines = append(lines, "# Generated by Mako setup wizard")
	lines = append(lines, "")
	
	for key, value := range existing {
		lines = append(lines, fmt.Sprintf("%s=%s", key, value))
	}
	
	content := strings.Join(lines, "\n") + "\n"
	return os.WriteFile(envPath, []byte(content), 0600) // 0600 for security
}

func createDefaultConfig() error {
	cfg := config.DefaultConfig()
	return cfg.Save()
}

func showCompletionMessage(defaultProvider string) {
	clearScreen()
	
	// Find provider info
	var providerColor string
	for _, p := range providers {
		if p.Name == defaultProvider {
			providerColor = p.Color
			break
		}
	}
	
	fmt.Printf("\n\n")
	fmt.Printf("  %s╔═══════════════════════════════════════════════════════╗%s\n", ColorGreen, ColorReset)
	fmt.Printf("  %s║                                                       ║%s\n", ColorGreen, ColorReset)
	fmt.Printf("  %s║              %s✨ Setup Complete! ✨%s                ║%s\n", ColorGreen, ColorBold, ColorGreen, ColorReset)
	fmt.Printf("  %s║                                                       ║%s\n", ColorGreen, ColorReset)
	fmt.Printf("  %s╚═══════════════════════════════════════════════════════╝%s\n", ColorGreen, ColorReset)
	fmt.Printf("\n")
	
	fmt.Printf("%sDefault provider: %s%s%s\n\n", ColorDimBlue, providerColor, defaultProvider, ColorReset)
	
	fmt.Printf("%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n", ColorLightBlue, ColorReset)
	fmt.Printf("%s  Quick Start Guide%s\n", ColorBold, ColorReset)
	fmt.Printf("%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n\n", ColorLightBlue, ColorReset)
	
	fmt.Printf("%sTry these commands:%s\n\n", ColorLightBlue, ColorReset)
	fmt.Printf("  %s▸ mako ask \"list files in current directory\"%s\n", ColorCyan, ColorReset)
	fmt.Printf("  %s▸ mako ask \"find large files over 100MB\"%s\n", ColorCyan, ColorReset)
	fmt.Printf("  %s▸ mako history%s\n", ColorCyan, ColorReset)
	fmt.Printf("  %s▸ mako stats%s\n\n", ColorCyan, ColorReset)
	
	fmt.Printf("%sManage providers:%s\n\n", ColorLightBlue, ColorReset)
	fmt.Printf("  %s▸ mako config providers%s       - View configured providers\n", ColorCyan, ColorReset)
	fmt.Printf("  %s▸ mako config switch <name>%s   - Switch active provider\n", ColorCyan, ColorReset)
	fmt.Printf("  %s▸ mako setup%s                  - Re-run this wizard\n\n", ColorCyan, ColorReset)
	
	fmt.Printf("%sDocumentation: %shttps://github.com/fabiobrug/mako%s\n\n", ColorDimBlue, ColorCyan, ColorReset)
	
	fmt.Printf("%sPress Enter to start Mako...%s", ColorLightBlue, ColorReset)
	bufio.NewReader(os.Stdin).ReadString('\n')
	
	clearScreen()
}

func clearScreen() {
	fmt.Print("\033[2J\033[H")
}

func showWelcomeBanner() {
	fmt.Printf("\n\n")
	fmt.Printf("  %s╔═══════════════════════════════════════════════════════╗%s\n", ColorCyan, ColorReset)
	fmt.Printf("  %s║                                                       ║%s\n", ColorCyan, ColorReset)
	fmt.Printf("  %s║          ███╗   ███╗ █████╗ ██╗  ██╗ ██████╗          ║%s\n", ColorCyan, ColorReset)
	fmt.Printf("  %s║          ████╗ ████║██╔══██╗██║ ██╔╝██╔═══██╗         ║%s\n", ColorCyan, ColorReset)
	fmt.Printf("  %s║          ██╔████╔██║███████║█████╔╝ ██║   ██║         ║%s\n", ColorCyan, ColorReset)
	fmt.Printf("  %s║          ██║╚██╔╝██║██╔══██║██╔═██╗ ██║   ██║         ║%s\n", ColorCyan, ColorReset)
	fmt.Printf("  %s║          ██║ ╚═╝ ██║██║  ██║██║  ██╗╚██████╔╝         ║%s\n", ColorCyan, ColorReset)
	fmt.Printf("  %s║          ╚═╝     ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝          ║%s\n", ColorCyan, ColorReset)
	fmt.Printf("  %s║                                                       ║%s\n", ColorCyan, ColorReset)
	fmt.Printf("  %s║           %sAI-Native Shell Orchestrator%s            ║%s\n", ColorCyan, ColorBold, ColorCyan, ColorReset)
	fmt.Printf("  %s║              %sWelcome! Let's get started.%s          ║%s\n", ColorCyan, ColorDimBlue, ColorCyan, ColorReset)
	fmt.Printf("  %s║                                                       ║%s\n", ColorCyan, ColorReset)
	fmt.Printf("  %s╚═══════════════════════════════════════════════════════╝%s\n", ColorCyan, ColorReset)
	fmt.Printf("\n")
}

func showSectionHeader(title string) {
	fmt.Printf("%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n", ColorLightBlue, ColorReset)
	fmt.Printf("%s  %s%s\n", ColorBold, title, ColorReset)
	fmt.Printf("%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n\n", ColorLightBlue, ColorReset)
}
