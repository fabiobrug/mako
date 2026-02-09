package ai

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type SystemContext struct {
	OS             string
	Shell          string
	CurrentDir     string
	RecentOutput   []string
	RecentCommands []string
	WorkingFiles   []string
}

func GetSystemContext(recentOutput []string) SystemContext {
	ctx := SystemContext{
		OS:           runtime.GOOS + "/" + runtime.GOARCH,
		Shell:        getShellName(),
		CurrentDir:   getCurrentDir(),
		RecentOutput: recentOutput,
	}

	ctx.WorkingFiles = detectWorkingFiles()

	return ctx
}

func GetEnhancedContext(recentOutput []string, recentCommands []string) SystemContext {
	ctx := GetSystemContext(recentOutput)
	ctx.RecentCommands = recentCommands
	return ctx
}

func getShellName() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "unknown"
	}

	parts := strings.Split(shell, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return shell
}

func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "unknown"
	}
	return dir
}

func CheckCommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

func detectWorkingFiles() []string {
	entries, err := os.ReadDir(".")
	if err != nil {
		return []string{}
	}

	var files []string
	for i, entry := range entries {
		if i >= 10 { // Limit to 10 files for context
			break
		}
		if !strings.HasPrefix(entry.Name(), ".") {
			files = append(files, entry.Name())
		}
	}

	return files
}

func AnalyzeRecentOutput(output []string) map[string]string {
	hints := make(map[string]string)

	for _, line := range output {
		if strings.Contains(strings.ToLower(line), "error") {
			hints["has_errors"] = "true"
		}
		if strings.Contains(strings.ToLower(line), "permission denied") {
			hints["needs_sudo"] = "true"
		}
		if strings.Contains(line, "command not found") {
			hints["missing_command"] = "true"
		}

		if strings.Contains(line, ".json") || strings.Contains(line, ".yml") {
			hints["working_with"] = "config files"
		}
		if strings.Contains(line, ".py") {
			hints["working_with"] = "python"
		}
		if strings.Contains(line, ".go") {
			hints["working_with"] = "golang"
		}
		if strings.Contains(line, ".js") || strings.Contains(line, ".ts") {
			hints["working_with"] = "javascript"
		}
	}

	return hints
}
