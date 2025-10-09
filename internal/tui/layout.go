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

// NewTile returns Tile, returned Tile is using specified TextPresenter.
func NewTile(presenter TextPresenter) *Tile {
	tile := &Tile{}
	tile.Init()
	tile.TextPresenter = presenter
	return tile
}

// NewTile returns flex Tile, returned Tile is using specified TextPresenter.
func NewFlexTile(presenter TextPresenter) *Tile {
	tile := NewTile(presenter)
	tile.FlexibleWidth = true
	tile.FlexibleHeight = true
	return tile
}

// NewTile returns fixed Tile, returned Tile is using specified TextPresenter.
func NewFixedTile(presenter TextPresenter, width, height int) *Tile {
	tile := NewTile(presenter)
	tile.MinimumWidth = width
	tile.MinimumHeight = height
	return tile
}

// NewEditTile returns editable TextBox as Tile.
func NewEditTile() *Tile {
	return NewTile(&presenter.EditTextPresenter{})
}

// NewLabelTile returns label text as Tile.
func NewLabelTile(text string) *Tile {
	return NewTile(&presenter.LabelTextPresenter{Text: text})
}

// NewLabelTile returns label text as Tile, and aligned center.
func NewCenteredLabelTile(text string) *Tile {
	return NewTile(&presenter.LabelTextPresenter{Text: text, AlignCenter: true})
}

// NewListTile returns list as Tile.
func NewListTile(items []string) *Tile {
	return NewTile(&presenter.ListTextPresenter{Items: items})
}

// NewListTile returns list as Tile.
func NewCustomListTile(items []string, cursorChar rune, prefix string) *Tile {
	return NewTile(&presenter.ListTextPresenter{
		Items:      items,
		CursorChar: cursorChar,
		Prefix:     prefix,
	})
}

// NewVerticalSeparator returns vertical separator.
func NewVerticalSeparator() *Tile {
	tile := NewTile(&presenter.VerticalSeparatorTextPresenter{})
	tile.MinimumWidth = 1
	tile.FlexibleHeight = true
	return tile
}

// NewHorizontalSeparator returns horizontal separator.
func NewHorizontalSeparator() *Tile {
	tile := NewTile(&presenter.HorizontalSeparatorTextPresenter{})
	tile.MinimumHeight = 1
	tile.FlexibleWidth = true
	return tile
}

// NewHBox returns horizontal box.
func NewHBox(controls ...Control) *Box {
	box := &Box{}
	box.Init(Horizontal)
	for _, ctrl := range controls {
		box.Controls = append(box.Controls, ctrl)
	}
	return box
}

// NewVBox returns vertical box.
func NewVBox(controls ...Control) *Box {
	box := &Box{}
	box.Init(Vertical)
	for _, ctrl := range controls {
		box.Controls = append(box.Controls, ctrl)
	}
	return box
}
