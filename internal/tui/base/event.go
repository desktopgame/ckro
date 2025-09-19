package base

import (
	"github.com/gdamore/tcell/v2"
)

// Event is notify a terminal updates.
type Event struct {
	runtime Runtime
	source  tcell.Event
}

// Init is initialize Event.
func (ev *Event) Init(runtime Runtime, source tcell.Event) {
	ev.runtime = runtime
	ev.source = source
}

// GetRuntime returns Runtime.
// each controls is able to access to stack of control through this object.
// in example, use when needs to controls show modal another controls.
func (ev *Event) GetRuntime() Runtime {
	return ev.runtime
}

// GetSource returns raw event of tcell.
func (ev *Event) GetSource() tcell.Event {
	return ev.source
}
