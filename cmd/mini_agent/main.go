package main

import (
	"context"
	"log"
	"mini-agent/internal/agent"
	"mini-agent/internal/channel"
	"mini-agent/internal/memory"
	"mini-agent/internal/provider"
	"mini-agent/internal/tool"
)

func main() {
	ctx := context.Background()

	// 1. Setup Provider (llama.cpp)
	llm := provider.NewLlamaCpp("http://localhost:8080/v1", "llama-3.2")

	// 2. Setup Memory
	// In-memory
	mem := memory.NewInMemory()
	// SQLite
	// mem, err := memory.NewSQLiteMemory("data.db")
	// if err != nil {
	// 	log.Fatalf("failed to open database: %v", err)
	// }

	// 3. Initialize Agent
	systemPrompt := "You are an agent, a sleek autonomous AI assistant. Be concise."
	myAgent := agent.New(llm, mem, systemPrompt)

	// Initialize and Register the Shell Tool
	shell := tool.NewShellTool()
	ft := tool.NewFileTool("./agent_workspace")
	myAgent.RegisterTool(shell)
	myAgent.RegisterTool(ft)

	// 4. Start Terminal Channel
	term := channel.NewTerminal()
	messages, _ := term.Listen(ctx)

	// 5. The Main Event Loop
	for msg := range messages {
		// Processing...
		response, err := myAgent.Process(ctx, msg.UserID, msg.Text)
		if err != nil {
			log.Printf("Error: %v", err)
			msg.Reply("I encountered an error processing that.")
			continue
		}

		// Send the response back to the terminal
		msg.Reply(response)
	}
}
