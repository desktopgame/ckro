package llm

import (
	"context"

	"github.com/openai/openai-go/v2"
)

type ChatEvent struct {
	EventBase[*openai.ChatCompletion]
	client     *openai.Client
	inputList  []openai.ChatCompletionMessageParamUnion
	model      string
	toolParams []openai.ChatCompletionToolUnionParam
	e          error
}

func (c *ChatEvent) Consume(ctx context.Context) {
	chatCompletion, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: c.inputList,
		Model:    c.model,
		Tools:    c.toolParams,
		ToolChoice: openai.ChatCompletionToolChoiceOptionUnionParam{
			OfAuto: openai.String("auto"),
		},
	})
	if err == nil {
		c.result = chatCompletion
		c.ch <- Complete
		return
	}
	c.e = err
	c.ch <- Error
}

func (c *ChatEvent) Cancel(ctx context.Context) {
	c.ch <- Cancel
}

func (c *ChatEvent) GetError() error {
	return c.e
}

func (c *ChatEvent) GetResult() *openai.ChatCompletion {
	return c.result
}
