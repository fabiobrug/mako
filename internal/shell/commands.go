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
	"github.com/fabiobrug/mako.git/internal/alias"
	"github.com/fabiobrug/mako.git/internal/database"
	"github.com/fabiobrug/mako.git/internal/safety"
)

var validator = safety.NewValidator()

// Global reference to ring buffer (will be set from main)
var recentOutputGetter func(int) []string

// SetRecentOutputGetter allows main to provide ring buffer access
func SetRecentOutputGetter(getter func(int) []string) {
	recentOutputGetter = getter
}

// readLineFromTTY reads a line of input from /dev/tty with the prefilled text
func readLineFromTTY(prefill string) (string, error) {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return "", err
	}
	defer tty.Close()

	// Write prefilled text
	tty.WriteString(prefill)

	// Read until newline - initialize with prefilled text
	result := []byte(prefill)
	buf := make([]byte, 1)
	
	for {
		n, err := tty.Read(buf)
		if err != nil {
			return "", err
		}
		if n > 0 {
			if buf[0] == '\n' || buf[0] == '\r' {
				break
			}
			// Handle backspace
			if buf[0] == 127 || buf[0] == 8 {
				if len(result) > 0 {
					result = result[:len(result)-1]
					// Visual backspace: move back, print space, move back again
					tty.WriteString("\b \b")
				}
			} else {
				result = append(result, buf[0])
				// Echo the character
				tty.Write(buf)
			}
		}
	}
	
	return string(result), nil
}

// wrapLine wraps a line of text to fit within maxWidth characters
func wrapLine(text string, maxWidth int) []string {
	if len(text) <= maxWidth {
		return []string{text}
	}

	var lines []string
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{text}
	}

	currentLine := ""
	for _, word := range words {
		// If adding this word would exceed maxWidth, start a new line
		if len(currentLine)+len(word)+1 > maxWidth {
			if currentLine != "" {
				lines = append(lines, currentLine)
				currentLine = word
			} else {
				// Single word is longer than maxWidth, split it
				lines = append(lines, word[:maxWidth])
				currentLine = word[maxWidth:]
			}
		} else {
			if currentLine == "" {
				currentLine = word
			} else {
				currentLine += " " + word
			}
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

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
			output, err := handleAsk(query, db)
			return true, output, err
		case "history":
			output, err := handleHistory(parts[2:], db)
			return true, output, err
		case "stats":
			output, err := handleStats(db)
			return true, output, err
		case "alias":
			output, err := handleAlias(parts[2:], db)
			return true, output, err
		case "help":
			return true, getHelpText(), nil
		case "v", "version":
			return true, fmt.Sprintf("v0.3.0\n"), nil
		case "draw":
			return true, getSharkArt(), nil
		case "clear":
			output, err := handleClear()
			return true, output, err
		default:
			return true, fmt.Sprintf("Unknown mako command: %s\n", parts[1]), nil
		}
	}

	return false, "", nil
}

func handleAsk(query string, db *database.DB) (string, error) {
	client, err := ai.NewGeminiClient()
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

	// Short delay to ensure main goroutine detects pause file
	time.Sleep(75 * time.Millisecond)

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
	time.Sleep(100 * time.Millisecond)

	// Handle choice
	switch choice {
	case "run":
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

			embedService, _ := ai.NewEmbeddingService()
			var embeddingBytes []byte
			if embedService != nil {
				vec, embedErr := embedService.Embed(safeCommand)
				if embedErr == nil {
					embeddingBytes = ai.VectorToBytes(vec)
				}
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
				writeTTY(fmt.Sprintf("\r\n%s Getting error explanation...%s\r\n", cyan, reset))

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

	case "explain":
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

	case "alternatives":
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

	case "edit":
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

			embedService, _ := ai.NewEmbeddingService()
			var embeddingBytes []byte
			if embedService != nil {
				vec, embedErr := embedService.Embed(safeCommand)
				if embedErr == nil {
					embeddingBytes = ai.VectorToBytes(vec)
				}
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

	case "copy":
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

	case "cancel":
		writeTTY(fmt.Sprintf("\r\n%sℹ Cancelled%s\r\n\r\n", gray, reset))
		
		// Still save the conversation turn (but mark as not executed)
		if conversation != nil {
			conversation.AddTurn(query, command, false)
			if err := conversation.Save(); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to save conversation: %v\n", err)
			}
		}
		
		return "", nil

	default:
		return "", nil
	}
}

func handleHistory(args []string, db *database.DB) (string, error) {
	lightBlue := "\033[38;2;93;173;226m"
	green := "\033[38;2;100;255;100m"
	red := "\033[38;2;255;100;100m"
	dimBlue := "\033[38;2;120;150;180m"
	gray := "\033[38;2;150;150;150m"
	reset := "\033[0m"
	
	if db == nil {
		return fmt.Sprintf("\n%s✗ Database not available%s\n\n", dimBlue, reset), nil
	}
	
	// Check for flags
	filterFailed := false
	filterSuccess := false
	interactive := false
	var filterArgs []string
	
	for _, arg := range args {
		if arg == "--failed" {
			filterFailed = true
		} else if arg == "--success" {
			filterSuccess = true
		} else if arg == "--interactive" || arg == "--browse" {
			interactive = true
		} else {
			filterArgs = append(filterArgs, arg)
		}
	}
	
	// Handle interactive mode
	if interactive {
		return handleInteractiveHistory(db, filterFailed, filterSuccess)
	}
	
	if len(filterArgs) > 0 && filterArgs[0] == "semantic" {
		if len(filterArgs) < 2 {
			return fmt.Sprintf("\n%sUsage:%s mako history semantic <query> [--failed|--success]\n\n", lightBlue, reset), nil
		}
		return handleSemanticHistory(strings.Join(filterArgs[1:], " "), db, filterFailed, filterSuccess)
	}
	
	if len(filterArgs) == 0 {
		var commands []database.Command
		var err error
		var title string
		
		if filterFailed {
			commands, err = db.GetCommandsByExitCode(false, 10)
			title = "Failed Commands"
		} else if filterSuccess {
			commands, err = db.GetCommandsByExitCode(true, 10)
			title = "Successful Commands"
		} else {
			commands, err = db.GetRecentCommands(10)
			title = "Recent Commands"
		}
		
		if err != nil {
			return "", err
		}
		if len(commands) == 0 {
			return fmt.Sprintf("\n%sNo command history yet%s\n\n", dimBlue, reset), nil
		}
		var output strings.Builder
		output.WriteString(fmt.Sprintf("\n%s╭─ %s%s\n", lightBlue, title, reset))
		for _, cmd := range commands {
			// Status icon
			statusIcon := fmt.Sprintf("%s✓%s", green, reset)
			if cmd.ExitCode != 0 {
				statusIcon = fmt.Sprintf("%s✗%s", red, reset)
			}
			
			// Format duration
			durationStr := fmt.Sprintf("%dms", cmd.Duration)
			if cmd.Duration >= 1000 {
				durationStr = fmt.Sprintf("%.1fs", float64(cmd.Duration)/1000.0)
			}
			
			output.WriteString(fmt.Sprintf("%s│%s  %s %s[%s]%s %s%-6s%s %s\n",
				lightBlue, reset,
				statusIcon,
				dimBlue, cmd.Timestamp.Format("15:04:05"), reset,
				gray, durationStr, reset,
				cmd.Command))
			
			// Add output preview if available
			if cmd.OutputPreview != "" {
				preview := cmd.OutputPreview
				if len(preview) > 60 {
					preview = preview[:60] + "..."
				}
				// Replace newlines with spaces for single-line preview
				preview = strings.ReplaceAll(preview, "\n", " ")
				preview = strings.ReplaceAll(preview, "\r", "")
				output.WriteString(fmt.Sprintf("%s│%s    %s↳ %s%s\n",
					lightBlue, reset,
					gray, preview, reset))
			}
		}
		output.WriteString(fmt.Sprintf("%s╰─%s\n\n", lightBlue, reset))
		return output.String(), nil
	}
	query := strings.Join(filterArgs, " ")
	commands, err := db.SearchCommands(query, 10)
	if err != nil {
		return "", err
	}
	if len(commands) == 0 {
		return fmt.Sprintf("\n%sNo commands found matching:%s %s\n\n", lightBlue, reset, query), nil
	}
	var output strings.Builder
	output.WriteString(fmt.Sprintf("\n%s╭─ Found %d commands matching '%s'%s\n", lightBlue, len(commands), query, reset))
	for _, cmd := range commands {
		// Status icon
		statusIcon := fmt.Sprintf("%s✓%s", green, reset)
		if cmd.ExitCode != 0 {
			statusIcon = fmt.Sprintf("%s✗%s", red, reset)
		}
		
		// Format duration
		durationStr := fmt.Sprintf("%dms", cmd.Duration)
		if cmd.Duration >= 1000 {
			durationStr = fmt.Sprintf("%.1fs", float64(cmd.Duration)/1000.0)
		}
		
		output.WriteString(fmt.Sprintf("%s│%s  %s %s[%s]%s %s%-6s%s %s\n",
			lightBlue, reset,
			statusIcon,
			dimBlue, cmd.Timestamp.Format("15:04:05"), reset,
			gray, durationStr, reset,
			cmd.Command))
		
		// Add output preview if available
		if cmd.OutputPreview != "" {
			preview := cmd.OutputPreview
			if len(preview) > 60 {
				preview = preview[:60] + "..."
			}
			preview = strings.ReplaceAll(preview, "\n", " ")
			preview = strings.ReplaceAll(preview, "\r", "")
			output.WriteString(fmt.Sprintf("%s│%s    %s↳ %s%s\n",
				lightBlue, reset,
				gray, preview, reset))
		}
	}
	output.WriteString(fmt.Sprintf("%s╰─%s\n\n", lightBlue, reset))
	return output.String(), nil
}

func handleSemanticHistory(query string, db *database.DB, filterFailed bool, filterSuccess bool) (string, error) {
	lightBlue := "\033[38;2;93;173;226m"
	green := "\033[38;2;100;255;100m"
	red := "\033[38;2;255;100;100m"
	dimBlue := "\033[38;2;120;150;180m"
	gray := "\033[38;2;150;150;150m"
	reset := "\033[0m"
	
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
	
	// Apply exit code filtering
	if filterFailed || filterSuccess {
		var filtered []database.Command
		for _, cmd := range commands {
			if filterFailed && cmd.ExitCode != 0 {
				filtered = append(filtered, cmd)
			} else if filterSuccess && cmd.ExitCode == 0 {
				filtered = append(filtered, cmd)
			}
		}
		commands = filtered
	}
	
	if len(commands) == 0 {
		return fmt.Sprintf("\n%sNo similar commands found for:%s %s\n\n", lightBlue, reset, query), nil
	}
	
	var output strings.Builder
	titleSuffix := ""
	if filterFailed {
		titleSuffix = " (failed only)"
	} else if filterSuccess {
		titleSuffix = " (successful only)"
	}
	output.WriteString(fmt.Sprintf("\n%s╭─ Found %d similar commands for '%s'%s%s\n", lightBlue, len(commands), query, titleSuffix, reset))
	for _, cmd := range commands {
		// Status icon
		statusIcon := fmt.Sprintf("%s✓%s", green, reset)
		if cmd.ExitCode != 0 {
			statusIcon = fmt.Sprintf("%s✗%s", red, reset)
		}
		
		// Format duration
		durationStr := fmt.Sprintf("%dms", cmd.Duration)
		if cmd.Duration >= 1000 {
			durationStr = fmt.Sprintf("%.1fs", float64(cmd.Duration)/1000.0)
		}
		
		output.WriteString(fmt.Sprintf("%s│%s  %s %s[%s]%s %s%-6s%s %s\n",
			lightBlue, reset,
			statusIcon,
			dimBlue, cmd.Timestamp.Format("15:04:05"), reset,
			gray, durationStr, reset,
			cmd.Command))
		
		// Add output preview if available
		if cmd.OutputPreview != "" {
			preview := cmd.OutputPreview
			if len(preview) > 60 {
				preview = preview[:60] + "..."
			}
			preview = strings.ReplaceAll(preview, "\n", " ")
			preview = strings.ReplaceAll(preview, "\r", "")
			output.WriteString(fmt.Sprintf("%s│%s    %s↳ %s%s\n",
				lightBlue, reset,
				gray, preview, reset))
		}
	}
	output.WriteString(fmt.Sprintf("%s╰─%s\n\n", lightBlue, reset))
	return output.String(), nil
}

func handleInteractiveHistory(db *database.DB, filterFailed bool, filterSuccess bool) (string, error) {
	lightBlue := "\033[38;2;93;173;226m"
	cyan := "\033[38;2;0;209;255m"
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

	// Get commands
	var commands []database.Command
	var err error
	if filterFailed {
		commands, err = db.GetCommandsByExitCode(false, 50)
	} else if filterSuccess {
		commands, err = db.GetCommandsByExitCode(true, 50)
	} else {
		commands, err = db.GetRecentCommands(50)
	}
	
	if err != nil {
		return "", err
	}
	if len(commands) == 0 {
		return fmt.Sprintf("\r\n%sNo commands in history%s\r\n\r\n", gray, reset), nil
	}

	// Build menu options
	menuArgs := []string{
		fmt.Sprintf("%sSelect a command:%s", lightBlue, reset),
	}

	for i, cmd := range commands {
		if i >= 20 { // Limit to 20 items for usability
			break
		}
		
		statusIcon := "✓"
		if cmd.ExitCode != 0 {
			statusIcon = "✗"
		}
		
		durationStr := fmt.Sprintf("%dms", cmd.Duration)
		if cmd.Duration >= 1000 {
			durationStr = fmt.Sprintf("%.1fs", float64(cmd.Duration)/1000.0)
		}
		
		label := fmt.Sprintf("%s [%s] %s - %s", 
			statusIcon,
			cmd.Timestamp.Format("15:04:05"),
			durationStr,
			cmd.Command)
		
		if len(label) > 70 {
			label = label[:67] + "..."
		}
		
		menuArgs = append(menuArgs, fmt.Sprintf("%s|%d", label, i))
	}
	
	menuArgs = append(menuArgs, "Cancel|cancel")

	// Pause PTY input
	pauseFile := filepath.Join(os.Getenv("HOME"), ".mako", "pause_input")
	os.WriteFile(pauseFile, []byte("1"), 0644)
	defer os.Remove(pauseFile)

	time.Sleep(75 * time.Millisecond)

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
	time.Sleep(100 * time.Millisecond)

	if choice == "cancel" {
		writeTTY(fmt.Sprintf("\r\n%sℹ Cancelled%s\r\n\r\n", gray, reset))
		return "", nil
	}

	// Parse selected index
	var selectedIdx int
	fmt.Sscanf(choice, "%d", &selectedIdx)
	
	if selectedIdx < 0 || selectedIdx >= len(commands) {
		writeTTY(fmt.Sprintf("\r\n%s✗ Invalid selection%s\r\n\r\n", red, reset))
		return "", nil
	}

	selectedCmd := commands[selectedIdx]

	// Show action menu
	actionMenuArgs := []string{
		fmt.Sprintf("%sWhat would you like to do?%s", lightBlue, reset),
		"Run this command again|run",
		"Copy to clipboard|copy",
		"View full output|output",
		"Cancel|cancel",
	}

	menuCmd = exec.Command(menuPath, actionMenuArgs...)
	menuCmd.Stderr = os.Stderr

	actionBytes, err := menuCmd.Output()
	if err != nil {
		return "", fmt.Errorf("menu failed: %w", err)
	}

	action := strings.TrimSpace(string(actionBytes))
	time.Sleep(100 * time.Millisecond)

	switch action {
	case "run":
		writeTTY(fmt.Sprintf("\r\n%s▸ Executing: %s%s%s\r\n\r\n", cyan, cyan, selectedCmd.Command, reset))

		cmd := exec.Command("bash", "-c", selectedCmd.Command)
		cmd.Stdin = os.Stdin

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		startTime := time.Now()
		execErr := cmd.Run()
		duration := time.Since(startTime).Milliseconds()

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

			safeCommand := validator.RedactSecrets(selectedCmd.Command)

			embedService, _ := ai.NewEmbeddingService()
			var embeddingBytes []byte
			if embedService != nil {
				vec, embedErr := embedService.Embed(safeCommand)
				if embedErr == nil {
					embeddingBytes = ai.VectorToBytes(vec)
				}
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

	case "copy":
		if clipboard.WriteAll(selectedCmd.Command) == nil {
			writeTTY(fmt.Sprintf("\r\n%s✓ Copied to clipboard!%s\r\n\r\n", green, reset))
		} else {
			writeTTY(fmt.Sprintf("\r\n%s✗ Failed to copy to clipboard%s\r\n\r\n", red, reset))
		}
		return "", nil

	case "output":
		if selectedCmd.OutputPreview != "" {
			writeTTY(fmt.Sprintf("\r\n%s╭─ Command Output%s\r\n", lightBlue, reset))
			output := selectedCmd.OutputPreview
			output = strings.ReplaceAll(output, "\n", "\r\n")
			lines := strings.Split(output, "\r\n")
			for _, line := range lines {
				if len(lines) <= 50 || line != "" { // Show all for small outputs
					writeTTY(fmt.Sprintf("%s│%s  %s\r\n", lightBlue, reset, line))
				}
			}
			writeTTY(fmt.Sprintf("%s╰─%s\r\n\r\n", lightBlue, reset))
		} else {
			writeTTY(fmt.Sprintf("\r\n%sNo output recorded%s\r\n\r\n", gray, reset))
		}
		return "", nil

	case "cancel":
		writeTTY(fmt.Sprintf("\r\n%sℹ Cancelled%s\r\n\r\n", gray, reset))
		return "", nil

	default:
		return "", nil
	}
}

func handleStats(db *database.DB) (string, error) {
	lightBlue := "\033[38;2;93;173;226m"
	cyan := "\033[38;2;0;209;255m"
	dimBlue := "\033[38;2;120;150;180m"
	reset := "\033[0m"
	if db == nil {
		return fmt.Sprintf("\n%s✗ Database not available%s\n\n", dimBlue, reset), nil
	}
	stats, err := db.GetStats()
	if err != nil {
		return "", err
	}
	var output strings.Builder
	output.WriteString(fmt.Sprintf("\n%s╭─ Mako Statistics%s\n", lightBlue, reset))
	output.WriteString(fmt.Sprintf("%s│%s  Total commands    %s%d%s\n", lightBlue, reset, cyan, stats["total_commands"], reset))
	output.WriteString(fmt.Sprintf("%s│%s  Commands today    %s%d%s\n", lightBlue, reset, cyan, stats["commands_today"], reset))
	output.WriteString(fmt.Sprintf("%s│%s  Avg duration      %s%.0fms%s\n", lightBlue, reset, cyan, stats["avg_duration_ms"], reset))
	output.WriteString(fmt.Sprintf("%s╰─%s\n\n", lightBlue, reset))
	return output.String(), nil
}

func handleClear() (string, error) {
	green := "\033[38;2;100;255;100m"
	reset := "\033[0m"
	
	if err := ai.ClearConversation(); err != nil {
		return "", fmt.Errorf("failed to clear conversation: %w", err)
	}
	
	return fmt.Sprintf("\n%s✓ Conversation history cleared%s\n\n", green, reset), nil
}

func handleAlias(args []string, db *database.DB) (string, error) {
	lightBlue := "\033[38;2;93;173;226m"
	cyan := "\033[38;2;0;209;255m"
	green := "\033[38;2;100;255;100m"
	red := "\033[38;2;255;100;100m"
	dimBlue := "\033[38;2;120;150;180m"
	gray := "\033[38;2;150;150;150m"
	reset := "\033[0m"

	store, err := alias.NewAliasStore()
	if err != nil {
		return "", err
	}

	if len(args) == 0 {
		return fmt.Sprintf("\n%sUsage:%s mako alias <save|list|delete|run> [args]\n\n", lightBlue, reset), nil
	}

	subcommand := args[0]

	switch subcommand {
	case "save":
		if len(args) < 3 {
			return fmt.Sprintf("\n%sUsage:%s mako alias save <name> <command> [--tags tag1,tag2,...]\n\n", lightBlue, reset), nil
		}
		name := args[1]
		
		// Parse command and tags
		var command string
		var tags []string
		
		// Find --tags flag if present
		tagsIdx := -1
		for i, arg := range args[2:] {
			if arg == "--tags" {
				tagsIdx = i + 2
				break
			}
		}
		
		if tagsIdx > 0 && tagsIdx < len(args)-1 {
			command = strings.Join(args[2:tagsIdx], " ")
			tagStr := strings.Join(args[tagsIdx+1:], " ")
			tags = strings.Split(tagStr, ",")
			// Trim spaces from tags
			for i, tag := range tags {
				tags[i] = strings.TrimSpace(tag)
			}
		} else {
			command = strings.Join(args[2:], " ")
			tags = []string{}
		}

		if err := store.Set(name, command, tags); err != nil {
			return "", err
		}

		tagStr := ""
		if len(tags) > 0 {
			tagStr = fmt.Sprintf(" %s[tags: %s]%s", dimBlue, strings.Join(tags, ", "), reset)
		}
		result := fmt.Sprintf("\r\n%s✓ Saved alias '%s':%s %s%s", green, name, reset, command, tagStr)
		
		// Check if parameters are missing (user might have forgotten to escape)
		if strings.Contains(command, "cmd/mako/") && !strings.Contains(command, "$") {
			result += fmt.Sprintf("\r\n%sℹ Note: If you want to use $1, $2, etc., wrap your command in single quotes:%s", gray, reset)
			result += fmt.Sprintf("\r\n%s  mako alias save %s '%s'%s", gray, name, strings.ReplaceAll(args[2], "\"", "'"), reset)
		}
		result += "\r\n\r\n"
		return result, nil

	case "list":
		// Check for --tag filter
		var filterTag string
		if len(args) > 1 && args[1] == "--tag" && len(args) > 2 {
			filterTag = args[2]
		}

		var aliases map[string]alias.AliasInfo
		if filterTag != "" {
			aliases = store.ListByTag(filterTag)
		} else {
			aliases = store.List()
		}

		if len(aliases) == 0 {
			if filterTag != "" {
				result := fmt.Sprintf("\r\n%sNo aliases with tag '%s'%s\r\n\r\n", dimBlue, filterTag, reset)
				return result, nil
			}
			result := fmt.Sprintf("\r\n%sNo aliases saved yet%s\r\n\r\n", dimBlue, reset)
			return result, nil
		}

		var output strings.Builder
		if filterTag != "" {
			output.WriteString(fmt.Sprintf("\r\n%s╭─ Aliases tagged '%s'%s\r\n", lightBlue, filterTag, reset))
		} else {
			output.WriteString(fmt.Sprintf("\r\n%s╭─ Saved Aliases%s\r\n", lightBlue, reset))
		}
		
		for name, info := range aliases {
			tagStr := ""
			if len(info.Tags) > 0 {
				tagStr = fmt.Sprintf(" %s[%s]%s", dimBlue, strings.Join(info.Tags, ", "), reset)
			}
			output.WriteString(fmt.Sprintf("%s│%s  %s%s%s → %s%s\r\n",
				lightBlue, reset,
				cyan, name, reset,
				info.Command, tagStr))
		}
		output.WriteString(fmt.Sprintf("%s╰─%s\r\n\r\n", lightBlue, reset))
		return output.String(), nil

	case "delete":
		if len(args) < 2 {
			return fmt.Sprintf("\r\n%sUsage:%s mako alias delete <name>\r\n\r\n", lightBlue, reset), nil
		}
		name := args[1]

		if err := store.Delete(name); err != nil {
			result := fmt.Sprintf("\r\n%s✗ %v%s\r\n\r\n", red, err, reset)
			return result, nil
		}

		result := fmt.Sprintf("\r\n%s✓ Deleted alias '%s'%s\r\n\r\n", green, name, reset)
		return result, nil

	case "run":
		if len(args) < 2 {
			return fmt.Sprintf("\n%sUsage:%s mako alias run <name> [args...]\n\n", lightBlue, reset), nil
		}
		name := args[1]
		aliasArgs := args[2:] // Extra arguments for parameter substitution

		command, ok := store.Get(name)
		if !ok {
			return fmt.Sprintf("\n%s✗ Alias '%s' not found%s\n\n", red, name, reset), nil
		}

		// Expand parameters ($1, $2, $@, etc.)
		command = alias.ExpandParameters(command, aliasArgs)

		// Execute the aliased command
		tty, _ := os.OpenFile("/dev/tty", os.O_WRONLY, 0)
		if tty != nil {
			defer tty.Close()
		}
		writeTTY := func(s string) {
			if tty != nil {
				fmt.Fprint(tty, s)
			}
		}

		writeTTY(fmt.Sprintf("\r\n%s╭─ Running Alias '%s'%s\r\n", lightBlue, name, reset))
		writeTTY(fmt.Sprintf("%s│%s  %s%s%s\r\n", lightBlue, reset, cyan, command, reset))
		writeTTY(fmt.Sprintf("%s╰─%s\r\n\r\n", lightBlue, reset))

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

			embedService, _ := ai.NewEmbeddingService()
			var embeddingBytes []byte
			if embedService != nil {
				vec, embedErr := embedService.Embed(safeCommand)
				if embedErr == nil {
					embeddingBytes = ai.VectorToBytes(vec)
				}
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

	case "export":
		if len(args) < 2 {
			return fmt.Sprintf("\r\n%sUsage:%s mako alias export <filepath>\r\n\r\n", lightBlue, reset), nil
		}
		exportPath := args[1]

		if err := store.ExportToFile(exportPath); err != nil {
			result := fmt.Sprintf("\r\n%s✗ Export failed: %v%s\r\n\r\n", red, err, reset)
			return result, nil
		}

		result := fmt.Sprintf("\r\n%s✓ Exported aliases to '%s'%s\r\n\r\n", green, exportPath, reset)
		return result, nil

	case "import":
		if len(args) < 2 {
			return fmt.Sprintf("\r\n%sUsage:%s mako alias import <filepath>\r\n\r\n", lightBlue, reset), nil
		}
		importPath := args[1]

		if err := store.ImportFromFile(importPath); err != nil {
			result := fmt.Sprintf("\r\n%s✗ Import failed: %v%s\r\n\r\n", red, err, reset)
			return result, nil
		}

		result := fmt.Sprintf("\r\n%s✓ Imported aliases from '%s'%s\r\n\r\n", green, importPath, reset)
		return result, nil

	default:
		return fmt.Sprintf("\n%sUnknown alias subcommand: %s%s\n\n", red, subcommand, reset), nil
	}
}

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
%s│%s  %smako stats%s                       Show statistics
%s│%s  %smako clear%s                       Clear conversation history
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
