package main

import (
	"github.com/desktopgame/ckro/internal/tui"
	"github.com/desktopgame/ckro/internal/tui/presenter"
)

type ModeLine struct {
	tile tui.Tile
}

func (m *ModeLine) Init() {
	m.tile = tui.Tile{}
	m.tile.Init()
	m.tile.FlexibleWidth = true
	m.tile.MinimumHeight = 1
	m.tile.TextPresenter = &presenter.LabelTextPresenter{}
}

func (m *ModeLine) Text(text string) {
	if p, ok := m.tile.TextPresenter.(*presenter.LabelTextPresenter); ok {
		p.Text = text
	}
}

func (m *ModeLine) Traverse(fm *tui.FocusManager) {
	m.tile.Traverse(fm)
}

func (m *ModeLine) Update() {
	m.tile.Update()
}

func (m *ModeLine) Draw(g *tui.Graphics) {
	m.tile.Draw(g)
}

func (m *ModeLine) MinimumSize(width int, height int) (Width int, Height int) {
	return m.tile.MinimumSize(width, height)
}

func (m *ModeLine) Move(x int, y int) {
	m.tile.Move(x, y)
}

func (m *ModeLine) Layout(w int, h int) {
	m.tile.Layout(w, h)
}

func (m *ModeLine) IsFlexibleWidth() bool {
	return m.tile.IsFlexibleWidth()
}

func (m *ModeLine) IsFlexibleHeight() bool {
	return m.tile.IsFlexibleHeight()
}
