package channel

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
)

type Terminal struct {
	Prompt string
}

func NewTerminal() Channel {
	return &Terminal{Prompt: ">> "}
}

// Listen starts a loop that reads from Stdin and sends to a channel
func (t *Terminal) Listen(ctx context.Context) (<-chan IncomingMessage, error) {
	out := make(chan IncomingMessage)
	scanner := bufio.NewScanner(os.Stdin)

	go func() {
		defer close(out)
		fmt.Println("--- Mini-Agent Terminal Session Started ---")
		fmt.Print(t.Prompt)

		for scanner.Scan() {
			text := strings.TrimSpace(scanner.Text())
			if text == "" {
				fmt.Print(t.Prompt)
				continue
			}

			if text == "/exit" || text == "/quit" {
				return
			}

			// Create the message object
			msg := IncomingMessage{
				Text:     text,
				UserID:   "local-user",
				Platform: "terminal",
				Reply: func(response string) error {
					fmt.Printf("\n[Agent]: %s\n\n", response)
					fmt.Print(t.Prompt)
					return nil
				},
			}

			select {
			case out <- msg:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out, nil
}

func (t *Terminal) Send(ctx context.Context, userID string, text string) error {
	fmt.Printf("\n[Notification]: %s\n", text)
	return nil
}
