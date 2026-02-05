package ai

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func GetSystemContext(recentOutput []string) SystemContext {
	return SystemContext{
		OS:           runtime.GOOS + "/" + runtime.GOARCH,
		Shell:        getShellName(),
		CurrentDir:   getCurrentDir(),
		RecentOutput: recentOutput,
	}
}

func getShellName() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "unknow"
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
		return "unknow"
	}
	return dir
}

func CheckCommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}
