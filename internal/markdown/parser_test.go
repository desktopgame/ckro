package markdown_test

import (
	"testing"

	"github.com/desktopgame/ckro/internal/markdown"
	"github.com/stretchr/testify/assert"
)

func TestList01(t *testing.T) {
	text := `
* item1
* item2
* item3
`
	doc := markdown.Parse(text)
	assert.Equal(t, len(doc.Blocks), 1)

	root := doc.Blocks[0]
	if lb, ok := root.(*markdown.ListBlock); ok {
		assert.Equal(t, len(lb.Items), 3)
	} else {
		assert.True(t, false)
	}
}
