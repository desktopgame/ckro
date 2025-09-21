package llm

import (
	"context"

	"github.com/openai/openai-go/v2"
)

type MessageEvent struct {
	EventBase[*openai.ChatCompletion]
}

func (m *MessageEvent) Consume(ctx context.Context) {
}

func (m *MessageEvent) Cancel(ctx context.Context) {
}

func (m *MessageEvent) GetError() error {
	return nil
}

func (m *MessageEvent) GetResult() *openai.ChatCompletion {
	return m.result
}
