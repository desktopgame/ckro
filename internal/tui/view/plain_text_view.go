package view

import (
	"github.com/desktopgame/ckro/internal/text"
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

type PlainTextView struct {
}

func (p *PlainTextView) Layout(textViewResolver TextViewResolver, e model.Element, x, y, w, h int) *TextLayout {
	return &TextLayout{
		Element:  e,
		Children: nil,
		Width:    text.DisplayWidth(e.GetText()),
		Height:   1,
	}
}

func (p *PlainTextView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	// x := 0
	// y := 0

	// def := tcell.StyleDefault
	style := ConvertStyle(textLayout.Element.GetStyle())
	clusters := text.GraphemeClusters(textLayout.Element.GetText())
	x := 0
	y := 0
	for _, cluster := range clusters {

		if cluster == "\t" {
			spaces := text.TabWidth - (x % text.TabWidth)
			for i := 0; i < spaces; i++ {
				renderer.SetContent(x+i, y, ' ', nil, tcell.StyleDefault)
			}
			x += spaces

		} else {
			runes := []rune(cluster)

			if len(runes) > 0 {
				mainRune := runes[0]
				var combining []rune

				// 残りのruneをcombining charactersとして設定
				if len(runes) > 1 {
					combining = runes[1:]
				}
				width := runewidth.RuneWidth(mainRune)

				renderer.SetContent(x, y, mainRune, combining, style)
				// 全角文字の場合、次のセルを空にする
				if width == 2 {
					x++
					renderer.SetContent(x, y, 0, nil, style)
				}
			}
			x++
		}
	}
}

func (p *PlainTextView) WidthWithTabStop(textViewResolver TextViewResolver, textLayout *TextLayout, column int) int {
	clusters := text.GraphemeClusters(textLayout.Element.GetText())
	totalWidth := 0
	for _, cluster := range clusters {
		width := 0
		if cluster == "\t" {
			width = text.TabWidth - (column % text.TabWidth)
		} else {
			width = runewidth.StringWidth(cluster)
		}
		column += width
		totalWidth += width
	}
	return totalWidth
}

func (p *PlainTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) (Width int, Height int) {
	return text.DisplayWidth(e.GetText()), 1
}
