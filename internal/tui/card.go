package tui

type Card struct {
	Controls []Control
	Index    int
}

func (c *Card) Init(controls []Control) {
	c.Controls = make([]Control, len(controls))
	copy(c.Controls, controls)
	c.Index = 0
}

func (c *Card) Target() Control {
	return c.Controls[c.Index]
}

func (c *Card) Prev() {
	c.Index--
	if c.Index < 0 {
		c.Index = len(c.Controls) - 1
	}
}

func (c *Card) Next() {
	c.Index++
	if c.Index >= len(c.Controls) {
		c.Index = 0
	}
}

func (c *Card) Traverse(fm *FocusManager) {
	c.Target().Traverse(fm)
}

func (c *Card) Update() {
	c.Target().Update()
}

func (c *Card) Draw(g *Graphics) {
	c.Target().Draw(g)
}

func (c *Card) MinimumSize(width int, height int) (Width int, Height int) {
	return c.Target().MinimumSize(width, height)
}

func (c *Card) Move(x int, y int) {
	c.Target().Move(x, y)
}

func (c *Card) Layout(w int, h int) {
	c.Target().Layout(w, h)
}

func (c *Card) IsFlexibleWidth() bool {
	return c.Target().IsFlexibleWidth()
}

func (c *Card) IsFlexibleHeight() bool {
	return c.Target().IsFlexibleHeight()
}
