package llm

import (
	"context"
)

type Status int

const (
	Complete Status = iota
	Error
	Cancel
)

type Event interface {
	Consume(ctx context.Context)
	Cancel(ctx context.Context)
	GetError() error
}

type EventBase[T any] struct {
	ch     chan Status
	result T
}
