package base

import (
	"github.com/gdamore/tcell/v2"
)

type Event struct {
	stackable Stackable
	source    tcell.Event
}

func (ev *Event) Init(stackable Stackable, source tcell.Event) {
	ev.stackable = stackable
	ev.source = source
}

func (ev *Event) GetStackable() Stackable {
	return ev.stackable
}

func (ev *Event) GetSource() tcell.Event {
	return ev.source
}
