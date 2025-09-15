package tui

type Control interface {
	MinimumSize() (Width int, Height int)
	Move(x int, y int)
	Layout(w int, h int)

	IsFlexibleWidth() bool
	IsFlexibleHeight() bool
}
