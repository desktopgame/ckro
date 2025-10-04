package tui

import (
	"github.com/desktopgame/ckro/internal/tui/presenter"
)

// Tile is implement of Control for TextBox.
type Tile struct {
	TextBox       *TextBox
	TextPresenter TextPresenter

	MinimumWidth   int
	MinimumHeight  int
	FlexibleWidth  bool
	FlexibleHeight bool
}

// Init is initialize Tile.
func (t *Tile) Init() {
	t.TextBox = &TextBox{}
	t.TextBox.Init()
	t.TextPresenter = &presenter.LabelTextPresenter{}
}

// Handle is delegate to TextPresenter.
func (t *Tile) Handle(ev Event) {
	t.TextPresenter.Handle(t.TextBox, ev.GetSource())
}

// Focus is switch cursor visibility as needed.
func (t *Tile) Focus(on bool) {
	if on {
		t.TextBox.ShowCursor = t.TextPresenter.ShowCursor()
	} else {
		t.TextBox.ShowCursor = false
	}
}

// Traverse is register TextBox to FocusManager.
func (t *Tile) Traverse(fm *FocusManager) {
	if t.TextPresenter.IsFocusable() {
		fm.Register(t)
	}
}

// Update is delegate to TextPresenter.
func (t *Tile) Update() {
	t.TextPresenter.Present(t.TextBox)
}

// Draw is delegate to TextBox.
func (t *Tile) Draw(g *Graphics) {
	t.TextBox.Draw(g)
}

// MinimumSize returns Tile fields to direct.
func (t *Tile) MinimumSize(width int, height int) (Width int, Height int) {
	return t.MinimumWidth, t.MinimumHeight
}

// Move is move the TextBox.
func (t *Tile) Move(x int, y int) {
	t.TextBox.X = x
	t.TextBox.Y = y
}

// Layout is resize the TextBox.
func (t *Tile) Layout(w int, h int) {
	if t.TextBox.Width != w || t.TextBox.Height != h {
		t.TextBox.Width = w
		t.TextBox.Height = h
		t.TextBox.CursorUpdate()
	}
}

// IsFlexibleWidth returns Tile fields to direct.
func (t *Tile) IsFlexibleWidth() bool {
	return t.FlexibleWidth
}

// IsFlexibleHeight returns Tile fields to direct.
func (t *Tile) IsFlexibleHeight() bool {
	return t.FlexibleHeight
}
