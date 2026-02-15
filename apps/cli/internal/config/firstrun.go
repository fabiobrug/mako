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

	// Create default config
	config := DefaultConfig()
	if err := config.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	// Step 1: AI Provider Setup
	fmt.Printf("%s━━━━━━━━━━━━━━━━━━━━━━━%s\n", dimBlue, reset)
	fmt.Printf("%sStep 1: AI Provider Setup%s\n", lightBlue, reset)
	fmt.Printf("%s━━━━━━━━━━━━━━━━━━━━━━━%s\n\n", dimBlue, reset)
	fmt.Printf("%sMako supports multiple AI providers:%s\n", dimBlue, reset)
	fmt.Printf("  • %sOllama%s - Local, free, private\n", cyan, reset)
	fmt.Printf("  • %sOpenAI%s - GPT models, best quality\n", cyan, reset)
	fmt.Printf("  • %sAnthropic%s - Claude models\n", cyan, reset)
	fmt.Printf("  • %sGemini%s - Google's models (free tier)\n", cyan, reset)
	fmt.Printf("  • %sDeepSeek / OpenRouter%s - Cost-effective alternatives\n\n", cyan, reset)
	
	fmt.Printf("%sTo configure your provider:%s\n", dimBlue, reset)
	fmt.Printf("  1. Navigate: %scd apps/cli%s\n", cyan, reset)
	fmt.Printf("  2. Copy config: %scp .env.example .env%s\n", cyan, reset)
	fmt.Printf("  3. Edit .env and set your provider and API key\n")
	fmt.Printf("  4. Or use: %smako config set api_key YOUR_KEY%s\n\n", cyan, reset)
	
	fmt.Printf("%sSetup guide: %shttps://github.com/fabiobrug/mako/blob/main/docs/SETUP.md%s\n\n", dimBlue, cyan, reset)

	// Step 2: Quick Tour
	fmt.Printf("%s━━━━━━━━━━━━━━━━%s\n", dimBlue, reset)
	fmt.Printf("%sStep 2: Quick Tour%s\n", lightBlue, reset)
	fmt.Printf("%s━━━━━━━━━━━━━━━━%s\n\n", dimBlue, reset)
	fmt.Printf("%sTry these commands:%s\n", dimBlue, reset)
	fmt.Printf("  %smako ask \"list files\"%s\n", cyan, reset)
	fmt.Printf("  %smako ask \"find large files\"%s\n", cyan, reset)
	fmt.Printf("  %smako history%s\n", cyan, reset)
	fmt.Printf("  %smako health%s - Check configuration\n\n", cyan, reset)

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


// ShowWelcomeMessage displays a welcome message on first run
func ShowWelcomeMessage() {
	cyan := "\033[38;2;0;209;255m"
	lightBlue := "\033[38;2;93;173;226m"
	reset := "\033[0m"

	fmt.Printf("\n%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n", lightBlue, reset)
	fmt.Printf("%s  Starting Mako v1.3.2%s\n", cyan, reset)
	fmt.Printf("%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n\n", lightBlue, reset)
}
