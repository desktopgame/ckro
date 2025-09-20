package tui

import (
	"github.com/desktopgame/ckro/internal/tui/presenter"
)

type GridCell struct {
	Control      Control
	StaticWidth  int
	StaticHeight int
	Width        int
	Height       int
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
			tile := Tile{}
			tile.Init()
			tile.FlexibleWidth = true
			tile.FlexibleHeight = true
			tile.TextPresenter = &presenter.FrameTextPresenter{}
			line = append(line, GridCell{
				Control: &tile,
			})
		}
		g.table = append(g.table, line)
	}
}

func (g *Grid) SetTile(row int, column int, staticWidth int, staticHeight int, presenter TextPresenter) *GridCell {
	if row < 0 || row >= g.rowCount || column < 0 || column >= g.columnCount {
		return nil
	}

	c := &g.table[row][column]

	tile := Tile{}
	tile.Init()
	tile.MinimumWidth = max(0, staticWidth)
	tile.MinimumHeight = max(0, staticHeight)
	tile.FlexibleWidth = staticWidth == 0
	tile.FlexibleHeight = staticHeight == 0
	tile.TextPresenter = presenter

	c.StaticWidth = staticWidth
	c.StaticHeight = staticHeight
	c.Control = &tile

	return c
}

func (g *Grid) SetControl(row int, column int, ctrl Control) *GridCell {
	if row < 0 || row >= g.rowCount || column < 0 || column >= g.columnCount {
		return nil
	}

	c := &g.table[row][column]
	c.Control = ctrl
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

func (g *Grid) HeightTable(h int) []int {
	yBorders := g.rowCount + 1

	heightTable := []int{}
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
		heightTable = append(heightTable, maxStaticHeight)
		sumHeight += maxStaticHeight
	}

	for i := 0; i < g.rowCount; i++ {
		if divRows > 0 {
			msh := heightTable[i]
			div := ((h - yBorders) - sumHeight) / divRows

			if div > msh {
				heightTable[i] = div
			}
		}
	}
	return heightTable
}

func (g *Grid) Traverse(fm *FocusManager) {
	for _, row := range g.table {
		for _, c := range row {
			c.Control.Traverse(fm)
		}
	}
}

func (g *Grid) Update() {
	for _, row := range g.table {
		for _, c := range row {
			c.Control.Update()
		}
	}
}

func (g *Grid) Draw(gg *Graphics) {
	for _, row := range g.table {
		for _, c := range row {
			c.Control.Draw(gg)
		}
	}
}

func (g *Grid) MinimumSize(width int, height int) (Width int, Height int) {
	xBorders := g.columnCount + 1
	yBorders := g.rowCount + 1
	sw, sh := g.StaticSize()
	return max(width, sw+xBorders), max(height, sh+yBorders)
}

func (g *Grid) Move(x int, y int) {
	g.x = x
	g.y = y
}

func (g *Grid) Layout(w int, h int) {
	xBorders := g.columnCount + 1
	yBorders := g.rowCount + 1

	for i := 0; i < g.rowCount; i++ {
		for j := 0; j < g.columnCount; j++ {
			gc := &g.table[i][j]
			mw, mh := gc.Control.MinimumSize(w, h)
			if gc.Control.IsFlexibleWidth() {
				gc.StaticWidth = 0
			} else {
				gc.StaticWidth = mw
			}
			if gc.Control.IsFlexibleHeight() {
				gc.StaticHeight = 0
			} else {
				gc.StaticHeight = mh
			}
		}
	}

	sw, _ := g.StaticSize()
	heightTable := g.HeightTable(h)

	useHeight := 0
	for i := 0; i < g.rowCount; i++ {
		useHeight += heightTable[i]
	}

	offsetY := g.y + 1
	yMod := max(0, h-useHeight-yBorders)
	for i := 0; i < g.rowCount; i++ {
		maxConsumeY := 0

		offsetX := g.x + 1
		xMod := 0
		if g.columnCount != g.StaticColumns(i) {
			xMod = ((w - xBorders) - sw) % (g.columnCount - g.StaticColumns(i))
		}
		maxHeight := 0
		for j := 0; j < g.columnCount; j++ {
			gc := g.table[i][j]
			gc.Control.Move(offsetX, offsetY)

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

			height := gc.StaticHeight
			consumeY := 0
			if height == 0 {
				height = heightTable[i]

				consumeY = max(0, min(yMod, yMod/(g.rowCount-g.StaticRows(j))))
				if consumeY == 0 && yMod > 0 {
					consumeY = yMod
				}
				if yMod > 0 {
					height += consumeY
				}
			}

			gc.Control.Layout(width, height)
			gc.Width = width
			gc.Height = height

			if height > maxHeight {
				maxHeight = height
				maxConsumeY = consumeY
			}

			offsetX += width + 1
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
	return true
}

func (g *Grid) IsFlexibleHeight() bool {
	return true
}
