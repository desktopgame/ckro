package tui

import "github.com/gdamore/tcell/v2"

type Stack struct {
	Layers []Layer
	Top    int
}

func (st *Stack) Init() {
	st.Layers = []Layer{}
}

func (st *Stack) Traverse(fm *FocusManager) {
	fm.Traverse(st.Layers[st.Top].Control)
}

func (st *Stack) Update() {
	if st.Top > 1 && st.Layers[st.Top].ClearBackground {
		st.Layers[st.Top].Control.Update()
	} else {
		for i := 0; i <= st.Top; i++ {
			st.Layers[i].Control.Update()
		}
	}
}

func (st *Stack) Draw(s tcell.Screen) {
	if st.Top > 1 && st.Layers[st.Top].ClearBackground {
		st.Layers[st.Top].Control.Draw(s)
	} else {
		for i := 0; i <= st.Top; i++ {
			st.Layers[i].Control.Draw(s)
		}
	}
}

func (st *Stack) Layout(width int, height int) {
	for i := 0; i <= st.Top; i++ {
		st.Layers[i].Control.Layout(width, height)
	}
}

func (st *Stack) MinimumSize(width int, height int) (Width int, Height int) {
	w := -1
	h := -1

	for i := 0; i <= st.Top; i++ {
		mw, mh := st.Layers[i].Control.MinimumSize(width, height)

		if mw > w {
			w = mw
		}
		if mh > h {
			h = mh
		}
	}
	return w, h
}
