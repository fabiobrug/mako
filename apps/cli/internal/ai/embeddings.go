package ai

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
)

type GeminiEmbeddingProvider struct {
	apiKey string
	model  string
	client *http.Client
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

	return &GeminiEmbeddingProvider{
		apiKey: apiKey,
		model:  model,
		client: &http.Client{},
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
