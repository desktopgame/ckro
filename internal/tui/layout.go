package tui

import "github.com/desktopgame/ckro/internal/tui/presenter"

// Frame wrapper
func WithFrame(ctrl Control) *Frame {
	fr := Frame{}
	fr.Control = ctrl
	return &fr
}

// Center wrapper
func WithCenter(ctrl Control, width int, height int) *Center {
	c := Center{}
	c.Control = ctrl
	c.Width = width
	c.Height = height
	return &c
}

// Layout composition helpers
func WithSeparators(orientation Orientation, controls ...Control) *Box {
	if len(controls) == 0 {
		return NewHBox()
	}

	var box *Box
	var separator *Tile

	if orientation == Horizontal {
		box = NewHBox()
		separator = NewVerticalSeparator()
	} else {
		box = NewVBox()
		separator = NewHorizontalSeparator()
	}

	for i, ctrl := range controls {
		if i > 0 {
			// Clone separator for each use
			sep := NewTile(separator.TextPresenter)
			sep.MinimumWidth = separator.MinimumWidth
			sep.MinimumHeight = separator.MinimumHeight
			sep.FlexibleWidth = separator.FlexibleWidth
			sep.FlexibleHeight = separator.FlexibleHeight
			box.Controls = append(box.Controls, sep)
		}
		box.Controls = append(box.Controls, ctrl)
	}

	return box
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

// Common layout patterns
func NewTextEditor() (*Tile, *Tile, *Box) {
	textArea := NewEditTile()
	textArea.MinimumWidth = 3
	textArea.FlexibleWidth = true
	textArea.FlexibleHeight = true
	textArea.TextBox.ShowCursor = true

	lineNumbers := NewTile(&presenter.LineNumberTextPresenter{
		TargetView: textArea.TextBox,
	})
	lineNumbers.MinimumWidth = 4
	lineNumbers.FlexibleHeight = true

	scrollBar := NewTile(&presenter.ScrollBarTextPresenter{
		TargetView: textArea.TextBox,
	})
	scrollBar.MinimumWidth = 1
	scrollBar.FlexibleHeight = true

	editorBox := NewHBox(
		lineNumbers,
		NewVerticalSeparator(),
		textArea,
		scrollBar,
	)

	return textArea, lineNumbers, editorBox
}

func NewFileTree(rootDir string, onFileOpen func(string)) *Tile {
	tree := NewTile(&presenter.TreeTextPresenter{
		RootDirectory: rootDir,
		OnFileOpen:    onFileOpen,
	})
	tree.MinimumWidth = 50
	tree.FlexibleHeight = true
	return tree
}

// Grid utilities
func NewGrid(rows, cols int) *Grid {
	grid := &Grid{}
	grid.Init(rows, cols)
	return grid
}
