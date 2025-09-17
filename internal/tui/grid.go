package tui

import (
	"github.com/desktopgame/ckro/internal/tui/presenter"
	"github.com/gdamore/tcell/v2"
)

type GridCell struct {
	TextBox       *TextBox
	TextPresenter TextPresenter
}

func (gc *GridCell) Handle(ev tcell.Event) {
	gc.TextPresenter.Handle(gc.TextBox, ev)
}

func (gc *GridCell) GetTextBox() *TextBox {
	return gc.TextBox
}

func (gc *GridCell) GetTextPresenter() TextPresenter {
	return gc.TextPresenter
}

type Grid struct {
	rowCount    int
	columnCount int
	table       [][]GridCell
	x           int
	y           int
	width       int
	height      int
}

func (g *Grid) Init(rowCount int, columnCount int) {
	g.rowCount = rowCount
	g.columnCount = columnCount

	for i := 0; i < rowCount; i++ {
		line := []GridCell{}
		for j := 0; j < columnCount; j++ {
			tb := TextBox{}
			tb.Init()

			tp := presenter.FrameTextPresenter{}
			line = append(line, GridCell{
				TextBox:       &tb,
				TextPresenter: &tp,
			})
		}
		g.table = append(g.table, line)
	}
}

func (g *Grid) Set(row int, column int, presenter TextPresenter) *GridCell {
	if row < 0 || row >= g.rowCount || column < 0 || column >= g.columnCount {
		return nil
	}

	c := g.table[row][column]
	if presenter != nil {
		c.TextPresenter = presenter
	}
	return &c
}

func (g *Grid) Traverse(fm *FocusManager) {
	for _, row := range g.table {
		for _, c := range row {
			fm.Register(&c)
		}
	}
}

func (g *Grid) Update() {
	for _, row := range g.table {
		for _, c := range row {
			c.TextPresenter.Present(c.TextBox)
		}
	}
}

func (g *Grid) Draw(s tcell.Screen) {
	for _, row := range g.table {
		for _, c := range row {
			c.TextBox.Draw(s)
		}
	}

	for x := g.x; x < g.x+g.width; x++ {
		s.SetContent(x, 0, '-', nil, tcell.StyleDefault)
		s.SetContent(x, g.y+g.height-1, '@', nil, tcell.StyleDefault)
	}

	for y := g.y; y < g.y+g.height; y++ {
		s.SetContent(0, y, '|', nil, tcell.StyleDefault)
		s.SetContent(g.width-1, y, '|', nil, tcell.StyleDefault)
	}
}

func (g *Grid) MinimumSize(width int, height int) (Width int, Height int) {
	return width, height
}

func (g *Grid) Move(x int, y int) {
	g.x = x
	g.y = y
}

func (g *Grid) Layout(w int, h int) {
	xBorders := g.columnCount + 1
	yBorders := g.rowCount + 1

	offsetY := g.y + 1
	yMod := (h - yBorders) % g.rowCount
	for i := 0; i < g.rowCount; i++ {
		consumeY := min(yMod, yMod/g.rowCount)

		offsetX := g.x + 1
		xMod := (w - xBorders) % g.columnCount
		for j := 0; j < g.columnCount; j++ {
			gc := g.table[i][j]
			gc.TextBox.X = offsetX
			gc.TextBox.Y = offsetY
			gc.TextBox.Width = (w - xBorders) / g.columnCount
			gc.TextBox.Height = (h - yBorders) / g.rowCount

			if xMod > 0 {
				consumeX := min(xMod, xMod/g.columnCount)
				gc.TextBox.Width += consumeX
				xMod -= consumeX
			}

			if yMod > 0 {
				gc.TextBox.Height += consumeY
			}

			offsetX += gc.TextBox.Width + 1
		}
		yMod -= consumeY
		offsetY += ((h - yBorders) / g.rowCount) + consumeY + 1
	}

	g.width = w
	g.height = h
}

func (g *Grid) IsFlexibleWidth() bool {
	return false
}

func (g *Grid) IsFlexibleHeight() bool {
	return false
}
