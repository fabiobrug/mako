package ai

import (
	"fmt"
	"net/http"
)

// DeepSeek uses OpenAI-compatible API format
// We can reuse the OpenAI provider with different base URL

// NewDeepSeekProvider creates a new DeepSeek provider
func NewDeepSeekProvider(cfg *ProviderConfig) (*OpenAIProvider, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("DeepSeek API key not found. Set LLM_API_KEY in .env")
	}
	
	model := cfg.Model
	if model == "" {
		model = "deepseek-chat" // Default model
	}
	
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "https://api.deepseek.com/v1"
	}
	
	return &OpenAIProvider{
		apiKey:  cfg.APIKey,
		model:   model,
		baseURL: baseURL,
		client:  &http.Client{},
	}, nil
}
