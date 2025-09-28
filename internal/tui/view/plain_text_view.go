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

func (p *PlainTextView) MoveLength(ctx Context, e model.Element) int {
	return text.GraphemeLength(ctx.GetText(e)) + 1 // include newline
}

func (p *PlainTextView) MoveUp(ctx Context, e model.Element, viewLocalPos int) int {
	return -1
}

func (p *PlainTextView) MoveDown(ctx Context, e model.Element, viewLocalPos int) int {
	return -1
}

func (p *PlainTextView) MoveLeft(ctx Context, e model.Element, viewLocalPos int) int {
	if viewLocalPos <= 0 {
		return -1
	}
	return viewLocalPos - 1
}

func (p *PlainTextView) MoveRight(ctx Context, e model.Element, viewLocalPos int) int {
	if viewLocalPos > p.MoveLength(ctx, e) {
		return -1
	}
	return viewLocalPos + 1
}

func (p *PlainTextView) ConvertPos(ctx Context, e model.Element, viewLocalPos int) (ViewLocalX int, ViewLocalY int) {
	return text.DisplayPos(ctx.GetText(e), viewLocalPos), 0
}

func (p *PlainTextView) ConvertModel(ctx Context, e model.Element, viewLocalPos int) model.Position {
	st := e.GetRange(0).StartPosition
	return model.Position{
		Row:    st.Row,
		Column: st.Column + viewLocalPos,
	}
}

func (p *PlainTextView) ConvertRelativeX(ctx Context, e model.Element, viewLocalPos int) int {
	return viewLocalPos
}

func (p *PlainTextView) MoveFirstLine(ctx Context, e model.Element, relX int) int {
	l := p.MoveLength(ctx, e)
	return min(relX, l-1)
}

func (p *PlainTextView) MoveLastLine(ctx Context, e model.Element, relX int) int {
	l := p.MoveLength(ctx, e)
	return min(relX, l-1)
}
