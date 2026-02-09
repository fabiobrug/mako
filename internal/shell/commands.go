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
			return true, fmt.Sprintf("v0.2.0\n"), nil
		case "draw":
			return true, getSharkArt(), nil
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

	command, err := client.GenerateCommand(query, context)
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
			"Edit before running|edit",
			"Copy to clipboard|copy",
			"Cancel|cancel",
		)
	} else {
		menuArgs = append(menuArgs,
			"Run command|run",
			"Explain what this does|explain",
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
		if stdout.Len() > 0 {
			output := stdout.String()
			output = strings.ReplaceAll(output, "\n", "\r\n")
			writeTTY(output)
		}
		if stderr.Len() > 0 {
			errOutput := stderr.String()
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
				Command:    safeCommand,
				Timestamp:  time.Now(),
				ExitCode:   exitCode,
				Duration:   duration,
				WorkingDir: workingDir,
				Embedding:  embeddingBytes,
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

	case "edit":
		writeTTY(fmt.Sprintf("\r\n%s▸ Edit command (press Enter when done):%s\r\n", cyan, reset))
		writeTTY(fmt.Sprintf("%s> %s", gray, reset))
		
		// Use a simple line editor
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
		if stdout.Len() > 0 {
			output := stdout.String()
			output = strings.ReplaceAll(output, "\n", "\r\n")
			writeTTY(output)
		}
		if stderr.Len() > 0 {
			errOutput := stderr.String()
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
				Command:    safeCommand,
				Timestamp:  time.Now(),
				ExitCode:   exitCode,
				Duration:   duration,
				WorkingDir: workingDir,
				Embedding:  embeddingBytes,
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
		return "", nil

	case "cancel":
		writeTTY(fmt.Sprintf("\r\n%sℹ Cancelled%s\r\n\r\n", gray, reset))
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
	if len(args) > 0 && args[0] == "semantic" {
		if len(args) < 2 {
			return fmt.Sprintf("\n%sUsage:%s mako history semantic <query>\n\n", lightBlue, reset), nil
		}
		return handleSemanticHistory(strings.Join(args[1:], " "), db)
	}
	if len(args) == 0 {
		commands, err := db.GetRecentCommands(10)
		if err != nil {
			return "", err
		}
		if len(commands) == 0 {
			return fmt.Sprintf("\n%sNo command history yet%s\n\n", dimBlue, reset), nil
		}
		var output strings.Builder
		output.WriteString(fmt.Sprintf("\n%s╭─ Recent Commands%s\n", lightBlue, reset))
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
		}
		output.WriteString(fmt.Sprintf("%s╰─%s\n\n", lightBlue, reset))
		return output.String(), nil
	}
	query := strings.Join(args, " ")
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
	}
	output.WriteString(fmt.Sprintf("%s╰─%s\n\n", lightBlue, reset))
	return output.String(), nil
}

func handleSemanticHistory(query string, db *database.DB) (string, error) {
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
	if len(commands) == 0 {
		return fmt.Sprintf("\n%sNo similar commands found for:%s %s\n\n", lightBlue, reset, query), nil
	}
	var output strings.Builder
	output.WriteString(fmt.Sprintf("\n%s╭─ Found %d similar commands for '%s'%s\n", lightBlue, len(commands), query, reset))
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
	}
	output.WriteString(fmt.Sprintf("%s╰─%s\n\n", lightBlue, reset))
	return output.String(), nil
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

func handleAlias(args []string, db *database.DB) (string, error) {
	lightBlue := "\033[38;2;93;173;226m"
	cyan := "\033[38;2;0;209;255m"
	green := "\033[38;2;100;255;100m"
	red := "\033[38;2;255;100;100m"
	dimBlue := "\033[38;2;120;150;180m"
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
			return fmt.Sprintf("\n%sUsage:%s mako alias save <name> <command>\n\n", lightBlue, reset), nil
		}
		name := args[1]
		command := strings.Join(args[2:], " ")

		if err := store.Set(name, command); err != nil {
			return "", err
		}

		return fmt.Sprintf("\n%s✓ Saved alias '%s':%s %s\n\n", green, name, reset, command), nil

	case "list":
		aliases := store.List()
		if len(aliases) == 0 {
			return fmt.Sprintf("\n%sNo aliases saved yet%s\n\n", dimBlue, reset), nil
		}

		var output strings.Builder
		output.WriteString(fmt.Sprintf("\n%s╭─ Saved Aliases%s\n", lightBlue, reset))
		for name, command := range aliases {
			output.WriteString(fmt.Sprintf("%s│%s  %s%s%s → %s\n",
				lightBlue, reset,
				cyan, name, reset,
				command))
		}
		output.WriteString(fmt.Sprintf("%s╰─%s\n\n", lightBlue, reset))
		return output.String(), nil

	case "delete":
		if len(args) < 2 {
			return fmt.Sprintf("\n%sUsage:%s mako alias delete <name>\n\n", lightBlue, reset), nil
		}
		name := args[1]

		if err := store.Delete(name); err != nil {
			return fmt.Sprintf("\n%s✗ %v%s\n\n", red, err, reset), nil
		}

		return fmt.Sprintf("\n%s✓ Deleted alias '%s'%s\n\n", green, name, reset), nil

	case "run":
		if len(args) < 2 {
			return fmt.Sprintf("\n%sUsage:%s mako alias run <name>\n\n", lightBlue, reset), nil
		}
		name := args[1]

		command, ok := store.Get(name)
		if !ok {
			return fmt.Sprintf("\n%s✗ Alias '%s' not found%s\n\n", red, name, reset), nil
		}

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
		if stdout.Len() > 0 {
			output := stdout.String()
			output = strings.ReplaceAll(output, "\n", "\r\n")
			writeTTY(output)
		}
		if stderr.Len() > 0 {
			errOutput := stderr.String()
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
				Command:    safeCommand,
				Timestamp:  time.Now(),
				ExitCode:   exitCode,
				Duration:   duration,
				WorkingDir: workingDir,
				Embedding:  embeddingBytes,
			})
		}

		if execErr != nil {
			writeTTY(fmt.Sprintf("\r\n%s✗ Command failed%s\r\n\r\n", red, reset))
		} else {
			writeTTY(fmt.Sprintf("\r\n%s✓ Command executed successfully%s\r\n\r\n", green, reset))
		}
		return "", nil

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
%s│%s  %smako history%s                     Show recent commands
%s│%s  %smako history <keyword>%s           Search by keyword
%s│%s  %smako history semantic <query>%s    Search by meaning
%s│%s  %smako alias save <name> <cmd>%s     Save a command alias
%s│%s  %smako alias list%s                  List all saved aliases
%s│%s  %smako alias run <name>%s            Run a saved alias
%s│%s  %smako alias delete <name>%s         Delete an alias
%s│%s  %smako stats%s                       Show statistics
%s│%s  %smako help%s                        Show this help
%s│%s  %smako version%s                     Show Mako version
%s│%s 
%s╰─%s %sRegular shell commands work normally!%s

`, lightBlue, reset,
		lightBlue, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
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
