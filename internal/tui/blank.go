package tui

type Blank struct {
}

func (b Blank) Traverse(fm *FocusManager) {
}

func (b Blank) Update() {
}

func (b Blank) Draw(g *Graphics) {
}

func (b Blank) MinimumSize(width int, height int) (Width int, Height int) {
	return 0, 0
}

func (b Blank) Move(x int, y int) {
}

func (b Blank) Layout(w int, h int) {
}

func (b Blank) IsFlexibleWidth() bool {
	return false
}

func (b Blank) IsFlexibleHeight() bool {
	return false
}
