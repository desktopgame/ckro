package tui

// Card is delegate to the which one of sub controls.
// uses to switch content like tab pane.
type Card struct {
	Controls []Control
	Index    int
}

// Init is initialize Card.
func (c *Card) Init(controls []Control) {
	c.Controls = make([]Control, len(controls))
	copy(c.Controls, controls)
	c.Index = 0
}

// Target returns current delegate target.
func (c *Card) Target() Control {
	return c.Controls[c.Index]
}

// Prev is show previous sub control.
func (c *Card) Prev() {
	c.Index--
	if c.Index < 0 {
		c.Index = len(c.Controls) - 1
	}
}

// Next is show next sub control.
func (c *Card) Next() {
	c.Index++
	if c.Index >= len(c.Controls) {
		c.Index = 0
	}
}

// Traverse is delegate to current delegate target
func (c *Card) Traverse(fm *FocusManager) {
	c.Target().Traverse(fm)
}

// Update is delegate to current delegate target
func (c *Card) Update() {
	c.Target().Update()
}

// Draw is delegate to current delegate target
func (c *Card) Draw(g *Graphics) {
	c.Target().Draw(g)
}

// MinimumSize is delegate to current delegate target
func (c *Card) MinimumSize(width int, height int) (Width int, Height int) {
	return c.Target().MinimumSize(width, height)
}

// Move is delegate to current delegate target
func (c *Card) Move(x int, y int) {
	c.Target().Move(x, y)
}

// Layout is delegate to current delegate target
func (c *Card) Layout(w int, h int) {
	c.Target().Layout(w, h)
}

// IsFlexibleWidth is delegate to current delegate target
func (c *Card) IsFlexibleWidth() bool {
	return c.Target().IsFlexibleWidth()
}

// IsFlexibleHeight is delegate to current delegate target
func (c *Card) IsFlexibleHeight() bool {
	return c.Target().IsFlexibleHeight()
}
