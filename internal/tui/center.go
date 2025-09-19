package tui

// Center is Control for centering.
type Center struct {
	Control Control
	Width   int
	Height  int
	x       int
	y       int
}

// Traverse is delegate to sub control.
func (c *Center) Traverse(fm *FocusManager) {
	c.Control.Traverse(fm)
}

// Update is delegate to sub control.
func (c *Center) Update() {
	c.Control.Update()
}

// Draw is delegate to sub control.
func (c *Center) Draw(g *Graphics) {
	c.Control.Draw(g)
}

// MinimumSize is provide Control minimum size with margin.
func (c *Center) MinimumSize(width int, height int) (Width int, Height int) {
	mw, mh := c.Control.MinimumSize(width, height)
	return max(mw, c.Width), max(mh, c.Height)
}

// Move is delegate to sub control.
func (c *Center) Move(x int, y int) {
	c.x = x
	c.y = y
}

// Layout is align center sub control to parent.
func (c *Center) Layout(width int, height int) {
	c.Control.Move(
		c.x+(width-c.Width)/2,
		c.y+(height-c.Height)/2,
	)
	c.Control.Layout(c.Width, c.Height)
}

// IsFlexibleWidth returns false.
func (c *Center) IsFlexibleWidth() bool {
	return false
}

// IsFlexibleHeight returns false.
func (c *Center) IsFlexibleHeight() bool {
	return false
}
