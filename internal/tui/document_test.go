package tui_test

import (
	"testing"

	"github.com/desktopgame/ckro/internal/tui"
)

func TestDocument(t *testing.T) {
	doc := tui.Document{}
	doc.Init()

	doc.InsertString("Hello")
	cursorColumn := doc.GetCursorColumn()
	if cursorColumn != 5 {
		t.Fatalf("got %q, want %q", cursorColumn, 5)
	}

	doc.RemoveChar()
	cursorColumn = doc.GetCursorColumn()
	if cursorColumn != 4 {
		t.Fatalf("got %q, want %q", cursorColumn, 4)
	}
}
