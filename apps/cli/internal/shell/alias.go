package shell

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/fabiobrug/mako.git/internal/ai"
	"github.com/fabiobrug/mako.git/internal/alias"
	"github.com/fabiobrug/mako.git/internal/database"
)

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
