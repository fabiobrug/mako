package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fabiobrug/mako.git/internal/config"
	"github.com/fabiobrug/mako.git/internal/retry"
)

type GeminiProvider struct {
	apiKey    string
	model     string
	client    *http.Client
	executor  *retry.ResilientExecutor
}

// NewGeminiProvider creates a new Gemini AI provider
func NewGeminiProvider(cfg *ProviderConfig) (*GeminiProvider, error) {
	apiKey := cfg.APIKey
	
	// Fallback to environment variable
	if apiKey == "" {
		apiKey = os.Getenv("GEMINI_API_KEY")
	}
	
	// Fallback to config file
	if apiKey == "" {
		fileCfg, err := config.LoadConfig()
		if err == nil && fileCfg.APIKey != "" {
			apiKey = fileCfg.APIKey
		}
	}
	
	if apiKey == "" {
		return nil, fmt.Errorf("Gemini API key not found. Set LLM_API_KEY in .env or use: mako config set api_key <your-key>")
	}
	
	// Default model if not specified
	model := cfg.Model
	if model == "" {
		model = "gemini-2.5-flash"
	}

	// Configure retry with exponential backoff
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

	// Configure circuit breaker
	cbConfig := &retry.CircuitBreakerConfig{
		MaxFailures: 5,
		Timeout:     30 * time.Second,
		MaxRequests: 1,
		OnStateChange: func(from, to retry.State) {
			if to == retry.StateOpen {
				fmt.Fprintf(os.Stderr, "⚠️  Gemini API circuit breaker opened due to repeated failures\n")
			} else if to == retry.StateClosed {
				fmt.Fprintf(os.Stderr, "✓ Gemini API circuit breaker closed - service recovered\n")
			}
		},
	}

	circuitBreaker := retry.NewCircuitBreaker(cbConfig)
	executor := retry.NewResilientExecutor(retryConfig, circuitBreaker)

	return &GeminiProvider{
		apiKey:   apiKey,
		model:    model,
		client:   &http.Client{Timeout: 30 * time.Second},
		executor: executor,
	}, nil
}

// NewGeminiClient creates a legacy Gemini client for backward compatibility
// Deprecated: Use NewGeminiProvider instead
func NewGeminiClient() (*GeminiProvider, error) {
	return NewGeminiProvider(&ProviderConfig{
		Provider: "gemini",
	})
}

func (g *GeminiProvider) GenerateCommand(userRequest string, systemCtx SystemContext) (string, error) {
	return g.GenerateCommandWithConversation(userRequest, systemCtx, nil)
}

// GenerateCommandWithConversation generates a command with conversation context
func (g *GeminiProvider) GenerateCommandWithConversation(userRequest string, systemCtx SystemContext, conversation *ConversationHistory) (string, error) {
	prompt := g.buildPromptWithConversation(userRequest, systemCtx, conversation)

	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{"text": prompt},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature":     0.1,
			"maxOutputTokens": 200,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Use resilient executor for retry + circuit breaker
	ctx := context.Background()
	command, err := g.executeWithRetry(ctx, func() (string, error) {
		return g.makeRequest(jsonData)
	})

	if err != nil {
		return "", err
	}

	return g.cleanCommand(command), nil
}

// executeWithRetry is a helper to execute API calls with retry and circuit breaker
func (g *GeminiProvider) executeWithRetry(ctx context.Context, fn func() (string, error)) (string, error) {
	return retry.DoWithResult(ctx, g.executor.Retry, func() (string, error) {
		var result string
		err := g.executor.CircuitBreaker.Execute(ctx, func() error {
			var innerErr error
			result, innerErr = fn()
			return innerErr
		})
		return result, err
	})
}

// makeRequest performs the actual HTTP request (extracted for retry/circuit breaker)
func (g *GeminiProvider) makeRequest(jsonData []byte) (string, error) {
	apiURL := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent", g.model)
	url := fmt.Sprintf("%s?key=%s", apiURL, g.apiKey)
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
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
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	return response.Candidates[0].Content.Parts[0].Text, nil
}

func (g *GeminiProvider) buildPromptWithConversation(userRequest string, systemCtx SystemContext, conversation *ConversationHistory) string {
	var promptBuild strings.Builder

	// Include conversation history if available
	if conversation != nil && conversation.IsActive() {
		promptBuild.WriteString(conversation.GetContext())
	}

	// Include user preferences if available
	if systemCtx.Preferences != nil {
		prefHints := systemCtx.Preferences.GetPreferenceHints()
		if prefHints != "" {
			promptBuild.WriteString(prefHints)
		}
	}

	return g.buildPromptCore(userRequest, systemCtx, &promptBuild)
}

func (g *GeminiProvider) buildPrompt(userRequest string, systemCtx SystemContext) string {
	var promptBuild strings.Builder
	return g.buildPromptCore(userRequest, systemCtx, &promptBuild)
}

func (g *GeminiProvider) buildPromptCore(userRequest string, systemCtx SystemContext, promptBuild *strings.Builder) string {
	promptBuild.WriteString("You are a shell command generator. Your ONLY job is to output a single shell command.\n\n")
	promptBuild.WriteString("RULES:\n")
	promptBuild.WriteString("- Output ONLY the command, nothing else\n")
	promptBuild.WriteString("- NO explanations, NO markdown, NO code blocks\n")
	promptBuild.WriteString("- Use proper flags and options for the task\n")
	promptBuild.WriteString("- The command must be safe and correct\n\n")

	promptBuild.WriteString(fmt.Sprintf("System: %s\n", systemCtx.OS))
	promptBuild.WriteString(fmt.Sprintf("Shell: %s\n", systemCtx.Shell))
	promptBuild.WriteString(fmt.Sprintf("Current directory: %s\n", systemCtx.CurrentDir))

	// NEW: Include project context for smart suggestions
	if systemCtx.Project != nil {
		projectHint := systemCtx.Project.GetProjectHint()
		if projectHint != "" {
			promptBuild.WriteString(fmt.Sprintf("Project type: %s\n", projectHint))
		}

		// Give AI specific hints about available commands
		if systemCtx.Project.TestCmd != "" {
			promptBuild.WriteString(fmt.Sprintf("Test command: %s\n", systemCtx.Project.TestCmd))
		}
		if systemCtx.Project.BuildCmd != "" {
			promptBuild.WriteString(fmt.Sprintf("Build command: %s\n", systemCtx.Project.BuildCmd))
		}
		if systemCtx.Project.RunCmd != "" {
			promptBuild.WriteString(fmt.Sprintf("Run command: %s\n", systemCtx.Project.RunCmd))
		}
	}

	// NEW: Include files in current directory
	if len(systemCtx.WorkingFiles) > 0 {
		promptBuild.WriteString(fmt.Sprintf("Files in directory: %s\n", strings.Join(systemCtx.WorkingFiles, ", ")))
	}

	// NEW: Include recent commands for context
	if len(systemCtx.RecentCommands) > 0 {
		promptBuild.WriteString("\nRecent commands:\n")
		for _, cmd := range systemCtx.RecentCommands {
			promptBuild.WriteString(fmt.Sprintf("  %s\n", cmd))
		}
	}

	// ENHANCED: Actually use the recent output!
	if len(systemCtx.RecentOutput) > 0 {
		promptBuild.WriteString("\nRecent terminal output:\n")

		// Analyze output for context hints
		hints := AnalyzeRecentOutput(systemCtx.RecentOutput)

		// Add intelligent context based on output
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
		startIdx := len(systemCtx.RecentOutput) - 5
		if startIdx < 0 {
			startIdx = 0
		}
		for _, line := range systemCtx.RecentOutput[startIdx:] {
			if strings.TrimSpace(line) != "" {
				promptBuild.WriteString(fmt.Sprintf("  %s\n", line))
			}
		}
	}

	promptBuild.WriteString("\nUser request: ")
	promptBuild.WriteString(userRequest)
	promptBuild.WriteString("\n\n")

	// NEW: Enhanced guidance for command composition
	promptBuild.WriteString("COMMAND COMPOSITION GUIDE:\n")
	promptBuild.WriteString("- Use pipes (|) to chain commands: cmd1 | cmd2\n")
	promptBuild.WriteString("- Use && for sequential execution: cmd1 && cmd2 (only run cmd2 if cmd1 succeeds)\n")
	promptBuild.WriteString("- Use || for alternatives: cmd1 || cmd2 (run cmd2 only if cmd1 fails)\n")
	promptBuild.WriteString("- Use ; for unconditional sequence: cmd1 ; cmd2 (always run both)\n")
	promptBuild.WriteString("- Combine multiple operations when needed for complex tasks\n")
	promptBuild.WriteString("- For filtering/processing: grep, awk, sed, sort, uniq, head, tail\n")
	promptBuild.WriteString("- For monitoring: watch, tail -f\n\n")

	promptBuild.WriteString("EXAMPLES:\n")
	promptBuild.WriteString("- Find errors and count: grep ERROR log.txt | wc -l\n")
	promptBuild.WriteString("- Find, count, show top 10: grep ERROR *.log | sort | uniq -c | sort -rn | head -10\n")
	promptBuild.WriteString("- Build and run: make build && ./app\n")
	promptBuild.WriteString("- Try command or fallback: command -v docker || echo \"Docker not installed\"\n\n")

	promptBuild.WriteString("Generate ONLY the shell command. No markdown, no backticks, no explanations. Just the raw command.")

	return promptBuild.String()
}

func (g *GeminiProvider) cleanCommand(command string) string {
	command = strings.TrimPrefix(command, "```bash")
	command = strings.TrimPrefix(command, "```sh")
	command = strings.TrimPrefix(command, "```")
	command = strings.TrimSuffix(command, "```")

	command = strings.TrimSpace(command)

	return command
}

func (g *GeminiProvider) ExplainError(failedCommand string, errorOutput string, systemCtx SystemContext) (string, error) {
	prompt := fmt.Sprintf(`Shell debugging assistant. Analyze this error briefly.

System: %s | Shell: %s | Dir: %s

Command: %s
Error: %s

Provide:
EXPLANATION: Brief 1-2 sentence explanation of the error
SUGGESTION: A corrected command (if applicable) or next steps

Be concise and actionable.`,
		systemCtx.OS,
		systemCtx.Shell,
		systemCtx.CurrentDir,
		failedCommand,
		errorOutput,
	)

	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{"text": prompt},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature":     0.3,
			"maxOutputTokens": 2048,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	// Use resilient executor for retry + circuit breaker
	ctx := context.Background()
	explanation, err := g.executeWithRetry(ctx, func() (string, error) {
		return g.makeRequestWithFinishReason(jsonData)
	})

	return explanation, err
}

// makeRequestWithFinishReason performs request and checks finish reason
func (g *GeminiProvider) makeRequestWithFinishReason(jsonData []byte) (string, error) {
	apiURL := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent", g.model)
	url := fmt.Sprintf("%s?key=%s", apiURL, g.apiKey)
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
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
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
			FinishReason string `json:"finishReason"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	// Check if response was truncated
	finishReason := response.Candidates[0].FinishReason
	if finishReason != "" && finishReason != "STOP" {
		fmt.Fprintf(os.Stderr, "Warning: API response truncated (reason: %s)\n", finishReason)
	}

	return response.Candidates[0].Content.Parts[0].Text, nil
}

// ExplainCommand generates a human-readable explanation of what a command does
func (g *GeminiProvider) ExplainCommand(command string, systemCtx SystemContext) (string, error) {
	prompt := fmt.Sprintf(`Explain this shell command in simple, clear terms.

System: %s | Shell: %s | Dir: %s

Command: %s

Provide a brief explanation (2-3 sentences) covering:
1. What the command does
2. What the key flags/options mean
3. Any potential side effects or warnings
4. **Security warnings** if the command has any security implications (destructive operations, permission changes, network access, etc.)

Be concise and user-friendly. If there are security concerns, highlight them clearly.`,
		systemCtx.OS,
		systemCtx.Shell,
		systemCtx.CurrentDir,
		command,
	)

	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{"text": prompt},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature":     0.3,
			"maxOutputTokens": 1024,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	// Use resilient executor for retry + circuit breaker
	ctx := context.Background()
	explanation, err := g.executeWithRetry(ctx, func() (string, error) {
		return g.makeRequest(jsonData)
	})

	return explanation, err
}

// SuggestAlternatives generates alternative commands that accomplish the same goal
func (g *GeminiProvider) SuggestAlternatives(command string, systemCtx SystemContext) (string, error) {
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
		systemCtx.OS,
		systemCtx.Shell,
		systemCtx.CurrentDir,
		command,
	)

	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{"text": prompt},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature":     0.5,
			"maxOutputTokens": 1024,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	// Use resilient executor for retry + circuit breaker
	ctx := context.Background()
	alternatives, err := g.executeWithRetry(ctx, func() (string, error) {
		return g.makeRequest(jsonData)
	})

	return alternatives, err
}
