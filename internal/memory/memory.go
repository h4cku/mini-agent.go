package memory

import (
	"context"
	"mini-agent/internal/provider"
)

// Memory handles short-term and long-term storage
type Memory interface {
	// Add adds a message to a specific session/conversation
	Add(ctx context.Context, sessionID string, msg provider.Message) error
	// GetRecent retrieves the last N messages for context
	GetRecent(ctx context.Context, sessionID string, limit int) ([]provider.Message, error)
	// Clear wipes a session
	Clear(ctx context.Context, sessionID string) error
}
