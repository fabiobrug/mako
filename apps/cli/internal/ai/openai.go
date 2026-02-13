package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type OpenAIProvider struct {
	apiKey  string
	model   string
	baseURL string
	client  *http.Client
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(cfg *ProviderConfig) (*OpenAIProvider, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("OpenAI API key not found. Set LLM_API_KEY in .env")
	}
	
	model := cfg.Model
	if model == "" {
		model = "gpt-4o-mini" // Default to cost-effective model
	}
	
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	
	return &OpenAIProvider{
		apiKey:  cfg.APIKey,
		model:   model,
		baseURL: baseURL,
		client:  &http.Client{},
	}, nil
}

func (o *OpenAIProvider) GenerateCommand(userRequest string, context SystemContext) (string, error) {
	return o.GenerateCommandWithConversation(userRequest, context, nil)
}

func (o *OpenAIProvider) GenerateCommandWithConversation(userRequest string, context SystemContext, conversation *ConversationHistory) (string, error) {
	prompt := o.buildPrompt(userRequest, context, conversation)
	
	messages := []map[string]interface{}{
		{
			"role":    "system",
			"content": "You are a shell command generator. Output ONLY the command, nothing else. NO explanations, NO markdown, NO code blocks.",
		},
		{
			"role":    "user",
			"content": prompt,
		},
	}
	
	requestBody := map[string]interface{}{
		"model":       o.model,
		"messages":    messages,
		"temperature": 0.1,
		"max_tokens":  200,
	}
	
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}
	
	url := fmt.Sprintf("%s/chat/completions", o.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", o.apiKey))
	
	resp, err := o.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}
	
	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}
	
	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}
	
	command := response.Choices[0].Message.Content
	command = o.cleanCommand(command)
	
	return command, nil
}

func (o *OpenAIProvider) ExplainError(failedCommand string, errorOutput string, context SystemContext) (string, error) {
	prompt := fmt.Sprintf(`Shell debugging assistant. Analyze this error briefly.

System: %s | Shell: %s | Dir: %s

Command: %s
Error: %s

Provide:
EXPLANATION: Brief 1-2 sentence explanation of the error
SUGGESTION: A corrected command (if applicable) or next steps

Be concise and actionable.`,
		context.OS,
		context.Shell,
		context.CurrentDir,
		failedCommand,
		errorOutput,
	)
	
	return o.sendRequest(prompt, 2048, 0.3)
}

func (o *OpenAIProvider) ExplainCommand(command string, context SystemContext) (string, error) {
	prompt := fmt.Sprintf(`Explain this shell command in simple, clear terms.

System: %s | Shell: %s | Dir: %s

Command: %s

Provide a brief explanation (2-3 sentences) covering:
1. What the command does
2. What the key flags/options mean
3. Any potential side effects or warnings
4. **Security warnings** if the command has any security implications

Be concise and user-friendly.`,
		context.OS,
		context.Shell,
		context.CurrentDir,
		command,
	)
	
	return o.sendRequest(prompt, 1024, 0.3)
}

func (o *OpenAIProvider) SuggestAlternatives(command string, context SystemContext) (string, error) {
	prompt := fmt.Sprintf(`Given this shell command, suggest 2-3 alternative ways to accomplish the same goal.

System: %s | Shell: %s | Dir: %s

Original Command: %s

Provide alternatives that:
1. Use different tools/approaches
2. May be safer, faster, or more efficient
3. Have different trade-offs (verbosity, portability, features)

Format each alternative as:
• [command] - brief explanation of difference/advantage

Be concise and practical.`,
		context.OS,
		context.Shell,
		context.CurrentDir,
		command,
	)
	
	return o.sendRequest(prompt, 1024, 0.5)
}

func (o *OpenAIProvider) sendRequest(prompt string, maxTokens int, temperature float64) (string, error) {
	messages := []map[string]interface{}{
		{
			"role":    "user",
			"content": prompt,
		},
	}
	
	requestBody := map[string]interface{}{
		"model":       o.model,
		"messages":    messages,
		"temperature": temperature,
		"max_tokens":  maxTokens,
	}
	
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}
	
	url := fmt.Sprintf("%s/chat/completions", o.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", o.apiKey))
	
	resp, err := o.client.Do(req)
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
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}
	
	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}
	
	return response.Choices[0].Message.Content, nil
}

func (o *OpenAIProvider) buildPrompt(userRequest string, context SystemContext, conversation *ConversationHistory) string {
	var promptBuild strings.Builder
	
	// Include conversation history if available
	if conversation != nil && conversation.IsActive() {
		promptBuild.WriteString(conversation.GetContext())
	}
	
	// Include user preferences if available
	if context.Preferences != nil {
		prefHints := context.Preferences.GetPreferenceHints()
		if prefHints != "" {
			promptBuild.WriteString(prefHints)
		}
	}
	
	promptBuild.WriteString(fmt.Sprintf("System: %s\n", context.OS))
	promptBuild.WriteString(fmt.Sprintf("Shell: %s\n", context.Shell))
	promptBuild.WriteString(fmt.Sprintf("Current directory: %s\n", context.CurrentDir))
	
	// Include project context
	if context.Project != nil {
		projectHint := context.Project.GetProjectHint()
		if projectHint != "" {
			promptBuild.WriteString(fmt.Sprintf("Project type: %s\n", projectHint))
		}
		
		if context.Project.TestCmd != "" {
			promptBuild.WriteString(fmt.Sprintf("Test command: %s\n", context.Project.TestCmd))
		}
		if context.Project.BuildCmd != "" {
			promptBuild.WriteString(fmt.Sprintf("Build command: %s\n", context.Project.BuildCmd))
		}
		if context.Project.RunCmd != "" {
			promptBuild.WriteString(fmt.Sprintf("Run command: %s\n", context.Project.RunCmd))
		}
	}
	
	// Include files in current directory
	if len(context.WorkingFiles) > 0 {
		promptBuild.WriteString(fmt.Sprintf("Files in directory: %s\n", strings.Join(context.WorkingFiles, ", ")))
	}
	
	// Include recent commands
	if len(context.RecentCommands) > 0 {
		promptBuild.WriteString("\nRecent commands:\n")
		for _, cmd := range context.RecentCommands {
			promptBuild.WriteString(fmt.Sprintf("  %s\n", cmd))
		}
	}
	
	// Include recent output analysis
	if len(context.RecentOutput) > 0 {
		promptBuild.WriteString("\nRecent terminal output:\n")
		
		hints := AnalyzeRecentOutput(context.RecentOutput)
		
		if hints["has_errors"] == "true" {
			promptBuild.WriteString("(Note: Recent output contains errors)\n")
		}
		if hints["needs_sudo"] == "true" {
			promptBuild.WriteString("(Note: Previous command had permission issues)\n")
		}
		if workingWith, ok := hints["working_with"]; ok {
			promptBuild.WriteString(fmt.Sprintf("(Note: User is working with %s)\n", workingWith))
		}
		
		// Include last 5 lines of output
		startIdx := len(context.RecentOutput) - 5
		if startIdx < 0 {
			startIdx = 0
		}
		for _, line := range context.RecentOutput[startIdx:] {
			if strings.TrimSpace(line) != "" {
				promptBuild.WriteString(fmt.Sprintf("  %s\n", line))
			}
		}
	}
	
	promptBuild.WriteString("\nUser request: ")
	promptBuild.WriteString(userRequest)
	
	return promptBuild.String()
}

func (o *OpenAIProvider) cleanCommand(command string) string {
	command = strings.TrimPrefix(command, "```bash")
	command = strings.TrimPrefix(command, "```sh")
	command = strings.TrimPrefix(command, "```")
	command = strings.TrimSuffix(command, "```")
	command = strings.TrimSpace(command)
	return command
}

// OpenAI Embedding Provider
type OpenAIEmbeddingProvider struct {
	apiKey  string
	model   string
	baseURL string
	client  *http.Client
}

func NewOpenAIEmbeddingProvider(cfg *ProviderConfig) (*OpenAIEmbeddingProvider, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("OpenAI API key not found")
	}
	
	model := cfg.Model
	if model == "" {
		model = "text-embedding-3-small" // Default embedding model
	}
	
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	
	return &OpenAIEmbeddingProvider{
		apiKey:  cfg.APIKey,
		model:   model,
		baseURL: baseURL,
		client:  &http.Client{},
	}, nil
}

func (o *OpenAIEmbeddingProvider) GenerateEmbedding(text string) ([]byte, error) {
	requestBody := map[string]interface{}{
		"input": text,
		"model": o.model,
	}
	
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}
	
	url := fmt.Sprintf("%s/embeddings", o.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", o.apiKey))
	
	resp, err := o.client.Do(req)
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
		Data []struct {
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
	}
	
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	
	if len(response.Data) == 0 {
		return nil, fmt.Errorf("no embedding in response")
	}
	
	return VectorToBytes(response.Data[0].Embedding), nil
}
