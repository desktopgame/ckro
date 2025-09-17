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

func (fr *Frame) Draw(s tcell.Screen) {
	for i := 0; i < fr.width+2; i++ {
		s.SetContent(fr.x+i, fr.y, '-', nil, tcell.StyleDefault)
		s.SetContent(fr.x+i, fr.y+fr.height-1, '-', nil, tcell.StyleDefault)
	}
	for i := 0; i < fr.height+2; i++ {
		s.SetContent(fr.x, fr.y+i, '|', nil, tcell.StyleDefault)
		s.SetContent(fr.x+fr.width-1, fr.y+i, '|', nil, tcell.StyleDefault)
	}
	fr.Control.Draw(s)
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
	mw, mh := fr.Control.MinimumSize(width-2, height-2)
	offsetX := fr.x + ((width - mw) / 2)
	offsetY := fr.y + ((height - mh) / 2)

	fr.Control.Move(offsetX, offsetY)
	fr.Control.Layout(mw, mh)
	fr.width = width
	fr.height = height
}

func (fr *Frame) IsFlexibleWidth() bool {
	return fr.Control.IsFlexibleWidth()
}

func (fr *Frame) IsFlexibleHeight() bool {
	return fr.Control.IsFlexibleHeight()
}
