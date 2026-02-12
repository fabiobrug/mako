package config

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

// IsFirstRun checks if this is the first time Mako is being run
func IsFirstRun() bool {
	configPath := GetConfigPath()
	_, err := os.Stat(configPath)
	return os.IsNotExist(err)
}

// RunFirstTimeSetup runs the interactive first-time setup wizard
func RunFirstTimeSetup() error {
	cyan := "\033[38;2;0;209;255m"
	lightBlue := "\033[38;2;93;173;226m"
	dimBlue := "\033[38;2;120;150;180m"
	reset := "\033[0m"

	// Clear screen and show welcome
	fmt.Print("\033[2J\033[H")
	fmt.Printf("\n%s  Mako - AI-Native Shell Orchestrator%s\n\n", cyan, reset)
	fmt.Printf("%sWelcome to Mako!%s\n\n", lightBlue, reset)
	fmt.Printf("%sMako is your AI-powered shell assistant.%s\n\n", dimBlue, reset)
	fmt.Printf("%sLet's get you set up...%s\n\n", lightBlue, reset)

	// Ensure .mako directory exists
	makoDir := GetMakoDir()
	if err := os.MkdirAll(makoDir, 0755); err != nil {
		return fmt.Errorf("failed to create .mako directory: %w", err)
	}

	// Step 1: API Key Setup
	fmt.Printf("%s━━━━━━━━━━━━━━━━%s\n", dimBlue, reset)
	fmt.Printf("%sStep 1: API Key%s\n", lightBlue, reset)
	fmt.Printf("%s━━━━━━━━━━━━━━━━%s\n\n", dimBlue, reset)
	fmt.Printf("%sMako uses Google's Gemini API (free tier available)%s\n", dimBlue, reset)
	fmt.Printf("%sGet your key: %shttps://ai.google.dev/%s\n\n", dimBlue, cyan, reset)

	apiKey, err := promptForAPIKey()
	if err != nil {
		return fmt.Errorf("failed to get API key: %w", err)
	}

	// Create and save config
	config := DefaultConfig()
	if apiKey != "" {
		config.APIKey = apiKey
		fmt.Printf("%s✓ API key saved!%s\n\n", lightBlue, reset)
	} else {
		fmt.Printf("%s⚠  You can set it later: mako config set api_key YOUR_KEY%s\n\n", dimBlue, reset)
	}

	if err := config.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	// Step 2: Quick Tour
	fmt.Printf("%s━━━━━━━━━━━━━━━━%s\n", dimBlue, reset)
	fmt.Printf("%sStep 2: Quick Tour%s\n", lightBlue, reset)
	fmt.Printf("%s━━━━━━━━━━━━━━━━%s\n\n", dimBlue, reset)
	fmt.Printf("%sTry these commands:%s\n", dimBlue, reset)
	fmt.Printf("  %smako ask \"list files\"%s\n", cyan, reset)
	fmt.Printf("  %smako ask \"find large files\"%s\n", cyan, reset)
	fmt.Printf("  %smako history%s\n\n", cyan, reset)

	// Step 3: Learn More
	fmt.Printf("%s━━━━━━━━━━━━━━━━━%s\n", dimBlue, reset)
	fmt.Printf("%sStep 3: Learn More%s\n", lightBlue, reset)
	fmt.Printf("%s━━━━━━━━━━━━━━━━━%s\n\n", dimBlue, reset)
	fmt.Printf("%sDocumentation: %shttps://github.com/fabiobrug/mako%s\n", dimBlue, cyan, reset)
	fmt.Printf("%sType %smako help%s for all commands%s\n\n", dimBlue, cyan, reset, reset)

	fmt.Printf("%sPress Enter to start Mako...%s", lightBlue, reset)
	fmt.Scanln()

	fmt.Print("\033[2J\033[H") // Clear screen
	return nil
}

// promptForAPIKey prompts the user to enter their Gemini API key
func promptForAPIKey() (string, error) {
	cyan := "\033[38;2;0;209;255m"
	reset := "\033[0m"

	fmt.Printf("%sEnter API key (or press Enter to skip): %s", cyan, reset)

	// Read password-style input (hidden)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // Add newline after password input
	
	if err != nil {
		// Fallback to regular input if terminal doesn't support hidden input
		var apiKey string
		fmt.Scanln(&apiKey)
		return strings.TrimSpace(apiKey), nil
	}

	return strings.TrimSpace(string(bytePassword)), nil
}

// ShowWelcomeMessage displays a welcome message on first run
func ShowWelcomeMessage() {
	cyan := "\033[38;2;0;209;255m"
	lightBlue := "\033[38;2;93;173;226m"
	reset := "\033[0m"

	fmt.Printf("\n%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n", lightBlue, reset)
	fmt.Printf("%s  Starting Mako v1.1.6%s\n", cyan, reset)
	fmt.Printf("%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n\n", lightBlue, reset)
}
