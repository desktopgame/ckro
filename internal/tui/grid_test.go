package tui_test

import (
	"testing"

	"github.com/desktopgame/ckro/internal/tui"
	"github.com/desktopgame/ckro/internal/tui/presenter"
)

func assertPos(t *testing.T, gc *tui.GridCell, x int, y int) {
	if gc.TextBox.X != x {
		t.Fatalf("got %d, want %d", gc.TextBox.X, x)
	}
	if gc.TextBox.Y != y {
		t.Fatalf("got %d, want %d", gc.TextBox.Y, y)
	}
}

func assertSize(t *testing.T, gc *tui.GridCell, w int, h int) {
	if gc.TextBox.Width != w {
		t.Fatalf("got %d, want %d", gc.TextBox.Width, w)
	}
	if gc.TextBox.Height != h {
		t.Fatalf("got %d, want %d", gc.TextBox.Height, h)
	}
}

func TestGrid(t *testing.T) {
	g := tui.Grid{}
	g.Init(2, 2)
	c00 := g.Set(0, 0, &presenter.EditTextPresenter{})
	c01 := g.Set(0, 1, &presenter.EditTextPresenter{})
	c10 := g.Set(1, 0, &presenter.EditTextPresenter{})
	c11 := g.Set(1, 1, &presenter.EditTextPresenter{})
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

func TestGrid1(t *testing.T) {
	g := tui.Grid{}
	g.Init(2, 2)
	c00 := g.Set(0, 0, &presenter.EditTextPresenter{})
	c01 := g.Set(0, 1, &presenter.EditTextPresenter{})
	c10 := g.Set(1, 0, &presenter.EditTextPresenter{})
	c11 := g.Set(1, 1, &presenter.EditTextPresenter{})
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
