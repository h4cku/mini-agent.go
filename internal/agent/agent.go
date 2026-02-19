package agent

import (
	"context"
	"fmt"
	"mini-agent/internal/memory"
	"mini-agent/internal/provider"
	"mini-agent/internal/tool"
)

type Agent struct {
	Provider provider.Provider
	Memory   memory.Memory
	Tools    map[string]tool.Tool
	System   string // The base "personality" prompt
}

func New(p provider.Provider, m memory.Memory, sys string) *Agent {
	return &Agent{
		Provider: p,
		Memory:   m,
		Tools:    make(map[string]tool.Tool),
		System:   sys,
	}
}

func (a *Agent) RegisterTool(t tool.Tool) {
	a.Tools[t.Name()] = t
}

func (a *Agent) Process(ctx context.Context, sessionID string, userInput string) (string, error) {
	// 1. Save user input to memory
	userMsg := provider.Message{Role: "user", Content: userInput}
	if err := a.Memory.Add(ctx, sessionID, userMsg); err != nil {
		return "", err
	}

	// 2. Prepare the Tool Definitions for the LLM
	var toolDefs []provider.ToolDefinition
	for _, t := range a.Tools {
		toolDefs = append(toolDefs, provider.ToolDefinition{
			Name:        t.Name(),
			Description: t.Description(),
			Parameters:  t.Definition(),
		})
	}

	// 3. Enter the Generation Loop (Max 5 turns to prevent infinite loops)
	for i := 0; i < 5; i++ {
		// Fetch full history (System + Chat History)
		history, _ := a.Memory.GetRecent(ctx, sessionID, 10)
		fullPrompt := append([]provider.Message{{Role: "system", Content: a.System}}, history...)

		// Ask the LLM what to do
		resp, err := a.Provider.Generate(ctx, fullPrompt, toolDefs)
		if err != nil {
			return "", err
		}

		// If the LLM just replied with text (no tool call), we are done
		if resp.ToolCallID == "" {
			a.Memory.Add(ctx, sessionID, resp)
			return resp.Content, nil
		}

		// 4. Handle Tool Execution
		// For simplicity, assuming LLM returns tool name in 'Content' or 'ToolCallID'
		// In a real OpenAI implementation, this is a specific struct field.
		targetTool, exists := a.Tools[resp.ToolCallID]
		if !exists {
			return "", fmt.Errorf("tool %s not found", resp.ToolCallID)
		}

		// Execute the tool (args usually come from the provider's response)
		result, err := targetTool.Execute(ctx, resp.Content)
		if err != nil {
			result = fmt.Sprintf("Error: %v", err)
		}

		// Add tool result to memory so the LLM can "see" the output
		toolResMsg := provider.Message{
			Role:    "tool",
			Content: result,
		}
		a.Memory.Add(ctx, sessionID, toolResMsg)
	}

	return "I'm sorry, I reached my maximum reasoning steps.", nil
}
