package tui

import (
	"github.com/desktopgame/ckro/internal/tui/view"
	"github.com/gdamore/tcell/v2"
)

// Clip is clipping specified region in Graphics.
type Clip struct {
	Graphics   *Graphics
	X          int
	Y          int
	Width      int
	Height     int
	FirstLineY int
}

// SetContent is set a character to specified cell, if not already settled.
// TODO: refactor
func (c *Clip) SetContent(x int, y int, primary rune, combining []rune, style tcell.Style) {
	if y < c.FirstLineY || y >= c.FirstLineY+c.Height {
		return
	}
	offset := y - c.FirstLineY
	c.Graphics.Draw(c.X+x, c.Y+offset, primary, combining, style)
}

// SetCursor is set a character to specified cell.
// TODO: refactor
func (c *Clip) SetCursor(x int, y int, primary rune, combining []rune, style tcell.Style) {
	c.Graphics.ForceDraw(c.X+x, c.Y+y, primary, combining, style)
}

func (c *Clip) Translate(offsetX int, offsetY int) view.Renderer {
	copy := *c
	copy.X += offsetX
	copy.Y += offsetY
	return &copy
}
