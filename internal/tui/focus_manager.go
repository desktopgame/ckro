package tui

import "github.com/gdamore/tcell/v2"

type FocusManager struct {
	tiles  []*Tile
	active int
}

func (fm *FocusManager) Init() {
	fm.tiles = []*Tile{}
	fm.active = 0
}

func (fm *FocusManager) Register(t *Tile) {
	fm.tiles = append(fm.tiles, t)
	t.TextBox.ShowCursor = false
}

func (fm *FocusManager) Traverse(ctrl Control) {
	fm.Init()
	ctrl.Traverse(fm)
}

func (fm *FocusManager) FocusPrev() {
	if len(fm.tiles) == 0 {
		return
	}

	fm.tiles[fm.active].TextBox.ShowCursor = false

	if fm.active > 0 {
		fm.active--
	} else {
		fm.active = len(fm.tiles) - 1
	}

	showCursor := fm.tiles[fm.active].TextPresenter.ShowCursor()
	fm.tiles[fm.active].TextBox.ShowCursor = showCursor
}

func (fm *FocusManager) FocusNext() {
	if len(fm.tiles) == 0 {
		return
	}

	fm.tiles[fm.active].TextBox.ShowCursor = false

	if fm.active < len(fm.tiles)-1 {
		fm.active++
	} else {
		fm.active = 0
	}

	showCursor := fm.tiles[fm.active].TextPresenter.ShowCursor()
	fm.tiles[fm.active].TextBox.ShowCursor = showCursor
}

func (fm *FocusManager) Handle(ev tcell.Event) {
	if len(fm.tiles) == 0 {
		return
	}

	if fm.active >= 0 {
		fm.tiles[fm.active].Handle(ev)
	}
}
