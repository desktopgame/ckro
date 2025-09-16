package tui

import "github.com/gdamore/tcell/v2"

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
}

func (g *Grid) Init(rowCount int, columnCount int) {
	g.rowCount = rowCount
	g.columnCount = columnCount

	for i := 0; i < rowCount; i++ {
		line := []GridCell{}
		for j := 0; j < columnCount; j++ {
			line = append(line, GridCell{})
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

func (g *Grid) Draw(s tcell.Screen) {}

func (g *Grid) MinimumSize(width int, height int) (Width int, Height int) {
	return width, height
}

func (g *Grid) Move(x int, y int) {
	g.x = x
	g.y = y
}

func (g *Grid) Layout(w int, h int) {}

func (g *Grid) IsFlexibleWidth() bool {
	return false
}

func (g *Grid) IsFlexibleHeight() bool {
	return false
}
