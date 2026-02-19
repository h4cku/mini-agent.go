package tool

import "context"

// Tool is a capability the agent can execute
type Tool interface {
	Name() string
	Description() string
	// Definition returns the JSON Schema for the LLM to understand arguments
	Definition() map[string]interface{}
	// Execute runs the actual logic with JSON-stringified arguments
	Execute(ctx context.Context, args string) (string, error)
}
