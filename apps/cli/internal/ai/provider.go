package ai

import (
	"fmt"
	"os"
	"strings"

	"github.com/fabiobrug/mako.git/internal/config"
)

// AIProvider defines the interface that all LLM providers must implement
type AIProvider interface {
	// GenerateCommand generates a shell command from a natural language request
	GenerateCommand(userRequest string, context SystemContext) (string, error)
	
	// GenerateCommandWithConversation generates a command with conversation history
	GenerateCommandWithConversation(userRequest string, context SystemContext, conversation *ConversationHistory) (string, error)
	
	// ExplainError provides an explanation and suggestion for a command error
	ExplainError(failedCommand string, errorOutput string, context SystemContext) (string, error)
	
	// ExplainCommand explains what a command does in human-readable terms
	ExplainCommand(command string, context SystemContext) (string, error)
	
	// SuggestAlternatives suggests alternative ways to accomplish the same goal
	SuggestAlternatives(command string, context SystemContext) (string, error)
}

// EmbeddingProvider defines the interface for embedding generation
type EmbeddingProvider interface {
	// GenerateEmbedding generates a vector embedding for the given text
	GenerateEmbedding(text string) ([]byte, error)
}

// ProviderConfig holds configuration for initializing a provider
type ProviderConfig struct {
	Provider string
	Model    string
	APIKey   string
	BaseURL  string
}

// LoadProviderConfig loads LLM provider configuration from environment or config file
func LoadProviderConfig() (*ProviderConfig, error) {
	// Try environment variables first
	provider := os.Getenv("LLM_PROVIDER")
	model := os.Getenv("LLM_MODEL")
	apiKey := os.Getenv("LLM_API_KEY")
	baseURL := os.Getenv("LLM_API_BASE")
	
	// Fallback to config file if env vars not set
	if provider == "" {
		cfg, err := config.LoadConfig()
		if err == nil {
			provider = cfg.LLMProvider
			model = cfg.LLMModel
			apiKey = cfg.APIKey
			baseURL = cfg.LLMBaseURL
		}
	}
	
	// Default to gemini for backward compatibility
	if provider == "" {
		provider = "gemini"
	}
	
	// Validate provider
	provider = strings.ToLower(strings.TrimSpace(provider))
	validProviders := map[string]bool{
		"openai":     true,
		"anthropic":  true,
		"openrouter": true,
		"gemini":     true,
		"deepseek":   true,
		"ollama":     true,
	}
	
	if !validProviders[provider] {
		return nil, fmt.Errorf("unsupported LLM provider: %s. Supported: openai, anthropic, openrouter, gemini, deepseek, ollama", provider)
	}
	
	return &ProviderConfig{
		Provider: provider,
		Model:    model,
		APIKey:   apiKey,
		BaseURL:  baseURL,
	}, nil
}

// LoadEmbeddingProviderConfig loads embedding provider configuration
func LoadEmbeddingProviderConfig() (*ProviderConfig, error) {
	// Try environment variables for embedding-specific config
	provider := os.Getenv("EMBEDDING_PROVIDER")
	model := os.Getenv("EMBEDDING_MODEL")
	apiKey := os.Getenv("EMBEDDING_API_KEY")
	baseURL := os.Getenv("EMBEDDING_API_BASE")
	
	// If no embedding-specific config, fall back to main LLM config
	if provider == "" {
		llmCfg, err := LoadProviderConfig()
		if err != nil {
			return nil, err
		}
		
		// Use the provider and API key from LLM config
		// but NOT the model (embedding models are different from text generation models)
		return &ProviderConfig{
			Provider: llmCfg.Provider,
			Model:    model, // Keep empty to use provider defaults
			APIKey:   llmCfg.APIKey,
			BaseURL:  baseURL,
		}, nil
	}
	
	provider = strings.ToLower(strings.TrimSpace(provider))
	
	return &ProviderConfig{
		Provider: provider,
		Model:    model,
		APIKey:   apiKey,
		BaseURL:  baseURL,
	}, nil
}

// NewAIProvider creates a new AI provider based on configuration
func NewAIProvider() (AIProvider, error) {
	cfg, err := LoadProviderConfig()
	if err != nil {
		return nil, err
	}
	
	switch cfg.Provider {
	case "gemini":
		return NewGeminiProvider(cfg)
	case "openai":
		return NewOpenAIProvider(cfg)
	case "anthropic":
		return NewAnthropicProvider(cfg)
	case "openrouter":
		return NewOpenRouterProvider(cfg)
	case "deepseek":
		return NewDeepSeekProvider(cfg)
	case "ollama":
		return NewOllamaProvider(cfg)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", cfg.Provider)
	}
}

// NewEmbeddingProvider creates a new embedding provider based on configuration
func NewEmbeddingProvider() (EmbeddingProvider, error) {
	cfg, err := LoadEmbeddingProviderConfig()
	if err != nil {
		return nil, err
	}
	
	switch cfg.Provider {
	case "gemini":
		return NewGeminiEmbeddingProvider(cfg)
	case "openai":
		return NewOpenAIEmbeddingProvider(cfg)
	case "ollama":
		return NewOllamaEmbeddingProvider(cfg)
	default:
		// For providers without embedding support, fall back to Gemini
		return NewGeminiEmbeddingProvider(&ProviderConfig{
			Provider: "gemini",
			APIKey:   os.Getenv("GEMINI_API_KEY"),
		})
	}
}
