package tui

import (
	"testing"

	"github.com/desktopgame/ckro/internal/tui/extensions/litemark"
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/desktopgame/ckro/internal/tui/view"
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

	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, len("1234あ"))
}

func TestTextBox02(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	tb.InsertString("1234\n1234\n123あ")

	assert.Equal(t, tb.bytePos.StartPosition.Row, 2)
	assert.Equal(t, tb.bytePos.StartPosition.Column, len("123あ"))
}

func TestTextBox03(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	tb.InsertString("1234\n1234\n123あ")
	tb.bytePos = view.CharacterReference{
		StartPosition: model.Position{
			Row:    0,
			Column: 0,
		},
		Bytes: 1,
	}
	tb.viewPosition = 0

	tb.MoveRight()

	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, len("1"))
}

func TestTextBox04(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	tb.InsertString("1234\n1234\n123あ")
	tb.bytePos = view.CharacterReference{
		StartPosition: model.Position{
			Row:    0,
			Column: 0,
		},
		Bytes: 1,
	}
	tb.viewPosition = 0

	tb.MoveRight() // 2
	tb.MoveRight() // 3
	tb.MoveRight() // 4
	tb.MoveRight() // NL

	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, len("1234"))
}

func TestTextBox05(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	tb.InsertString("1234\n1234\n123あ")
	tb.bytePos = view.CharacterReference{
		StartPosition: model.Position{
			Row:    0,
			Column: 0,
		},
		Bytes: 1,
	}
	tb.viewPosition = 0

	tb.MoveRight() // 2
	tb.MoveRight() // 3
	tb.MoveRight() // 4
	tb.MoveRight() // NL

	tb.InsertString("\n")

	assert.Equal(t, tb.bytePos.StartPosition.Row, 1)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)
}

func TestTextBox06(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	tb.InsertString("#")
	tb.InsertString(" ")

	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 2)

	tb.InsertString("AAA")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 5)
}

func TestTextBox07(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	tb.InsertString("1234\n\n\n```\n123\n```")
	tb.bytePos = view.CharacterReference{
		StartPosition: model.Position{
			Row:    0,
			Column: 0,
		},
		Bytes: 1,
	}
	tb.viewPosition = 0

	tb.InsertString("\n")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 1)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)
}

func TestTextBox08(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	tb.InsertString("\n\n\n\n\n")

	assert.Equal(t, tb.bytePos.StartPosition.Row, 5)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)

	tb.MoveUp()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 4)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)
}

func TestTextBox09(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	tb.InsertString("\n\n\n\n\n")

	assert.Equal(t, tb.bytePos.StartPosition.Row, 5)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)

	tb.MoveLeft()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 4)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)
}

func TestTextBox10(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	tb.InsertString("\n\n\n\n\n")

	assert.Equal(t, tb.bytePos.StartPosition.Row, 5)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)

	tb.MoveLeft()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 4)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)

	tb.MoveRight()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 5)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)
}

func TestTextBox11(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10

	tb.InsertString("a")
	tb.InsertString("\n")

	tb.InsertString("a")
	tb.InsertString("\n")

	tb.InsertString("a")

	tb.MoveUp()

	assert.Equal(t, tb.bytePos.StartPosition.Row, 1)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 1)
}

func TestTextBox12(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10

	tb.InsertString("a")
	tb.InsertString("\n")

	tb.InsertString("a")
	tb.InsertString("\n")

	tb.InsertString("a")

	tb.MoveUp()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 1)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 1)

	tb.MoveDown()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 2)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 1)
}

func TestTextBox13(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10

	tb.InsertString("a")
	tb.InsertString("\n")

	tb.InsertString("a")
	tb.InsertString("\n")

	tb.InsertString("a")

	tb.MoveDown()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 2)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 1)
}

func TestTextBox14(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10

	tb.InsertString("a")
	tb.InsertString("\n")

	tb.InsertString("a")
	tb.InsertString("\n")

	tb.InsertString("a")

	tb.MoveUp()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 1)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 1)

	tb.MoveUp()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 1)

	tb.MoveUp()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 1)
}

func TestTextBox15(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc
	tb.InsertString("1234\n\n\n```\n123\n```\n\n")
	tb.bytePos = view.CharacterReference{
		StartPosition: model.Position{
			Row:    0,
			Column: 0,
		},
		Bytes: 1,
	}
	tb.viewPosition = 0

	tb.MoveDown()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 1)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)

	tb.MoveDown()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 2)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)

	tb.MoveDown()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 4)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)

	tb.MoveDown()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 6)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)
}

func TestTextBox16(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc
	tb.InsertString("1234\n\n\n```\n123\n```\n\n")
	tb.bytePos = view.CharacterReference{
		StartPosition: model.Position{
			Row:    0,
			Column: 0,
		},
		Bytes: 1,
	}
	tb.viewPosition = 0

	tb.MoveDown()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 1)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)

	tb.MoveDown()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 2)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)

	tb.MoveDown()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 4)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)

	tb.MoveRight()
	tb.MoveRight()
	tb.MoveRight()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 4)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 3)

	tb.InsertString("\n")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 5)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)
}

func TestTextBox17(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc
	tb.InsertString("1234\n\n\n```\n123\n```\n\n")
	tb.bytePos = view.CharacterReference{
		StartPosition: model.Position{
			Row:    0,
			Column: 0,
		},
		Bytes: 1,
	}
	tb.viewPosition = 0

	tb.MoveDown()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 1)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)

	tb.MoveDown()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 2)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)

	tb.MoveDown()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 4)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)

	tb.MoveRight()
	tb.MoveRight()
	tb.MoveRight()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 4)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 3)

	tb.InsertString("\n")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 5)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)

	tb.InsertString("\n")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 6)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)
}

func TestTextBox18(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc
	tb.InsertString("1234")

	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 4)

	tb.RemoveChar()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 3)
}

func TestTextBox19(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc
	tb.InsertString("\n\n\n")

	assert.Equal(t, tb.bytePos.StartPosition.Row, 3)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)

	tb.RemoveChar()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 2)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)
}

func TestTextBox20(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc
	tb.InsertString("#")
	tb.InsertString(" ")
	tb.InsertString("H")

	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 3)

	tb.RemoveChar()
}

func TestTextBox21(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc
	tb.InsertString("#")
	tb.InsertString(" ")
	tb.InsertString("H")

	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 3)

	tb.MoveLeft()
	tb.MoveLeft()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 2)

	tb.MoveRight()
	tb.MoveRight()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 3)
}

func TestTextBox22(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc
	tb.InsertString("\n")
	tb.InsertString("#")
	tb.InsertString(" ")
	tb.InsertString("H")

	assert.Equal(t, tb.bytePos.StartPosition.Row, 1)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 3)

	tb.RemoveChar()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 1)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)

	tb.RemoveChar()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)
}

func TestTextBox23(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc
	tb.InsertString("\n")
	tb.InsertString("#")
	tb.InsertString(" ")
	tb.InsertString("Heading")
	tb.InsertString("\n")
	tb.InsertString("\n")
	tb.InsertString("#")
	tb.InsertString(" ")
	tb.InsertString("Heading2")

	assert.Equal(t, tb.bytePos.StartPosition.Row, 3)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 10)

	tb.MoveLeft()
	tb.MoveLeft()
	tb.MoveLeft()
	tb.MoveLeft()
	tb.MoveLeft()
	tb.MoveLeft()
	tb.MoveLeft()
	tb.MoveLeft() // H
	assert.Equal(t, tb.bytePos.StartPosition.Row, 3)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 2)

	tb.RemoveChar()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 2)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)
}

func TestTextBox24(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 20
	tb.Height = 20
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc
	tb.InsertString("\n")
	tb.InsertString("#")
	tb.InsertString(" ")
	tb.InsertString("Heading")
	tb.InsertString("\n")
	tb.InsertString("\n")
	tb.InsertString("##")
	tb.InsertString(" ")
	tb.InsertString("Heading2")

	assert.Equal(t, tb.bytePos.StartPosition.Row, 3)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 11)

	tb.MoveLeft()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 3)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 10)

	tb.MoveLeft()
	tb.MoveLeft()
	tb.MoveLeft()
	tb.MoveLeft()
	tb.MoveLeft()
	tb.MoveLeft()
	tb.MoveLeft() // H
	assert.Equal(t, tb.bytePos.StartPosition.Row, 3)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 3)

	tb.RemoveChar()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 2)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)
}

func TestTextBox25(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 20
	tb.Height = 20
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc
	tb.InsertString("#")
	tb.InsertString(" ")
	tb.InsertString("H")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 3)

	tb.InsertString("e")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 4)
}

func TestTextBox26(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 20
	tb.Height = 20
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc
	tb.InsertString("\n")
	tb.MoveLeft()

	tb.InsertString("##")
	tb.InsertString(" ")
	tb.InsertString("H")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 4)

	tb.InsertString("H")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 5)
}

func TestTextBox27(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 20
	tb.Height = 20
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc

	tb.InsertString("##")
	tb.InsertString(" ")
	tb.InsertString("H")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 4)

	tb.InsertString("H")
	tb.InsertString("H")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 6)

	tb.MoveLeft()
	tb.MoveLeft()
	tb.MoveLeft()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 3)

	tb.InsertString("\n")
}

func TestTextBox28(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 20
	tb.Height = 20
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc

	tb.InsertString("##")
	tb.InsertString(" ")
	tb.InsertString("H")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 4)

	tb.InsertString("H")
	tb.InsertString("H")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 6)

	tb.MoveLeft()
	tb.MoveLeft()
	tb.MoveLeft()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 3)

	tb.InsertString("\n")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 1)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 3)

	tb.InsertString("\n")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 2)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 3)
}

func TestTextBox29(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc
	tb.InsertString("1234\n\n\n```\n123\nabcd\n```\n\n")
	tb.bytePos = view.CharacterReference{
		StartPosition: model.Position{
			Row:    0,
			Column: 0,
		},
		Bytes: 1,
	}
	tb.viewPosition = 0

	tb.MoveDown()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 1)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)

	tb.MoveDown()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 2)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)

	tb.MoveDown()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 4)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)

	tb.MoveRight()
	tb.MoveRight() // 3
	assert.Equal(t, tb.bytePos.StartPosition.Row, 4)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 2)

	tb.MoveRight() // NL
	assert.Equal(t, tb.bytePos.StartPosition.Row, 4)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 3)

	tb.MoveRight() // a
	tb.MoveRight()
	tb.MoveRight()
	tb.MoveRight() // d
	assert.Equal(t, tb.bytePos.StartPosition.Row, 5)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 3)

	tb.MoveRight() // NL
	assert.Equal(t, tb.bytePos.StartPosition.Row, 5)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 4)

	tb.MoveRight()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 7)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)
}

func TestTextBox30(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc
	tb.InsertString("1234\n\n\n```\n123\nabcd\n```\n\n")
	tb.bytePos = view.CharacterReference{
		StartPosition: model.Position{
			Row:    0,
			Column: 0,
		},
		Bytes: 1,
	}
	tb.viewPosition = 0

	tb.MoveDown()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 1)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)

	tb.MoveDown()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 2)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)

	tb.MoveDown()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 4)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)

	tb.MoveRight()
	tb.MoveRight() // 3
	assert.Equal(t, tb.bytePos.StartPosition.Row, 4)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 2)

	tb.MoveRight() // NL
	assert.Equal(t, tb.bytePos.StartPosition.Row, 4)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 3)

	tb.MoveRight() // a
	tb.MoveRight()
	tb.MoveRight()
	tb.MoveRight() // d
	assert.Equal(t, tb.bytePos.StartPosition.Row, 5)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 3)

	tb.MoveRight() // NL
	assert.Equal(t, tb.bytePos.StartPosition.Row, 5)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 4)

	tb.MoveRight()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 7)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)

	tb.RemoveChar()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 6)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 2)
}

func TestTextBox31(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc
	tb.InsertString("1234\n\n\n```\n123\nabcd\n```\n\n")
	tb.bytePos = view.CharacterReference{
		StartPosition: model.Position{
			Row:    0,
			Column: 0,
		},
		Bytes: 1,
	}
	tb.viewPosition = 0

	tb.MoveDown()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 1)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)

	tb.MoveDown()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 2)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)

	tb.MoveDown()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 4)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)

	tb.RemoveChar()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 3)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 2)
}

func TestTextBox32(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc

	tb.MoveRight()
	tb.MoveRight()
	tb.InsertString("Hello")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 2)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 5)
}

func TestTextBox33(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc

	tb.MoveRight()
	tb.MoveRight()
	tb.InsertString("Hello")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 2)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 5)

	tb.MoveRight()
	tb.MoveRight()
	tb.RemoveChar()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 2)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 5)
}

func TestTextBox34(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc
	tb.InsertString("1\n---\n")
	tb.bytePos = view.CharacterReference{
		StartPosition: model.Position{
			Row:    0,
			Column: 0,
		},
		Bytes: 1,
	}
	tb.viewPosition = 0

	tb.MoveRight()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 1)

	tb.MoveRight()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 1)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 3)

	tb.MoveRight()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 2)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)

	tb.RemoveChar()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 1)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 2)
}

func TestTextBox35(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc
	tb.InsertString("*")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 1)

	tb.InsertString("a")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 2)

	tb.InsertString("*")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 3)
}

func TestTextBox36(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc
	tb.InsertString("**")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 2)

	tb.InsertString("aa")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 4)

	tb.InsertString("*")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 5)
}

func TestTextBox37(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc
	tb.InsertString("**")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 2)

	tb.InsertString("aa")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 4)

	tb.InsertString("*")
	tb.InsertString("*")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 6)

	tb.InsertString("\n")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 1)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)
}

func TestTextBox38(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc
	tb.InsertString("**")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 2)

	tb.InsertString("aa")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 4)

	tb.InsertString("*")
	tb.InsertString("*")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 6)

	tb.RemoveChar()
	tb.RemoveChar()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)
}

func TestTextBox39(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc
	tb.InsertString("abc")
	tb.InsertString(" ")

	tb.InsertString("**")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 6)

	tb.InsertString("xyz")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 9)

	tb.InsertString("*")
	tb.InsertString("*")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 11)
}

func TestTextBox40(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc
	tb.InsertString("abc")
	tb.InsertString(" ")

	tb.InsertString("**")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 6)

	tb.InsertString("xyz")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 9)

	tb.InsertString("*")
	tb.InsertString("*")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 11)

	tb.RemoveChar()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 8)
}

func TestTextBox41(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc
	tb.InsertString("abc")
	tb.InsertString(" ")

	tb.InsertString("**")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 6)

	tb.InsertString("xyz")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 9)

	tb.InsertString("*")
	tb.InsertString("*")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 11)

	tb.RemoveChar()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 8)

	tb.RemoveChar()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 7)

	tb.RemoveChar()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 4)
}

func TestTextBox42(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc
	tb.InsertString("abc")
	tb.InsertString(" ")

	tb.InsertString("**")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 6)

	tb.InsertString("xyz")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 9)

	tb.InsertString("*")
	tb.InsertString("*")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 11)

	tb.MoveLeft()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 8)

	tb.InsertString("V")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 9)
}

func TestTextBox43(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc
	tb.InsertString("*")
	tb.InsertString("abcd")
	tb.InsertString("*")

	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 6)

	tb.RemoveChar()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 4)
	assert.Equal(t, tb.viewPosition, 3)
}

func TestTextBox44(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc
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

	assert.Equal(t, tb.bytePos.StartPosition.Row, 1)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 4)

	tb.MoveLeft()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 1)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 3)

	tb.MoveLeft()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 1)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 2)

	tb.MoveLeft()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 1)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 1)

	tb.MoveLeft()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 1)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)

	tb.MoveLeft()
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 7)

	tb.RemoveChar() // g
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 6)

	tb.RemoveChar() // n
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 5)

	tb.RemoveChar() // a
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 4)

	tb.RemoveChar() // l
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 2)
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
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 68 + 2
	tb.Height = 10
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc
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
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 68 + 2
	tb.Height = 10
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc
	tb.InsertString("\n# \nB")

	assert.Equal(t, tb.bytePos.StartPosition.Row, 2)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 1)
}

func TestTextBoxFind01(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc
	tb.InsertString("aa\nbb\nccc")

	tb.FindPrev("b\ncc")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 1)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 1)

	tb.FindPrev("a\nb")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 1)
	assert.Equal(t, tb.bytePos.Bytes, 1)

	tb.MoveLineEnd()
	tb.MoveDown()
	tb.MoveLineEnd()
	tb.MoveDown()
	tb.MoveLineEnd()
	tb.MoveDown()
	tb.FindPrev("a\nbb\nccc")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 1)
	assert.Equal(t, tb.bytePos.Bytes, 1)
}

func TestTextBoxFind02(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc
	tb.InsertString("aa\nbb\nccc")

	tb.FindPrev("cc")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 2)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 1)
	assert.Equal(t, tb.bytePos.Bytes, 1)

	tb.FindPrev("bb")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 1)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)
	assert.Equal(t, tb.bytePos.Bytes, 1)
}

func TestTextBoxFind03(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc
	tb.InsertString("aa\nbb\nccc")

	tb.MoveReset()

	tb.FindNext("a\nbb")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 1)
	assert.Equal(t, tb.bytePos.Bytes, 1)

	tb.FindNext("bb")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 1)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)
	assert.Equal(t, tb.bytePos.Bytes, 1)

	tb.FindNext("b\nc")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 1)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 1)
	assert.Equal(t, tb.bytePos.Bytes, 1)
}

func TestTextBoxFind04(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc
	tb.InsertString("あいう\nbb\nかきく")

	tb.MoveReset()

	tb.FindNext("あいう")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)
	assert.Equal(t, tb.bytePos.Bytes, len("あ"))

	tb.FindNext("う\nbb")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, len("あい"))
	assert.Equal(t, tb.bytePos.Bytes, len("う"))

	assert.False(t, tb.FindNext("う"))

	tb.FindNext("か")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 2)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)
	assert.Equal(t, tb.bytePos.Bytes, len("か"))
}

func TestTextBoxFind05(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc
	tb.InsertString("あいう\nbb\nかきく")

	tb.FindPrev("bb\nか")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 1)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)
	assert.Equal(t, tb.bytePos.Bytes, 1)

	assert.False(t, tb.FindPrev("あいう\nb"))

	tb.FindPrev("あいう\n")
	assert.Equal(t, tb.bytePos.StartPosition.Row, 0)
	assert.Equal(t, tb.bytePos.StartPosition.Column, 0)
	assert.Equal(t, tb.bytePos.Bytes, len("あ"))
}

func TestTextBoxReplace01(t *testing.T) {
	tb := TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 10
	tb.Height = 10
	newDoc := &litemark.StyledDocument{}
	newDoc.Init()
	tb.Document = newDoc
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
