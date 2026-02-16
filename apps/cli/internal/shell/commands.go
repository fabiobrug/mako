package shell

import (
	"fmt"
	"strings"

	"github.com/fabiobrug/mako.git/internal/cache"
	"github.com/fabiobrug/mako.git/internal/database"
	"github.com/fabiobrug/mako.git/internal/safety"
)

var validator = safety.NewValidator()

// Global reference to ring buffer (will be set from main)
var recentOutputGetter func(int) []string

// Global reference to embedding cache (will be set from main)
var embeddingCache *cache.EmbeddingCache

// SetRecentOutputGetter allows main to provide ring buffer access
func SetRecentOutputGetter(getter func(int) []string) {
	recentOutputGetter = getter
}

// SetEmbeddingCache allows main to provide cache access
func SetEmbeddingCache(cache *cache.EmbeddingCache) {
	embeddingCache = cache
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
			// Support contextual help like "mako help quickstart" or "mako help --alias"
			if len(parts) > 2 {
				topic := strings.ToLower(parts[2])
				helpText := getContextualHelp(topic)
				if helpText != "" {
					return true, helpText, nil
				}
			}
			return true, getHelpText(), nil
		case "v", "version":
			return true, fmt.Sprintf("v1.3.4\n"), nil
		case "draw":
			return true, getSharkArt(), nil
		case "clear":
			output, err := handleClear()
			return true, output, err
		case "health":
			output, err := handleHealth(db)
			return true, output, err
		case "export":
			output, err := handleExport(parts[2:], db)
			return true, output, err
		case "import":
			output, err := handleImport(parts[2:], db)
			return true, output, err
		case "sync":
			output, err := handleSync(db)
			return true, output, err
		case "config":
			output, err := handleConfig(parts[2:])
			return true, output, err
		case "update":
			output, err := handleUpdate(parts[2:])
			return true, output, err
		case "completion":
			output, err := handleCompletion(parts[2:])
			return true, output, err
		case "uninstall":
			output, err := handleUninstall()
			return true, output, err
		default:
			return true, fmt.Sprintf("Unknown mako command: %s\n", parts[1]), nil
		}
	}

	return false, "", nil
}
