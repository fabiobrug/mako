package health

import (
	"fmt"
	"os"

	"github.com/fabiobrug/mako.git/internal/cache"
	"github.com/fabiobrug/mako.git/internal/config"
	"github.com/fabiobrug/mako.git/internal/database"
)

// HealthStatus represents the health status of a component
type HealthStatus string

const (
	StatusOK      HealthStatus = "OK"
	StatusWarning HealthStatus = "WARNING"
	StatusError   HealthStatus = "ERROR"
)

// ComponentHealth represents health of a single component
type ComponentHealth struct {
	Name    string
	Status  HealthStatus
	Message string
	Details map[string]interface{}
}

// HealthReport contains overall health information
type HealthReport struct {
	Components   []ComponentHealth
	Suggestions  []string
	OverallOK    bool
}

// Checker performs health checks
type Checker struct {
	db           *database.DB
	embeddingCache *cache.EmbeddingCache
	apiKey       string
}

// NewChecker creates a new health checker
func NewChecker(db *database.DB, embeddingCache *cache.EmbeddingCache, apiKey string) *Checker {
	return &Checker{
		db:             db,
		embeddingCache: embeddingCache,
		apiKey:         apiKey,
	}
}

// Check performs all health checks
func (c *Checker) Check() (*HealthReport, error) {
	report := &HealthReport{
		Components:  make([]ComponentHealth, 0),
		Suggestions: make([]string, 0),
		OverallOK:   true,
	}

	// Check database
	dbHealth := c.checkDatabase()
	report.Components = append(report.Components, dbHealth)
	if dbHealth.Status != StatusOK {
		report.OverallOK = false
	}

	// Check API key
	apiHealth := c.checkAPIKey()
	report.Components = append(report.Components, apiHealth)
	if apiHealth.Status != StatusOK {
		report.OverallOK = false
	}

	// Check cache
	cacheHealth := c.checkCache()
	report.Components = append(report.Components, cacheHealth)
	if cacheHealth.Status == StatusError {
		report.OverallOK = false
	}

	// Check disk space
	diskHealth := c.checkDiskSpace()
	report.Components = append(report.Components, diskHealth)
	if diskHealth.Status == StatusError {
		report.OverallOK = false
	}

	// Check embedding provider
	embeddingHealth := c.checkEmbeddingProvider()
	report.Components = append(report.Components, embeddingHealth)
	if embeddingHealth.Status == StatusError {
		report.OverallOK = false
	}

	// Generate suggestions
	report.Suggestions = c.generateSuggestions(report.Components)

	return report, nil
}

// checkDatabase verifies database health
func (c *Checker) checkDatabase() ComponentHealth {
	health := ComponentHealth{
		Name:    "Database",
		Details: make(map[string]interface{}),
	}

	// Get command count
	count, err := c.db.GetCommandCount()
	if err != nil {
		health.Status = StatusError
		health.Message = fmt.Sprintf("Failed to query database: %v", err)
		return health
	}
	health.Details["command_count"] = count

	// Get database size
	size, err := c.db.GetDatabaseSize()
	if err != nil {
		health.Status = StatusWarning
		health.Message = "Could not determine database size"
	} else {
		health.Details["size_bytes"] = size
		health.Details["size_mb"] = size / (1024 * 1024)
	}

	// Check for corruption (simple check)
	_, err = c.db.GetRecentCommands(1)
	if err != nil {
		health.Status = StatusError
		health.Message = "Database may be corrupted"
		return health
	}

	health.Status = StatusOK
	health.Message = fmt.Sprintf("%d commands, %d MB", count, size/(1024*1024))
	return health
}

// checkAPIKey verifies API key is set
func (c *Checker) checkAPIKey() ComponentHealth {
	health := ComponentHealth{
		Name:    "AI Provider",
		Details: make(map[string]interface{}),
	}

	// Get provider configuration from environment variables first
	provider := os.Getenv("LLM_PROVIDER")
	model := os.Getenv("LLM_MODEL")
	apiKey := os.Getenv("LLM_API_KEY")
	baseURL := os.Getenv("LLM_API_BASE")
	
	// If not in env, try loading from config file
	if provider == "" || apiKey == "" {
		if cfg, err := config.LoadConfig(); err == nil {
			if provider == "" && cfg.LLMProvider != "" {
				provider = cfg.LLMProvider
			}
			if model == "" && cfg.LLMModel != "" {
				model = cfg.LLMModel
			}
			if apiKey == "" && cfg.APIKey != "" {
				apiKey = cfg.APIKey
			}
			if baseURL == "" && cfg.LLMBaseURL != "" {
				baseURL = cfg.LLMBaseURL
			}
		}
	}
	
	// Legacy fallback
	if apiKey == "" {
		apiKey = os.Getenv("GEMINI_API_KEY")
	}
	
	// Default values
	if provider == "" {
		provider = "gemini"
	}
	if model == "" {
		model = "default"
	}
	
	health.Details["provider"] = provider
	health.Details["model"] = model

	// Ollama doesn't require API key
	if provider == "ollama" {
		if baseURL == "" {
			baseURL = "http://localhost:11434"
		}
		health.Status = StatusOK
		health.Message = fmt.Sprintf("Using %s (local, model: %s)", provider, model)
		health.Details["base_url"] = baseURL
		return health
	}

	// Check API key for cloud providers
	if apiKey == "" {
		health.Status = StatusError
		health.Message = fmt.Sprintf("API key not set for %s provider", provider)
		health.Details["help"] = "Set LLM_API_KEY in .env or use: cp .env.example .env"
		return health
	}

	// Basic validation (not a real API call to save quota)
	if len(apiKey) < 20 {
		health.Status = StatusWarning
		health.Message = fmt.Sprintf("Using %s - API key looks invalid (too short)", provider)
	} else {
		health.Status = StatusOK
		health.Message = fmt.Sprintf("Using %s (model: %s)", provider, model)
		health.Details["key_length"] = len(apiKey)
	}

	return health
}

// checkCache verifies embedding cache performance
func (c *Checker) checkCache() ComponentHealth {
	health := ComponentHealth{
		Name:    "Cache",
		Details: make(map[string]interface{}),
	}

	if c.embeddingCache == nil {
		health.Status = StatusWarning
		health.Message = "Cache not initialized"
		return health
	}

	stats := c.embeddingCache.Stats()
	health.Details["size"] = stats.Size
	health.Details["max_size"] = stats.MaxSize
	health.Details["hits"] = stats.Hits
	health.Details["misses"] = stats.Misses
	health.Details["hit_rate"] = fmt.Sprintf("%.1f%%", stats.HitRate*100)
	health.Details["memory_mb"] = stats.MemoryUsed / (1024 * 1024)

	// Evaluate cache performance
	if stats.HitRate < 0.4 && stats.Hits+stats.Misses > 100 {
		health.Status = StatusWarning
		health.Message = fmt.Sprintf("Low hit rate: %.1f%% (recommended: >60%%)", stats.HitRate*100)
	} else if stats.HitRate >= 0.6 || stats.Hits+stats.Misses < 100 {
		health.Status = StatusOK
		health.Message = fmt.Sprintf("Hit rate: %.1f%%", stats.HitRate*100)
	} else {
		health.Status = StatusOK
		health.Message = fmt.Sprintf("Hit rate: %.1f%%", stats.HitRate*100)
	}

	return health
}

// checkDiskSpace verifies available disk space
func (c *Checker) checkDiskSpace() ComponentHealth {
	health := ComponentHealth{
		Name:    "Disk Space",
		Details: make(map[string]interface{}),
	}

	// Get home directory
	home, err := os.UserHomeDir()
	if err != nil {
		health.Status = StatusWarning
		health.Message = "Could not check disk space"
		return health
	}

	makoDir := home + "/.mako"
	
	// Get directory size
	var totalSize int64
	err = os.MkdirAll(makoDir, 0755)
	if err == nil {
		// Walk directory and sum file sizes
		entries, err := os.ReadDir(makoDir)
		if err == nil {
			for _, entry := range entries {
				if !entry.IsDir() {
					info, err := entry.Info()
					if err == nil {
						totalSize += info.Size()
					}
				}
			}
		}
	}

	sizeMB := totalSize / (1024 * 1024)
	health.Details["size_mb"] = sizeMB

	// Simple threshold check (100MB limit as mentioned in requirements)
	const maxSizeMB = 100
	if sizeMB > maxSizeMB {
		health.Status = StatusWarning
		health.Message = fmt.Sprintf("%d MB / %d MB limit exceeded", sizeMB, maxSizeMB)
	} else {
		health.Status = StatusOK
		health.Message = fmt.Sprintf("%d MB / %d MB limit", sizeMB, maxSizeMB)
	}

	return health
}

// checkEmbeddingProvider verifies embedding provider configuration
func (c *Checker) checkEmbeddingProvider() ComponentHealth {
	health := ComponentHealth{
		Name:    "Embedding Provider",
		Details: make(map[string]interface{}),
	}

	// Get embedding provider configuration from environment variables
	embeddingProvider := os.Getenv("EMBEDDING_PROVIDER")
	embeddingModel := os.Getenv("EMBEDDING_MODEL")
	embeddingAPIKey := os.Getenv("EMBEDDING_API_KEY")
	embeddingBaseURL := os.Getenv("EMBEDDING_API_BASE")
	
	// Get LLM provider as fallback
	llmProvider := os.Getenv("LLM_PROVIDER")
	llmAPIKey := os.Getenv("LLM_API_KEY")
	
	// If not in env, try loading from config file
	if llmProvider == "" || llmAPIKey == "" {
		if cfg, err := config.LoadConfig(); err == nil {
			if llmProvider == "" && cfg.LLMProvider != "" {
				llmProvider = cfg.LLMProvider
			}
			if llmAPIKey == "" && cfg.APIKey != "" {
				llmAPIKey = cfg.APIKey
			}
		}
	}
	
	// Legacy fallback
	if llmAPIKey == "" {
		llmAPIKey = os.Getenv("GEMINI_API_KEY")
	}
	
	// Determine which provider is being used for embeddings
	if embeddingProvider == "" {
		embeddingProvider = llmProvider
		if embeddingProvider == "" {
			embeddingProvider = "gemini" // default
		}
	}
	
	health.Details["provider"] = embeddingProvider
	
	// Set default models based on provider
	defaultModel := embeddingModel
	if defaultModel == "" {
		switch embeddingProvider {
		case "gemini":
			defaultModel = "text-embedding-005"
		case "openai":
			defaultModel = "text-embedding-3-small"
		case "ollama":
			defaultModel = "nomic-embed-text"
		default:
			defaultModel = "default"
		}
	}
	health.Details["model"] = defaultModel
	
	// Check if provider requires API key
	if embeddingProvider == "ollama" {
		// Ollama doesn't require API key
		if embeddingBaseURL == "" {
			embeddingBaseURL = "http://localhost:11434"
		}
		health.Status = StatusOK
		health.Message = fmt.Sprintf("Using %s (local, model: %s)", embeddingProvider, defaultModel)
		health.Details["base_url"] = embeddingBaseURL
		health.Details["semantic_search"] = "enabled"
		return health
	}
	
	// Check API key for cloud providers
	actualAPIKey := embeddingAPIKey
	if actualAPIKey == "" {
		actualAPIKey = llmAPIKey
	}
	
	if actualAPIKey == "" {
		health.Status = StatusError
		health.Message = fmt.Sprintf("API key not set for %s embedding provider", embeddingProvider)
		health.Details["help"] = "Set EMBEDDING_API_KEY or LLM_API_KEY in .env"
		health.Details["semantic_search"] = "disabled"
		return health
	}
	
	// Basic validation
	if len(actualAPIKey) < 20 {
		health.Status = StatusWarning
		health.Message = fmt.Sprintf("Using %s - API key looks invalid (too short)", embeddingProvider)
		health.Details["semantic_search"] = "may not work"
	} else {
		health.Status = StatusOK
		health.Message = fmt.Sprintf("Using %s (model: %s)", embeddingProvider, defaultModel)
		health.Details["key_length"] = len(actualAPIKey)
		health.Details["semantic_search"] = "enabled"
	}
	
	return health
}

// generateSuggestions creates optimization suggestions
func (c *Checker) generateSuggestions(components []ComponentHealth) []string {
	suggestions := make([]string, 0)

	// Check cache performance
	for _, comp := range components {
		if comp.Name == "Cache" {
			if hitRate, ok := comp.Details["hit_rate"].(string); ok {
				// Parse hit rate
				var rate float64
				fmt.Sscanf(hitRate, "%f%%", &rate)
				
				if rate < 60 && comp.Details["hits"].(int64)+comp.Details["misses"].(int64) > 100 {
					suggestions = append(suggestions, "Cache hit rate is low - consider increasing cache size")
				} else if rate >= 80 {
					suggestions = append(suggestions, "Cache is working well")
				}
			}
		}

		if comp.Name == "Database" {
			if sizeMB, ok := comp.Details["size_mb"].(int64); ok {
				if sizeMB > 200 {
					suggestions = append(suggestions, "Consider archiving commands older than 1 year")
				}
			}

			if count, ok := comp.Details["command_count"].(int64); ok {
				if count > 100000 {
					suggestions = append(suggestions, "Large command history - consider periodic cleanup")
				}
			}
		}

		if comp.Name == "Disk Space" {
			if sizeMB, ok := comp.Details["size_mb"].(int64); ok {
				if sizeMB > 80 {
					suggestions = append(suggestions, "Approaching disk space limit - consider cleanup")
				}
			}
		}
	}

	if len(suggestions) == 0 {
		suggestions = append(suggestions, "Everything looks good!")
	}

	return suggestions
}

// FormatReport formats the health report for display
func FormatReport(report *HealthReport) string {
	var output string
	output += "▸ Mako Health Check\r\n\r\n"

	for _, comp := range report.Components {
		statusSymbol := "✓"
		if comp.Status == StatusWarning {
			statusSymbol = "⚠"
		} else if comp.Status == StatusError {
			statusSymbol = "✗"
		}

		output += fmt.Sprintf("%s %s: %s\r\n", statusSymbol, comp.Name, comp.Message)
	}

	output += "\r\nPerformance Tips:\r\n"
	for _, suggestion := range report.Suggestions {
		output += fmt.Sprintf("- %s\r\n", suggestion)
	}

	return output
}
