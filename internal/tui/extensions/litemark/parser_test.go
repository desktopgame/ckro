package litemark_test

import (
	"strings"
	"testing"

	"github.com/desktopgame/ckro/internal/tui/extensions/litemark"
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
	assert.Equal(t, litemark.GetText(&r, tx.LineIndex, p1.Spans[0]), "this text is ")

	it := tx.Inlines[1].(*litemark.Italic)
	assert.Equal(t, it.Spans[0].StartColumn, 13)
	assert.Equal(t, it.Spans[0].EndColumn, 23)
	assert.Equal(t, litemark.GetText(&r, tx.LineIndex, it.Spans[0]), "*litemark*")
	assert.Equal(t, litemark.GetText(&r, tx.LineIndex, it.Spans[1]), "litemark")

	p2 := tx.Inlines[2].(*litemark.PlainText)
	assert.Equal(t, p2.Spans[0].StartColumn, 23)
	assert.Equal(t, p2.Spans[0].EndColumn, 44)
	assert.Equal(t, litemark.GetText(&r, tx.LineIndex, p2.Spans[0]), ", this is dialect of ")

	bd := tx.Inlines[3].(*litemark.Bold)
	assert.Equal(t, bd.Spans[0].StartColumn, 44)
	assert.Equal(t, bd.Spans[0].EndColumn, 56)
	assert.Equal(t, litemark.GetText(&r, tx.LineIndex, bd.Spans[0]), "**markdown**")
	assert.Equal(t, litemark.GetText(&r, tx.LineIndex, bd.Spans[1]), "markdown")

	p3 := tx.Inlines[4].(*litemark.PlainText)
	assert.Equal(t, p3.Spans[0].StartColumn, 56)
	assert.Equal(t, p3.Spans[0].EndColumn, 57)
	assert.Equal(t, litemark.GetText(&r, tx.LineIndex, p3.Spans[0]), ".")
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
	assert.Equal(t, litemark.GetText(&r, tx.LineIndex, link.Spans[1]), "link")
	assert.Equal(t, litemark.GetText(&r, tx.LineIndex, link.Spans[2]), "https://www.google.com/?hl=ja")

	image := tx.Inlines[3].(*litemark.Image)
	assert.Equal(t, image.Spans[0].StartColumn, 55)
	assert.Equal(t, image.Spans[0].EndColumn, 55+19)
	assert.Equal(t, litemark.GetText(&r, tx.LineIndex, image.Spans[1]), "image")
	assert.Equal(t, litemark.GetText(&r, tx.LineIndex, image.Spans[2]), "image.png")
}

func Test05(t *testing.T) {
	lines := []string{
		"| Header 1 | Header 2 |",
		"|----------|----------|",
		"| Cell 1   | Cell 2   |",
		"| Cell 3   | Cell 4   |",
	}

	r := litemark.StringReader{
		Source: lines,
	}
	blocks := litemark.Parse(&r)

	tbl := blocks[0].(*litemark.Table)
	assert.Equal(t, tbl.LineIndex, 0)
	assert.Equal(t, tbl.LineCount, 4)

	h1 := tbl.Headers[0]
	h1Span := h1.Inlines[0].BaseInline().Spans[0]
	assert.Equal(t, h1Span.StartColumn, 1)
	assert.Equal(t, h1Span.EndColumn, 11)

	h2 := tbl.Headers[1]
	h2Span := h2.Inlines[0].BaseInline().Spans[0]
	assert.Equal(t, h2Span.StartColumn, 12)
	assert.Equal(t, h2Span.EndColumn, 22)

	row1 := tbl.Rows[0]
	row1col1 := row1.Columns[0]
	row1col1Span := row1col1.Inlines[0].BaseInline().Spans[0]
	assert.Equal(t, row1col1Span.StartColumn, 1)
	assert.Equal(t, row1col1Span.EndColumn, 11)
}
