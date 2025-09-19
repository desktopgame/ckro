package tui

import "github.com/gdamore/tcell/v2"

type Window struct {
	stack        Stack
	focusManager FocusManager
	width        int
	height       int
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

func (w *Window) Pop() {
	if len(w.stack.Layers) > 0 {
		callable := w.stack.Layers[len(w.stack.Layers)-1].OnPop
		if callable != nil {
			callable(w)
		}

		w.stack.Layers = w.stack.Layers[:len(w.stack.Layers)-1]
		w.stack.Top--
		w.stack.Traverse(&w.focusManager)

		w.focusManager.Grab()
	}
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

func (w *Window) Frame(s tcell.Screen, width int, height int) {
	mw, mh := w.stack.MinimumSize(width, height)

	if mw <= width && mh <= height {
		w.stack.Update()
		w.stack.Draw(s)
	}
}

func (w *Window) FocusPrev() {
	w.focusManager.FocusPrev()
}

func (w *Window) FocusNext() {
	w.focusManager.FocusNext()
}

func (w *Window) Resize(width int, height int) {
	if width > 0 && height > 0 {
		w.stack.Layout(width, height)
	}
	w.width = width
	w.height = height
}

func (w *Window) Handle(ev Event) {
	w.focusManager.Handle(ev)
}
