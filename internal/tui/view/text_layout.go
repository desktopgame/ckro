package view

import "github.com/desktopgame/ckro/internal/tui/model"

// TextLayout is representation a layout at a certain point in time.
type TextLayout struct {
	Element  model.Element
	Children []*TextLayout

	RelativeX int
	RelativeY int

	Width         int
	Height        int
	MinimumWidth  int
	MinimumHeight int
	WidthTable    []int
}
