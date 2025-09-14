package tui

import "github.com/gdamore/tcell/v2"

type Clip struct {
	Screen tcell.Screen
	X      int
	Y      int
	Width  int
	Height int
}

func (c Clip) SetContent(x int, y int, primary rune, combining []rune, style tcell.Style) {
	c.Screen.SetContent(c.X+x, c.Y+y, primary, combining, style)
}
