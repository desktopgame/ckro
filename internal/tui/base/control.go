package base

import "github.com/gdamore/tcell/v2"

type Control interface {
	Traverse(fm *FocusManager)
	Update()
	Draw(s tcell.Screen)

	MinimumSize(width int, height int) (Width int, Height int)
	Move(x int, y int)
	Layout(w int, h int)

	IsFlexibleWidth() bool
	IsFlexibleHeight() bool
}
