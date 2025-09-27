package view

import (
	"github.com/desktopgame/ckro/internal/text"
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

type PlainTextView struct {
}

func (p *PlainTextView) Layout(ctx Context, textLayout *TextLayout, x, y, w, h int) {
	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
}

func (p *PlainTextView) Draw(ctx Context, textLayout *TextLayout, renderer Renderer) {
	// x := 0
	// y := 0

	// def := tcell.StyleDefault
	clusters := text.GraphemeClusters(ctx.GetText(textLayout.Element))
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

				renderer.SetContent(x, y, mainRune, combining, tcell.StyleDefault)
				// 全角文字の場合、次のセルを空にする
				if width == 2 {
					x++
					renderer.SetContent(x, y, 0, nil, tcell.StyleDefault)
				}
			}
			x++
		}
	}
}

func (p *PlainTextView) MinimumSize(ctx Context, e model.Element, width int, height int) *TextLayout {
	return &TextLayout{
		Element:       e,
		MinimumWidth:  text.DisplayWidth(ctx.GetText(e)),
		MinimumHeight: 1,
	}
}
