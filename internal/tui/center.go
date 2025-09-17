package tui

import (
	"github.com/gdamore/tcell/v2"
)

type Center struct {
	Control Control
	Width   int
	Height  int
	x       int
	y       int
}

func (c *Center) Traverse(fm *FocusManager) {
	c.Control.Traverse(fm)
}

func (c *Center) Update() {
	c.Control.Update()
}

func (c *Center) Draw(s tcell.Screen) {
	c.Control.Draw(s)
}

func (c *Center) MinimumSize(width int, height int) (Width int, Height int) {
	mw, mh := c.Control.MinimumSize(width, height)
	return max(mw, width), max(mh, height)
}

func (c *Center) Move(x int, y int) {
	c.x = x
	c.y = y
}

func (c *Center) Layout(width int, height int) {
	c.Control.Move(
		c.x+(width-c.Width)/2,
		c.y+(height-c.Height)/2,
	)
	c.Control.Layout(c.Width, c.Height)
}

func (c *Center) IsFlexibleWidth() bool {
	return false
}

func (c *Center) IsFlexibleHeight() bool {
	return false
}
