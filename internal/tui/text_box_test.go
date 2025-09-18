package tui_test

import (
	"testing"

	"github.com/desktopgame/ckro/internal/tui"
)

func TestTextBox(t *testing.T) {
	msg :=
		`
Hello, world1
👨‍👩‍👧‍👦
Hello, world2
あいうえお
`

	tb := tui.TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 20
	tb.Height = 6
	tb.Document.InsertString(msg)
	tb.ShowCursor = true

	tb.CursorUpdate()

	cursorRow := tb.Document.GetCursorRow()
	if cursorRow != 5 {
		t.Fatalf("got %q, want %q", cursorRow, 5)
	}

	tb.Document.InsertLine()
	tb.CursorUpdate()

	row := tb.GetScrollY()
	if row != 1 {
		t.Fatalf("got %q, want %q", row, 1)
	}

}

func TestCursor(t *testing.T) {
	tb := tui.TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 20
	tb.Height = 6
	tb.Document.InsertString("12345678901234567890")
	tb.ShowCursor = true

	tb.CursorUpdate()

	col, _, _, _ := tb.CursorPosition()
	if col != 0 {
		t.Fatalf("got %q, want %q", col, 0)
	}
}

func TestScroll(t *testing.T) {
	tb := tui.TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 5
	tb.Height = 2
	tb.Document.InsertString("123456789")
	tb.ShowCursor = true

	tb.CursorUpdate()

	col, row, _, _ := tb.CursorPosition()
	if row-tb.GetScrollY() != 1 {
		t.Fatalf("got %q, want %q", row, 1)
	}
	if col != 4 {
		t.Fatalf("got %q, want %q", col, 4)
	}

	tb.Document.InsertString("0")
	tb.CursorUpdate()

	col, row, _, _ = tb.CursorPosition()
	if row-tb.GetScrollY() != 1 {
		t.Fatalf("got %q, want %q", row-tb.GetScrollY(), 1)
	}
	if col != 0 {
		t.Fatalf("got %q, want %q", col, 0)
	}
}

func TestWrap(t *testing.T) {
	tb := tui.TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 5
	tb.Height = 2
	tb.Document.InsertString("1234あ")
	tb.ShowCursor = true

	tb.CursorUpdate()

	col, row, _, _ := tb.CursorPosition()
	if row-tb.GetScrollY() != 1 {
		t.Fatalf("got %q, want %q", row-tb.GetScrollY(), 1)
	}
	if col != 2 {
		t.Fatalf("got %q, want %q", col, 2)
	}
}

func TestWrapWithTab(t *testing.T) {
	tb := tui.TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 5
	tb.Height = 2
	tb.Document.InsertString("12345\t")
	tb.ShowCursor = true

	tb.CursorUpdate()

	col, _, _, _ := tb.CursorPosition()
	if col != 4 {
		t.Fatalf("got %q, want %q", col, 4)
	}

	tb.Document.MoveLeft()

	tb.CursorUpdate()

	col, _, _, _ = tb.CursorPosition()
	if col != 0 {
		t.Fatalf("got %q, want %q", col, 0)
	}
}
