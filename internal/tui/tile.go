package tui

type Tile struct {
	TextBox *TextBox

	MinimumWidth    int
	MinimumHeight   int
	FlexibleWidth   bool
	FlexibleHeigght bool
}

func (t *Tile) Init() {
	t.TextBox = &TextBox{}
	t.TextBox.Init()
}

func (t *Tile) MinimumSize() (Width int, Height int) {
	return t.MinimumWidth, t.MinimumHeight
}

func (t *Tile) Move(x int, y int) {
	t.TextBox.X = x
	t.TextBox.Y = y
}

func (t *Tile) Layout(w int, h int) {
	t.TextBox.Width = w
	t.TextBox.Height = h
}

func (t *Tile) IsFlexibleWidth() bool {
	return t.FlexibleWidth
}

func (t *Tile) IsFlexibleHeight() bool {
	return t.FlexibleHeigght
}
