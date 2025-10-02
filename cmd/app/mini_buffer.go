package main

import (
	"github.com/desktopgame/ckro/internal/tui"
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/desktopgame/ckro/internal/tui/presenter"
	"github.com/gdamore/tcell/v2"
)

type MiniBuffer struct {
	tile tui.Tile
	ch   chan string
}

func (m *MiniBuffer) Init() {
	m.Close()
	m.tile = tui.Tile{}
	m.tile.Init()
	m.tile.FlexibleWidth = true
	m.tile.MinimumHeight = 1
	m.tile.TextPresenter = &presenter.EditTextPresenter{}
	m.ch = make(chan string)
}

func (m *MiniBuffer) onSubmit(s string) {
	m.ch <- s
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

func (m *MiniBuffer) Close() {
	if m.ch != nil {
		close(m.ch)
		m.ch = nil
	}
}

func (m *MiniBuffer) Handle(ev tui.Event) {
	switch e := ev.GetSource().(type) {
	case *tcell.EventKey:
		switch e.Key() {
		case tcell.KeyEnter:
			tb := m.tile.TextBox
			sg := tb.Document.Read(model.Range{
				StartPosition: model.Position{
					Row:    0,
					Column: 0,
				},
				EndPosition: model.Position{
					Row:    0,
					Column: tb.Document.GetLineBytes(tb.Document.GetLineCount() - 1),
				},
			})
			m.onSubmit(sg.GetLine(0))
			m.tile.TextBox.Document.Clear()
			m.tile.TextBox.CursorReset()
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
