package safety

import (
	"fmt"
	"regexp"
	"strings"
)

// CommandRisk represents the danger level of a command
type CommandRisk int

const (
	RiskSafe CommandRisk = iota
	RiskLow
	RiskMedium
	RiskHigh
	RiskCritical
)

// ValidationResult contains risk assessment for a command
type ValidationResult struct {
	Risk    CommandRisk
	Reasons []string
	Safe    bool
}

// Validator checks commands for safety
type Validator struct {
	criticalPatterns   []Pattern
	highRiskPatterns   []Pattern
	mediumRiskPatterns []Pattern
	secretPatterns     []*regexp.Regexp
}

type Pattern struct {
	Regex   *regexp.Regexp
	Message string
}

// NewValidator creates a new command validator
func NewValidator() *Validator {
	return &Validator{
		criticalPatterns: []Pattern{
			{regexp.MustCompile(`rm\s+(-rf?|--recursive)\s+/\s*$`), "Recursive delete of root directory"},
			{regexp.MustCompile(`rm\s+(-rf?|--recursive)\s+/\w+\s*$`), "Recursive delete of top-level directory"},
			{regexp.MustCompile(`dd\s+if=/dev/zero\s+of=/dev/[sh]d`), "Disk wipe command"},
			{regexp.MustCompile(`mkfs\.\w+\s+/dev/`), "Filesystem formatting"},
			{regexp.MustCompile(`:\(\)\{\s*:\|:&\s*\};:`), "Fork bomb"},
			{regexp.MustCompile(`curl.*\|\s*bash`), "Piping untrusted script to shell"},
			{regexp.MustCompile(`wget.*\|\s*sh`), "Piping untrusted script to shell"},
		},
		highRiskPatterns: []Pattern{
			{regexp.MustCompile(`rm\s+(-rf?|--recursive).*\*`), "Recursive delete with wildcard"},
			{regexp.MustCompile(`sudo\s+rm\s+-rf`), "Sudo recursive delete"},
			{regexp.MustCompile(`chmod\s+-R\s+777`), "Recursive permission change to 777"},
			{regexp.MustCompile(`chown\s+-R.*\s+/`), "Recursive ownership change"},
			{regexp.MustCompile(`>\s*/dev/sd[a-z]`), "Writing directly to disk device"},
			{regexp.MustCompile(`/dev/null\s+&$`), "Background job with output to /dev/null"},
		},
		mediumRiskPatterns: []Pattern{
			{regexp.MustCompile(`sudo\s+`), "Sudo command (elevated privileges)"},
			{regexp.MustCompile(`rm\s+-r`), "Recursive delete"},
			{regexp.MustCompile(`docker\s+system\s+prune\s+-a`), "Docker cleanup removes all unused images"},
			{regexp.MustCompile(`git\s+push\s+--force`), "Force push (rewrites history)"},
			{regexp.MustCompile(`npm\s+install\s+-g`), "Global npm install"},
		},
		secretPatterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)(password|passwd|pwd)\s*=\s*['"]?[^\s'"]+`),
			regexp.MustCompile(`(?i)(api[_-]?key|apikey|access[_-]?key)\s*=\s*['"]?[^\s'"]+`),
			regexp.MustCompile(`(?i)(token|auth[_-]?token)\s*=\s*['"]?[^\s'"]+`),
			regexp.MustCompile(`(?i)(secret|secret[_-]?key)\s*=\s*['"]?[^\s'"]+`),
			regexp.MustCompile(`Bearer\s+[A-Za-z0-9\-._~+/]+=*`),
			regexp.MustCompile(`ghp_[A-Za-z0-9]{36}`),                      // GitHub PAT
			regexp.MustCompile(`xox[baprs]-[0-9]{10,13}-[A-Za-z0-9]{24,}`), // Slack token
			regexp.MustCompile(`sk-[A-Za-z0-9]{48}`),                       // OpenAI API key
			regexp.MustCompile(`AIza[0-9A-Za-z\-_]{35}`),                   // Google API key
			regexp.MustCompile(`(?i)-p\s+[^\s]+`),                          // -p password flags
			regexp.MustCompile(`://[^:@/]+:([^@]+)@`),                      // URLs with credentials
		},
	}
}

// ValidateCommand checks if a command is safe to execute
func (v *Validator) ValidateCommand(command string) ValidationResult {
	result := ValidationResult{
		Risk:    RiskSafe,
		Reasons: []string{},
		Safe:    true,
	}

	// Check critical patterns
	for _, pattern := range v.criticalPatterns {
		if pattern.Regex.MatchString(command) {
			result.Risk = RiskCritical
			result.Reasons = append(result.Reasons, pattern.Message)
			result.Safe = false
		}
	}

	// Check high risk patterns (only if not already critical)
	if result.Risk < RiskCritical {
		for _, pattern := range v.highRiskPatterns {
			if pattern.Regex.MatchString(command) {
				if result.Risk < RiskHigh {
					result.Risk = RiskHigh
				}
				result.Reasons = append(result.Reasons, pattern.Message)
			}
		}
	}

	// Check medium risk patterns (only if not higher risk)
	if result.Risk < RiskHigh {
		for _, pattern := range v.mediumRiskPatterns {
			if pattern.Regex.MatchString(command) {
				if result.Risk < RiskMedium {
					result.Risk = RiskMedium
				}
				result.Reasons = append(result.Reasons, pattern.Message)
			}
		}
	}

	// Medium, high, and critical risk commands are not safe
	// They should all show warnings to the user
	if result.Risk >= RiskMedium {
		result.Safe = false
	}

	return result
}

// RedactSecrets replaces sensitive information with asterisks
func (v *Validator) RedactSecrets(command string) string {
	redacted := command

	for _, pattern := range v.secretPatterns {
		redacted = pattern.ReplaceAllStringFunc(redacted, func(match string) string {
			// Handle URL credentials: ://user:password@host -> ://user:***@host
			if strings.Contains(match, "://") && strings.Contains(match, "@") {
				return "://user:***@"
			}

			// Handle -p password flag: -p mypassword -> -p ***
			if strings.HasPrefix(match, "-p ") {
				return "-p ***"
			}

			// Keep the key name, redact the value for key=value pairs
			parts := strings.SplitN(match, "=", 2)
			if len(parts) == 2 {
				return parts[0] + "=***"
			}

			// For Bearer tokens
			if strings.HasPrefix(match, "Bearer ") {
				return "Bearer ***"
			}

			// For standalone tokens (GitHub/Slack/OpenAI) - full redaction
			return "***"
		})
	}

	return redacted
}

// GetRiskColor returns ANSI color code for risk level
func (v *Validator) GetRiskColor(risk CommandRisk) string {
	switch risk {
	case RiskCritical:
		return "\033[38;2;255;50;50m" // Bright red
	case RiskHigh:
		return "\033[38;2;255;150;0m" // Orange
	case RiskMedium:
		return "\033[38;2;255;200;0m" // Yellow
	case RiskLow:
		return "\033[38;2;150;150;150m" // Gray
	default:
		return "\033[38;2;100;255;100m" // Green
	}
}

// GetRiskLabel returns human-readable risk level
func (v *Validator) GetRiskLabel(risk CommandRisk) string {
	switch risk {
	case RiskCritical:
		return " CRITICAL DANGER"
	case RiskHigh:
		return " HIGH RISK"
	case RiskMedium:
		return " MEDIUM RISK"
	case RiskLow:
		return " LOW RISK"
	default:
		return "✓ SAFE"
	}
}

// FormatWarning creates a formatted warning message
func (v *Validator) FormatWarning(result ValidationResult) string {
	if result.Safe {
		return ""
	}

	cyan := "\033[38;2;0;209;255m"
	lightBlue := "\033[38;2;93;173;226m"
	riskColor := v.GetRiskColor(result.Risk)
	reset := "\033[0m"

	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("\r\n%s╭─ Safety Warning%s\r\n", lightBlue, reset))
	msg.WriteString(fmt.Sprintf("%s│%s\r\n", lightBlue, reset))
	msg.WriteString(fmt.Sprintf("%s│%s  %s%s%s\r\n", lightBlue, reset, riskColor, v.GetRiskLabel(result.Risk), reset))
	msg.WriteString(fmt.Sprintf("%s│%s\r\n", lightBlue, reset))

	for _, reason := range result.Reasons {
		msg.WriteString(fmt.Sprintf("%s│%s  %s• %s%s\r\n", lightBlue, reset, riskColor, reason, reset))
	}

	msg.WriteString(fmt.Sprintf("%s│%s\r\n", lightBlue, reset))

	if result.Risk == RiskCritical {
		msg.WriteString(fmt.Sprintf("%s│%s  %sThis command is BLOCKED for your safety.%s\r\n", lightBlue, reset, riskColor, reset))
	} else if result.Risk == RiskHigh {
		msg.WriteString(fmt.Sprintf("%s│%s  %sPlease review carefully before confirming.%s\r\n", lightBlue, reset, riskColor, reset))
	} else {
		msg.WriteString(fmt.Sprintf("%s│%s  %sPlease confirm before running.%s\r\n", lightBlue, reset, cyan, reset))
	}

	msg.WriteString(fmt.Sprintf("%s╰─%s\r\n", lightBlue, reset))

	return msg.String()
}
