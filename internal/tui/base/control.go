package base

// Control is graphical unit for compose a screen.
// Control can be contain another Control, but no distinction to interface by the see outer.
type Control interface {
	// Traverse is register sub controls to FocusManager by the order by focus.
	Traverse(fm *FocusManager)

	// Update is prepare for draw to terminal.
	Update()

	// Draw is rendering Control as text, to terminal.
	Draw(g *Graphics)

	// MinimumSize is provide size of needs to print Control.
	MinimumSize(width int, height int) (Width int, Height int)

	// Move is move a Control to specified position.
	Move(x int, y int)

	// Layout is relayout sub controls within specified size.
	Layout(w int, h int)

	// IsFlexibleWidth is returns true if Control is flexible on horizontal.
	IsFlexibleWidth() bool

	// IsFlexibleHeight is returns true if Control is flexible on vertical.
	IsFlexibleHeight() bool
}
