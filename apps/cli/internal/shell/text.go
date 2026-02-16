package shell

import (
	"fmt"
	"strings"
)

func getHelpText() string {
	lightBlue := "\033[38;2;93;173;226m"
	cyan := "\033[38;2;0;209;255m"
	dimBlue := "\033[38;2;120;150;180m"
	reset := "\033[0m"
	return fmt.Sprintf(`
%s╭─ Mako Commands%s
%s│%s
%s│%s  %smako ask <question>%s              Generate command from natural language
%s│%s  
%s│%s  %smako history%s                     Show recent commands
%s│%s  %smako history <keyword>%s           Search by keyword
%s│%s  %smako history semantic <query>%s    Search by meaning
%s│%s  %smako history --failed%s            Show only failed commands
%s│%s  %smako history --success%s           Show only successful commands
%s│%s  %smako history --interactive%s       Browse history interactively
%s│%s  
%s│%s  %smako alias save <name> <cmd>%s     Save a command alias
%s│%s  %smako alias list [--tag <tag>]%s    List all saved aliases
%s│%s  %smako alias run <name> [args]%s     Run a saved alias with parameters
%s│%s  %smako alias delete <name>%s         Delete an alias
%s│%s  %smako alias export <file>%s         Export aliases to file
%s│%s  %smako alias import <file>%s         Import aliases from file
%s│%s  
%s│%s  %smako config list%s                 Show all configuration settings
%s│%s  %smako config get <key>%s            Get configuration value
%s│%s  %smako config set <key> <value>%s    Set configuration value
%s│%s  %smako config reset%s                Reset to default configuration
%s│%s  
%s│%s  %smako update check%s                Check for updates
%s│%s  %smako update install%s              Install latest version
%s│%s  
%s│%s  %smako stats%s                       Show statistics
%s│%s  %smako health%s                      Check Mako health and performance
%s│%s  %smako export [--last N] > file%s    Export command history to JSON
%s│%s  %smako import <file>%s               Import commands from JSON
%s│%s  %smako sync%s                        Sync bash history to Mako
%s│%s  
%s│%s  %smako clear%s                       Clear conversation history
%s│%s  %smako completion <bash|zsh|fish>%s  Generate shell completion script
%s│%s  %smako help%s                        Show this help
%s│%s  %smako version%s                     Show Mako version
%s│%s 
%s╰─%s %sRegular shell commands work normally!%s

`, lightBlue, reset,
		lightBlue, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset,
		lightBlue, reset, dimBlue, reset)
}

func getSharkArt() string {
	return fmt.Sprintln(`
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣴⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⣾⣿⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⡀⣀⢀⣀⣀⣀⣀⣀⣀⣀⣤⣤⣤⠤⠤⠤⠤⠤⠄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡀⣼⣿⣿⣿⣿⢿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⢰⣿⣿⣿⣿⡉⠉⠉⠉⠉⠉⠉⠉⠉⠉⣾⣿⣆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠛⠻⠻⡿⢿⣿⣿⣇⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣶⣶⠀⠀
⠈⢹⣿⣿⣿⣿⣿⣦⣼⣷⣦⣀⠀⠀⠀⠈⠛⠉⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠘⠛⢧⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣼⣿⣿⠀⠀⠀
⠀⠀⠛⢿⣿⣿⣿⣿⣿⣿⣿⣿⣷⣦⣤⣀⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣰⣿⣿⣿⡿⠇⠀⠀⠀
⠀⠀⠀⠈⠘⠿⣿⣿⣿⣿⣿⣿⠛⢿⠛⠟⠛⠋⠋⠉⠋⠙⠛⠳⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣠⣿⣿⣿⣿⣿⠃⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠈⠉⠻⡿⡇⢀⣤⣤⣶⣶⣶⣶⣏⠉⠉⠉⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣰⣤⣆⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠠⠶⠿⣄⠀⠀⠀⠀⠀⠀⠀⢤⣾⣿⣿⣿⣿⣿⠙⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠱⢾⡿⣿⣿⣿⣿⣿⣿⣿⣦⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣶⣿⣿⣿⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠐⠢⠤⢤⣀⣀⣀⣀⣀⣀⣀⣀⣀⣠⣤⣤⣤⣤⣤⣭⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠁⠛⠻⡿⣿⣿⣿⣿⣿⣿⣷⣶⣤⣤⣄⣀⣀⣀⣠⣤⣶⣾⣿⣿⣿⣿⣿⣶⣤⣄⣀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣀⣠⣤⣶⣿⣿⣿⣿⣿⣿⢿⠟⠿⠛⠛⠛⠙⠿⣿⣿⣿⣿⣿⣾⡀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠁⠙⠛⠿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⣶⣶⣆⣀⣀⣀⣿⠿⢿⣿⣿⣿⡿⠿⠟⠟⠛⠋⠈⠀⠀⠀⠀⠀⠀⠀⠙⠿⣿⣿⣿⣿⣷⣆⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣤⣶⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⢿⣿⣿⣿⣿⣿⣷⣶⣦⣬⣭⣄⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠛⢿⢿⣿⣿⣀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣴⣿⣿⣿⣿⣿⡏⠿⠙⠃⠋⠀⠀⠙⠛⠘⠟⠿⠻⠟⢿⠿⡿⠿⡿⢽⠿⡿⠻⠷⠆⠀⠁⠉⠀⠘⠋⠛⠘⠛⠿⠻⠿⠿⠿⣿⢶⣶⣤⣤⣄⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠋⠻⠻⠦⠤
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣴⡿⠿⠟⠛⠁⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠉⠉⠙⠛⠓⠒⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀`)
}

// getContextualHelp returns help for specific topics
func getContextualHelp(topic string) string {
	cyan := "\033[38;2;0;209;255m"
	lightBlue := "\033[38;2;93;173;226m"
	dimBlue := "\033[38;2;120;150;180m"
	reset := "\033[0m"

	// Remove -- prefix if present
	topic = strings.TrimPrefix(topic, "--")

	switch topic {
	case "quickstart", "quick", "start":
		return fmt.Sprintf(`
%s╭─ Mako Quickstart%s
%s│%s
%s│%s  %sGet started with Mako in 3 steps:%s
%s│%s
%s│%s  %s1. Generate commands from natural language%s
%s│%s     %smako ask "list files sorted by size"%s
%s│%s
%s│%s  %s2. Search your history%s
%s│%s     %smako history docker%s                # Text search
%s│%s     %smako history semantic "containers"%s  # Semantic search
%s│%s
%s│%s  %s3. Save frequent commands as aliases%s
%s│%s     %smako alias save deploy "git push && make deploy"%s
%s│%s     %smako alias run deploy%s
%s│%s
%s│%s  %sTips:%s
%s│%s  • Use %smako health%s to check configuration
%s│%s  • Use %smako config list%s to view settings
%s│%s  • Use %smako stats%s to see usage statistics
%s│%s
%s╰─%s %sRun 'mako help' for full command list%s

`, lightBlue, reset,
			lightBlue, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset,
			lightBlue, reset, dimBlue, reset)

	case "alias", "aliases":
		return fmt.Sprintf(`
%s╭─ Mako Aliases%s
%s│%s
%s│%s  %sAliases let you save and reuse commands:%s
%s│%s
%s│%s  %smako alias save <name> <command>%s  Save a command
%s│%s  %smako alias list%s                   List all aliases
%s│%s  %smako alias list --tag <tag>%s       List by tag
%s│%s  %smako alias run <name>%s             Run an alias
%s│%s  %smako alias delete <name>%s          Delete an alias
%s│%s  %smako alias export <file>%s          Export to file
%s│%s  %smako alias import <file>%s          Import from file
%s│%s
%s│%s  %sExamples:%s
%s│%s  %smako alias save deploy "git push && make build"%s
%s│%s  %smako alias save backup "tar -czf backup.tar.gz ."%s
%s│%s  %smako alias run deploy%s
%s│%s
%s│%s  %sTagging:%s
%s│%s  You can organize aliases with tags in the description:
%s│%s  %smako alias save deploy "git push" --tag git%s
%s│%s
%s╰─%s %sRun 'mako help' for full command list%s

`, lightBlue, reset,
			lightBlue, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset,
			lightBlue, reset, dimBlue, reset)

	case "history":
		return fmt.Sprintf(`
%s╭─ Mako History%s
%s│%s
%s│%s  %sSearch and browse your command history:%s
%s│%s
%s│%s  %smako history%s                       Recent commands
%s│%s  %smako history <keyword>%s             Text search
%s│%s  %smako history semantic <query>%s      Semantic search (by meaning)
%s│%s  %smako history --failed%s              Only failed commands
%s│%s  %smako history --success%s             Only successful commands
%s│%s  %smako history --interactive%s         Browse interactively
%s│%s
%s│%s  %sWhat is semantic search?%s
%s│%s  Find commands by describing what you want, not exact text:
%s│%s  %smako history semantic "show containers"%s
%s│%s  Finds: docker ps, docker container ls, kubectl get pods, etc.
%s│%s
%s│%s  %sSync bash history:%s
%s│%s  %smako sync%s  Import your existing bash history
%s│%s
%s╰─%s %sRun 'mako help' for full command list%s

`, lightBlue, reset,
			lightBlue, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset,
			lightBlue, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset,
			lightBlue, reset, dimBlue, reset)

	case "config", "configuration":
		return fmt.Sprintf(`
%s╭─ Mako Configuration%s
%s│%s
%s│%s  %sManage Mako settings:%s
%s│%s
%s│%s  %smako config list%s               Show all settings
%s│%s  %smako config get <key>%s          Get a value
%s│%s  %smako config set <key> <value>%s  Set a value
%s│%s  %smako config reset%s              Reset to defaults
%s│%s
%s│%s  %sKey settings:%s
%s│%s  • api_key         Your AI provider API key
%s│%s  • llm_provider    AI provider (gemini, openai, anthropic, etc.)
%s│%s  • llm_model       Model to use for command generation
%s│%s  • cache_size      Embedding cache size
%s│%s  • auto_update     Check for updates on startup
%s│%s
%s│%s  %sHealth check:%s
%s│%s  %smako health%s  Validate your configuration
%s│%s
%s╰─%s %sRun 'mako help' for full command list%s

`, lightBlue, reset,
			lightBlue, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset,
			lightBlue, reset,
			lightBlue, reset,
			lightBlue, reset,
			lightBlue, reset,
			lightBlue, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset,
			lightBlue, reset, dimBlue, reset)

	case "embedding", "embeddings", "semantic":
		return fmt.Sprintf(`
%s╭─ Understanding Embeddings%s
%s│%s
%s│%s  %sWhat are embeddings?%s
%s│%s  Numerical representations of text that capture meaning.
%s│%s  Similar meanings have similar embeddings.
%s│%s
%s│%s  %sWhy does Mako use them?%s
%s│%s  To enable semantic search - find commands by meaning,
%s│%s  not just exact text matches.
%s│%s
%s│%s  %sConfiguration:%s
%s│%s  By default, uses the same provider as your LLM.
%s│%s
%s│%s  Default embedding models:
%s│%s  • Gemini: text-embedding-005
%s│%s  • OpenAI: text-embedding-3-small
%s│%s  • Ollama: nomic-embed-text (local, free)
%s│%s
%s│%s  %sUse local embeddings (free & private):%s
%s│%s  Set in your .env file:
%s│%s  %sEMBEDDING_PROVIDER=ollama%s
%s│%s  %sEMBEDDING_MODEL=nomic-embed-text%s
%s│%s
%s│%s  %sCheck configuration:%s
%s│%s  %smako health%s       Check embedding provider status
%s│%s  %smako config list%s  View current configuration
%s│%s
%s╰─%s %sRun 'mako help' for full command list%s

`, lightBlue, reset,
			lightBlue, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset,
			lightBlue, reset,
			lightBlue, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset,
			lightBlue, reset,
			lightBlue, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset,
			lightBlue, reset,
			lightBlue, reset,
			lightBlue, reset,
			lightBlue, reset,
			lightBlue, reset,
			lightBlue, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset, cyan, reset,
			lightBlue, reset,
			lightBlue, reset, dimBlue, reset)

	default:
		return "" // Return empty string to show full help
	}
}

func handleUninstall() (string, error) {
	cyan := "\033[38;2;0;209;255m"
	lightBlue := "\033[38;2;93;173;226m"
	dimBlue := "\033[38;2;120;150;180m"
	reset := "\033[0m"

	output := fmt.Sprintf("\r\n%sUninstall Mako%s\r\n", cyan, reset)
	output += fmt.Sprintf("%s━━━━━━━━━━━━━━━━━━━━━━%s\r\n\r\n", dimBlue, reset)
	output += fmt.Sprintf("%sTo uninstall Mako, run:%s\r\n\r\n", lightBlue, reset)
	output += fmt.Sprintf("  %scurl -sSL https://get-mako.sh/uninstall.sh | bash%s\r\n\r\n", cyan, reset)
	output += fmt.Sprintf("%sOr manually:%s\r\n\r\n", lightBlue, reset)
	output += fmt.Sprintf("  %s# Remove binaries%s\r\n", dimBlue, reset)
	output += "  sudo rm /usr/local/bin/mako /usr/local/bin/mako-menu\r\n\r\n"
	output += fmt.Sprintf("  %s# Remove configuration%s\r\n", dimBlue, reset)
	output += "  rm -rf ~/.mako\r\n\r\n"
	output += fmt.Sprintf("  %s# Remove completions (optional)%s\r\n", dimBlue, reset)
	output += "  rm /etc/bash_completion.d/mako\r\n"
	output += "  rm ~/.zsh/completions/_mako\r\n"
	output += "  rm ~/.config/fish/completions/mako.fish\r\n\r\n"

	return output, nil
}
