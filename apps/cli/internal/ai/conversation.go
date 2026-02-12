package ai

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	conversationFile    = ".mako/conversation.json"
	maxConversationTurns = 5 // Keep last 5 exchanges
	conversationTimeout  = 5 * time.Minute // Auto-clear after 5 minutes of inactivity
)

// ConversationTurn represents a single user request and AI response
type ConversationTurn struct {
	Timestamp   time.Time `json:"timestamp"`
	UserRequest string    `json:"user_request"`
	AIResponse  string    `json:"ai_response"`
	Executed    bool      `json:"executed"` // Whether the command was executed
}

// ConversationHistory manages the conversation state
type ConversationHistory struct {
	Turns        []ConversationTurn `json:"turns"`
	LastActivity time.Time          `json:"last_activity"`
	SessionID    string             `json:"session_id"`
}

// LoadConversation loads the conversation history from disk
func LoadConversation() (*ConversationHistory, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	convPath := filepath.Join(homeDir, conversationFile)
	
	data, err := os.ReadFile(convPath)
	if err != nil {
		if os.IsNotExist(err) {
			// No conversation exists, start fresh
			return &ConversationHistory{
				Turns:        []ConversationTurn{},
				LastActivity: time.Now(),
				SessionID:    generateSessionID(),
			}, nil
		}
		return nil, fmt.Errorf("failed to read conversation: %w", err)
	}

	var conv ConversationHistory
	if err := json.Unmarshal(data, &conv); err != nil {
		// Corrupted file, start fresh
		return &ConversationHistory{
			Turns:        []ConversationTurn{},
			LastActivity: time.Now(),
			SessionID:    generateSessionID(),
		}, nil
	}

	// Check if conversation has timed out
	if time.Since(conv.LastActivity) > conversationTimeout {
		return &ConversationHistory{
			Turns:        []ConversationTurn{},
			LastActivity: time.Now(),
			SessionID:    generateSessionID(),
		}, nil
	}

	return &conv, nil
}

// SaveConversation saves the conversation history to disk
func (c *ConversationHistory) Save() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	convPath := filepath.Join(homeDir, conversationFile)
	
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(convPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	c.LastActivity = time.Now()

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal conversation: %w", err)
	}

	if err := os.WriteFile(convPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write conversation: %w", err)
	}

	return nil
}

// AddTurn adds a new conversation turn and maintains the limit
func (c *ConversationHistory) AddTurn(userRequest, aiResponse string, executed bool) {
	turn := ConversationTurn{
		Timestamp:   time.Now(),
		UserRequest: userRequest,
		AIResponse:  aiResponse,
		Executed:    executed,
	}

	c.Turns = append(c.Turns, turn)

	// Keep only the last N turns
	if len(c.Turns) > maxConversationTurns {
		c.Turns = c.Turns[len(c.Turns)-maxConversationTurns:]
	}

	c.LastActivity = time.Now()
}

// GetContext returns a formatted string of recent conversation for AI context
func (c *ConversationHistory) GetContext() string {
	if len(c.Turns) == 0 {
		return ""
	}

	var context string
	context += "CONVERSATION HISTORY (most recent at bottom):\n"

	for i, turn := range c.Turns {
		executedMarker := ""
		if turn.Executed {
			executedMarker = " ✓"
		}
		context += fmt.Sprintf("%d. User: %s\n", i+1, turn.UserRequest)
		context += fmt.Sprintf("   AI: %s%s\n", turn.AIResponse, executedMarker)
	}

	context += "\nBuild upon this conversation context. The user's current request may be:\n"
	context += "- A refinement of the previous command\n"
	context += "- A follow-up question about the same topic\n"
	context += "- A new request (if clearly unrelated)\n\n"

	return context
}

// Clear removes all conversation history
func (c *ConversationHistory) Clear() {
	c.Turns = []ConversationTurn{}
	c.LastActivity = time.Now()
	c.SessionID = generateSessionID()
}

// IsActive returns whether there's an active conversation
func (c *ConversationHistory) IsActive() bool {
	return len(c.Turns) > 0 && time.Since(c.LastActivity) < conversationTimeout
}

// ClearConversation clears the conversation file
func ClearConversation() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	convPath := filepath.Join(homeDir, conversationFile)
	
	// Remove the file
	if err := os.Remove(convPath); err != nil {
		if os.IsNotExist(err) {
			return nil // Already cleared
		}
		return fmt.Errorf("failed to clear conversation: %w", err)
	}

	return nil
}

func generateSessionID() string {
	return fmt.Sprintf("session_%d", time.Now().Unix())
}

// GetLastCommand returns the last AI-generated command, if any
func (c *ConversationHistory) GetLastCommand() string {
	if len(c.Turns) == 0 {
		return ""
	}
	return c.Turns[len(c.Turns)-1].AIResponse
}

// GetLastUserRequest returns the last user request, if any
func (c *ConversationHistory) GetLastUserRequest() string {
	if len(c.Turns) == 0 {
		return ""
	}
	return c.Turns[len(c.Turns)-1].UserRequest
}
