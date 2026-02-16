package shell

import (
	"fmt"

	"github.com/fabiobrug/mako.git/internal/ai"
	"github.com/fabiobrug/mako.git/internal/database"
	"github.com/fabiobrug/mako.git/internal/health"
)

func handleHealth(db *database.DB) (string, error) {
	if db == nil {
		return "\r\n✗ Database not available - cannot perform health check\r\n\r\n", nil
	}
	
	// Create health checker (it will get provider config from environment)
	checker := health.NewChecker(db, embeddingCache, "")
	
	// Run health check
	report, err := checker.Check()
	if err != nil {
		return "", fmt.Errorf("health check failed: %w", err)
	}
	
	// Format and return report
	output := health.FormatReport(report)
	return output, nil
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
	var output string
	output += fmt.Sprintf("\n%s╭─ Mako Statistics%s\n", lightBlue, reset)
	output += fmt.Sprintf("%s│%s  Total commands    %s%d%s\n", lightBlue, reset, cyan, stats["total_commands"], reset)
	output += fmt.Sprintf("%s│%s  Commands today    %s%d%s\n", lightBlue, reset, cyan, stats["commands_today"], reset)
	output += fmt.Sprintf("%s│%s  Avg duration      %s%.0fms%s\n", lightBlue, reset, cyan, stats["avg_duration_ms"], reset)
	output += fmt.Sprintf("%s╰─%s\n\n", lightBlue, reset)
	return output, nil
}

func handleClear() (string, error) {
	green := "\033[38;2;100;255;100m"
	reset := "\033[0m"
	
	if err := ai.ClearConversation(); err != nil {
		return "", fmt.Errorf("failed to clear conversation: %w", err)
	}
	
	return fmt.Sprintf("\n%s✓ Conversation history cleared%s\n\n", green, reset), nil
}
