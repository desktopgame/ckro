package view

import "github.com/desktopgame/ckro/internal/tui/model"

type TextView interface {
	Layout(e model.Element, width int) *TextLayout
	Draw(textLayout *TextLayout, renderer Renderer, x int, y int)

	Width(textLayout *TextLayout, row int) int
	Height(textLayout *TextLayout) int
}
