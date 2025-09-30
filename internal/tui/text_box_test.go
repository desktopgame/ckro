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
	assert.Equal(t, tb.bytePos.StartPosition.Column, 2)

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
