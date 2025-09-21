package llm

import (
	"context"
)

type ErrorEvent struct {
	EventBase[error]
	e error
}

func (e *ErrorEvent) Consume(ctx context.Context) {
}

func (e *ErrorEvent) Cancel(ctx context.Context) {
}

func (e *ErrorEvent) GetError() error {
	return e.e
}

func (e *ErrorEvent) GetResult() error {
	return e.e
}
