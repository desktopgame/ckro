package llm

import (
	"context"
)

type ConfirmEvent struct {
	EventBase[struct{}]
	Approve bool
}

func (c *ConfirmEvent) Consume(ctx context.Context) {
	c.ch <- Complete
}

func (c *ConfirmEvent) Cancel(ctx context.Context) {
	c.ch <- Cancel
}

func (c *ConfirmEvent) GetResult() struct{} {
	return struct{}{}
}
