package litemark_test

import (
	"strings"
	"testing"

	"github.com/desktopgame/ckro/internal/litemark"
	"github.com/stretchr/testify/assert"
)

func Test01(t *testing.T) {
	text := `
# H1
## H2
### H3
`
	text = strings.Trim(text, " \t\n")

	r := litemark.StringReader{
		Source: strings.Split(text, "\n"),
	}
	blocks := litemark.Parse(&r)

	h1 := blocks[0].(*litemark.Heading)
	assert.Equal(t, h1.Level, 1)

	h2 := blocks[1].(*litemark.Heading)
	assert.Equal(t, h2.Level, 2)

	h3 := blocks[2].(*litemark.Heading)
	assert.Equal(t, h3.Level, 3)
}

func Test02(t *testing.T) {
	lines := []string{
		"```lang",
		"L1",
		"L2",
		"L3",
		"```",
	}

	r := litemark.StringReader{
		Source: lines,
	}
	blocks := litemark.Parse(&r)

	cb := blocks[0].(*litemark.CodeBlock)
	assert.Equal(t, cb.LineIndex, 0)
	assert.Equal(t, cb.LineCount, 5)
}

func Test03(t *testing.T) {
	text := `
this text is *litemark*, this is dialect of **markdown**.
`
	text = strings.Trim(text, " \t\n")

	r := litemark.StringReader{
		Source: strings.Split(text, "\n"),
	}
	blocks := litemark.Parse(&r)

	tx := blocks[0].(*litemark.Text)

	p1 := tx.Inlines[0].(*litemark.PlainText)
	assert.Equal(t, p1.StartColumn, 0)
	assert.Equal(t, p1.EndColumn, 13)

	it := tx.Inlines[1].(*litemark.Italic)
	assert.Equal(t, it.StartColumn, 13)
	assert.Equal(t, it.EndColumn, 23)

	p2 := tx.Inlines[2].(*litemark.PlainText)
	assert.Equal(t, p2.StartColumn, 23)
	assert.Equal(t, p2.EndColumn, 44)

	bd := tx.Inlines[3].(*litemark.Bold)
	assert.Equal(t, bd.StartColumn, 44)
	assert.Equal(t, bd.EndColumn, 56)

	p3 := tx.Inlines[4].(*litemark.PlainText)
	assert.Equal(t, p3.StartColumn, 56)
	assert.Equal(t, p3.EndColumn, 57)
}

func Test04(t *testing.T) {
	text := `
this is [link](https://www.google.com/?hl=ja), this is ![image](image.png)
`
	text = strings.Trim(text, " \t\n")

	r := litemark.StringReader{
		Source: strings.Split(text, "\n"),
	}
	blocks := litemark.Parse(&r)

	tx := blocks[0].(*litemark.Text)

	link := tx.Inlines[1].(*litemark.Link)
	assert.Equal(t, link.StartColumn, 8)
	assert.Equal(t, link.EndColumn, 8+37)

	image := tx.Inlines[3].(*litemark.Image)
	assert.Equal(t, image.StartColumn, 55)
	assert.Equal(t, image.EndColumn, 55+19)
}
