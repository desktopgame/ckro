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

func TestInsert(t *testing.T) {
	buf := model.Buffer{}
	buf.Init()
	buf.InsertString(0, 0, "123456789")
	at, err := buf.InsertString(0, 5, "Hello\nHello")
	assert.Nil(t, err)
	assert.Equal(t, at.Row, 1)
	assert.Equal(t, at.Column, 5)
}

func TestInsertEmptyLine(t *testing.T) {
	buf := model.Buffer{}
	buf.Init()
	buf.InsertString(0, 0, "123456789")
	at, err := buf.InsertString(0, 5, "Hello\n\n")
	assert.Nil(t, err)
	assert.Equal(t, at.Row, 2)
	assert.Equal(t, at.Column, 0)
}

func TestTrack01(t *testing.T) {
	buf := model.Buffer{}
	buf.Init()

	tr := buf.CreateTrack(0, 0)
	assert.Equal(t, tr.Position.Row, 0)
	assert.Equal(t, tr.Position.Column, 0)

	buf.InsertString(0, 0, "123456789")
	assert.Equal(t, tr.Position.Row, 0)
	assert.Equal(t, tr.Position.Column, 9)

	buf.RemoveString(0, 0, 9)
	assert.Equal(t, tr.Lost, false)
}

func TestTrack02(t *testing.T) {
	buf := model.Buffer{}
	buf.Init()

	buf.InsertString(0, 0, "123456789")

	tr := buf.CreateTrack(0, 1)

	buf.RemoveString(0, 2, 7)
	assert.Equal(t, tr.Position.Row, 0)
	assert.Equal(t, tr.Position.Column, 1)

	buf.InsertString(0, 0, "ABC")
	assert.Equal(t, tr.Position.Row, 0)
	assert.Equal(t, tr.Position.Column, 4)
}

func TestTrack03(t *testing.T) {
	buf := model.Buffer{}
	buf.Init()

	buf.InsertString(0, 0, "123456789")

	tr := buf.CreateTrack(0, 1)

	buf.RemoveString(0, 2, 7)
	assert.Equal(t, tr.Position.Row, 0)
	assert.Equal(t, tr.Position.Column, 1)

	buf.InsertString(0, 0, "ABC\n")
	assert.Equal(t, tr.Position.Row, 1)
	assert.Equal(t, tr.Position.Column, 1)
}

func TestTrack04(t *testing.T) {
	buf := model.Buffer{}
	buf.Init()

	buf.InsertString(0, 0, "123456789")

	tr := buf.CreateTrack(0, 1)

	buf.RemoveString(0, 2, 7)
	assert.Equal(t, tr.Position.Row, 0)
	assert.Equal(t, tr.Position.Column, 1)

	buf.InsertString(0, 0, "ABC\nAAA")
	assert.Equal(t, tr.Position.Row, 1)
	assert.Equal(t, tr.Position.Column, 4)
}

func TestTrack05(t *testing.T) {
	buf := model.Buffer{}
	buf.Init()

	buf.InsertString(0, 0, "123456789")

	tr := buf.CreateTrack(0, 1)

	buf.RemoveString(0, 2, 7)
	assert.Equal(t, tr.Position.Row, 0)
	assert.Equal(t, tr.Position.Column, 1)

	buf.InsertString(0, 0, "ABC\nAAA\nAAA")
	assert.Equal(t, tr.Position.Row, 2)
	assert.Equal(t, tr.Position.Column, 4)
}
