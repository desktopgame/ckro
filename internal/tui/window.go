package tui

import (
	"github.com/gdamore/tcell/v2"
)

// Window is basic implementation of base.Runtime
type Window struct {
	g               Graphics
	screen          tcell.Screen
	stack           Stack
	focusManager    FocusManager
	width           int
	height          int
	backgroundTasks int
}

// Init is initialize Window.
func (w *Window) Init(s tcell.Screen, width int, height int) {
	w.g = Graphics{}
	w.screen = s
	w.g.Init(s)
	w.g.Resize(width, height)
}

// Push is push the layer to stack.
func (w *Window) Push(layer Layer) {
	w.stack.Layers = append(w.stack.Layers, layer)
	w.stack.Top = len(w.stack.Layers) - 1
	w.stack.Traverse(&w.focusManager)

	if w.width > 0 && w.height > 0 {
		w.stack.Layout(w.width, w.height)
	}

	w.focusManager.Grab()
}

// Pop is pop the layer from stack.
func (w *Window) Pop(returnCode int) {
	if len(w.stack.Layers) > 0 {
		l := w.stack.Layers[len(w.stack.Layers)-1]

		w.stack.Layers = w.stack.Layers[:len(w.stack.Layers)-1]
		w.stack.Top--
		w.stack.Traverse(&w.focusManager)
		w.focusManager.Grab()

		callable := l.OnPop
		if callable != nil {
			callable(returnCode)
		}
	}
}

// BeginBackground is increment number of background tasks.
// Runtime is ignore events while zero than bigger of number of background tasks.
func (w *Window) BeginBackground() {
	w.backgroundTasks++
}

// EndBackground is decrement number of background tasks.
func (w *Window) EndBackground() {
	w.backgroundTasks--
}

// Repaint is do request rerender terminal.
func (w *Window) Repaint() {
	w.screen.PostEvent(tcell.NewEventInterrupt(RepaintMessage{}))
}

// DoInBackground returns true if number of background tasks is zero than bigger.
func (w *Window) DoInBackground() bool {
	return w.backgroundTasks > 0
}

// Top returns Control of top layer.
func (w *Window) Top() Control {
	if w.stack.Top >= 0 {
		return w.stack.Layers[w.stack.Top].Control
	}
	return nil
}

// GetLayerCount returns count of layers.
func (w *Window) GetLayerCount() int {
	return len(w.stack.Layers)
}

// Blit is render the layer into terminal.
func (w *Window) Blit(width int, height int) {
	mw, mh := w.stack.MinimumSize(width, height)

	if mw <= width && mh <= height {
		w.g.Clear()
		w.stack.Update()
		w.stack.Draw(&w.g)
	}
}

// FocusPrev is move focus to previous.
func (w *Window) FocusPrev() {
	w.focusManager.FocusPrev()
}

// FocusNext is move focus to next.
func (w *Window) FocusNext() {
	w.focusManager.FocusNext()
}

// Resize is resize layers.
func (w *Window) Resize(width int, height int) {
	w.g.Resize(width, height)
	if width > 0 && height > 0 {
		w.stack.Layout(width, height)
	}
	w.width = width
	w.height = height
}

// Handle is process the event by focused control.
func (w *Window) Handle(ev Event) {
	w.focusManager.Handle(ev)
}
