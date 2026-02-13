package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type OllamaProvider struct {
	model   string
	baseURL string
	client  *http.Client
}

// NewOllamaProvider creates a new Ollama provider for local LLM inference
func NewOllamaProvider(cfg *ProviderConfig) (*OllamaProvider, error) {
	model := cfg.Model
	if model == "" {
		model = "llama3.2" // Default model
	}
	
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	
	// Verify Ollama is running by checking the API
	client := &http.Client{}
	resp, err := client.Get(fmt.Sprintf("%s/api/tags", baseURL))
	if err != nil {
		return nil, fmt.Errorf("Ollama not reachable at %s. Make sure Ollama is running: https://ollama.ai", baseURL)
	}
	resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ollama API returned status %d. Make sure Ollama is running", resp.StatusCode)
	}
	
	return &OllamaProvider{
		model:   model,
		baseURL: baseURL,
		client:  client,
	}, nil
}

func (o *OllamaProvider) GenerateCommand(userRequest string, context SystemContext) (string, error) {
	return o.GenerateCommandWithConversation(userRequest, context, nil)
}

func (o *OllamaProvider) GenerateCommandWithConversation(userRequest string, context SystemContext, conversation *ConversationHistory) (string, error) {
	prompt := o.buildPrompt(userRequest, context, conversation)
	
	requestBody := map[string]interface{}{
		"model":  o.model,
		"prompt": prompt,
		"stream": false,
		"options": map[string]interface{}{
			"temperature": 0.1,
			"num_predict": 200,
		},
	}
	
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}
	
	url := fmt.Sprintf("%s/api/generate", o.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
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
		Response string `json:"response"`
		Done     bool   `json:"done"`
	}
	
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}
	
	command := o.cleanCommand(response.Response)
	return command, nil
}

func (o *OllamaProvider) ExplainError(failedCommand string, errorOutput string, context SystemContext) (string, error) {
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

func (o *OllamaProvider) ExplainCommand(command string, context SystemContext) (string, error) {
	prompt := fmt.Sprintf(`Explain this shell command in simple, clear terms.

System: %s | Shell: %s | Dir: %s

Command: %s

Provide a brief explanation (2-3 sentences) covering:
1. What the command does
2. What the key flags/options mean
3. Any potential side effects or warnings
4. Security warnings if the command has any security implications

Be concise and user-friendly.`,
		context.OS,
		context.Shell,
		context.CurrentDir,
		command,
	)
	
	return o.sendRequest(prompt, 1024, 0.3)
}

func (o *OllamaProvider) SuggestAlternatives(command string, context SystemContext) (string, error) {
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

func (o *OllamaProvider) sendRequest(prompt string, maxTokens int, temperature float64) (string, error) {
	requestBody := map[string]interface{}{
		"model":  o.model,
		"prompt": prompt,
		"stream": false,
		"options": map[string]interface{}{
			"temperature": temperature,
			"num_predict": maxTokens,
		},
	}
	
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}
	
	url := fmt.Sprintf("%s/api/generate", o.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	
	req.Header.Set("Content-Type", "application/json")
	
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
		Response string `json:"response"`
	}
	
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}
	
	return response.Response, nil
}

func (o *OllamaProvider) buildPrompt(userRequest string, context SystemContext, conversation *ConversationHistory) string {
	var promptBuild strings.Builder
	
	promptBuild.WriteString("You are a shell command generator. Your ONLY job is to output a single shell command.\n\n")
	promptBuild.WriteString("RULES:\n")
	promptBuild.WriteString("- Output ONLY the command, nothing else\n")
	promptBuild.WriteString("- NO explanations, NO markdown, NO code blocks\n")
	promptBuild.WriteString("- Use proper flags and options for the task\n")
	promptBuild.WriteString("- The command must be safe and correct\n\n")
	
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
	promptBuild.WriteString("\n\nGenerate ONLY the shell command. No markdown, no backticks, no explanations. Just the raw command.")
	
	return promptBuild.String()
}

func (o *OllamaProvider) cleanCommand(command string) string {
	command = strings.TrimPrefix(command, "```bash")
	command = strings.TrimPrefix(command, "```sh")
	command = strings.TrimPrefix(command, "```")
	command = strings.TrimSuffix(command, "```")
	command = strings.TrimSpace(command)
	return command
}

// Ollama Embedding Provider
type OllamaEmbeddingProvider struct {
	model   string
	baseURL string
	client  *http.Client
}

func NewOllamaEmbeddingProvider(cfg *ProviderConfig) (*OllamaEmbeddingProvider, error) {
	model := cfg.Model
	if model == "" {
		model = "nomic-embed-text" // Default embedding model for Ollama
	}
	
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	
	return &OllamaEmbeddingProvider{
		model:   model,
		baseURL: baseURL,
		client:  &http.Client{},
	}, nil
}

func (o *OllamaEmbeddingProvider) GenerateEmbedding(text string) ([]byte, error) {
	requestBody := map[string]interface{}{
		"model":  o.model,
		"prompt": text,
	}
	
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}
	
	url := fmt.Sprintf("%s/api/embeddings", o.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	
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
		Embedding []float32 `json:"embedding"`
	}
	
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	
	if len(response.Embedding) == 0 {
		return nil, fmt.Errorf("no embedding in response")
	}
	
	return VectorToBytes(response.Embedding), nil
}
