package main

import (
	"github.com/desktopgame/ckro/internal/tui"
	"github.com/desktopgame/ckro/internal/tui/presenter"
)

type MiniBuffer struct {
	tile tui.Tile
}

func (m *MiniBuffer) Init() {
	m.tile = tui.Tile{}
	m.tile.Init()
	m.tile.FlexibleWidth = true
	m.tile.MinimumHeight = 1
	m.tile.TextPresenter = &presenter.EditTextPresenter{}
}

func (m *MiniBuffer) Traverse(fm *tui.FocusManager) {
	m.tile.Traverse(fm)
}

func (m *MiniBuffer) Update() {
	m.tile.Update()
}

func (m *MiniBuffer) Draw(g *tui.Graphics) {
	m.tile.Draw(g)
}

func (m *MiniBuffer) MinimumSize(width int, height int) (Width int, Height int) {
	return m.tile.MinimumSize(width, height)
}

func (m *MiniBuffer) Move(x int, y int) {
	m.tile.Move(x, y)
}

func (m *MiniBuffer) Layout(w int, h int) {
	m.tile.Layout(w, h)
}

func (m *MiniBuffer) IsFlexibleWidth() bool {
	return m.tile.IsFlexibleWidth()
}

func (m *MiniBuffer) IsFlexibleHeight() bool {
	return m.tile.IsFlexibleHeight()
}
