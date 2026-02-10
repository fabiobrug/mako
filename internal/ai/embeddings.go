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

const embedAPIURL = "https://generativelanguage.googleapis.com/v1beta/models/text-embedding-004:embedContent"

type EmbeddingService struct {
	apiKey string
	client *http.Client
}

func NewEmbeddingService() (*EmbeddingService, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY not set")
	}

	return &EmbeddingService{
		apiKey: apiKey,
		client: &http.Client{},
	}, nil
}

func (e *EmbeddingService) Embed(text string) ([]float32, error) {
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
// This method satisfies the database.EmbeddingService interface
func (e *EmbeddingService) GenerateEmbedding(text string) ([]byte, error) {
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
