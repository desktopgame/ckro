package main

import (
	"github.com/desktopgame/ckro/internal/tui"
	"github.com/desktopgame/ckro/internal/tui/presenter"
	"github.com/gdamore/tcell/v2"
)

type MiniBuffer struct {
	tile     tui.Tile
	OnSubmit func(string)
}

func (m *MiniBuffer) Init(onSubmit func(string)) {
	m.tile = tui.Tile{}
	m.tile.Init()
	m.tile.FlexibleWidth = true
	m.tile.MinimumHeight = 1
	m.tile.TextPresenter = &presenter.EditTextPresenter{}
	m.OnSubmit = onSubmit
}

func (m *MiniBuffer) SetEditable(edidtable bool) {
	if edit, ok := m.tile.TextPresenter.(*presenter.EditTextPresenter); ok {
		edit.ReadOnly = !edidtable
	}
}

func (m *MiniBuffer) Editable() {
	m.SetEditable(true)
}

func (m *MiniBuffer) ReadOnly() {
	m.SetEditable(false)
}

func (m *MiniBuffer) Handle(ev tui.Event) {
	switch e := ev.GetSource().(type) {
	case *tcell.EventKey:
		switch e.Key() {
		case tcell.KeyEnter:
			if m.OnSubmit != nil {
				m.OnSubmit(m.tile.TextBox.GetDocument().GetBuffer().GetLineAt(0).GetContent())
				m.tile.TextBox.Document.Init()
				m.tile.TextBox.CursorReset()
			}
			return
		}
	}
	m.tile.Handle(ev)
}

func (m *MiniBuffer) Focus(on bool) {
	m.tile.Focus(on)
}

func (m *MiniBuffer) Traverse(fm *tui.FocusManager) {
	fm.Register(m)
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
