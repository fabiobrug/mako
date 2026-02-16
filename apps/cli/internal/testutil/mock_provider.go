package testutil

import (
	"fmt"
	"strings"
)

// MockAIProvider is a test implementation of AI provider
type MockAIProvider struct {
	GenerateFunc  func(prompt, context string) (string, error)
	GenerateError error
	GenerateResult string
	CallCount     int
}

func NewMockAIProvider(result string, err error) *MockAIProvider {
	return &MockAIProvider{
		GenerateResult: result,
		GenerateError:  err,
	}
}

func (m *MockAIProvider) GenerateCommand(prompt, context string) (string, error) {
	m.CallCount++
	if m.GenerateFunc != nil {
		return m.GenerateFunc(prompt, context)
	}
	if m.GenerateError != nil {
		return "", m.GenerateError
	}
	return m.GenerateResult, nil
}

func (m *MockAIProvider) Name() string {
	return "mock"
}

// MockEmbeddingProvider is a test implementation of embedding provider
type MockEmbeddingProvider struct {
	EmbedFunc    func(text string) ([]float32, error)
	EmbedError   error
	EmbedResult  []float32
	CallCount    int
}

func NewMockEmbeddingProvider(result []float32, err error) *MockEmbeddingProvider {
	return &MockEmbeddingProvider{
		EmbedResult: result,
		EmbedError:  err,
	}
}

func (m *MockEmbeddingProvider) GenerateEmbedding(text string) ([]float32, error) {
	m.CallCount++
	if m.EmbedFunc != nil {
		return m.EmbedFunc(text)
	}
	if m.EmbedError != nil {
		return nil, m.EmbedError
	}
	// Return a deterministic embedding based on text length
	if m.EmbedResult != nil {
		return m.EmbedResult, nil
	}
	// Generate simple embedding: first 4 floats are based on text properties
	embedding := make([]float32, 768)
	embedding[0] = float32(len(text))
	embedding[1] = float32(strings.Count(text, " "))
	embedding[2] = float32(strings.Count(text, "-"))
	embedding[3] = float32(len(strings.Fields(text)))
	return embedding, nil
}

func (m *MockEmbeddingProvider) Name() string {
	return "mock-embedding"
}

// PredictableAIProvider returns specific commands based on prompts
type PredictableAIProvider struct {
	Responses map[string]string
	CallCount int
}

func NewPredictableAIProvider() *PredictableAIProvider {
	return &PredictableAIProvider{
		Responses: map[string]string{
			"list files":              "ls -lh",
			"list all files":          "ls -lha",
			"show current directory":  "pwd",
			"disk usage":              "df -h",
			"find large files":        "find . -type f -size +100M",
			"git status":              "git status",
			"create directory test":   "mkdir test",
			"remove file test.txt":    "rm test.txt",
		},
	}
}

func (p *PredictableAIProvider) GenerateCommand(prompt, context string) (string, error) {
	p.CallCount++
	prompt = strings.ToLower(strings.TrimSpace(prompt))
	
	if cmd, ok := p.Responses[prompt]; ok {
		return cmd, nil
	}
	
	// Default response
	return fmt.Sprintf("echo 'Generated command for: %s'", prompt), nil
}

func (p *PredictableAIProvider) Name() string {
	return "predictable"
}
