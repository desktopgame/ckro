package llm

import (
	"context"

	"github.com/openai/openai-go/v2"
)

type ChatEvent struct {
	EventBase[*openai.ChatCompletion]
	client    *openai.Client
	inputList []openai.ChatCompletionMessageParamUnion
	model     string
}

func (c *ChatEvent) Consume(ctx context.Context) {
	chatCompletion, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: c.inputList,
		Model:    c.model,
	})
	if err == nil {
		c.result = chatCompletion
		c.ch <- Complete
	}
}

func (c *ChatEvent) Cancel(ctx context.Context) {
	c.ch <- Cancel
}

func (c *ChatEvent) GetResult() *openai.ChatCompletion {
	return c.result
}
