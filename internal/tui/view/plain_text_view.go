package view

import (
	"github.com/desktopgame/ckro/internal/text"
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

type PlainTextView struct {
}

func (p *PlainTextView) Layout(e model.Element, width int) *TextLayout {
	return &TextLayout{
		Element:  e,
		Children: nil,
	}
}

func (p *PlainTextView) Draw(textLayout *TextLayout, renderer Renderer, x int, y int) {
	// x := 0
	// y := 0
	def := tcell.StyleDefault
	clusters := text.GraphemeClusters(textLayout.Element.GetText())
	for _, cluster := range clusters {

		if cluster == "\t" {
			spaces := text.TabWidth - (x % text.TabWidth)
			for i := 0; i < spaces; i++ {
				renderer.SetContent(x+i, y, ' ', nil, def)
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

				renderer.SetContent(x, y, mainRune, combining, def)
				// 全角文字の場合、次のセルを空にする
				if width == 2 {
					x++
					renderer.SetContent(x, y, 0, nil, def)
				}
			}
			x++
		}
	}
}

func (p *PlainTextView) Width(textLayout *TextLayout, row int) int {
	return text.DisplayWidth(textLayout.Element.GetText())
}

func (p *PlainTextView) Height(textLayout *TextLayout) int {
	return 1
}
