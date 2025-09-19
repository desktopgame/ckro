package tui

import "github.com/gdamore/tcell/v2"

// Clip is clipping specified region in Graphics.
type Clip struct {
	Graphics *Graphics
	X        int
	Y        int
	Width    int
	Height   int
}

// SetContent is set a character to specified cell, if not already settled.
// TODO: refactor
func (c Clip) SetContent(x int, y int, primary rune, combining []rune, style tcell.Style) {
	c.Graphics.Draw(c.X+x, c.Y+y, primary, combining, style)
}

// SetCursor is set a character to specified cell.
// TODO: refactor
func (c Clip) SetCursor(x int, y int, primary rune, combining []rune, style tcell.Style) {
	c.Graphics.ForceDraw(c.X+x, c.Y+y, primary, combining, style)
}
