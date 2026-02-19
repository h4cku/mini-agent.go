package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type LlamaCppProvider struct {
	BaseURL string
	Model   string
}

func NewLlamaCpp(baseURL string, model string) Provider {
	return &LlamaCppProvider{
		BaseURL: baseURL, // e.g., "http://localhost:8080/v1"
		Model:   model,
	}
}

func (l *LlamaCppProvider) Name() string {
	return "llama.cpp"
}

func (l *LlamaCppProvider) Generate(ctx context.Context, history []Message, tools []ToolDefinition) (Message, error) {
	// Llama.cpp uses the standard OpenAI payload format
	type wrappedTool struct {
		Type     string         `json:"type"`
		Function ToolDefinition `json:"function"`
	}

	var formattedTools []wrappedTool
	for _, t := range tools {
		formattedTools = append(formattedTools, wrappedTool{
			Type:     "function", // <--- THIS is what was missing
			Function: t,
		})
	}

	payload := map[string]any{
		"model":    l.Model,
		"messages": history,
		"tools":    formattedTools,
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", l.BaseURL+"/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return Message{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Message{}, err
	}
	defer resp.Body.Close()

	// Handle standard OpenAI-style response
	var openAIResp struct {
		Choices []struct {
			Message struct {
				Role      string `json:"role"`
				Content   string `json:"content"`
				ToolCalls []struct {
					ID       string `json:"id"`
					Function struct {
						Name      string `json:"name"`
						Arguments string `json:"arguments"`
					} `json:"function"`
				} `json:"tool_calls"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return Message{}, err
	}

	if len(openAIResp.Choices) == 0 {
		return Message{}, fmt.Errorf("empty response from llama.cpp")
	}

	// Transform the response into our internal Message format
	choice := openAIResp.Choices[0].Message
	msg := Message{
		Role:    choice.Role,
		Content: choice.Content,
	}

	// If the model wants to call a tool, we map it to our ToolCallID
	if len(choice.ToolCalls) > 0 {
		msg.ToolCallID = choice.ToolCalls[0].Function.Name
		msg.Content = choice.ToolCalls[0].Function.Arguments
	}

	return msg, nil
}
