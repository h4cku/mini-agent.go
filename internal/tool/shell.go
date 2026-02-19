package tool

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type ShellTool struct {
	AllowedCommands []string // Optional: restrict to 'ls', 'echo', etc.
}

func NewShellTool() Tool {
	return &ShellTool{}
}

func (s *ShellTool) Name() string {
	return "execute_shell"
}

func (s *ShellTool) Description() string {
	return "Executes a shell command on the host system. Use this to check files, network, or system status."
}

func (s *ShellTool) Definition() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"command": map[string]any{
				"type":        "string",
				"description": "The full shell command to run (e.g., 'ls -la')",
			},
		},
		"required": []string{"command"},
	}
}

func (s *ShellTool) Execute(ctx context.Context, args string) (string, error) {
	fmt.Println(args)
	// 1. Parse the JSON arguments provided by the LLM
	var params struct {
		Command string `json:"command"`
	}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return "", fmt.Errorf("invalid arguments: %v", err)
	}

	// 2. Security Check (Basic example: block destructive commands)
	forbidden := []string{"rm -rf /", "mkfs", "shutdown"}
	for _, f := range forbidden {
		if strings.Contains(params.Command, f) {
			return "Error: Command is forbidden for security reasons.", nil
		}
	}

	// 3. Execute the command
	// We use 'sh -c' to allow for pipes and redirects
	cmd := exec.CommandContext(ctx, "sh", "-c", params.Command)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Sprintf("Error: %v\nOutput: %s", err, string(output)), nil
	}

	return string(output), nil
}
