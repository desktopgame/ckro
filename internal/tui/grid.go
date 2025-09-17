package tui

import (
	"github.com/desktopgame/ckro/internal/tui/presenter"
	"github.com/gdamore/tcell/v2"
)

type GridCell struct {
	TextBox       *TextBox
	TextPresenter TextPresenter
	StaticWidth   int
	StaticHeight  int
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

func (g *Grid) Set(row int, column int, staticWidth int, staticHeight int, presenter TextPresenter) *GridCell {
	if row < 0 || row >= g.rowCount || column < 0 || column >= g.columnCount {
		return nil
	}

	c := &g.table[row][column]
	c.StaticWidth = staticWidth
	c.StaticHeight = staticHeight
	if presenter != nil {
		c.TextPresenter = presenter
	}
	return c
}

func (g *Grid) StaticSize() (Width int, Height int) {
	sw := 0
	for i := 0; i < g.columnCount; i++ {
		msw := g.MaxStaticWidth(i)
		sw += msw
	}

	sh := 0
	for i := 0; i < g.rowCount; i++ {
		msh := g.MaxStaticHeight(i)
		sh += msh
	}

	return sw, sh
}

func (g *Grid) MaxStaticWidth(column int) int {
	w := 0
	for i := 0; i < g.rowCount; i++ {
		gc := g.table[i][column]

		if gc.StaticWidth > w {
			w = gc.StaticWidth
		}
	}
	return w
}

func (g *Grid) MaxStaticHeight(row int) int {
	h := 0
	for i := 0; i < g.columnCount; i++ {
		gc := g.table[row][i]

		if gc.StaticHeight > h {
			h = gc.StaticHeight
		}
	}
	return h
}

func (g *Grid) StaticRows(column int) int {
	c := 0
	for i := 0; i < g.rowCount; i++ {
		if g.table[i][column].StaticHeight > 0 {
			c++
		}
	}
	return c
}

func (g *Grid) StaticColumns(row int) int {
	c := 0
	for i := 0; i < g.columnCount; i++ {
		if g.table[row][i].StaticWidth > 0 {
			c++
		}
	}
	return c
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
	/*
		xBorders := g.columnCount + 1
		yBorders := g.rowCount + 1

		offsetY := g.y + 1
		yMod := (g.height - yBorders) % g.rowCount
		for i := 0; i < g.rowCount; i++ {

			for x := g.x; x < g.x+g.width; x++ {
				s.SetContent(x, offsetY-1, '-', nil, tcell.StyleDefault)
			}

			consumeY := max(0, min(yMod, yMod/g.rowCount))
			if consumeY == 0 && yMod > 0 {
				consumeY = yMod
			}

			if consumeY < 0 {
				consumeY = 0
			}
			yMod -= consumeY
			offsetY += ((g.height - yBorders) / g.rowCount) + consumeY + 1
		}

		offsetX := g.x + 1
		xMod := (g.width - xBorders) % g.columnCount
		for i := 0; i < g.columnCount; i++ {
			for j := 0; j < g.columnCount; j++ {
				for y := g.y; y < g.y+g.height; y++ {
					s.SetContent(offsetX-1, y, '|', nil, tcell.StyleDefault)
				}

				width := (g.width - xBorders) / g.columnCount

				if xMod > 0 {
					consumeX := max(0, min(xMod, xMod/g.columnCount))
					if consumeX == 0 && xMod > 0 {
						consumeX = xMod
					}

					width += consumeX
					xMod -= consumeX
				}

				offsetX += width + 1
			}
		}
	*/

	for _, row := range g.table {
		for _, c := range row {
			c.TextBox.Draw(s)
		}
	}

	/*
		for x := g.x; x < g.x+g.width; x++ {
			s.SetContent(x, g.y, '-', nil, tcell.StyleDefault)
			s.SetContent(x, g.y+g.height-1, '-', nil, tcell.StyleDefault)
		}

		for y := g.y; y < g.y+g.height; y++ {
			s.SetContent(g.width-1, y, '|', nil, tcell.StyleDefault)
		}

		s.SetContent(g.x, g.y, '*', nil, tcell.StyleDefault)
		s.SetContent(g.x+g.width-1, g.y, '*', nil, tcell.StyleDefault)
		s.SetContent(g.x, g.y+g.height-1, '*', nil, tcell.StyleDefault)
		s.SetContent(g.x+g.width-1, g.y+g.height-1, '*', nil, tcell.StyleDefault)
	*/
}

func (g *Grid) MinimumSize(width int, height int) (Width int, Height int) {
	sw, sh := g.StaticSize()
	return max(width, sw), max(height, sh)
}

func (g *Grid) Move(x int, y int) {
	g.x = x
	g.y = y
}

func (g *Grid) Layout(w int, h int) {
	xBorders := g.columnCount + 1
	yBorders := g.rowCount + 1

	sw, sh := g.StaticSize()

	maxStaticHeightTable := []int{}
	divRows := 0
	sumHeight := 0
	for i := 0; i < g.rowCount; i++ {
		maxStaticHeight := 0
		dynamicCount := 0
		for j := 0; j < g.columnCount; j++ {
			gc := g.table[i][j]

			if gc.StaticHeight > maxStaticHeight {
				maxStaticHeight = gc.StaticHeight
			} else if gc.StaticHeight == 0 {
				dynamicCount++
			}
		}
		if dynamicCount > 0 {
			maxStaticHeight = 0
			divRows++
		}
		maxStaticHeightTable = append(maxStaticHeightTable, maxStaticHeight)
		sumHeight += maxStaticHeight
	}

	for i := 0; i < g.rowCount; i++ {
		msh := maxStaticHeightTable[i]
		div := ((h - yBorders) - sumHeight) / divRows

		if div > msh {
			maxStaticHeightTable[i] = div
		}
	}

	yyMod := 0
	for i := 0; i < g.rowCount; i++ {
		yyMod += maxStaticHeightTable[i]
	}

	offsetY := g.y + 1
	yMod := max(0, h-yyMod-yBorders)
	for i := 0; i < g.rowCount; i++ {
		maxConsumeY := 0

		offsetX := g.x + 1
		xMod := ((w - xBorders) - sw) % (g.columnCount - g.StaticColumns(i))
		maxHeight := 0
		for j := 0; j < g.columnCount; j++ {
			gc := g.table[i][j]
			gc.TextBox.X = offsetX
			gc.TextBox.Y = offsetY

			width := gc.StaticWidth
			if width == 0 {
				width = ((w - xBorders) - sw) / (g.columnCount - g.StaticColumns(i))

				if xMod > 0 {
					consumeX := max(0, min(xMod, xMod/(g.columnCount-g.StaticColumns(i))))
					if consumeX == 0 && xMod > 0 {
						consumeX = xMod
					}

					width += consumeX
					xMod -= consumeX
				}
			}
			gc.TextBox.Width = width

			height := gc.StaticHeight
			consumeY := 0
			if height == 0 {
				height = ((h - yBorders) - sh) / (g.rowCount - g.StaticRows(j))

				consumeY = max(0, min(yMod, yMod/(g.rowCount-g.StaticRows(j))))
				if consumeY == 0 && yMod > 0 {
					consumeY = yMod
				}
				if yMod > 0 {
					height += consumeY
				}
			}
			gc.TextBox.Height = height

			if height > maxHeight {
				maxHeight = height
				maxConsumeY = consumeY
			}

			offsetX += gc.TextBox.Width + 1
		}
		if maxConsumeY < 0 {
			maxConsumeY = 0
		}
		yMod -= maxConsumeY
		offsetY += maxHeight + 1
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
