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

// Stack utilities
func NewStack(controls ...Control) *Stack {
	stack := &Stack{}
	stack.Init()
	for _, ctrl := range controls {
		stack.Layers = append(stack.Layers, ctrl)
	}
	if len(controls) > 0 {
		stack.Top = len(controls) - 1
	}
	return stack
}

// Quick layout builders
func QuickEditor() (*Tile, *Box) {
	textArea, _, editorBox := NewTextEditor()
	return textArea, editorBox
}

func QuickForm(labelWidth int, pairs ...struct{ Label, Input string }) *Grid {
	grid := NewGrid(len(pairs), 2)

	for i, pair := range pairs {
		// ラベル（固定サイズ）
		label := NewFixedTile(&presenter.LabelTextPresenter{Text: pair.Label, AlignCenter: true}, labelWidth, 5)
		grid.SetControl(i, 0, label)

		// 入力フィールド（フレーム付き、固定高さ）
		input := NewEditTile()
		input.FlexibleWidth = true
		input.FlexibleHeight = false
		input.MinimumHeight = 3
		input.TextBox.ShowCursor = true
		framedInput := WithFrame(input)
		grid.SetControl(i, 1, framedInput)
	}

	return grid
}

func QuickDialog(title string, content Control, buttonLabels ...string) *Box {
	titleTile := NewCenteredLabelTile(title)
	titleTile.FlexibleWidth = true
	titleTile.MinimumHeight = 1

	buttons := NewHBox()
	for _, label := range buttonLabels {
		btn := NewFixedTile(&presenter.LabelTextPresenter{
			Text:        label,
			AlignCenter: true,
		}, len(label)+4, 3)
		buttons.Controls = append(buttons.Controls, btn)
		if len(buttons.Controls) > 1 {
			// Add spacing between buttons
			buttons.Controls = append(buttons.Controls[:len(buttons.Controls)-1],
				NewFixedTile(&presenter.FrameTextPresenter{}, 2, 1),
				buttons.Controls[len(buttons.Controls)-1])
		}
	}

	return NewVBox(
		titleTile,
		NewHorizontalSeparator(),
		content,
		NewHorizontalSeparator(),
		buttons,
	)
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

func HSplit(controls ...Control) *Box {
	return WithSeparators(Horizontal, controls...)
}

func VSplit(controls ...Control) *Box {
	return WithSeparators(Vertical, controls...)
}

// Command Palette utility
func QuickCommandPalette(commands []string, onExecute func(string)) *CommandPalette {
	palette := &CommandPalette{}
	palette.Init(commands, onExecute)
	return palette
}
