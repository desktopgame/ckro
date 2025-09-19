package base

import (
	"github.com/gdamore/tcell/v2"
)

type Event struct {
	runtime Runtime
	source  tcell.Event
}

func (ev *Event) Init(runtime Runtime, source tcell.Event) {
	ev.runtime = runtime
	ev.source = source
}

func (ev *Event) GetRuntime() Runtime {
	return ev.runtime
}

func (ev *Event) GetSource() tcell.Event {
	return ev.source
}
