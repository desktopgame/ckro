package tui

import "github.com/gdamore/tcell/v2"

type Clip struct {
	Graphics *Graphics
	X        int
	Y        int
	Width    int
	Height   int
}

func (c Clip) SetContent(x int, y int, primary rune, combining []rune, style tcell.Style) {
	c.Graphics.Draw(c.X+x, c.Y+y, primary, combining, style)
}

func (c Clip) SetCursor(x int, y int, primary rune, combining []rune, style tcell.Style) {
	c.Graphics.ForceDraw(c.X+x, c.Y+y, primary, combining, style)
}
