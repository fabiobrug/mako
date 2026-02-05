package parser

import (
	"strings"
)

type CommandInfo struct {
	RawLine    string
	Command    string
	Args       []string
	IsPrompt   bool
	IsComplete bool
}

func ParseLine(line string) CommandInfo {
	trimmed := strings.TrimSpace(line)

	info := CommandInfo{
		RawLine: line,
	}

	if isPromptLine(trimmed) {
		info.IsPrompt = true
		return info
	}

	if trimmed == "" {
		return info
	}

	parts := strings.Fields(trimmed)
	if len(parts) > 0 {
		info.Command = parts[0]
		if len(parts) > 1 {
			info.Args = parts[1:]
		}
		info.IsComplete = true
	}

	return info
}

func isPromptLine(line string) bool {
	prompts := []string{
		"$",
		"#",
		">",
		"❯",  
	}

	for _, prompt := range prompts {
		if strings.HasSuffix(line, prompt) || strings.Contains(line, prompt+" ") {
			return true
		}
	}

	if strings.Contains(line, "@") && (strings.Contains(line, "$") || strings.Contains(line, "#")) {
		return true
	}

	return false
}

func IsIgnoredCommand(cmd string) bool {
	ignored := []string{
		"exit",
		"clear",
		"",
	}

	for _, ig := range ignored {
		if cmd == ig {
			return true
		}
	}

	return false
}
