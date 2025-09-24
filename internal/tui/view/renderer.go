package view

import "github.com/gdamore/tcell/v2"

type Renderer interface {
	SetContent(x int, y int, primary rune, combining []rune, style tcell.Style)

	Translate(offsetX int, offsetY int) Renderer
}
