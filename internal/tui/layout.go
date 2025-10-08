package tui

import "github.com/desktopgame/ckro/internal/tui/presenter"

// WithFrame returns Frame as wrapper of specified Control.
func WithFrame(ctrl Control) *Frame {
	fr := Frame{}
	fr.Control = ctrl
	return &fr
}

// WithCenter returns Center as wrapper of specified Control.
func WithCenter(ctrl Control, width int, height int) *Center {
	c := Center{}
	c.Control = ctrl
	c.Width = width
	c.Height = height
	return &c
}

// Tile creation utilities
func NewTile(presenter TextPresenter) *Tile {
	tile := &Tile{}
	tile.Init()
	tile.TextPresenter = presenter
	return tile
}

func NewFlexTile(presenter TextPresenter) *Tile {
	tile := NewTile(presenter)
	tile.FlexibleWidth = true
	tile.FlexibleHeight = true
	return tile
}

func NewFixedTile(presenter TextPresenter, width, height int) *Tile {
	tile := NewTile(presenter)
	tile.MinimumWidth = width
	tile.MinimumHeight = height
	return tile
}

func NewEditTile() *Tile {
	return NewTile(&presenter.EditTextPresenter{})
}

func NewLabelTile(text string) *Tile {
	return NewTile(&presenter.LabelTextPresenter{Text: text})
}

func NewCenteredLabelTile(text string) *Tile {
	return NewTile(&presenter.LabelTextPresenter{Text: text, AlignCenter: true})
}

func NewListTile(items []string) *Tile {
	return NewTile(&presenter.ListTextPresenter{Items: items})
}

func NewCustomListTile(items []string, cursorChar rune, prefix string) *Tile {
	return NewTile(&presenter.ListTextPresenter{
		Items:      items,
		CursorChar: cursorChar,
		Prefix:     prefix,
	})
}

// Separator utilities
func NewVerticalSeparator() *Tile {
	tile := NewTile(&presenter.VerticalSeparatorTextPresenter{})
	tile.MinimumWidth = 1
	tile.FlexibleHeight = true
	return tile
}

func NewHorizontalSeparator() *Tile {
	tile := NewTile(&presenter.HorizontalSeparatorTextPresenter{})
	tile.MinimumHeight = 1
	tile.FlexibleWidth = true
	return tile
}

// Box creation utilities
func NewHBox(controls ...Control) *Box {
	box := &Box{}
	box.Init(Horizontal)
	for _, ctrl := range controls {
		box.Controls = append(box.Controls, ctrl)
	}
	return box
}

func NewVBox(controls ...Control) *Box {
	box := &Box{}
	box.Init(Vertical)
	for _, ctrl := range controls {
		box.Controls = append(box.Controls, ctrl)
	}
	return box
}
