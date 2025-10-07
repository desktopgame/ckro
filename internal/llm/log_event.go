package llm

import (
	"context"

	"github.com/openai/openai-go/v2"
)

type LogResult struct {
}

type LogEvent struct {
	Body openai.ChatCompletionMessageParamUnion
	EventBase[LogResult]
}

func (l *LogEvent) Consume(ctx context.Context) {
}

func (l *LogEvent) Cancel(ctx context.Context) {
}

func (l *LogEvent) GetError() error {
	return nil
}

func (l *LogEvent) GetResult() LogResult {
	return LogResult{}
}
