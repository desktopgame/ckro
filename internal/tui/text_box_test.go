package tui_test

import (
	"testing"

	"github.com/desktopgame/ckro/internal/tui"
)

func TestTextBox(t *testing.T) {
	// 画面に一言
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
	tb.HasFocus = true

	cursorRow := tb.Document.GetCursorRow()
	if cursorRow != 5 {
		t.Fatalf("got %q, want %q", cursorRow, 5)
	}

	tb.CursorUpdate()

	row := tb.ScrollY
	if row != 1 {
		t.Fatalf("got %q, want %q", row, 1)
	}

}
