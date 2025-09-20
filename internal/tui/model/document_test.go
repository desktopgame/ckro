package model_test

import (
	"strings"
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

func TestTooManyLines(t *testing.T) {
	doc := model.Document{}
	doc.Init()

	firstLine := "こんにちは！今日はどんなご用件でしょうか？😊"
	doc.InsertString(firstLine)
	doc.InsertLine()

	content := `以下は、デバッグやテストに使えるサンプルとして、複数行にわたるテキストです。  
（必要に応じて内容を調整してください）

"""
Line 1: This is the first line of a multi-line text sample.
Line 2: Here we add some more content, perhaps with numbers like 12345 or symbols #!$.
Line 3: The third line can contain a short sentence, e.g., "Debugging mode activated."
Line 4: In this line we might include a JSON snippet for structure testing:
{
    "key1": "value1",
    "key2": [1, 2, 3],
    "nested": {
        "innerKey": "innerValue"
    }
}
Line 5: Finally, end with a closing statement or marker like END_OF_TEXT.
"""

これでご要望に合っているでしょうか？必要があればさらに長くしたり、特定の形式（例えば YAML や CSV）に変換してお渡しすることも可能です。
`

	doc.InsertString(content)

	lines := strings.Split(firstLine+"\n"+content, "\n")
	buf := doc.GetBuffer()
	assert.Equal(t, lines[0], firstLine)
	assert.Equal(t, lines[0], buf.GetLineAt(0).GetContent())
	assert.Equal(t, lines[1], "以下は、デバッグやテストに使えるサンプルとして、複数行にわたるテキストです。  ")
	assert.Equal(t, lines[1], buf.GetLineAt(1).GetContent())

	for i := 0; i < buf.GetLineCount(); i++ {
		assert.Equal(t, lines[i], buf.GetLineAt(i).GetContent(), "%d", i)
	}
}
