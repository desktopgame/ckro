package tui

// Blank is most simply implement of Control.
// Blank is provide empty region, and not draw content, and not update.
type Blank struct {
}

// Traverse is do nothinng.
func (b Blank) Traverse(fm *FocusManager) {
}

// Update is do nothinng.
func (b Blank) Update() {
}

// Draw is do nothinng.
func (b Blank) Draw(g *Graphics) {
}

// MinimumSize returns zero.
func (b Blank) MinimumSize(width int, height int) (Width int, Height int) {
	return 0, 0
}

// Move is do nothinng.
func (b Blank) Move(x int, y int) {
}

// Layout is do nothinng.
func (b Blank) Layout(w int, h int) {
}

// MinimumSize returns false.
func (b Blank) IsFlexibleWidth() bool {
	return false
}

// MinimumSize returns false.
func (b Blank) IsFlexibleHeight() bool {
	return false
}
