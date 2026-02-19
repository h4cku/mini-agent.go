package provider

import "context"

// Message represents a single turn in a conversation
type Message struct {
	Role       string `json:"role"` // system, user, assistant, tool
	Content    string `json:"content"`
	ToolCallID string `json:"tool_call_id,omitempty"` // Optional: for linking tool results
}

// Provider defines how to talk to an LLM
type Provider interface {
	// Generate takes a history and returns the next response
	Generate(ctx context.Context, history []Message, tools []ToolDefinition) (Message, error)

	// Name returns the provider identifier (e.g., "openai-gpt4o")
	Name() string
}

// ToolDefinition is what we send to the LLM so it knows what it can do
type ToolDefinition struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Type        string         `json:"type"`
	Parameters  map[string]any `json:"parameters"` // JSON Schema
}
