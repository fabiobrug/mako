package config

import (
	"fmt"
	"os"
)

// IsFirstRun checks if this is the first time Mako is being run
func IsFirstRun() bool {
	configPath := GetConfigPath()
	_, err := os.Stat(configPath)
	return os.IsNotExist(err)
}

// RunFirstTimeSetup runs the interactive first-time setup wizard
// This is now handled by the onboarding package
func RunFirstTimeSetup() error {
	// Import is done in main.go to avoid circular dependencies
	// The actual wizard implementation is in internal/onboarding/wizard.go
	
	// Ensure .mako directory exists
	makoDir := GetMakoDir()
	if err := os.MkdirAll(makoDir, 0755); err != nil {
		return err
	}
	
	// The wizard is called from main.go
	return nil
}


// ShowWelcomeMessage displays a welcome message on first run
func ShowWelcomeMessage() {
	cyan := "\033[38;2;0;209;255m"
	lightBlue := "\033[38;2;93;173;226m"
	reset := "\033[0m"

	fmt.Printf("\n%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n", lightBlue, reset)
	fmt.Printf("%s  Starting Mako v1.3.7%s\n", cyan, reset)
	fmt.Printf("%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n\n", lightBlue, reset)
}
