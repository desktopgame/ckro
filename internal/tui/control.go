package tui

import "github.com/gdamore/tcell/v2"

type Control interface {
	Update()
	Traverse(fm *FocusManager)
	Draw(s tcell.Screen)

	MinimumSize(width int, height int) (Width int, Height int)
	Move(x int, y int)
	Layout(w int, h int)

	IsFlexibleWidth() bool
	IsFlexibleHeight() bool
}
