package tui

import "github.com/desktopgame/ckro/internal/tui/model"

type TextLayout struct {
	Element  model.Element
	Children []*TextLayout

	Width  int
	Height int
}
