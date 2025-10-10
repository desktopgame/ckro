package tui_test

import (
	"strings"
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

func TestTextBox01(t *testing.T) {
	tb := newPlainTextBox(10, 10)
	tb.InsertString("1234あ")
	mustBytePos(t, tb, 0, len("1234あ"))
}

func TestTextBox02(t *testing.T) {
	tb := newPlainTextBox(10, 10)
	tb.InsertString("1234\n1234\n123あ")
	mustBytePos(t, tb, 2, len("123あ"))
}

func TestTextBox03(t *testing.T) {
	tb := newPlainTextBox(10, 10)
	tb.InsertString("1234\n1234\n123あ")
	tb.MoveTextStart()

	tb.MoveRight()
	mustBytePos(t, tb, 0, len("1"))
}

func TestTextBox04(t *testing.T) {
	tb := newPlainTextBox(10, 10)
	tb.InsertString("1234\n1234\n123あ")
	tb.MoveTextStart()

	tb.MoveRight() // 2
	tb.MoveRight() // 3
	tb.MoveRight() // 4
	tb.MoveRight() // NL
	mustBytePos(t, tb, 0, len("1234"))
}

func TestTextBox05(t *testing.T) {
	tb := newPlainTextBox(10, 10)
	tb.InsertString("1234\n1234\n123あ")
	tb.MoveTextStart()

	tb.MoveRight() // 2
	tb.MoveRight() // 3
	tb.MoveRight() // 4
	tb.MoveRight() // NL

	tb.InsertString("\n")
	mustBytePos(t, tb, 1, 0)
}

func TestTextBox06(t *testing.T) {
	tb := newStyledTextBox(10, 10)
	tb.InsertString("#")
	tb.InsertString(" ")

	mustBytePos(t, tb, 0, 2)

	tb.InsertString("AAA")
	mustBytePos(t, tb, 0, 5)
}

func TestTextBox07(t *testing.T) {
	tb := newPlainTextBox(10, 10)
	tb.InsertString("1234\n\n\n```\n123\n```")
	tb.MoveTextStart()

	tb.InsertString("\n")
	mustBytePos(t, tb, 1, 0)
}

func TestTextBox08(t *testing.T) {
	tb := newPlainTextBox(10, 10)
	tb.InsertString("\n\n\n\n\n")
	mustBytePos(t, tb, 5, 0)

	tb.MoveUp()
	mustBytePos(t, tb, 4, 0)
}

func TestTextBox09(t *testing.T) {
	tb := newPlainTextBox(10, 10)
	tb.InsertString("\n\n\n\n\n")
	mustBytePos(t, tb, 5, 0)

	tb.MoveLeft()
	mustBytePos(t, tb, 4, 0)
}

func TestTextBox10(t *testing.T) {
	tb := newPlainTextBox(10, 10)
	tb.InsertString("\n\n\n\n\n")
	mustBytePos(t, tb, 5, 0)

	tb.MoveLeft()
	mustBytePos(t, tb, 4, 0)

	tb.MoveRight()
	mustBytePos(t, tb, 5, 0)
}

func TestTextBox11(t *testing.T) {
	tb := newPlainTextBox(10, 10)

	tb.InsertString("a")
	tb.InsertString("\n")

	tb.InsertString("a")
	tb.InsertString("\n")

	tb.InsertString("a")

	tb.MoveUp()
	mustBytePos(t, tb, 1, 1)
}

func TestTextBox12(t *testing.T) {
	tb := newPlainTextBox(10, 10)

	tb.InsertString("a")
	tb.InsertString("\n")

	tb.InsertString("a")
	tb.InsertString("\n")

	tb.InsertString("a")

	tb.MoveUp()
	mustBytePos(t, tb, 1, 1)

	tb.MoveDown()
	mustBytePos(t, tb, 2, 1)
}

func TestTextBox13(t *testing.T) {
	tb := newPlainTextBox(10, 10)

	tb.InsertString("a")
	tb.InsertString("\n")

	tb.InsertString("a")
	tb.InsertString("\n")

	tb.InsertString("a")

	tb.MoveDown()
	mustBytePos(t, tb, 2, 1)
}

func TestTextBox14(t *testing.T) {
	tb := newPlainTextBox(10, 10)

	tb.InsertString("a")
	tb.InsertString("\n")

	tb.InsertString("a")
	tb.InsertString("\n")

	tb.InsertString("a")

	tb.MoveUp()
	mustBytePos(t, tb, 1, 1)

	tb.MoveUp()
	mustBytePos(t, tb, 0, 1)

	tb.MoveUp()
	mustBytePos(t, tb, 0, 0)
}

func TestTextBox15(t *testing.T) {
	tb := newStyledTextBox(10, 10)
	tb.InsertString("1234\n\n\n```\n123\n```\n\n")
	tb.MoveTextStart()

	tb.MoveDown()
	mustBytePos(t, tb, 1, 0)

	tb.MoveDown()
	mustBytePos(t, tb, 2, 0)

	tb.MoveDown()
	mustBytePos(t, tb, 4, 0)

	tb.MoveDown()
	mustBytePos(t, tb, 6, 0)
}

func TestTextBox16(t *testing.T) {
	tb := newStyledTextBox(10, 10)
	tb.InsertString("1234\n\n\n```\n123\n```\n\n")
	tb.MoveTextStart()

	tb.MoveDown()
	mustBytePos(t, tb, 1, 0)

	tb.MoveDown()
	mustBytePos(t, tb, 2, 0)

	tb.MoveDown()
	mustBytePos(t, tb, 4, 0)

	tb.MoveRight()
	tb.MoveRight()
	tb.MoveRight()
	mustBytePos(t, tb, 4, 3)

	tb.InsertString("\n")
	mustBytePos(t, tb, 5, 0)
}

func TestTextBox17(t *testing.T) {
	tb := newStyledTextBox(10, 10)
	tb.InsertString("1234\n\n\n```\n123\n```\n\n")
	tb.MoveTextStart()

	tb.MoveDown()
	mustBytePos(t, tb, 1, 0)

	tb.MoveDown()
	mustBytePos(t, tb, 2, 0)

	tb.MoveDown()
	mustBytePos(t, tb, 4, 0)

	tb.MoveRight()
	tb.MoveRight()
	tb.MoveRight()
	mustBytePos(t, tb, 4, 3)

	tb.InsertString("\n")
	mustBytePos(t, tb, 5, 0)

	tb.InsertString("\n")
	mustBytePos(t, tb, 6, 0)
}

func TestTextBox18(t *testing.T) {
	tb := newStyledTextBox(10, 10)
	tb.InsertString("1234")
	mustBytePos(t, tb, 0, 4)

	tb.RemoveChar()
	mustBytePos(t, tb, 0, 3)
}

func TestTextBox19(t *testing.T) {
	tb := newStyledTextBox(10, 10)
	tb.InsertString("\n\n\n")
	mustBytePos(t, tb, 3, 0)

	tb.RemoveChar()
	mustBytePos(t, tb, 2, 0)
}

func TestTextBox20(t *testing.T) {
	tb := newStyledTextBox(10, 10)
	tb.InsertString("#")
	tb.InsertString(" ")
	tb.InsertString("H")
	mustBytePos(t, tb, 0, 3)

	tb.RemoveChar()
}

func TestTextBox21(t *testing.T) {
	tb := newStyledTextBox(10, 10)
	tb.InsertString("#")
	tb.InsertString(" ")
	tb.InsertString("H")
	mustBytePos(t, tb, 0, 3)

	tb.MoveLeft()
	tb.MoveLeft()
	mustBytePos(t, tb, 0, 2)

	tb.MoveRight()
	tb.MoveRight()
	mustBytePos(t, tb, 0, 3)
}

func TestTextBox22(t *testing.T) {
	tb := newStyledTextBox(10, 10)
	tb.InsertString("\n")
	tb.InsertString("#")
	tb.InsertString(" ")
	tb.InsertString("H")
	mustBytePos(t, tb, 1, 3)

	tb.RemoveChar()
	mustBytePos(t, tb, 1, 0)

	tb.RemoveChar()
	mustBytePos(t, tb, 0, 0)
}

func TestTextBox23(t *testing.T) {
	tb := newStyledTextBox(10, 10)
	tb.InsertString("\n")
	tb.InsertString("#")
	tb.InsertString(" ")
	tb.InsertString("Heading")
	tb.InsertString("\n")
	tb.InsertString("\n")
	tb.InsertString("#")
	tb.InsertString(" ")
	tb.InsertString("Heading2")
	mustBytePos(t, tb, 3, 10)

	tb.MoveLeft()
	tb.MoveLeft()
	tb.MoveLeft()
	tb.MoveLeft()
	tb.MoveLeft()
	tb.MoveLeft()
	tb.MoveLeft()
	tb.MoveLeft() // H
	mustBytePos(t, tb, 3, 2)

	tb.RemoveChar()
	mustBytePos(t, tb, 2, 2)
}

func TestTextBox24(t *testing.T) {
	tb := newStyledTextBox(20, 20)
	tb.InsertString("\n")
	tb.InsertString("#")
	tb.InsertString(" ")
	tb.InsertString("Heading")
	tb.InsertString("\n")
	tb.InsertString("\n")
	tb.InsertString("##")
	tb.InsertString(" ")
	tb.InsertString("Heading2")
	mustBytePos(t, tb, 3, 11)

	tb.MoveLeft()
	mustBytePos(t, tb, 3, 10)

	tb.MoveLeft()
	tb.MoveLeft()
	tb.MoveLeft()
	tb.MoveLeft()
	tb.MoveLeft()
	tb.MoveLeft()
	tb.MoveLeft() // H
	mustBytePos(t, tb, 3, 3)

	tb.RemoveChar()
	mustBytePos(t, tb, 2, 3)
}

func TestTextBox25(t *testing.T) {
	tb := newStyledTextBox(20, 20)
	tb.InsertString("#")
	tb.InsertString(" ")
	tb.InsertString("H")
	mustBytePos(t, tb, 0, 3)

	tb.InsertString("e")
	mustBytePos(t, tb, 0, 4)
}

func TestTextBox26(t *testing.T) {
	tb := newStyledTextBox(20, 20)
	tb.InsertString("\n")
	tb.MoveLeft()
	mustBytePos(t, tb, 0, 0)

	tb.InsertString("##")
	tb.InsertString(" ")
	tb.InsertString("H")
	mustBytePos(t, tb, 0, 4)

	tb.InsertString("H")
	mustBytePos(t, tb, 0, 5)
}

func TestTextBox27(t *testing.T) {
	tb := newStyledTextBox(20, 20)

	tb.InsertString("##")
	tb.InsertString(" ")
	tb.InsertString("H")
	mustBytePos(t, tb, 0, 4)

	tb.InsertString("H")
	tb.InsertString("H")
	mustBytePos(t, tb, 0, 6)

	tb.MoveLeft()
	tb.MoveLeft()
	tb.MoveLeft()
	mustBytePos(t, tb, 0, 3)

	tb.InsertString("\n")
}

func TestTextBox28(t *testing.T) {
	tb := newStyledTextBox(20, 20)

	tb.InsertString("##")
	tb.InsertString(" ")
	tb.InsertString("H")
	mustBytePos(t, tb, 0, 4)

	tb.InsertString("H")
	tb.InsertString("H")
	mustBytePos(t, tb, 0, 6)

	tb.MoveLeft()
	tb.MoveLeft()
	tb.MoveLeft()
	mustBytePos(t, tb, 0, 3)

	tb.InsertString("\n")
	mustBytePos(t, tb, 1, 3)

	tb.InsertString("\n")
	mustBytePos(t, tb, 2, 3)
}

func TestTextBox29(t *testing.T) {
	tb := newStyledTextBox(10, 10)
	tb.InsertString("1234\n\n\n```\n123\nabcd\n```\n\n")
	tb.MoveTextStart()

	tb.MoveDown()
	mustBytePos(t, tb, 1, 0)

	tb.MoveDown()
	mustBytePos(t, tb, 2, 0)

	tb.MoveDown()
	mustBytePos(t, tb, 4, 0)

	tb.MoveRight()
	tb.MoveRight() // 3
	mustBytePos(t, tb, 4, 2)

	tb.MoveRight() // NL
	mustBytePos(t, tb, 4, 3)

	tb.MoveRight() // a
	tb.MoveRight()
	tb.MoveRight()
	tb.MoveRight() // d
	mustBytePos(t, tb, 5, 3)

	tb.MoveRight() // NL
	mustBytePos(t, tb, 5, 4)

	tb.MoveRight()
	mustBytePos(t, tb, 7, 0)
}

func TestTextBox30(t *testing.T) {
	tb := newStyledTextBox(10, 10)
	tb.InsertString("1234\n\n\n```\n123\nabcd\n```\n\n")
	tb.MoveTextStart()

	tb.MoveDown()
	mustBytePos(t, tb, 1, 0)

	tb.MoveDown()
	mustBytePos(t, tb, 2, 0)

	tb.MoveDown()
	mustBytePos(t, tb, 4, 0)

	tb.MoveRight()
	tb.MoveRight() // 3
	mustBytePos(t, tb, 4, 2)

	tb.MoveRight() // NL
	mustBytePos(t, tb, 4, 3)

	tb.MoveRight() // a
	tb.MoveRight()
	tb.MoveRight()
	tb.MoveRight() // d
	mustBytePos(t, tb, 5, 3)

	tb.MoveRight() // NL
	mustBytePos(t, tb, 5, 4)

	tb.MoveRight()
	mustBytePos(t, tb, 7, 0)

	tb.RemoveChar()
	mustBytePos(t, tb, 6, 2)
}

func TestTextBox31(t *testing.T) {
	tb := newStyledTextBox(10, 10)
	tb.InsertString("1234\n\n\n```\n123\nabcd\n```\n\n")
	tb.MoveTextStart()

	tb.MoveDown()
	mustBytePos(t, tb, 1, 0)

	tb.MoveDown()
	mustBytePos(t, tb, 2, 0)

	tb.MoveDown()
	mustBytePos(t, tb, 4, 0)

	tb.RemoveChar()
	mustBytePos(t, tb, 3, 2)
}

func TestTextBox32(t *testing.T) {
	tb := newStyledTextBox(10, 10)

	tb.MoveRight()
	tb.MoveRight()
	tb.InsertString("Hello")
	mustBytePos(t, tb, 2, 5)
}

func TestTextBox33(t *testing.T) {
	tb := newStyledTextBox(10, 10)

	tb.MoveRight()
	tb.MoveRight()
	tb.InsertString("Hello")
	mustBytePos(t, tb, 2, 5)

	tb.MoveRight()
	tb.MoveRight()
	tb.RemoveChar()
	mustBytePos(t, tb, 2, 5)
}

func TestTextBox34(t *testing.T) {
	tb := newStyledTextBox(10, 10)
	tb.InsertString("1\n---\n")
	tb.MoveTextStart()

	tb.MoveRight()
	mustBytePos(t, tb, 0, 1)

	tb.MoveRight()
	mustBytePos(t, tb, 1, 3)

	tb.MoveRight()
	mustBytePos(t, tb, 2, 0)

	tb.RemoveChar()
	mustBytePos(t, tb, 1, 2)
}

func TestTextBox35(t *testing.T) {
	tb := newStyledTextBox(10, 10)
	tb.InsertString("*")
	mustBytePos(t, tb, 0, 1)

	tb.InsertString("a")
	mustBytePos(t, tb, 0, 2)

	tb.InsertString("*")
	mustBytePos(t, tb, 0, 3)
}

func TestTextBox36(t *testing.T) {
	tb := newStyledTextBox(10, 10)
	tb.InsertString("**")
	mustBytePos(t, tb, 0, 2)

	tb.InsertString("aa")
	mustBytePos(t, tb, 0, 4)

	tb.InsertString("*")
	mustBytePos(t, tb, 0, 5)
}

func TestTextBox37(t *testing.T) {
	tb := newStyledTextBox(10, 10)
	tb.InsertString("**")
	mustBytePos(t, tb, 0, 2)

	tb.InsertString("aa")
	mustBytePos(t, tb, 0, 4)

	tb.InsertString("*")
	tb.InsertString("*")
	mustBytePos(t, tb, 0, 6)

	tb.InsertString("\n")
	mustBytePos(t, tb, 1, 0)
}

func TestTextBox38(t *testing.T) {
	tb := newStyledTextBox(10, 10)
	tb.InsertString("**")
	mustBytePos(t, tb, 0, 2)

	tb.InsertString("aa")
	mustBytePos(t, tb, 0, 4)

	tb.InsertString("*")
	tb.InsertString("*")
	mustBytePos(t, tb, 0, 6)

	tb.RemoveChar()
	tb.RemoveChar()
	mustBytePos(t, tb, 0, 0)
}

func TestTextBox39(t *testing.T) {
	tb := newStyledTextBox(10, 10)
	tb.InsertString("abc")
	tb.InsertString(" ")

	tb.InsertString("**")
	mustBytePos(t, tb, 0, 6)

	tb.InsertString("xyz")
	mustBytePos(t, tb, 0, 9)

	tb.InsertString("*")
	tb.InsertString("*")
	mustBytePos(t, tb, 0, 11)
}

func TestTextBox40(t *testing.T) {
	tb := newStyledTextBox(10, 10)
	tb.InsertString("abc")
	tb.InsertString(" ")

	tb.InsertString("**")
	mustBytePos(t, tb, 0, 6)

	tb.InsertString("xyz")
	mustBytePos(t, tb, 0, 9)

	tb.InsertString("*")
	tb.InsertString("*")
	mustBytePos(t, tb, 0, 11)

	tb.RemoveChar()
	mustBytePos(t, tb, 0, 10)
}

func TestTextBox41(t *testing.T) {
	tb := newStyledTextBox(10, 10)
	tb.InsertString("abc")
	tb.InsertString(" ")

	tb.InsertString("**")
	mustBytePos(t, tb, 0, 6)

	tb.InsertString("xyz")
	mustBytePos(t, tb, 0, 9)

	tb.InsertString("*")
	tb.InsertString("*")
	mustBytePos(t, tb, 0, 11)

	tb.RemoveChar()
	mustBytePos(t, tb, 0, 10)

	tb.RemoveChar()
	mustBytePos(t, tb, 0, 9)

	tb.RemoveChar()
	mustBytePos(t, tb, 0, 4)
}

func TestTextBox42(t *testing.T) {
	tb := newStyledTextBox(10, 10)
	tb.InsertString("abc")
	tb.InsertString(" ")

	tb.InsertString("**")
	mustBytePos(t, tb, 0, 6)

	tb.InsertString("xyz")
	mustBytePos(t, tb, 0, 9)

	tb.InsertString("*")
	tb.InsertString("*")
	mustBytePos(t, tb, 0, 11)

	tb.MoveLeft()
	mustBytePos(t, tb, 0, 8)

	tb.InsertString("V")
	mustBytePos(t, tb, 0, 9)
}

func TestTextBox43(t *testing.T) {
	tb := newStyledTextBox(10, 10)
	tb.InsertString("*")
	tb.InsertString("abcd")
	tb.InsertString("*")
	mustBytePos(t, tb, 0, 6)

	tb.RemoveChar()
	mustBytePos(t, tb, 0, 5)
	assert.Equal(t, tb.GetViewPosition(), 3)
}

func TestTextBox44(t *testing.T) {
	tb := newStyledTextBox(10, 10)
	tb.InsertString("```lang")
	tb.InsertString("\n")
	tb.InsertString("abcd")
	tb.InsertString("\n")
	tb.InsertString("```")

	if tb.FindPrev("abcd") {
		tb.MoveRight()
		tb.MoveRight()
		tb.MoveRight()
		tb.MoveRight()
	}

	mustBytePos(t, tb, 1, 4)

	tb.MoveLeft()
	mustBytePos(t, tb, 1, 3)

	tb.MoveLeft()
	mustBytePos(t, tb, 1, 2)

	tb.MoveLeft()
	mustBytePos(t, tb, 1, 1)

	tb.MoveLeft()
	mustBytePos(t, tb, 1, 0)

	tb.MoveLeft()
	mustBytePos(t, tb, 0, 7)

	tb.RemoveChar() // g
	mustBytePos(t, tb, 0, 6)

	tb.RemoveChar() // n
	mustBytePos(t, tb, 0, 5)

	tb.RemoveChar() // a
	mustBytePos(t, tb, 0, 4)

	tb.RemoveChar() // l
	mustBytePos(t, tb, 0, 2)
}

func TestTextBox45(t *testing.T) {
	// tb := TextBox{}
	// tb.Init()
	// tb.X = 0
	// tb.Y = 0
	// tb.Width = 68 + 2
	// tb.Height = 10
	// newDoc := &litemark.StyledDocument{}
	// newDoc.Init()
	// tb.Document = newDoc
	// tb.InsertString("{{{\n")
	// tb.InsertString("私はOpenAIによって開発された大規模言語モデル、ChatGPT（Chat Generative Pre‑trained Transformer）です。質問に答えたり、情報を整理したり、アイデアを提案したりするのが得意です。何か知りたいことや相談したいことがあれば、お気軽にどうぞ！\n")
	// tb.InsertString("}}}")
	//
	// for i := 0; i < 17; i++ {
	// 	tb.MoveLeft()
	// }
	// x, y, _, _ := tb.CursorPosition()
	// y = y - tb.GetScrollY()
	// assert.Equal(t, x, 65)
	// assert.Equal(t, y, 3)
}

func TestTextBox46(t *testing.T) {
	tb := newStyledTextBox(20, 20)
	tb.InsertString("**BOLD**")

	for i := 0; i < 4; i++ {
		tb.MoveLeft()
	}
	tb.InsertString("\n")

	r := model.Range{
		StartPosition: model.Position{
			Row:    0,
			Column: 0,
		},
		EndPosition: model.Position{
			Row:    0,
			Column: tb.Document.GetLineBytes(0),
		},
	}
	sg := tb.Document.Read(r)
	assert.Equal(t, len(sg.GetLine(0)), 0)
}

func TestTextBox47(t *testing.T) {
	tb := newStyledTextBox(20, 20)
	tb.InsertString("\n# \nB")
	mustBytePos(t, tb, 2, 1)
}

func TestTextBox48(t *testing.T) {
	tb := newStyledTextBox(20, 20)
	tb.InsertString("**aa** **bb**")
	mustBytePos(t, tb, 0, 13)
}

func TestTextBox49(t *testing.T) {
	tb := newStyledTextBox(20, 20)
	tb.InsertString("\n**aa** **bb**")
	mustBytePos(t, tb, 1, 13)
}

func TestTextBox50(t *testing.T) {
	tb := newStyledTextBox(20, 20)
	tb.InsertString("**aa** **bb**\n")
	mustBytePos(t, tb, 1, 0)
}

func TestTextBox51(t *testing.T) {
	tb := newStyledTextBox(20, 20)
	tb.InsertString("\n")

	tb.MoveLeft()
	mustBytePos(t, tb, 0, 0)

	tb.InsertString("\n**aa** **bb**")
	mustBytePos(t, tb, 1, 13)
}

func TestTextBox52(t *testing.T) {
	tb := newStyledTextBox(20, 20)
	tb.InsertString("\n")

	tb.MoveLeft()
	mustBytePos(t, tb, 0, 0)

	tb.InsertString("**aa** a\nb **bb**")
	mustBytePos(t, tb, 1, 8)
}

func TestTextBox53(t *testing.T) {
	tb := newStyledTextBox(20, 20)
	tb.InsertString("abc   def")
	mustBytePos(t, tb, 0, 9)

	tb.MoveLeft()
	tb.MoveLeft()
	tb.MoveLeft()
	tb.MoveLeft()
	tb.MoveLeft()
	mustBytePos(t, tb, 0, 4)

	tb.InsertString("**aa**")
	mustBytePos(t, tb, 0, 10)
}

func TestTextBox54(t *testing.T) {
	tb := newStyledTextBox(20, 20)
	tb.InsertString("abcd")
	mustBytePos(t, tb, 0, 4)

	tb.MoveLeft()
	tb.RemoveChar()
}

func TestTextBox55(t *testing.T) {
	tb := newStyledTextBox(20, 20)
	tb.InsertString("abcd")
	tb.SelectionStart()
	tb.MoveLeft()
	tb.MoveLeft()

	ts := tb.GetSelection()
	first, last := ts.Ordered()
	r := model.Range{
		StartPosition: first.StartPosition,
		EndPosition:   last.StartPosition,
	}
	sg := tb.Document.Read(r)
	assert.Equal(t, sg.GetLine(0), "cd")

	tb.RemoveSelection()
	tb.SelectionEnd()
	mustBytePos(t, tb, 0, 2)
	assert.Equal(t, tb.GetViewPosition(), 2)
}

func TestTextBox56(t *testing.T) {
	tb := newStyledTextBox(20, 20)
	tb.InsertString("abcd\nefgh\nijkl")
	tb.MoveLeft()
	tb.MoveLeft()
	tb.SelectionStart()
	tb.MoveUp()
	tb.MoveUp()

	ts := tb.GetSelection()
	first, last := ts.Ordered()
	r := model.Range{
		StartPosition: first.StartPosition,
		EndPosition:   last.StartPosition,
	}
	sg := tb.Document.Read(r)
	assert.Equal(t, sg.GetLine(0), "cd")
	assert.Equal(t, sg.GetLine(1), "efgh")
	assert.Equal(t, sg.GetLine(2), "ij")

	tb.RemoveSelection()
	tb.SelectionEnd()
	mustBytePos(t, tb, 0, 2)
	assert.Equal(t, tb.GetViewPosition(), 2)
}

func TestTextBox57(t *testing.T) {
	tb := newStyledTextBox(20, 20)
	tb.InsertString("あいうeお\nかきくけこ\nijkl")
	tb.MoveLeft()
	tb.MoveLeft()
	tb.SelectionStart()
	tb.MoveUp()
	tb.MoveUp()

	ts := tb.GetSelection()
	first, last := ts.Ordered()
	r := model.Range{
		StartPosition: first.StartPosition,
		EndPosition:   last.StartPosition,
	}
	sg := tb.Document.Read(r)
	assert.Equal(t, sg.GetLine(0), "うeお")
	assert.Equal(t, sg.GetLine(1), "かきくけこ")
	assert.Equal(t, sg.GetLine(2), "ij")

	tb.RemoveSelection()
	tb.SelectionEnd()
	mustBytePos(t, tb, 0, len("あい"))
	assert.Equal(t, tb.GetViewPosition(), 2)
}

func TestTextBox58(t *testing.T) {
	tb := newStyledTextBox(20, 20)
	tb.Document.ReplaceAll(strings.NewReader("{{{\n```\nabc\n```\n\n}}}"))
	tb.MoveTextEnd()
	tb.RemoveChar()

	mustBytePos(t, tb, 3, 2)
}

func TestTextBox59(t *testing.T) {
	tb := newStyledTextBox(20, 20)
	tb.Document.ReplaceAll(strings.NewReader("{{{\nabc\n\n}}}"))
	tb.MoveTextStart()
	tb.RemoveChar()

	mustBytePos(t, tb, 0, 2)
}

func TestTextBox60(t *testing.T) {
	tb := newStyledTextBox(20, 20)
	tb.Document.ReplaceAll(strings.NewReader("{{{\nabc\n}}}\n\n"))
	tb.MoveTextEnd()
	tb.RemoveChar()

	mustBytePos(t, tb, 2, 2)
}

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
