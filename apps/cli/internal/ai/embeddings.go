package ai

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fabiobrug/mako.git/internal/retry"
)

type GeminiEmbeddingProvider struct {
	apiKey   string
	model    string
	client   *http.Client
	executor *retry.ResilientExecutor
}

// NewGeminiEmbeddingProvider creates a new Gemini embedding provider
func NewGeminiEmbeddingProvider(cfg *ProviderConfig) (*GeminiEmbeddingProvider, error) {
	apiKey := cfg.APIKey
	
	// Fallback to environment variable
	if apiKey == "" {
		apiKey = os.Getenv("GEMINI_API_KEY")
	}
	
	if apiKey == "" {
		return nil, fmt.Errorf("Gemini API key not found. Set GEMINI_API_KEY or EMBEDDING_API_KEY")
	}
	
	// Default model if not specified
	model := cfg.Model
	if model == "" {
		model = "gemini-embedding-001"
	}

	// Configure retry with exponential backoff for embeddings
	retryConfig := &retry.Config{
		MaxAttempts:  3,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     5 * time.Second,
		Multiplier:   2.0,
		Jitter:       true,
		RetryableErrors: func(err error) bool {
			// Don't retry on authentication errors or invalid requests
			if strings.Contains(err.Error(), "status 401") ||
				strings.Contains(err.Error(), "status 403") ||
				strings.Contains(err.Error(), "status 400") {
				return false
			}
			// Retry on timeouts, rate limits, and server errors
			return strings.Contains(err.Error(), "status 429") ||
				strings.Contains(err.Error(), "status 500") ||
				strings.Contains(err.Error(), "status 502") ||
				strings.Contains(err.Error(), "status 503") ||
				strings.Contains(err.Error(), "status 504") ||
				strings.Contains(err.Error(), "timeout") ||
				strings.Contains(err.Error(), "connection")
		},
	}

	// Configure circuit breaker for embeddings
	cbConfig := &retry.CircuitBreakerConfig{
		MaxFailures: 5,
		Timeout:     30 * time.Second,
		MaxRequests: 1,
		OnStateChange: func(from, to retry.State) {
			// Silent for embeddings to avoid cluttering output
		},
	}

	circuitBreaker := retry.NewCircuitBreaker(cbConfig)
	executor := retry.NewResilientExecutor(retryConfig, circuitBreaker)

	return &GeminiEmbeddingProvider{
		apiKey:   apiKey,
		model:    model,
		client:   &http.Client{Timeout: 30 * time.Second},
		executor: executor,
	}, nil
}

// EmbeddingService is deprecated, kept for backward compatibility
type EmbeddingService = GeminiEmbeddingProvider

// NewEmbeddingService creates a legacy embedding service
// Deprecated: Use NewEmbeddingProvider instead
func NewEmbeddingService() (*GeminiEmbeddingProvider, error) {
	return NewGeminiEmbeddingProvider(&ProviderConfig{
		Provider: "gemini",
	})
}

func (e *GeminiEmbeddingProvider) Embed(text string) ([]float32, error) {
	requestBody := map[string]interface{}{
		"content": map[string]interface{}{
			"parts": []map[string]interface{}{
				{"text": text},
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	// Use resilient executor for retry + circuit breaker
	ctx := context.Background()
	embedding, err := e.executeWithRetry(ctx, func() ([]float32, error) {
		return e.makeEmbeddingRequest(jsonData)
	})

	return embedding, err
}

// executeWithRetry is a helper to execute API calls with retry and circuit breaker
func (e *GeminiEmbeddingProvider) executeWithRetry(ctx context.Context, fn func() ([]float32, error)) ([]float32, error) {
	return retry.DoWithResult(ctx, e.executor.Retry, func() ([]float32, error) {
		var result []float32
		err := e.executor.CircuitBreaker.Execute(ctx, func() error {
			var innerErr error
			result, innerErr = fn()
			return innerErr
		})
		return result, err
	})
}

// makeEmbeddingRequest performs the actual HTTP request for embeddings
func (e *GeminiEmbeddingProvider) makeEmbeddingRequest(jsonData []byte) ([]float32, error) {
	embedAPIURL := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:embedContent", e.model)
	url := fmt.Sprintf("%s?key=%s", embedAPIURL, e.apiKey)
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var response struct {
		Embedding struct {
			Values []float32 `json:"values"`
		} `json:"embedding"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Embedding.Values, nil
}

// GenerateEmbedding generates an embedding and returns it as bytes
// This method satisfies the EmbeddingProvider interface
func (e *GeminiEmbeddingProvider) GenerateEmbedding(text string) ([]byte, error) {
	vec, err := e.Embed(text)
	if err != nil {
		return nil, err
	}
	return VectorToBytes(vec), nil
}

func VectorToBytes(vec []float32) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, vec)
	return buf.Bytes()
}

func BytesToVector(data []byte) ([]float32, error) {
	if len(data)%4 != 0 {
		return nil, fmt.Errorf("invalid vector data length")
	}

	vec := make([]float32, len(data)/4)
	buf := bytes.NewReader(data)

	err := binary.Read(buf, binary.LittleEndian, &vec)
	if err != nil {
		return nil, err
	}

	return vec, nil
}

func CosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct, normA, normB float64

	for i := range a {
		dotProduct += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return float32(dotProduct / (math.Sqrt(normA) * math.Sqrt(normB)))
}
