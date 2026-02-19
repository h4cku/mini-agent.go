package tool

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FileTool struct {
	WorkspaceRoot string
}

func NewFileTool(root string) Tool {
	// Ensure the root is an absolute path
	absRoot, _ := filepath.Abs(root)
	os.MkdirAll(absRoot, 0755)
	return &FileTool{WorkspaceRoot: absRoot}
}

func (f *FileTool) Name() string { return "manage_files" }

func (f *FileTool) Description() string {
	return "Read, write, or list files in the workspace. Operations: 'read', 'write', 'list'."
}

func (f *FileTool) Definition() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"operation": map[string]any{"type": "string", "enum": []string{"read", "write", "list"}},
			"path":      map[string]any{"type": "string", "description": "Relative path inside workspace"},
			"content":   map[string]any{"type": "string", "description": "Content for write operation"},
		},
		"required": []string{"operation", "path"},
	}
}

// validatePath prevents "path traversal" (e.g., ../../etc/passwd)
func (f *FileTool) validatePath(path string) (string, error) {
	finalPath := filepath.Join(f.WorkspaceRoot, path)
	if !strings.HasPrefix(finalPath, f.WorkspaceRoot) {
		return "", fmt.Errorf("security error: attempt to access path outside workspace")
	}
	return finalPath, nil
}

func (f *FileTool) Execute(ctx context.Context, args string) (string, error) {
	var p struct {
		Op      string `json:"operation"`
		Path    string `json:"path"`
		Content string `json:"content"`
	}
	if err := json.Unmarshal([]byte(args), &p); err != nil {
		return "", err
	}

	safePath, err := f.validatePath(p.Path)
	if err != nil {
		return err.Error(), nil
	}

	switch p.Op {
	case "read":
		data, err := os.ReadFile(safePath)
		if err != nil {
			return fmt.Sprintf("Error reading file: %v", err), nil
		}
		return string(data), nil

	case "write":
		err := os.WriteFile(safePath, []byte(p.Content), 0644)
		if err != nil {
			return fmt.Sprintf("Error writing file: %v", err), nil
		}
		return "File written successfully.", nil

	case "list":
		entries, err := os.ReadDir(safePath)
		if err != nil {
			return fmt.Sprintf("Error listing dir: %v", err), nil
		}
		var names []string
		for _, e := range entries {
			names = append(names, e.Name())
		}
		return strings.Join(names, "\n"), nil

	default:
		return "Unknown operation", nil
	}
}
