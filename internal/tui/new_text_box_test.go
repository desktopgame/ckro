package tui_test

import (
	"bufio"
	"io"
	"slices"
	"strings"
	"testing"

	"github.com/desktopgame/ckro/internal/text"
	"github.com/desktopgame/ckro/internal/tui"
	"github.com/desktopgame/ckro/internal/tui/base"
	"github.com/desktopgame/ckro/internal/tui/extensions/litemark"
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

//
// Testing Library
//

func newPlainTextBox(width int, height int) *tui.TextBox {
	tb := tui.TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = width
	tb.Height = height
	tb.ShowCursor = true

	doc := model.PlainDocument{}
	doc.Init()
	tb.Document = &doc

	engine := tui.PlainTextEngine{}
	tb.TextEngine = &engine
	return &tb
}

func newStyledTextBox(width int, height int) *tui.TextBox {
	tb := tui.TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = width
	tb.Height = height
	tb.ShowCursor = true

	doc := litemark.StyledDocument{}
	doc.Init()
	tb.Document = &doc

	engine := tui.LitemarkEngine{}
	tb.TextEngine = &engine
	return &tb
}

func parseToCell(r io.Reader, width int, height int) []tcell.SimCell {
	sc := bufio.NewScanner(r)

	table := make([]tcell.SimCell, width*height)
	spaceBytes := []byte{' '}
	spaceRunes := []rune{' '}
	for i := 0; i < width*height; i++ {
		table[i].Bytes = spaceBytes
		table[i].Runes = spaceRunes
	}

	ch := 0
	for sc.Scan() {
		line := sc.Text()
		z := html.NewTokenizer(strings.NewReader(line))

		depth := 0
		x := 0

	L:
		for {
			tt := z.Next()
			switch tt {
			case html.ErrorToken:
				break L
			case html.TextToken:
				if depth > 0 {
				} else {
					tagText := string(z.Text())
					clusters := text.GraphemeClusters(tagText)
					for _, cluster := range clusters {
						cell := tcell.SimCell{}
						runes := []rune(cluster)
						mainRune := runes[0]

						if runewidth.RuneWidth(mainRune) == 2 {
							cell.Bytes = []byte(string(runes[0]))
							cell.Style = tcell.StyleDefault
							cell.Runes = []rune{runes[0]}
							table[ch*height+x] = cell
							x++

							cell.Bytes = []byte(string(runes[1:]))
							cell.Style = tcell.StyleDefault
							cell.Runes = runes[1:]
							table[ch*height+x] = cell
							x++
						} else {
							cell.Bytes = []byte(cluster)
							cell.Style = tcell.StyleDefault
							cell.Runes = []rune(cluster)
							table[ch*height+x] = cell
							x++
						}
					}
				}
			case html.StartTagToken, html.EndTagToken:
				tn, _ := z.TagName()
				if len(tn) > 0 {
					if tt == html.StartTagToken {
						depth++
					} else {
						depth--
					}
				}
			case html.SelfClosingTagToken:
				tn, _ := z.TagName()
				if string(tn) == "c" {
					cell := tcell.SimCell{}
					cell.Bytes = []byte{' '}
					cell.Style = tcell.StyleDefault.Reverse(true)
					cell.Runes = []rune{' '}
					table[ch*height+x] = cell
					x++
				}
			}
		}

		ch++
	}
	return table
}

func renderAndCompare(t *testing.T, tb *tui.TextBox, expected io.Reader) {
	screen := tcell.NewSimulationScreen("UTF-8")
	screen.Init()
	screen.SetSize(tb.Width, tb.Height)

	g := base.Graphics{}
	g.Init(screen)
	g.Resize(tb.Width, tb.Height)

	tb.Draw(&g)
	screen.Show()

	cells, _, _ := screen.GetContents()
	assert.True(t, deepEquals(cells, parseToCell(expected, tb.Width, tb.Height)))
}

func deepEquals(a []tcell.SimCell, b []tcell.SimCell) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		aCell := a[i]
		bCell := b[i]

		aFg, aBg, aAttr := aCell.Style.Decompose()
		bFg, bBg, bAttr := bCell.Style.Decompose()

		if aFg != bFg {
			return false
		}
		if aBg != bBg {
			return false
		}
		if aAttr != bAttr {
			return false
		}

		if !slices.Equal(aCell.Bytes, bCell.Bytes) {
			return false
		}
		if !slices.Equal(aCell.Runes, bCell.Runes) {
			return false
		}
	}
	return true
}

//
// Tests
//

func TestTextBoxN01(t *testing.T) {
	tb := newPlainTextBox(10, 10)
	tb.InsertString("1234あ")

	renderAndCompare(t, tb, strings.NewReader("1234あ<c/>"))
}
