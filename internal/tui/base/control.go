package base

type Control interface {
	Traverse(fm *FocusManager)
	Update()
	Draw(g *Graphics)

	MinimumSize(width int, height int) (Width int, Height int)
	Move(x int, y int)
	Layout(w int, h int)

	IsFlexibleWidth() bool
	IsFlexibleHeight() bool
}
