package tui

import (
	"github.com/desktopgame/ckro/internal/tui/presenter"
	"github.com/gdamore/tcell/v2"
)

type Tile struct {
	TextBox       *TextBox
	TextPresenter TextPresenter

	MinimumWidth   int
	MinimumHeight  int
	FlexibleWidth  bool
	FlexibleHeight bool
}

func (t *Tile) Init() {
	t.TextBox = &TextBox{}
	t.TextBox.Init()
	t.TextPresenter = &presenter.LabelTextPresenter{}
}

func (t *Tile) Update() {
	t.TextPresenter.Present(t.TextBox)
}

func (t *Tile) Handle(ev tcell.Event) {
	t.TextPresenter.Handle(t.TextBox, ev)
}

func (t *Tile) Draw(s tcell.Screen) {
	t.TextBox.Draw(s)
}

func (t *Tile) MinimumSize(width int, height int) (Width int, Height int) {
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
	return t.FlexibleHeight
}
