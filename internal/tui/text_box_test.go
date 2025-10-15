package tui_test

import (
	"testing"

	"github.com/desktopgame/ckro/internal/tui"
	"github.com/desktopgame/ckro/internal/tui/extensions/litemark"
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/stretchr/testify/assert"
)

//
// Testing Library
//

func newPlainTextBox(width int, height int) *tui.TextBox {
	tb := tui.TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = width
	tb.Height = height
	tb.ShowCursor = true

	doc := model.PlainDocument{}
	doc.Init()
	tb.Document = &doc

	engine := tui.PlainTextViewResolver{}
	tb.ViewResolver = &engine
	return &tb
}

func newStyledTextBox(width int, height int) *tui.TextBox {
	tb := tui.TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = width
	tb.Height = height
	tb.ShowCursor = true

	doc := litemark.StyledDocument{}
	doc.Init()
	tb.Document = &doc

	engine := tui.LitemarkTextViewResolver{}
	tb.ViewResolver = &engine
	return &tb
}

func mustBytePos(t *testing.T, tb *tui.TextBox, row, column int) {
	bPos := tb.GetBytePosition()
	assert.Equal(t, bPos.StartPosition.Row, row)
	assert.Equal(t, bPos.StartPosition.Column, column)
}

//
// Tests
//

func TestTextBoxFind01(t *testing.T) {
	tb := newStyledTextBox(10, 10)
	tb.InsertString("aa\nbb\nccc")

	tb.FindPrev("b\ncc")
	mustBytePos(t, tb, 1, 1)

	tb.FindPrev("a\nb")
	mustBytePos(t, tb, 0, 1)
	assert.Equal(t, tb.GetBytePosition().Bytes, 1)

	tb.MoveLineEnd()
	tb.MoveDown()
	tb.MoveLineEnd()
	tb.MoveDown()
	tb.MoveLineEnd()
	tb.MoveDown()
	tb.FindPrev("a\nbb\nccc")
	mustBytePos(t, tb, 0, 1)
	assert.Equal(t, tb.GetBytePosition().Bytes, 1)
}

func TestTextBoxFind02(t *testing.T) {
	tb := newStyledTextBox(10, 10)
	tb.InsertString("aa\nbb\nccc")

	tb.FindPrev("cc")
	mustBytePos(t, tb, 2, 1)
	assert.Equal(t, tb.GetBytePosition().Bytes, 1)

	tb.FindPrev("bb")
	mustBytePos(t, tb, 1, 0)
	assert.Equal(t, tb.GetBytePosition().Bytes, 1)
}

func TestTextBoxFind03(t *testing.T) {
	tb := newStyledTextBox(10, 10)
	tb.InsertString("aa\nbb\nccc")

	tb.MoveReset()

	tb.FindNext("a\nbb")
	mustBytePos(t, tb, 0, 1)
	assert.Equal(t, tb.GetBytePosition().Bytes, 1)

	tb.FindNext("bb")
	mustBytePos(t, tb, 1, 0)
	assert.Equal(t, tb.GetBytePosition().Bytes, 1)

	tb.FindNext("b\nc")
	mustBytePos(t, tb, 1, 1)
	assert.Equal(t, tb.GetBytePosition().Bytes, 1)
}

func TestTextBoxFind04(t *testing.T) {
	tb := newStyledTextBox(10, 10)
	tb.InsertString("あいう\nbb\nかきく")

	tb.MoveReset()

	tb.FindNext("あいう")
	mustBytePos(t, tb, 0, 0)
	assert.Equal(t, tb.GetBytePosition().Bytes, len("あ"))

	tb.FindNext("う\nbb")
	mustBytePos(t, tb, 0, len("あい"))
	assert.Equal(t, tb.GetBytePosition().Bytes, len("う"))

	assert.False(t, tb.FindNext("う"))

	tb.FindNext("か")
	mustBytePos(t, tb, 2, 0)
	assert.Equal(t, tb.GetBytePosition().Bytes, len("か"))
}

func TestTextBoxFind05(t *testing.T) {
	tb := newStyledTextBox(10, 10)
	tb.InsertString("あいう\nbb\nかきく")

	tb.FindPrev("bb\nか")
	mustBytePos(t, tb, 1, 0)
	assert.Equal(t, tb.GetBytePosition().Bytes, 1)

	assert.False(t, tb.FindPrev("あいう\nb"))

	tb.FindPrev("あいう\n")
	mustBytePos(t, tb, 0, 0)
	assert.Equal(t, tb.GetBytePosition().Bytes, len("あ"))
}

func TestTextBoxReplace01(t *testing.T) {
	tb := newStyledTextBox(10, 10)
	tb.InsertString("あいう\nbb\nかきく")

	assert.True(t, tb.FindPrev("あいう"))
	tb.Replace(len("あいう"), "ABC")
	sg := tb.Document.Read(model.Range{
		StartPosition: model.Position{
			Row:    0,
			Column: 0,
		},
		EndPosition: model.Position{
			Row:    0,
			Column: tb.Document.GetLineBytes(0),
		},
	})
	assert.Equal(t, sg.GetLine(0), "ABC")
}
