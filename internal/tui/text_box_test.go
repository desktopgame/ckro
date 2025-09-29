package tui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextBox01(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	tb.InsertString("1234あ")
	tb.ShowCursor = true

	assert.Equal(t, tb.bytePos.Row, 0)
	assert.Equal(t, tb.bytePos.Column, len("1234")+1)
}

func TestTextBox02(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	tb.InsertString("1234\n1234\n123あ")
	tb.ShowCursor = true

	assert.Equal(t, tb.bytePos.Row, 2)
	assert.Equal(t, tb.bytePos.Column, len("123")+1)
}
