package model_test

import (
	"testing"

	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/stretchr/testify/assert"
)

func TestDocument(t *testing.T) {
	doc := model.Document{}
	doc.Init()

	doc.InsertString("Hello")
	cursorColumn := doc.GetCursorColumn()
	assert.Equal(t, cursorColumn, 5)

	doc.RemoveChar()
	cursorColumn = doc.GetCursorColumn()
	assert.Equal(t, cursorColumn, 4)
}

func TestReplace(t *testing.T) {
	doc := model.Document{}
	doc.Init()

	doc.InsertString("Hello1")
	doc.InsertLine()

	doc.InsertString("Hello2")
	doc.InsertLine()

	doc.InsertString("Hello3")
	doc.InsertLine()

	doc.MoveReset()
	assert.True(t, doc.FindNext("Hello2"))
	assert.True(t, doc.Replace(6, "Hello\nNewLine\nNewLine"))
}
