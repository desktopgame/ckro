package base

import "github.com/gdamore/tcell/v2"

// Graphics is abstractions to operation for terminal.
// Graphics has bitmap for avoiding overdraw, also used to hide overwrap region by above layer.
type Graphics struct {
	screen tcell.Screen
	bitmap [][]bool
	width  int
	height int
}

// Init is initialize Graphics.
func (g *Graphics) Init(screen tcell.Screen) {
	g.screen = screen
}

// Clear is clear a bitmap.
func (g *Graphics) Clear() {
	if g.bitmap != nil {
		for i := 0; i < g.height; i++ {
			for j := 0; j < g.width; j++ {
				g.bitmap[i][j] = false
			}
		}
	}
}

// Draw is set a character to specified cell, if not already settled.
func (g *Graphics) Draw(x int, y int, primary rune, combine []rune, style tcell.Style) {
	if !g.bitmap[y][x] {
		g.screen.SetContent(x, y, primary, combine, style)
		g.bitmap[y][x] = true
	}
}

// ForceDraw is set a character to specified cell.
func (g *Graphics) ForceDraw(x int, y int, primary rune, combine []rune, style tcell.Style) {
	g.screen.SetContent(x, y, primary, combine, style)
	g.bitmap[y][x] = true
}

// Glass is mark the specified cell as used.
func (g *Graphics) Glass(x int, y int) {
	g.bitmap[y][x] = true
}

// GlassRange is mark the cells within specified range as used.
func (g *Graphics) GlassRange(x int, y int, width int, height int) {
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			g.Glass(x+j, y+i)
		}
	}
}

// Resize is restructure bitmap.
func (g *Graphics) Resize(width int, height int) {
	bitmap := [][]bool{}
	for i := 0; i < height; i++ {
		row := make([]bool, width)
		bitmap = append(bitmap, row)
	}
	g.bitmap = bitmap
	g.width = width
	g.height = height
}
