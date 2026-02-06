package shell

import (
	"fmt"
	"strings"

	"github.com/fabiobrug/mako.git/internal/ai"
	"github.com/fabiobrug/mako.git/internal/database"
)

func InterceptCommand(line string, db *database.DB) (bool, string, error) {
	trimmed := strings.TrimSpace(line)

	if strings.HasPrefix(trimmed, "mako ") {
		parts := strings.Fields(trimmed)
		if len(parts) < 2 {
			return true, "Usage: mako <command>\n", nil
		}

		switch parts[1] {
		case "ask":
			if len(parts) < 3 {
				return true, "Usage: mako ask <question>\n", nil
			}
			query := strings.Join(parts[2:], " ")
			output, err := handleAsk(query)
			return true, output, err

		case "history":
			output, err := handleHistory(parts[2:], db)
			return true, output, err

		case "stats":
			output, err := handleStats(db)
			return true, output, err

		case "help":
			return true, getHelpText(), nil

		case "v", "version":
			return true, fmt.Sprintf("v0.1.0\n"), nil

		default:
			return true, fmt.Sprintf("Unknown mako command: %s\n", parts[1]), nil
		}
	}

	return false, "", nil
}

func handleAsk(query string) (string, error) {
	client, err := ai.NewGeminiClient()
	if err != nil {
		return "", err
	}

	context := ai.GetSystemContext([]string{})
	command, err := client.GenerateCommand(query, context)
	if err != nil {
		return "", err
	}

	output := fmt.Sprintf("\n\n Suggested command:\n   %s\n\n", command)
	output += "  Note: Commands generated inside Mako are shown but not auto-executed.\n"
	output += "Copy and paste the command above to run it.\n\n"

	return output, nil
}

func handleHistory(args []string, db *database.DB) (string, error) {
	if db == nil {
		return "\nDatabase not available\n\n", nil
	}

	if len(args) > 0 && args[0] == "semantic" {
		if len(args) < 2 {
			return "\nUsage: mako history semantic <query>\n\n", nil
		}
		return handleSemanticHistory(strings.Join(args[1:], " "), db)
	}

	if len(args) == 0 {
		commands, err := db.GetRecentCommands(10)
		if err != nil {
			return "", err
		}

		if len(commands) == 0 {
			return "\nNo command history yet.\n\n", nil
		}

		var output strings.Builder
		output.WriteString("\n\n Recent commands:\n\n")
		for _, cmd := range commands {
			output.WriteString(fmt.Sprintf("[%s] %s\n",
				cmd.Timestamp.Format("15:04:05"),
				cmd.Command))
		}
		output.WriteString("\n")

		return output.String(), nil
	}

	query := strings.Join(args, " ")
	commands, err := db.SearchCommands(query, 10)
	if err != nil {
		return "", err
	}

	if len(commands) == 0 {
		return fmt.Sprintf("\nNo commands found matching: %s\n\n", query), nil
	}

	var output strings.Builder
	output.WriteString(fmt.Sprintf("\n\n Found %d commands:\n\n", len(commands)))
	for _, cmd := range commands {
		output.WriteString(fmt.Sprintf("[%s] %s\n",
			cmd.Timestamp.Format("15:04:05"),
			cmd.Command))
	}
	output.WriteString("\n")

	return output.String(), nil
}

func handleSemanticHistory(query string, db *database.DB) (string, error) {
	embedService, err := ai.NewEmbeddingService()
	if err != nil {
		return "", err
	}

	queryVec, err := embedService.Embed(query)
	if err != nil {
		return "", err
	}

	queryBytes := ai.VectorToBytes(queryVec)
	commands, err := db.SearchCommandsSemantic(queryBytes, 10, 0.5)
	if err != nil {
		return "", err
	}

	if len(commands) == 0 {
		return fmt.Sprintf("\nNo semantically similar commands found for: %s\n\n", query), nil
	}

	var output strings.Builder
	output.WriteString(fmt.Sprintf("\n\n Found %d similar commands:\n\n", len(commands)))
	for _, cmd := range commands {
		output.WriteString(fmt.Sprintf("[%s] %s\n",
			cmd.Timestamp.Format("15:04:05"),
			cmd.Command))
	}
	output.WriteString("\n")

	return output.String(), nil
}

func handleStats(db *database.DB) (string, error) {
	if db == nil {
		return "\nDatabase not available\n\n", nil
	}

	stats, err := db.GetStats()
	if err != nil {
		return "", err
	}

	var output strings.Builder
	output.WriteString("\n\n Mako Statistics:\n\n")
	output.WriteString(fmt.Sprintf("  Total commands: %d\n", stats["total_commands"]))
	output.WriteString(fmt.Sprintf("  Commands today: %d\n", stats["commands_today"]))
	output.WriteString(fmt.Sprintf("  Avg duration: %.0fms\n\n", stats["avg_duration_ms"]))

	return output.String(), nil
}

func getHelpText() string {
	return `

 Mako Commands (inside Mako shell):

  mako ask <question>              Generate command from natural language
  mako history                     Show recent commands
  mako history <keyword>           Search by keyword
  mako history semantic <query>    Search by meaning
  mako stats                       Show statistics
  mako help                        Show this help
  mako version                     Show the current Mako version

Regular shell commands work normally!

`
}
