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
	assert.Equal(t, litemark.GetText(&r, h1.LineIndex, h1.Span), "H1")

	h2 := blocks[1].(*litemark.Heading)
	assert.Equal(t, h2.Level, 2)
	assert.Equal(t, litemark.GetText(&r, h2.LineIndex, h2.Span), "H2")

	h3 := blocks[2].(*litemark.Heading)
	assert.Equal(t, h3.Level, 3)
	assert.Equal(t, litemark.GetText(&r, h3.LineIndex, h3.Span), "H3")
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
	assert.Equal(t, p1.Spans[0].StartColumn, 0)
	assert.Equal(t, p1.Spans[0].EndColumn, 13)

	it := tx.Inlines[1].(*litemark.Italic)
	assert.Equal(t, it.Spans[0].StartColumn, 13)
	assert.Equal(t, it.Spans[0].EndColumn, 23)

	p2 := tx.Inlines[2].(*litemark.PlainText)
	assert.Equal(t, p2.Spans[0].StartColumn, 23)
	assert.Equal(t, p2.Spans[0].EndColumn, 44)

	bd := tx.Inlines[3].(*litemark.Bold)
	assert.Equal(t, bd.Spans[0].StartColumn, 44)
	assert.Equal(t, bd.Spans[0].EndColumn, 56)

	p3 := tx.Inlines[4].(*litemark.PlainText)
	assert.Equal(t, p3.Spans[0].StartColumn, 56)
	assert.Equal(t, p3.Spans[0].EndColumn, 57)
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
	assert.Equal(t, link.Spans[0].StartColumn, 8)
	assert.Equal(t, link.Spans[0].EndColumn, 8+37)

	image := tx.Inlines[3].(*litemark.Image)
	assert.Equal(t, image.Spans[0].StartColumn, 55)
	assert.Equal(t, image.Spans[0].EndColumn, 55+19)
}
