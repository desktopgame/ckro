package tui

import "github.com/gdamore/tcell/v2"

type Stack struct {
	Layers []Control
	Top    int
}

func (st *Stack) Init() {
	st.Layers = []Control{}
}

func (st *Stack) Traverse(fm *FocusManager) {
	fm.Traverse(st.Layers[st.Top])
}

func (st *Stack) Update() {
	for i := 0; i <= st.Top; i++ {
		st.Layers[i].Update()
	}
}

func (st *Stack) Draw(s tcell.Screen) {
	for i := 0; i <= st.Top; i++ {
		st.Layers[i].Draw(s)
	}
}

func (st *Stack) Layout(width int, height int) {
	for i := 0; i <= st.Top; i++ {
		st.Layers[i].Layout(width, height)
	}
}

func (st *Stack) MinimumSize(width int, height int) (Width int, Height int) {
	w := -1
	h := -1

	for i := 0; i <= st.Top; i++ {
		mw, mh := st.Layers[i].MinimumSize(width, height)

		if mw > w {
			w = mw
		}
		if mh > h {
			h = mh
		}
	}
	return w, h
}
