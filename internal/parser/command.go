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

// ValidatePipeline checks if a command pipeline is syntactically valid
func ValidatePipeline(command string) bool {
	// Basic validation rules:
	// 1. Should not start or end with pipe operators
	// 2. Should not have consecutive operators without commands
	// 3. Should have balanced quotes

	trimmed := strings.TrimSpace(command)
	if trimmed == "" {
		return false
	}

	// Check for operators at start/end
	invalidStarts := []string{"|", "&&", "||", ";"}
	for _, op := range invalidStarts {
		if strings.HasPrefix(trimmed, op) || strings.HasSuffix(trimmed, op) {
			return false
		}
	}

	// Check for consecutive pipes/operators
	invalidPatterns := []string{"||", "| |", "| &&", "&& |", "; |"}
	for _, pattern := range invalidPatterns {
		if strings.Contains(trimmed, pattern) {
			return false
		}
	}

	// Check balanced quotes
	if !hasBalancedQuotes(trimmed) {
		return false
	}

	return true
}

// hasBalancedQuotes checks if quotes are balanced
func hasBalancedQuotes(s string) bool {
	singleQuotes := 0
	doubleQuotes := 0
	escaped := false

	for _, char := range s {
		if escaped {
			escaped = false
			continue
		}

		switch char {
		case '\\':
			escaped = true
		case '\'':
			singleQuotes++
		case '"':
			doubleQuotes++
		}
	}

	return singleQuotes%2 == 0 && doubleQuotes%2 == 0
}

// IsPipeline checks if a command contains pipeline operators
func IsPipeline(command string) bool {
	pipelineOperators := []string{"|", "&&", "||", ";"}
	for _, op := range pipelineOperators {
		if strings.Contains(command, op) {
			return true
		}
	}
	return false
}

// GetPipelineComplexity returns a score for how complex the pipeline is
func GetPipelineComplexity(command string) int {
	complexity := 0
	
	// Count pipe operators
	complexity += strings.Count(command, "|")
	
	// Count boolean operators
	complexity += strings.Count(command, "&&")
	complexity += strings.Count(command, "||")
	
	// Count semicolons
	complexity += strings.Count(command, ";")
	
	// Bonus for nested commands (subshells)
	complexity += strings.Count(command, "$(")
	complexity += strings.Count(command, "`")
	
	return complexity
}
