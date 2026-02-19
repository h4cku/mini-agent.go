package channel

import "context"

// IncomingMessage represents a message from a user on any platform
type IncomingMessage struct {
	Text     string
	UserID   string
	Platform string             // "telegram", "terminal", etc.
	Reply    func(string) error // Helper to reply directly to this message
}

// Channel defines an input/output method for the user
type Channel interface {
	// Listen starts the long-polling or webhook server
	Listen(ctx context.Context) (<-chan IncomingMessage, error)
	// Send pushes a proactive message to a specific user
	Send(ctx context.Context, userID string, text string) error
}
