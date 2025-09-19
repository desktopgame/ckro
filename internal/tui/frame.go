package tui

import (
	"github.com/gdamore/tcell/v2"
)

type Frame struct {
	Control Control
	x       int
	y       int
	width   int
	height  int
}

func (fr *Frame) Traverse(fm *FocusManager) {
	fr.Control.Traverse(fm)
}

func (fr *Frame) Update() {
	fr.Control.Update()
}

func (fr *Frame) Draw(g *Graphics) {
	for i := 1; i < fr.width-1; i++ {
		g.Draw(fr.x+i, fr.y, '-', nil, tcell.StyleDefault)
		g.Draw(fr.x+i, fr.y+fr.height-1, '-', nil, tcell.StyleDefault)
	}
	for i := 1; i < fr.height-1; i++ {
		g.Draw(fr.x, fr.y+i, '|', nil, tcell.StyleDefault)
		g.Draw(fr.x+fr.width-1, fr.y+i, '|', nil, tcell.StyleDefault)
	}
	g.Draw(fr.x, fr.y, '*', nil, tcell.StyleDefault)
	g.Draw(fr.x, fr.y+fr.height-1, '*', nil, tcell.StyleDefault)
	g.Draw(fr.x+fr.width-1, fr.y, '*', nil, tcell.StyleDefault)
	g.Draw(fr.x+fr.width-1, fr.y+fr.height-1, '*', nil, tcell.StyleDefault)
	fr.Control.Draw(g)
	g.GlassRange(fr.x, fr.y, fr.width, fr.height)
}

func (fr *Frame) MinimumSize(width int, height int) (Width int, Height int) {
	mw, mh := fr.Control.MinimumSize(width-2, height-2)
	return mw + 2, mh + 2
}

func (fr *Frame) Move(x int, y int) {
	fr.x = x
	fr.y = y
}

func (fr *Frame) Layout(width int, height int) {
	offsetX := fr.x + 1
	offsetY := fr.y + 1

	fr.Control.Move(offsetX, offsetY)
	fr.Control.Layout(width-2, height-2)
	fr.width = width
	fr.height = height
}

func (fr *Frame) IsFlexibleWidth() bool {
	return fr.Control.IsFlexibleWidth()
}

func (fr *Frame) IsFlexibleHeight() bool {
	return fr.Control.IsFlexibleHeight()
}
