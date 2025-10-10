package view

import (
	"github.com/gdamore/tcell/v2"
)

// Renderer is terminal graphics backend.
// only tcell supported in currently.
type Renderer interface {
	// SetContent put rune at specified position.
	SetContent(x int, y int, primary rune, combining []rune, style tcell.Style)

	// Translate returns renderer of after coordinate translate.
	// region will be reseted.
	Translate(offsetX int, offsetY int) Renderer

	// Translate returns renderer of after clip.
	Region(width int, height int) Renderer
}
