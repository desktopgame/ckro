package tui_test

import (
	"testing"

	"github.com/desktopgame/ckro/internal/tui"
	"github.com/desktopgame/ckro/internal/tui/presenter"
)

func assertPos(t *testing.T, gc *tui.GridCell, x int, y int) {
	if tile, ok := gc.Control.(*tui.Tile); ok {
		if tile.TextBox.X != x {
			t.Fatalf("got %d, want %d", tile.TextBox.X, x)
		}
		if tile.TextBox.Y != y {
			t.Fatalf("got %d, want %d", tile.TextBox.Y, y)
		}
	}
}

func assertSize(t *testing.T, gc *tui.GridCell, w int, h int) {
	if tile, ok := gc.Control.(*tui.Tile); ok {
		if tile.TextBox.Width != w {
			t.Fatalf("got %d, want %d", tile.TextBox.Width, w)
		}
		if tile.TextBox.Height != h {
			t.Fatalf("got %d, want %d", tile.TextBox.Height, h)
		}
	}
}

func TestGrid01(t *testing.T) {
	g := tui.Grid{}
	g.Init(2, 2)
	c00 := g.SetTile(0, 0, 0, 0, &presenter.EditTextPresenter{})
	c01 := g.SetTile(0, 1, 0, 0, &presenter.EditTextPresenter{})
	c10 := g.SetTile(1, 0, 0, 0, &presenter.EditTextPresenter{})
	c11 := g.SetTile(1, 1, 0, 0, &presenter.EditTextPresenter{})
	g.Layout(11, 11)

	/*
		0----------------------------------------------------
		1
		2
		3
		4
		5----------------------------------------------------
		6
		7
		8
		9
		10---------------------------------------------------
	*/

	assertPos(t, c00, 1, 1)
	assertPos(t, c01, 6, 1)
	assertPos(t, c10, 1, 6)
	assertPos(t, c11, 6, 6)
}

func TestGrid02(t *testing.T) {
	g := tui.Grid{}
	g.Init(2, 2)
	c00 := g.SetTile(0, 0, 0, 0, &presenter.EditTextPresenter{})
	c01 := g.SetTile(0, 1, 0, 0, &presenter.EditTextPresenter{})
	c10 := g.SetTile(1, 0, 0, 0, &presenter.EditTextPresenter{})
	c11 := g.SetTile(1, 1, 0, 0, &presenter.EditTextPresenter{})
	g.Layout(12, 12)

	/*
		0----------------------------------------------------
		1
		2
		3
		4
		5
		6----------------------------------------------------
		7
		8
		9
		10
		11---------------------------------------------------
	*/

	assertPos(t, c00, 1, 1)
	assertSize(t, c00, 5, 5)
	assertPos(t, c01, 7, 1)
	assertPos(t, c10, 1, 7)
	assertPos(t, c11, 7, 7)
}

func TestGrid03(t *testing.T) {
	g := tui.Grid{}
	g.Init(2, 1)
	c00 := g.SetTile(0, 0, 0, 1, &presenter.EditTextPresenter{})
	c01 := g.SetTile(1, 0, 0, 0, &presenter.EditTextPresenter{})
	g.Layout(10, 10)

	assertPos(t, c00, 1, 1)
	assertSize(t, c00, 8, 1)
	assertPos(t, c01, 1, 3)
}

func TestGrid04(t *testing.T) {
	g := tui.Grid{}
	g.Init(2, 2)
	c00 := g.SetTile(0, 0, 0, 3, &presenter.FrameTextPresenter{})
	c01 := g.SetTile(1, 0, 0, 0, &presenter.FrameTextPresenter{})
	c10 := g.SetTile(0, 1, 0, 0, &presenter.FrameTextPresenter{})
	c11 := g.SetTile(1, 1, 0, 0, &presenter.FrameTextPresenter{})
	g.Layout(20, 20)

	assertPos(t, c00, 1, 1)
	assertPos(t, c01, 1, 11)
	assertPos(t, c10, 11, 1)
	assertPos(t, c11, 11, 11)
}
