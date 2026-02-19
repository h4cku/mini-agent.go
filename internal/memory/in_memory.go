package memory

import (
	"context"
	"mini-agent/internal/provider"
	"sync"
)

type InMemory struct {
	// Map of sessionID to a slice of messages
	storage map[string][]provider.Message
	mu      sync.RWMutex
}

func NewInMemory() Memory {
	return &InMemory{
		storage: make(map[string][]provider.Message),
	}
}

// Add appends a message to the history of a specific session
func (m *InMemory) Add(ctx context.Context, sessionID string, msg provider.Message) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.storage[sessionID] = append(m.storage[sessionID], msg)
	return nil
}

// GetRecent retrieves the last N messages for context
func (m *InMemory) GetRecent(ctx context.Context, sessionID string, limit int) ([]provider.Message, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	history, exists := m.storage[sessionID]
	if !exists {
		return []provider.Message{}, nil
	}

	// If history is longer than limit, slice it
	if len(history) > limit {
		return history[len(history)-limit:], nil
	}

	return history, nil
}

// Clear wipes the session history
func (m *InMemory) Clear(ctx context.Context, sessionID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.storage, sessionID)
	return nil
}
