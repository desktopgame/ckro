package tui

import "github.com/desktopgame/ckro/internal/tui/model"

type TextView interface {
	Layout(e model.Element, width int) *TextLayout
	Draw(e model.Element, textLayout TextLayout, g *Graphics)

	Width(e model.Element, textLayout TextLayout, row int) int
	Height(e model.Element, textLayout TextLayout) int
}
