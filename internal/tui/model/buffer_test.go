package model_test

import (
	"testing"

	"github.com/desktopgame/ckro/internal/tui/model"
)

func TestLine(t *testing.T) {
	line := model.Line{}
	line.AppendString("Hello")

	content := line.GetContent()
	if content != "Hello" {
		t.Fatalf("got %q, want %q", content, "Hello")
	}

	line.AppendString(", world!")
	content = line.GetContent()
	if content != "Hello, world!" {
		t.Fatalf("got %q, want %q", content, "Hello, world!")
	}

	line.Remove(0, 3)
	content = line.GetContent()
	if content != "lo, world!" {
		t.Fatalf("got %q, want %q", content, "lo, world!")
	}

	line.Remove(4, 2)
	content = line.GetContent()
	if content != "lo, rld!" {
		t.Fatalf("got %q, want %q", content, "lo, rld!")
	}
}

func TestBuffer(t *testing.T) {
	buf := model.Buffer{}
	buf.Init()
	buf.InsertString(0, 0, "Line1\nLine2")

	lc := buf.GetLineCount()
	if lc != 2 {
		t.Fatalf("got %q, want %q", lc, 2)
	}
}
