package base

type Focusable interface {
	Handle(ev Event)
	Focus(on bool)
}

type FocusableTree interface {
	Focusable

	SubFocusFirst()
	SubFocusPrev() bool
	SubFocusNext() bool
	SubFocusLast()
}

type FocusManager struct {
	tiles  []Focusable
	active int
	tree   FocusableTree
}

func (fm *FocusManager) Init() {
	fm.tiles = []Focusable{}
	fm.active = 0
}

func (fm *FocusManager) Register(t Focusable) {
	fm.tiles = append(fm.tiles, t)
	// t.GetTextBox().ShowCursor = false
	t.Focus(false)
}

func (fm *FocusManager) Traverse(ctrl Control) {
	fm.Init()
	ctrl.Traverse(fm)
}

func (fm *FocusManager) Grab() {
	if len(fm.tiles) == 0 {
		return
	}

	// showCursor := fm.tiles[fm.active].GetTextPresenter().ShowCursor()
	// fm.tiles[fm.active].GetTextBox().ShowCursor = showCursor
	fm.tiles[fm.active].Focus(true)
}

func (fm *FocusManager) FocusPrev() {
	if len(fm.tiles) == 0 {
		return
	}

	if fm.tree != nil {
		if !fm.tree.SubFocusPrev() {
			fm.tree = nil
		}
		return
	}

	// fm.tiles[fm.active].GetTextBox().ShowCursor = false
	fm.tiles[fm.active].Focus(false)

	if fm.active > 0 {
		fm.active--
	} else {
		fm.active = len(fm.tiles) - 1
	}

	// showCursor := fm.tiles[fm.active].GetTextPresenter().ShowCursor()
	// fm.tiles[fm.active].GetTextBox().ShowCursor = showCursor
	fm.tiles[fm.active].Focus(true)

	if tree, ok := fm.tiles[fm.active].(FocusableTree); ok {
		tree.SubFocusLast()
		fm.tree = tree
	}
}

func (fm *FocusManager) FocusNext() {
	if len(fm.tiles) == 0 {
		return
	}

	if fm.tree != nil {
		if !fm.tree.SubFocusNext() {
			fm.tree = nil
		}
		return
	}

	// fm.tiles[fm.active].GetTextBox().ShowCursor = false
	fm.tiles[fm.active].Focus(false)

	if fm.active < len(fm.tiles)-1 {
		fm.active++
	} else {
		fm.active = 0
	}

	// showCursor := fm.tiles[fm.active].GetTextPresenter().ShowCursor()
	// fm.tiles[fm.active].GetTextBox().ShowCursor = showCursor
	fm.tiles[fm.active].Focus(true)

	if tree, ok := fm.tiles[fm.active].(FocusableTree); ok {
		tree.SubFocusFirst()
		fm.tree = tree
	}
}

func (fm *FocusManager) Handle(ev Event) {
	if len(fm.tiles) == 0 {
		return
	}

	if fm.active >= 0 {
		fm.tiles[fm.active].Handle(ev)
	}
}
