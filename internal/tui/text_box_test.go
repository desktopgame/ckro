package tui

import (
	"testing"

	"github.com/desktopgame/ckro/internal/tui/model"
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

func TestTextBox03(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	tb.InsertString("1234\n1234\n123あ")
	tb.bytePos = model.Position{
		Row:    0,
		Column: 0,
	}
	tb.viewPosition = 0
	tb.ShowCursor = true

	tb.MoveRight()

	assert.Equal(t, tb.bytePos.Row, 0)
	assert.Equal(t, tb.bytePos.Column, len("1"))
}

func TestTextBox04(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	tb.InsertString("1234\n1234\n123あ")
	tb.bytePos = model.Position{
		Row:    0,
		Column: 0,
	}
	tb.viewPosition = 0
	tb.ShowCursor = true

	tb.MoveRight() // 2
	tb.MoveRight() // 3
	tb.MoveRight() // 4
	tb.MoveRight() // NL

	assert.Equal(t, tb.bytePos.Row, 0)
	assert.Equal(t, tb.bytePos.Column, len("1234"))
}
