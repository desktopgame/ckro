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
	clusters := text.GraphemeClusters(ctx.GetText(textLayout.Element))
	x := 0

	viewLine := 0
	width := textLayout.Width
	for _, cluster := range clusters {
		runes := []rune(cluster)

		if cluster == "\t" {
			w := text.TabWidth - (x % text.TabWidth)
			if x+w > width {
				viewLine++
				x = 0
			}
			for i := 0; i < w; i++ {
				renderer.SetContent(x+i, viewLine, ' ', nil, tcell.StyleDefault)
			}
			x += w
		} else if len(runes) > 0 {
			mainRune := runes[0]
			w := runewidth.RuneWidth(mainRune)

			if x+w > width {
				viewLine++
				x = 0
			}
			renderer.SetContent(x, viewLine, mainRune, runes[1:], tcell.StyleDefault)
			if w == 2 {
				x++
				renderer.SetContent(x, viewLine, 0, nil, tcell.StyleDefault)
			}
			x++
		}
	}
}

func (p *PlainTextView) MinimumSize(ctx Context, e model.Element, width int, height int) *TextLayout {
	line := ctx.GetText(e)
	lineWidth := text.DisplayWidth(line)

	if lineWidth <= width {
		return &TextLayout{
			Element:       e,
			MinimumWidth:  width,
			MinimumHeight: 1,
		}
	}

	x := 0
	viewLine := 1

	clusters := text.GraphemeClusters(line)
	for _, cluster := range clusters {
		runes := []rune(cluster)

		if cluster == "\t" {
			w := text.TabWidth - (x % text.TabWidth)
			if x+w > width {
				viewLine++
				x = 0
			}
			x += w
		} else if len(runes) > 0 {
			mainRune := runes[0]
			w := runewidth.RuneWidth(mainRune)

			if x+w > width {
				viewLine++
				x = 0
			}
			x += w
		}
	}

	return &TextLayout{
		Element:       e,
		MinimumWidth:  width,
		MinimumHeight: viewLine,
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
	if viewLocalPos >= p.MoveLength(ctx, e)-1 {
		return -1
	}
	return viewLocalPos + 1
}

func (p *PlainTextView) ConvertPos(ctx Context, textLayout *TextLayout, viewLocalPos int) (ViewLocalX int, ViewLocalY int) {
	line := ctx.GetText(textLayout.Element)

	x := 0
	clusterCount := 0
	viewLine := 0

	width := textLayout.Width
	clusters := text.GraphemeClusters(line)
	if len(clusters) == 0 && viewLocalPos == 0 {
		return 0, 0
	}
	for _, cluster := range clusters {
		runes := []rune(cluster)
		if clusterCount == viewLocalPos {
			return x, viewLine
		}

		if cluster == "\t" {
			w := text.TabWidth - (x % text.TabWidth)
			if x+w > width {
				viewLine++
				x = 0
			}
			x += w
		} else if len(runes) > 0 {
			mainRune := runes[0]
			w := runewidth.RuneWidth(mainRune)

			if x+w > width {
				viewLine++
				x = 0
			}
			x += w
		}
		clusterCount++
	}
	return x, viewLine
}

func (p *PlainTextView) ConvertModel(ctx Context, textLayout *TextLayout, viewLocalPos int) CharacterReference {
	e := textLayout.Element
	r := e.GetRange(0)
	str := ctx.Document.Read(r).GetLine(0)
	st := r.StartPosition

	graphemes := text.GraphemeLength(str)
	if viewLocalPos >= graphemes {
		return CharacterReference{
			StartPosition: model.Position{
				Row:    st.Row,
				Column: st.Column + len(str),
			},
			Bytes: 0,
		}
	}

	bPos, bLen := text.GraphemeToByteRange(str, viewLocalPos)
	return CharacterReference{
		StartPosition: model.Position{
			Row:    st.Row,
			Column: st.Column + bPos,
		},
		Bytes: bLen,
	}
}

func (p *PlainTextView) ConvertViewLocalPos(ctx Context, textLayout *TextLayout, bytePos model.Position) int {
	e := textLayout.Element
	r := e.GetRange(0)
	str := ctx.Document.Read(r).GetLine(bytePos.Row - r.StartPosition.Row)

	bytes := 0
	clusters := text.GraphemeClusters(str)
	for i, cluster := range clusters {
		if bytes == bytePos.Column {
			return i
		}
		bytes += len(cluster)
	}
	return p.MoveLength(ctx, textLayout.Element)
}

func (p *PlainTextView) ConvertRelativeX(ctx Context, textLayout *TextLayout, viewLocalPos int) int {
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
