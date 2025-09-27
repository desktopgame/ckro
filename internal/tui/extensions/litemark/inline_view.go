package litemark

import (
	"github.com/desktopgame/ckro/internal/text"
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/desktopgame/ckro/internal/tui/view"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

type InlineView struct {
}

func (il *InlineView) Layout(ctx view.Context, textLayout *view.TextLayout, x, y, w, h int) {
	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
}

func (il *InlineView) Draw(ctx view.Context, textLayout *view.TextLayout, renderer view.Renderer) {
	// x := 0
	// y := 0

	// def := tcell.StyleDefault
	style := tcell.StyleDefault
	clusters := text.GraphemeClusters(ctx.GetSegment(textLayout.Element, 1).GetLine(0))
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

func (il *InlineView) WidthWithTabStop(ctx view.Context, textLayout *view.TextLayout, column int) int {
	clusters := text.GraphemeClusters(ctx.GetSegment(textLayout.Element, 1).GetLine(0))
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

func (il *InlineView) MinimumSize(ctx view.Context, e model.Element, width int, height int) *view.TextLayout {
	return &view.TextLayout{
		Element:       e,
		MinimumWidth:  text.DisplayWidth(ctx.GetSegment(e, 1).GetLine(0)),
		MinimumHeight: 1,
	}
}
