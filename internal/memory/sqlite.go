package memory

import (
	"context"
	"database/sql"
	"fmt"
	"mini-agent/internal/provider"

	_ "modernc.org/sqlite" // Pure Go SQLite driver
)

type SQLiteMemory struct {
	db *sql.DB
}

func NewSQLiteMemory(dbPath string) (Memory, error) {
	// 1. Open the database file (it will be created if it doesn't exist)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// 2. Create the messages table if it doesn't exist
	schema := `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id TEXT NOT NULL,
		role TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_session ON messages(session_id);`

	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return &SQLiteMemory{db: db}, nil
}

func (s *SQLiteMemory) Add(ctx context.Context, sessionID string, msg provider.Message) error {
	query := `INSERT INTO messages (session_id, role, content) VALUES (?, ?, ?)`
	_, err := s.db.ExecContext(ctx, query, sessionID, msg.Role, msg.Content)
	return err
}

func (s *SQLiteMemory) GetRecent(ctx context.Context, sessionID string, limit int) ([]provider.Message, error) {
	// We select the most recent messages, then reverse them to maintain chronological order
	query := `
		SELECT role, content FROM (
			SELECT role, content, id FROM messages 
			WHERE session_id = ? 
			ORDER BY id DESC LIMIT ?
		) ORDER BY id ASC`

	rows, err := s.db.QueryContext(ctx, query, sessionID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []provider.Message
	for rows.Next() {
		var msg provider.Message
		if err := rows.Scan(&msg.Role, &msg.Content); err != nil {
			return nil, err
		}
		history = append(history, msg)
	}
	return history, nil
}

func (s *SQLiteMemory) Clear(ctx context.Context, sessionID string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM messages WHERE session_id = ?", sessionID)
	return err
}
