package model_test

import (
	"testing"

	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/stretchr/testify/assert"
)

func TestLine(t *testing.T) {
	line := model.Line{}
	line.AppendString("Hello")

	content := line.GetContent()
	assert.Equal(t, content, "Hello")

	line.AppendString(", world!")
	content = line.GetContent()
	assert.Equal(t, content, "Hello, world!")

	line.Remove(0, 3)
	content = line.GetContent()
	assert.Equal(t, content, "lo, world!")

	line.Remove(4, 2)
	content = line.GetContent()
	assert.Equal(t, content, "lo, rld!")
}

func TestBuffer(t *testing.T) {
	buf := model.Buffer{}
	buf.Init()
	buf.InsertString(0, 0, "Line1\nLine2")

	lc := buf.GetLineCount()
	assert.Equal(t, lc, 2)
}
