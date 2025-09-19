package base

// Focusable has distinction of focus on and off, and capture a Event.
type Focusable interface {
	Handle(ev Event)
	Focus(on bool)
}

// FocusableTree is manage focus to sub controls.
// if implemented this interface, parent controls stealing event, but movable focus to each children.
// parent controls can be delegate event to children as needed.
type FocusableTree interface {
	Focusable

	// SubFocusFirst is focus to first children.
	SubFocusFirst()

	// SubFocusPrev is focus to previous children.
	SubFocusPrev() bool

	// SubFocusNext is focus to next children.
	SubFocusNext() bool

	// SubFocusLast is focus to last children.
	SubFocusLast()
}

// FocusManager is manage focus order in controls.
type FocusManager struct {
	tiles  []Focusable
	active int
	tree   FocusableTree
}

// Init is initialize FocusManager.
func (fm *FocusManager) Init() {
	fm.tiles = []Focusable{}
	fm.active = 0
}

// Register is Focusable entry add  to table.
func (fm *FocusManager) Register(t Focusable) {
	fm.tiles = append(fm.tiles, t)
	t.Focus(false)
}

// Traverse is add all controls to table by the order of focus.
func (fm *FocusManager) Traverse(ctrl Control) {
	fm.Init()
	ctrl.Traverse(fm)
}

// Grab is try to activate focus on current object.
func (fm *FocusManager) Grab() {
	if len(fm.tiles) == 0 {
		return
	}

	fm.tiles[fm.active].Focus(true)

	if tree, ok := fm.tiles[fm.active].(FocusableTree); ok {
		tree.SubFocusFirst()
		fm.tree = tree
	}
}

// FocusPrev is move focus to previous.
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

	fm.tiles[fm.active].Focus(false)

	if fm.active > 0 {
		fm.active--
	} else {
		fm.active = len(fm.tiles) - 1
	}

	fm.tiles[fm.active].Focus(true)

	if tree, ok := fm.tiles[fm.active].(FocusableTree); ok {
		tree.SubFocusLast()
		fm.tree = tree
	}
}

// FocusNext is move focus to next.
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

	fm.tiles[fm.active].Focus(false)

	if fm.active < len(fm.tiles)-1 {
		fm.active++
	} else {
		fm.active = 0
	}

	fm.tiles[fm.active].Focus(true)

	if tree, ok := fm.tiles[fm.active].(FocusableTree); ok {
		tree.SubFocusFirst()
		fm.tree = tree
	}
}

// Handle is process the event by focused control.
func (fm *FocusManager) Handle(ev Event) {
	if len(fm.tiles) == 0 {
		return
	}

	if fm.active >= 0 {
		fm.tiles[fm.active].Handle(ev)
	}
}
