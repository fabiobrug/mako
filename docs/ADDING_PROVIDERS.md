# Adding New AI Providers to Mako

This guide explains how to add support for new AI providers to Mako.

## Overview

Mako uses a provider interface pattern that makes it easy to add new AI providers. Each provider must implement the `AIProvider` interface defined in `internal/ai/provider.go`.

## Provider Interface

Every AI provider must implement these methods:

```go
type AIProvider interface {
    GenerateCommand(userRequest string, context SystemContext) (string, error)
    GenerateCommandWithConversation(userRequest string, context SystemContext, conversation *ConversationHistory) (string, error)
    ExplainError(failedCommand string, errorOutput string, context SystemContext) (string, error)
    ExplainCommand(command string, context SystemContext) (string, error)
    SuggestAlternatives(command string, context SystemContext) (string, error)
}
```

## Step-by-Step Guide

### 1. Create a New Provider File

Create a new file in `apps/cli/internal/ai/` named after your provider (e.g., `myprovider.go`).

```go
package ai

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strings"
)

type MyProvider struct {
    apiKey  string
    model   string
    baseURL string
    client  *http.Client
}
```

### 2. Implement the Constructor

Create a constructor function that accepts a `ProviderConfig`:

```go
func NewMyProvider(cfg *ProviderConfig) (*MyProvider, error) {
    if cfg.APIKey == "" {
        return nil, fmt.Errorf("MyProvider API key not found. Set LLM_API_KEY in .env")
    }
    
    model := cfg.Model
    if model == "" {
        model = "default-model-name" // Your provider's default model
    }
    
    baseURL := cfg.BaseURL
    if baseURL == "" {
        baseURL = "https://api.myprovider.com/v1" // Your provider's API URL
    }
    
    return &MyProvider{
        apiKey:  cfg.APIKey,
        model:   model,
        baseURL: baseURL,
        client:  &http.Client{},
    }, nil
}
```

### 3. Implement Required Methods

#### GenerateCommand

This is the core method for generating shell commands:

```go
func (m *MyProvider) GenerateCommand(userRequest string, context SystemContext) (string, error) {
    return m.GenerateCommandWithConversation(userRequest, context, nil)
}
```

#### GenerateCommandWithConversation

This method includes conversation history:

```go
func (m *MyProvider) GenerateCommandWithConversation(userRequest string, context SystemContext, conversation *ConversationHistory) (string, error) {
    // Build the prompt with context
    prompt := m.buildPrompt(userRequest, context, conversation)
    
    // Make API request to your provider
    // ... (see examples below)
    
    // Clean and return the command
    command := m.cleanCommand(response.Text)
    return command, nil
}
```

#### ExplainError, ExplainCommand, SuggestAlternatives

These methods provide additional functionality:

```go
func (m *MyProvider) ExplainError(failedCommand string, errorOutput string, context SystemContext) (string, error) {
    prompt := fmt.Sprintf(`Shell debugging assistant. Analyze this error briefly.
Command: %s
Error: %s

Provide:
EXPLANATION: Brief explanation of the error
SUGGESTION: A corrected command or next steps`, failedCommand, errorOutput)
    
    return m.sendRequest(prompt, 2048, 0.3)
}

func (m *MyProvider) ExplainCommand(command string, context SystemContext) (string, error) {
    prompt := fmt.Sprintf(`Explain this shell command: %s`, command)
    return m.sendRequest(prompt, 1024, 0.3)
}

func (m *MyProvider) SuggestAlternatives(command string, context SystemContext) (string, error) {
    prompt := fmt.Sprintf(`Suggest 2-3 alternatives for: %s`, command)
    return m.sendRequest(prompt, 1024, 0.5)
}
```

### 4. Helper Methods

Implement helper methods for your provider:

```go
func (m *MyProvider) buildPrompt(userRequest string, context SystemContext, conversation *ConversationHistory) string {
    var promptBuild strings.Builder
    
    // Include conversation history if available
    if conversation != nil && conversation.IsActive() {
        promptBuild.WriteString(conversation.GetContext())
    }
    
    // Include system context
    promptBuild.WriteString(fmt.Sprintf("System: %s\n", context.OS))
    promptBuild.WriteString(fmt.Sprintf("Shell: %s\n", context.Shell))
    promptBuild.WriteString(fmt.Sprintf("Directory: %s\n", context.CurrentDir))
    
    // Include user request
    promptBuild.WriteString("\nUser request: ")
    promptBuild.WriteString(userRequest)
    
    return promptBuild.String()
}

func (m *MyProvider) cleanCommand(command string) string {
    // Remove common markdown artifacts
    command = strings.TrimPrefix(command, "```bash")
    command = strings.TrimPrefix(command, "```sh")
    command = strings.TrimPrefix(command, "```")
    command = strings.TrimSuffix(command, "```")
    command = strings.TrimSpace(command)
    return command
}

func (m *MyProvider) sendRequest(prompt string, maxTokens int, temperature float64) (string, error) {
    // Implement your provider's API call here
    // See existing providers for examples
}
```

### 5. Register the Provider

Add your provider to the factory function in `internal/ai/provider.go`:

```go
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
    case "myprovider":  // Add your provider here
        return NewMyProvider(cfg)
    // ... other providers
    default:
        return nil, fmt.Errorf("unsupported provider: %s", cfg.Provider)
    }
}
```

### 6. Update Configuration

Add your provider to the list of valid providers in `LoadProviderConfig()`:

```go
validProviders := map[string]bool{
    "openai":     true,
    "anthropic":  true,
    "gemini":     true,
    "myprovider": true,  // Add here
    // ... other providers
}
```

### 7. Update Documentation

Add your provider to:
- `.env.example` - Add configuration template
- `docs/SETUP.md` - Add setup instructions
- `README.md` - Add to supported providers table

## API Format Examples

### REST API with JSON

Most providers use REST APIs with JSON. Here's a template:

```go
func (m *MyProvider) sendRequest(prompt string, maxTokens int, temperature float64) (string, error) {
    requestBody := map[string]interface{}{
        "model":       m.model,
        "prompt":      prompt,
        "max_tokens":  maxTokens,
        "temperature": temperature,
    }
    
    jsonData, err := json.Marshal(requestBody)
    if err != nil {
        return "", err
    }
    
    url := fmt.Sprintf("%s/completions", m.baseURL)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return "", err
    }
    
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", m.apiKey))
    
    resp, err := m.client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    
    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
    }
    
    var response struct {
        Text string `json:"text"`
    }
    
    if err := json.Unmarshal(body, &response); err != nil {
        return "", err
    }
    
    return response.Text, nil
}
```

### OpenAI-Compatible APIs

If your provider is OpenAI-compatible, you can reuse the OpenAI provider:

```go
func NewMyProvider(cfg *ProviderConfig) (*OpenAIProvider, error) {
    if cfg.APIKey == "" {
        return nil, fmt.Errorf("API key not found")
    }
    
    model := cfg.Model
    if model == "" {
        model = "my-default-model"
    }
    
    baseURL := cfg.BaseURL
    if baseURL == "" {
        baseURL = "https://api.myprovider.com/v1"
    }
    
    // Reuse OpenAI provider implementation
    return &OpenAIProvider{
        apiKey:  cfg.APIKey,
        model:   model,
        baseURL: baseURL,
        client:  &http.Client{},
    }, nil
}
```

## Embedding Support (Optional)

If your provider supports embeddings, implement the `EmbeddingProvider` interface:

```go
type EmbeddingProvider interface {
    GenerateEmbedding(text string) ([]byte, error)
}
```

Example:

```go
type MyEmbeddingProvider struct {
    apiKey  string
    model   string
    baseURL string
    client  *http.Client
}

func NewMyEmbeddingProvider(cfg *ProviderConfig) (*MyEmbeddingProvider, error) {
    // Similar to NewMyProvider
}

func (m *MyEmbeddingProvider) GenerateEmbedding(text string) ([]byte, error) {
    // Call your provider's embedding API
    // Return the embedding as bytes using VectorToBytes()
    
    vec := []float32{...} // Your embedding vector
    return VectorToBytes(vec), nil
}
```

Register in `NewEmbeddingProvider()`:

```go
func NewEmbeddingProvider() (EmbeddingProvider, error) {
    cfg, err := LoadEmbeddingProviderConfig()
    if err != nil {
        return nil, err
    }
    
    switch cfg.Provider {
    case "myprovider":
        return NewMyEmbeddingProvider(cfg)
    // ... other providers
    }
}
```

## Testing Your Provider

### 1. Unit Tests

Create a test file `myprovider_test.go`:

```go
package ai

import (
    "testing"
)

func TestMyProvider_GenerateCommand(t *testing.T) {
    // Mock or test with real API
    provider := &MyProvider{
        apiKey:  "test-key",
        model:   "test-model",
        baseURL: "https://api.myprovider.com",
        client:  &http.Client{},
    }
    
    context := SystemContext{
        OS:         "linux",
        Shell:      "bash",
        CurrentDir: "/home/user",
    }
    
    cmd, err := provider.GenerateCommand("list files", context)
    if err != nil {
        t.Fatalf("GenerateCommand failed: %v", err)
    }
    
    if cmd == "" {
        t.Fatal("Expected non-empty command")
    }
}
```

### 2. Integration Tests

Test with real API:

```bash
# Set your API key
export LLM_PROVIDER=myprovider
export LLM_API_KEY=your-test-key
export LLM_MODEL=test-model

# Build
cd apps/cli
go build -o mako ./cmd/mako

# Test
./mako
# Inside: mako ask "list files"
```

## Best Practices

1. **Error Handling**: Provide clear error messages
   - API key missing
   - Invalid model name
   - Network errors
   - Rate limits

2. **Model Defaults**: Choose sensible defaults
   - Balance cost and quality
   - Prefer smaller models for simple tasks

3. **Timeouts**: Set appropriate HTTP timeouts
   ```go
   client: &http.Client{
       Timeout: 30 * time.Second,
   }
   ```

4. **Rate Limiting**: Handle rate limit errors gracefully
   ```go
   if resp.StatusCode == 429 {
       return "", fmt.Errorf("rate limit exceeded, please try again later")
   }
   ```

5. **Context Passing**: Use all available context for better results
   - System info (OS, shell, directory)
   - Recent commands and output
   - Project type detection
   - User preferences

6. **Command Cleaning**: Remove markdown and formatting artifacts
   - Strip code block markers
   - Remove explanations
   - Trim whitespace

## Examples to Study

Look at these existing providers for reference:

- **`gemini.go`** - Google's Gemini API
- **`openai.go`** - OpenAI's ChatGPT API
- **`ollama.go`** - Local Ollama API
- **`anthropic.go`** - Anthropic's Claude API

## Common Pitfalls

1. ❌ **Not cleaning the output** - LLMs often return markdown
2. ❌ **Ignoring conversation history** - Results in poor follow-up commands
3. ❌ **Not handling errors** - Leads to poor user experience
4. ❌ **Hardcoding values** - Use configuration instead
5. ❌ **Not validating API keys** - Check early, fail fast

## Getting Help

- Check existing provider implementations
- Ask in [GitHub Discussions](https://github.com/fabiobrug/mako/discussions)
- Open an issue if you need help

## Contributing

Once your provider is working:

1. Test thoroughly with different scenarios
2. Update documentation
3. Add configuration examples to `.env.example`
4. Submit a Pull Request

We welcome contributions! 🎉
