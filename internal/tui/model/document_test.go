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
