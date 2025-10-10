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
	offsetX    int
	offsetY    int
	regionW    int
	regionH    int
}

// SetContent is set a character to specified cell, if not already settled.
// TODO: refactor
func (c *Clip) SetContent(x int, y int, primary rune, combining []rune, style tcell.Style) {
	y += c.offsetY
	if y < c.FirstLineY || y >= c.FirstLineY+c.Height {
		return
	}
	offset := y - c.FirstLineY
	c.Graphics.Draw(c.offsetX+c.X+x, c.Y+offset, primary, combining, style)
}

// SetCursor is set a character to specified cell.
// TODO: refactor
func (c *Clip) SetCursor(x int, y int, primary rune, combining []rune, style tcell.Style) {
	c.Graphics.ForceDraw(c.offsetX+c.X+x, c.offsetY+c.Y+y, primary, combining, style)
}

func (c *Clip) Translate(offsetX int, offsetY int) view.Renderer {
	copy := *c
	copy.offsetX += offsetX
	copy.offsetY += offsetY
	copy.regionW = 0
	copy.regionH = 0
	return &copy
}

func (c *Clip) Region(width int, height int) view.Renderer {
	copy := *c
	copy.regionW = width
	copy.regionH = height
	return &copy
}
