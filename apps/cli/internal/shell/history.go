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
)

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
	
	embedService, err := ai.NewEmbeddingProvider()
	if err != nil {
		return "", err
	}
	queryBytes, err := embedService.GenerateEmbedding(query)
	if err != nil {
		return "", err
	}
	commands, err := db.SearchCommandsSemantic(query, queryBytes, 10, 0.5)
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
