package shell

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/fabiobrug/mako.git/internal/ai"
	"github.com/fabiobrug/mako.git/internal/database"
	"github.com/fabiobrug/mako.git/internal/safety"
)

func handleAsk(query string, db *database.DB) (string, error) {
	client, err := ai.NewAIProvider()
	if err != nil {
		return "", err
	}

	// Load conversation history
	conversation, err := ai.LoadConversation()
	if err != nil {
		// Log error but continue without conversation
		fmt.Fprintf(os.Stderr, "Warning: Failed to load conversation: %v\n", err)
		conversation = nil
	}

	// ENHANCED: Get recent output and commands for context
	var recentOutput []string
	if recentOutputGetter != nil {
		recentOutput = recentOutputGetter(10) // Last 10 lines
	}

	var recentCommands []string
	if db != nil {
		commands, err := db.GetRecentCommands(5)
		if err == nil {
			for _, cmd := range commands {
				recentCommands = append(recentCommands, cmd.Command)
			}
		}
	}

	// Build enhanced context
	context := ai.GetEnhancedContext(recentOutput, recentCommands)

	// Generate command with conversation history
	command, err := client.GenerateCommandWithConversation(query, context, conversation)
	if err != nil {
		return "", err
	}

	// Clean bracketed paste markers
	command = strings.ReplaceAll(command, "\x1b[200~", "")
	command = strings.ReplaceAll(command, "\x1b[201~", "")
	command = strings.ReplaceAll(command, "[200~", "")
	command = strings.ReplaceAll(command, "[201~", "")
	command = strings.TrimPrefix(command, "~")
	command = strings.TrimSuffix(command, "~")
	command = strings.TrimSpace(command)

	// Safety validation
	validationResult := validator.ValidateCommand(command)

	cyan := "\033[38;2;0;209;255m"
	lightBlue := "\033[38;2;93;173;226m"
	green := "\033[38;2;100;255;100m"
	red := "\033[38;2;255;100;100m"
	gray := "\033[38;2;150;150;150m"
	reset := "\033[0m"

	tty, _ := os.OpenFile("/dev/tty", os.O_WRONLY, 0)
	if tty != nil {
		defer tty.Close()
	}
	writeTTY := func(s string) {
		if tty != nil {
			fmt.Fprint(tty, s)
		}
	}

	// Display command
	output := fmt.Sprintf("\r\n%s╭─ Generated Command%s\r\n", lightBlue, reset)
	output += fmt.Sprintf("%s│%s  %s%s%s\r\n", lightBlue, reset, cyan, command, reset)
	output += fmt.Sprintf("%s╰─%s\r\n", lightBlue, reset)
	writeTTY(output)

	// Block critical commands
	if validationResult.Risk == safety.RiskCritical {
		writeTTY(validator.FormatWarning(validationResult))
		writeTTY(fmt.Sprintf("\r\n%s✗ Command blocked for safety%s\r\n\r\n", red, reset))
		return "", nil
	}

	// Show warnings for risky commands
	if !validationResult.Safe {
		writeTTY(validator.FormatWarning(validationResult))
	}

	// Pause PTY input BEFORE any delays to ensure immediate effect
	pauseFile := filepath.Join(os.Getenv("HOME"), ".mako", "pause_input")
	os.WriteFile(pauseFile, []byte("1"), 0644)
	defer os.Remove(pauseFile)

	// Longer delay for zsh which may have more aggressive input buffering
	time.Sleep(150 * time.Millisecond)

	// Menu options
	menuArgs := []string{
		fmt.Sprintf("%sWhat would you like to do?%s", lightBlue, reset),
	}

	if !validationResult.Safe {
		menuArgs = append(menuArgs,
			"Confirm and run|run",
			"Explain what this does|explain",
			"Suggest alternatives|alternatives",
			"Edit before running|edit",
			"Copy to clipboard|copy",
			"Cancel|cancel",
		)
	} else {
		menuArgs = append(menuArgs,
			"Run command|run",
			"Explain what this does|explain",
			"Suggest alternatives|alternatives",
			"Edit before running|edit",
			"Copy to clipboard|copy",
			"Cancel|cancel",
		)
	}

	// Call menu
	var menuPath string
	possiblePaths := []string{
		"./mako-menu",
		filepath.Join(filepath.Dir(os.Args[0]), "mako-menu"),
		"mako-menu",
	}

	for _, path := range possiblePaths {
		if absPath, err := filepath.Abs(path); err == nil {
			if _, err := os.Stat(absPath); err == nil {
				menuPath = absPath
				break
			}
		}
	}

	if menuPath == "" {
		menuPath = "mako-menu"
	}

	menuCmd := exec.Command(menuPath, menuArgs...)
	menuCmd.Stderr = os.Stderr

	choiceBytes, err := menuCmd.Output()
	if err != nil {
		return "", fmt.Errorf("menu failed: %w", err)
	}

	choice := strings.TrimSpace(string(choiceBytes))
	
	// Give time for pause_file removal and terminal state to settle
	time.Sleep(150 * time.Millisecond)

	// Handle choice
	switch choice {
	case "run":
		return handleAskRun(query, command, db, client, conversation, context, writeTTY, cyan, lightBlue, green, red, gray, reset)

	case "explain":
		return handleAskExplain(command, client, context, writeTTY, cyan, lightBlue, red, reset)

	case "alternatives":
		return handleAskAlternatives(command, client, context, writeTTY, cyan, lightBlue, red, reset)

	case "edit":
		return handleAskEdit(query, command, db, writeTTY, cyan, lightBlue, green, red, gray, reset)

	case "copy":
		return handleAskCopy(query, command, conversation, writeTTY, green, red, reset)

	case "cancel":
		return handleAskCancel(query, command, conversation, writeTTY, gray, reset)

	default:
		return "", nil
	}
}

func handleAskRun(query, command string, db *database.DB, client ai.AIProvider, conversation *ai.ConversationHistory, context ai.SystemContext, writeTTY func(string), cyan, lightBlue, green, red, gray, reset string) (string, error) {
	writeTTY(fmt.Sprintf("\r\n%s▸ Executing...%s\r\n\r\n", cyan, reset))

	cmd := exec.Command("bash", "-c", command)
	cmd.Stdin = os.Stdin

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	startTime := time.Now()
	execErr := cmd.Run()
	duration := time.Since(startTime).Milliseconds()

	// Output with proper line endings
	outputStr := ""
	if stdout.Len() > 0 {
		outputStr = stdout.String()
		output := strings.ReplaceAll(outputStr, "\n", "\r\n")
		writeTTY(output)
	}
	if stderr.Len() > 0 {
		errOutput := stderr.String()
		if outputStr != "" {
			outputStr += "\n" + errOutput
		} else {
			outputStr = errOutput
		}
		errOutput = strings.ReplaceAll(errOutput, "\n", "\r\n")
		writeTTY(errOutput)
	}

	// Save to database
	if db != nil {
		workingDir, _ := os.Getwd()
		exitCode := 0
		if execErr != nil {
			if exitErr, ok := execErr.(*exec.ExitError); ok {
				exitCode = exitErr.ExitCode()
			} else {
				exitCode = 1
			}
		}

		safeCommand := validator.RedactSecrets(command)

		embedService, _ := ai.NewEmbeddingProvider()
		var embeddingBytes []byte
		if embedService != nil {
			embeddingBytes, _ = embedService.GenerateEmbedding(safeCommand)
		}

		db.SaveCommand(database.Command{
			Command:       safeCommand,
			Timestamp:     time.Now(),
			ExitCode:      exitCode,
			Duration:      duration,
			WorkingDir:    workingDir,
			OutputPreview: outputStr,
			Embedding:     embeddingBytes,
		})
	}

	// NEW: Auto-explain errors
	if execErr != nil {
		writeTTY(fmt.Sprintf("\r\n%s✗ Command failed%s\r\n", red, reset))

		// Offer error explanation
		if stderr.Len() > 0 {
			writeTTY(fmt.Sprintf("\r\n%s▸ Getting error explanation...%s\r\n", cyan, reset))

			explanation, explainErr := client.ExplainError(command, stderr.String(), context)
			if explainErr == nil && strings.TrimSpace(explanation) != "" {
				writeTTY(fmt.Sprintf("\r\n%s╭─ Error Analysis%s\r\n", lightBlue, reset))

				// Split into lines and display with proper formatting
				lines := strings.Split(explanation, "\n")
				for _, line := range lines {
					line = strings.TrimSpace(line)
					if line != "" {
						// Wrap long lines to fit terminal width
						wrappedLines := wrapLine(line, 76) // 76 chars to account for "│  " prefix
						for _, wrappedLine := range wrappedLines {
							writeTTY(fmt.Sprintf("%s│%s  %s\r\n", lightBlue, reset, wrappedLine))
						}
					}
				}

				writeTTY(fmt.Sprintf("%s╰─%s\r\n", lightBlue, reset))
			} else if explainErr != nil {
				writeTTY(fmt.Sprintf("%s⚠ Could not get explanation: %v%s\r\n", gray, explainErr, reset))
			}
		}
		writeTTY("\r\n")
	} else {
		writeTTY(fmt.Sprintf("\r\n%s✓ Command executed successfully%s\r\n\r\n", green, reset))
	}

	// Save conversation turn
	if conversation != nil {
		conversation.AddTurn(query, command, true)
		if err := conversation.Save(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to save conversation: %v\n", err)
		}
	}

	// Learn from executed command
	if context.Preferences != nil {
		context.Preferences.LearnFromCommand(command)
		if err := context.Preferences.Save(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to save preferences: %v\n", err)
		}
	}

	return "", nil
}

func handleAskExplain(command string, client ai.AIProvider, context ai.SystemContext, writeTTY func(string), cyan, lightBlue, red, reset string) (string, error) {
	writeTTY(fmt.Sprintf("\r\n%s▸ Getting explanation...%s\r\n", cyan, reset))
	
	explanation, explainErr := client.ExplainCommand(command, context)
	if explainErr != nil {
		writeTTY(fmt.Sprintf("\r\n%s✗ Failed to get explanation: %v%s\r\n\r\n", red, explainErr, reset))
		return "", nil
	}

	writeTTY(fmt.Sprintf("\r\n%s╭─ Command Explanation%s\r\n", lightBlue, reset))
	lines := strings.Split(explanation, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			wrappedLines := wrapLine(line, 76)
			for _, wrappedLine := range wrappedLines {
				writeTTY(fmt.Sprintf("%s│%s  %s\r\n", lightBlue, reset, wrappedLine))
			}
		}
	}
	writeTTY(fmt.Sprintf("%s╰─%s\r\n\r\n", lightBlue, reset))
	return "", nil
}

func handleAskAlternatives(command string, client ai.AIProvider, context ai.SystemContext, writeTTY func(string), cyan, lightBlue, red, reset string) (string, error) {
	writeTTY(fmt.Sprintf("\r\n%s▸ Getting alternative suggestions...%s\r\n", cyan, reset))
	
	alternatives, altErr := client.SuggestAlternatives(command, context)
	if altErr != nil {
		writeTTY(fmt.Sprintf("\r\n%s✗ Failed to get alternatives: %v%s\r\n\r\n", red, altErr, reset))
		return "", nil
	}

	writeTTY(fmt.Sprintf("\r\n%s╭─ Alternative Commands%s\r\n", lightBlue, reset))
	lines := strings.Split(alternatives, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			wrappedLines := wrapLine(line, 76)
			for _, wrappedLine := range wrappedLines {
				writeTTY(fmt.Sprintf("%s│%s  %s\r\n", lightBlue, reset, wrappedLine))
			}
		}
	}
	writeTTY(fmt.Sprintf("%s╰─%s\r\n\r\n", lightBlue, reset))
	return "", nil
}

func handleAskEdit(query, command string, db *database.DB, writeTTY func(string), cyan, lightBlue, green, red, gray, reset string) (string, error) {
	writeTTY(fmt.Sprintf("\r\n%s▸ Edit command (press Enter when done):%s\r\n", cyan, reset))
	writeTTY(fmt.Sprintf("%s> %s", gray, reset))
	
	// Use simple line editor
	editedCommand, editErr := readLineFromTTY(command)
	if editErr != nil {
		writeTTY(fmt.Sprintf("\r\n%s✗ Edit failed: %v%s\r\n\r\n", red, editErr, reset))
		return "", nil
	}

	editedCommand = strings.TrimSpace(editedCommand)
	if editedCommand == "" {
		writeTTY(fmt.Sprintf("\r\n%sℹ Cancelled%s\r\n\r\n", gray, reset))
		return "", nil
	}

	// Show edited command
	writeTTY(fmt.Sprintf("\r\n%s╭─ Edited Command%s\r\n", lightBlue, reset))
	writeTTY(fmt.Sprintf("%s│%s  %s%s%s\r\n", lightBlue, reset, cyan, editedCommand, reset))
	writeTTY(fmt.Sprintf("%s╰─%s\r\n", lightBlue, reset))

	// Execute edited command
	command = editedCommand // Update command variable
	writeTTY(fmt.Sprintf("\r\n%s▸ Executing...%s\r\n\r\n", cyan, reset))

	cmd := exec.Command("bash", "-c", command)
	cmd.Stdin = os.Stdin

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	startTime := time.Now()
	execErr := cmd.Run()
	duration := time.Since(startTime).Milliseconds()

	// Output with proper line endings
	outputStr := ""
	if stdout.Len() > 0 {
		outputStr = stdout.String()
		output := strings.ReplaceAll(outputStr, "\n", "\r\n")
		writeTTY(output)
	}
	if stderr.Len() > 0 {
		errOutput := stderr.String()
		if outputStr != "" {
			outputStr += "\n" + errOutput
		} else {
			outputStr = errOutput
		}
		errOutput = strings.ReplaceAll(errOutput, "\n", "\r\n")
		writeTTY(errOutput)
	}

	// Save to database
	if db != nil {
		workingDir, _ := os.Getwd()
		exitCode := 0
		if execErr != nil {
			if exitErr, ok := execErr.(*exec.ExitError); ok {
				exitCode = exitErr.ExitCode()
			} else {
				exitCode = 1
			}
		}

		safeCommand := validator.RedactSecrets(command)

		embedService, _ := ai.NewEmbeddingProvider()
		var embeddingBytes []byte
		if embedService != nil {
			embeddingBytes, _ = embedService.GenerateEmbedding(safeCommand)
		}

		db.SaveCommand(database.Command{
			Command:       safeCommand,
			Timestamp:     time.Now(),
			ExitCode:      exitCode,
			Duration:      duration,
			WorkingDir:    workingDir,
			OutputPreview: outputStr,
			Embedding:     embeddingBytes,
		})
	}

	if execErr != nil {
		writeTTY(fmt.Sprintf("\r\n%s✗ Command failed%s\r\n\r\n", red, reset))
	} else {
		writeTTY(fmt.Sprintf("\r\n%s✓ Command executed successfully%s\r\n\r\n", green, reset))
	}
	return "", nil
}

func handleAskCopy(query, command string, conversation *ai.ConversationHistory, writeTTY func(string), green, red, reset string) (string, error) {
	if clipboard.WriteAll(command) == nil {
		writeTTY(fmt.Sprintf("\r\n%s✓ Copied to clipboard!%s\r\n\r\n", green, reset))
	} else {
		writeTTY(fmt.Sprintf("\r\n%s✗ Failed to copy to clipboard%s\r\n\r\n", red, reset))
	}
	
	// Save conversation turn (but mark as not executed since only copied)
	if conversation != nil {
		conversation.AddTurn(query, command, false)
		if err := conversation.Save(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to save conversation: %v\n", err)
		}
	}
	
	return "", nil
}

func handleAskCancel(query, command string, conversation *ai.ConversationHistory, writeTTY func(string), gray, reset string) (string, error) {
	writeTTY(fmt.Sprintf("\r\n%sℹ Cancelled%s\r\n\r\n", gray, reset))
	
	// Still save the conversation turn (but mark as not executed)
	if conversation != nil {
		conversation.AddTurn(query, command, false)
		if err := conversation.Save(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to save conversation: %v\n", err)
		}
	}
	
	return "", nil
}
