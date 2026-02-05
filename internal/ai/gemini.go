package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const geminiAPIURL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent"

type GeminiClient struct {
	apiKey string
	client *http.Client
}

func NewGeminiClient() (*GeminiClient, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	return &GeminiClient{
		apiKey: apiKey,
		client: &http.Client{},
	}, nil
}

type SystemContext struct {
	OS           string
	Shell        string
	CurrentDir   string
	RecentOutput []string
}

func (g *GeminiClient) GenerateCommand(userRequest string, context SystemContext) (string, error) {
	prompt := g.buildPrompt(userRequest, context)

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

	url := fmt.Sprintf("%s?key=%s", geminiAPIURL, g.apiKey)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "applications/json")

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

	command := response.Candidates[0].Content.Parts[0].Text
	command = g.cleanCommand(command)

	return command, nil
}

func (g *GeminiClient) buildPrompt(userRequest string, context SystemContext) string {
	var promptBuild strings.Builder

	promptBuild.WriteString("You are a shell command generator. Your ONLY job is to output a single shell command.\n\n")
	promptBuild.WriteString("RULES:\n")
	promptBuild.WriteString("- Output ONLY the command, nothing else\n")
	promptBuild.WriteString("- NO explanations, NO markdown, NO code blocks\n")
	promptBuild.WriteString("- Use proper flags and options for the task\n")
	promptBuild.WriteString("- The command must be safe and correct\n\n")

	promptBuild.WriteString(fmt.Sprintf("System: %s\n", context.OS))
	promptBuild.WriteString(fmt.Sprintf("Shell: %s\n", context.Shell))
	promptBuild.WriteString(fmt.Sprintf("Current directory: %s\n\n", context.CurrentDir))

	if len(context.RecentOutput) > 0 {
		promptBuild.WriteString("\nRecent terminal output:\n")
		for _, line := range context.RecentOutput {
			promptBuild.WriteString(fmt.Sprintf(" %s\n", line))
		}
	}

	promptBuild.WriteString("\nUser request: ")
	promptBuild.WriteString(userRequest)
	promptBuild.WriteString("\n\nGenerate ONLY the shell command. No markdown, no backticks, no explanations. Just the raw command.")

	return promptBuild.String()
}

func (g *GeminiClient) cleanCommand(command string) string {
	command = strings.TrimPrefix(command, "```bash")
	command = strings.TrimPrefix(command, "```sh")
	command = strings.TrimPrefix(command, "```")
	command = strings.TrimSuffix(command, "```")

	command = strings.TrimSpace(command)

	return command
}

func (g *GeminiClient) ExplainError(failedCommand string, errorOutput string, context SystemContext) (string, error) {
	prompt := fmt.Sprintf(`You are a helpful shell debugging assistant.

System: %s
Shell: %s
Current directory: %s

Failed command: %s

Error output:
%s

Provide a brief explanation of what went wrong and suggest a corrected command.
Format your response as:
EXPLANATION: <brief explanation>
SUGGESTION: <corrected command>`,
		context.OS,
		context.Shell,
		context.CurrentDir,
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
			"maxOutputTokens": 300,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s?key=%s", geminiAPIURL, g.apiKey)
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
		return "", err
	}

	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	return response.Candidates[0].Content.Parts[0].Text, nil
}
