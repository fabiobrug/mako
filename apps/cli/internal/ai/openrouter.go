package ai

import (
	"fmt"
	"net/http"
)

// OpenRouter uses OpenAI-compatible API format
// We can reuse the OpenAI provider with different base URL

// NewOpenRouterProvider creates a new OpenRouter provider
// OpenRouter provides access to multiple AI models through a single API
func NewOpenRouterProvider(cfg *ProviderConfig) (*OpenAIProvider, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("OpenRouter API key not found. Set LLM_API_KEY in .env")
	}
	
	model := cfg.Model
	if model == "" {
		model = "deepseek/deepseek-chat" // Default to cost-effective model
	}
	
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "https://openrouter.ai/api/v1"
	}
	
	return &OpenAIProvider{
		apiKey:  cfg.APIKey,
		model:   model,
		baseURL: baseURL,
		client:  &http.Client{},
	}, nil
}
