package tui

// Stack is stack of layer.
type Stack struct {
	Layers []Layer
	Top    int
}

// Init is initialize Stack.
func (st *Stack) Init() {
	st.Layers = []Layer{}
}

// Traverse is delegate to top layer.
func (st *Stack) Traverse(fm *FocusManager) {
	fm.Traverse(st.Layers[st.Top].Control)
}

// Update is update all layers or update top layer only, if top layer's ClearBackground is true.
func (st *Stack) Update() {
	if st.Top > 0 && st.Layers[st.Top].ClearBackground {
		st.Layers[st.Top].Control.Update()
	} else {
		for i := 0; i <= st.Top; i++ {
			st.Layers[i].Control.Update()
		}
	}
}

// Draw is draw all layers or draw top layer only, if top layer's ClearBackground is true.
func (st *Stack) Draw(g *Graphics) {
	if st.Top > 0 && st.Layers[st.Top].ClearBackground {
		st.Layers[st.Top].Control.Draw(g)
	} else {
		for i := st.Top; i >= 0; i-- {
			st.Layers[i].Control.Draw(g)
		}
	}
}

// Layout is layout all layers.
func (st *Stack) Layout(width int, height int) {
	for i := 0; i <= st.Top; i++ {
		st.Layers[i].Control.Layout(width, height)
	}
}

// MinimumSize returns maximum of minimum size of each layers.
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
