package tui

import "github.com/gdamore/tcell/v2"

type Window struct {
	stack        Stack
	focusManager FocusManager
	width        int
	height       int
}

func (w *Window) Push(ctrl Control) {
	w.stack.Layers = append(w.stack.Layers, ctrl)
	w.stack.Top = len(w.stack.Layers) - 1
	w.stack.Traverse(&w.focusManager)

	if w.width > 0 && w.height > 0 {
		w.stack.Layout(w.width, w.height)
	}

	w.focusManager.Grab()
}

func (w *Window) Pop() {
	if len(w.stack.Layers) > 0 {
		w.stack.Layers = w.stack.Layers[:len(w.stack.Layers)-1]
		w.stack.Top--
		w.stack.Traverse(&w.focusManager)

		w.focusManager.Grab()
	}
}

func (w *Window) Frame(s tcell.Screen, width int, height int) {
	mw, mh := w.stack.MinimumSize(width, height)

	if mw <= width && mh <= height {
		w.stack.Update()
		w.stack.Draw(s)
	}
}

func (w *Window) Resize(width int, height int) {
	if width > 0 && height > 0 {
		w.stack.Layout(width, height)
	}
	w.width = width
	w.height = height
}

func (w *Window) Handle(ev tcell.Event) {
	w.focusManager.Handle(ev)
}
