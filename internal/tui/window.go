package tui

import (
	"github.com/gdamore/tcell/v2"
)

type Window struct {
	g               Graphics
	screen          tcell.Screen
	stack           Stack
	focusManager    FocusManager
	width           int
	height          int
	backgroundTasks int
}

func (w *Window) Init(s tcell.Screen, width int, height int) {
	w.g = Graphics{}
	w.screen = s
	w.g.Init(s)
	w.g.Resize(width, height)
}

func (w *Window) Push(layer Layer) {
	w.stack.Layers = append(w.stack.Layers, layer)
	w.stack.Top = len(w.stack.Layers) - 1
	w.stack.Traverse(&w.focusManager)

	if w.width > 0 && w.height > 0 {
		w.stack.Layout(w.width, w.height)
	}

	w.focusManager.Grab()
}

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

func (w *Window) BeginBackground() {
	w.backgroundTasks++
}

func (w *Window) EndBackground() {
	w.backgroundTasks--
}

func (w *Window) Repaint() {
	w.screen.PostEvent(tcell.NewEventInterrupt(RepaintMessage{}))
}

func (w *Window) DoInBackground() bool {
	return w.backgroundTasks > 0
}

func (w *Window) Top() Control {
	if w.stack.Top >= 0 {
		return w.stack.Layers[w.stack.Top].Control
	}
	return nil
}

func (w *Window) GetLayerCount() int {
	return len(w.stack.Layers)
}

func (w *Window) Blit(width int, height int) {
	mw, mh := w.stack.MinimumSize(width, height)

	if mw <= width && mh <= height {
		w.g.Clear()
		w.stack.Update()
		w.stack.Draw(&w.g)
	}
}

func (w *Window) FocusPrev() {
	w.focusManager.FocusPrev()
}

func (w *Window) FocusNext() {
	w.focusManager.FocusNext()
}

func (w *Window) Resize(width int, height int) {
	w.g.Resize(width, height)
	if width > 0 && height > 0 {
		w.stack.Layout(width, height)
	}
	w.width = width
	w.height = height
}

func (w *Window) Handle(ev Event) {
	w.focusManager.Handle(ev)
}
