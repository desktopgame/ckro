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
	inlineElement := textLayout.Element.(*InlineElement)
	style := tcell.StyleDefault

	if inlineElement.IsBold {
		style = style.Bold(true).Foreground(tcell.ColorDarkRed)
	}

	if inlineElement.IsItalic {
		style = style.Italic(true)
	}

	if inlineElement.IsUnderline {
		style = style.Underline(true)
	}

	if fg, ok := inlineElement.Foreground.TryValue(); ok {
		style = style.Foreground(fg)
	}

	if bg, ok := inlineElement.Background.TryValue(); ok {
		style = style.Background(bg)
	}

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

func (il *InlineView) MoveLength(ctx view.Context, textLayout *view.TextLayout) int {
	line := ctx.GetSegment(textLayout.Element, 1).GetLine(0)
	return text.GraphemeLength(line)
}

func (il *InlineView) MoveUp(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	return -1
}

func (il *InlineView) MoveDown(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	return -1
}

func (il *InlineView) MoveLeft(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	if viewLocalPos <= 0 {
		return -1
	}
	return viewLocalPos - 1
}

func (il *InlineView) MoveRight(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	if viewLocalPos >= il.MoveLength(ctx, textLayout)-1 {
		return -1
	}
	return viewLocalPos + 1
}

func (il *InlineView) ConvertPos(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (ViewLocalX int, ViewLocalY int) {
	e := textLayout.Element
	return text.DisplayPos(ctx.GetText(e), viewLocalPos), 0
}

func (il *InlineView) ConvertModel(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) view.CharacterReference {
	e := textLayout.Element
	inlineElement := e.(*InlineElement)
	r := e.GetRange(0)
	st := r.StartPosition
	bPos, bLen := text.GraphemeToByteRange(ctx.GetSegment(e, 1).GetLine(0), viewLocalPos)
	return view.CharacterReference{
		StartPosition: model.Position{
			Row:    st.Row,
			Column: st.Column + inlineElement.Pad + bPos,
		},
		Bytes: bLen,
	}
}

func (il *InlineView) ConvertViewLocalPos(ctx view.Context, textLayout *view.TextLayout, bytePos model.Position) int {
	e := textLayout.Element
	inlineElement := e.(*InlineElement)
	r := e.GetRange(0)
	str := ctx.Document.Read(r).GetLine(bytePos.Row - r.StartPosition.Row)

	if bytePos.Column == e.GetRange(1).StartPosition.Column {
		return 0
	}

	bytes := inlineElement.Pad
	clusters := text.GraphemeClusters(str)
	for i, cluster := range clusters {
		if i <= inlineElement.Pad {
			continue
		}
		if bytes == bytePos.Column-r.StartPosition.Column {
			return i - inlineElement.Pad - 1
		}
		bytes += len(cluster)
	}
	return il.MoveLength(ctx, textLayout) - 1
}

func (il *InlineView) ShouldBeforeInsertionNewLineOnLineBegin(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) bool {
	ile := textLayout.Element.(*InlineElement)
	hasColor := ile.Foreground.IsSome() || ile.Background.IsSome()
	if ile.IsBold || ile.IsItalic || ile.IsUnderline || hasColor {
		return viewLocalPos == 0
	}
	return false
}

func (il *InlineView) ShouldRemoveWithLine(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) bool {
	return false
}

func (il *InlineView) ShouldRemoveWithSpecifiedColumnAfter(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (model.Position, bool) {
	return model.Position{}, false
}
