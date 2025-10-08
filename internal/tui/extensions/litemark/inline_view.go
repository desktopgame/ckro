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

func (il *InlineView) DrawWithTabStop(ctx view.Context, textLayout *view.TextLayout, renderer view.Renderer, column int) int {
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

	selectStyle := tcell.StyleDefault.Reverse(true)

	r := textLayout.Element.GetRange(1)
	sg := ctx.Document.Read(r)
	clusters := text.GraphemeClusters(sg.GetLine(0))
	x := 0
	y := 0
	at := r.StartPosition
	sel := ctx.TextSelection
	for _, cluster := range clusters {
		next := at
		next.Column += len(cluster)

		if cluster == "\t" {
			spaces := text.TabWidth - (column % text.TabWidth)
			tabStyle := tcell.StyleDefault
			if sel.Contain(at) {
				tabStyle = selectStyle
			}
			for i := 0; i < spaces; i++ {
				renderer.SetContent(x+i, y, ' ', nil, tabStyle)
			}
			x += spaces
			column += spaces

		} else {
			runes := []rune(cluster)

			charStyle := style
			if sel.Contain(at) {
				charStyle = selectStyle
			}

			if len(runes) > 0 {
				mainRune := runes[0]
				var combining []rune

				// 残りのruneをcombining charactersとして設定
				if len(runes) > 1 {
					combining = runes[1:]
				}
				width := runewidth.RuneWidth(mainRune)

				renderer.SetContent(x, y, mainRune, combining, charStyle)
				// 全角文字の場合、次のセルを空にする
				if width == 2 {
					x++
					column++
					renderer.SetContent(x, y, 0, nil, charStyle)
				}
			}
			x++
			column++
		}
		at = next
	}
	return column
}

func (il *InlineView) Draw(ctx view.Context, textLayout *view.TextLayout, renderer view.Renderer) {
	il.DrawWithTabStop(ctx, textLayout, renderer, 0)
}

func (il *InlineView) WidthWithTabStop(ctx view.Context, e model.Element, column int) int {
	clusters := text.GraphemeClusters(ctx.GetSegment(e, 1).GetLine(0))
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
	return text.DisplayPos(ctx.GetSegment(e, 1).GetLine(0), viewLocalPos), 0
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
	r := e.GetRange(1)

	if bytePos.Row < r.StartPosition.Row {
		return 0
	}
	if bytePos.Row > r.EndPosition.Row {
		return il.MoveLength(ctx, textLayout)
	}

	if bytePos.Column < r.StartPosition.Column {
		return 0
	}
	if bytePos.Column >= r.EndPosition.Column {
		return il.MoveLength(ctx, textLayout)
	}

	r = e.GetRange(0)
	str := ctx.Document.Read(r).GetLine(bytePos.Row - r.StartPosition.Row)

	bytes := inlineElement.Pad
	clusters := text.GraphemeClusters(str)
	if inlineElement.Pad == 0 {
		for i, cluster := range clusters {
			if bytes == bytePos.Column-r.StartPosition.Column {
				return i
			}
			bytes += len(cluster)
		}
	} else {
		for i, cluster := range clusters {
			if i < inlineElement.Pad {
				continue
			}
			if bytes == bytePos.Column-r.StartPosition.Column {
				return i - inlineElement.Pad
			}
			bytes += len(cluster)
		}
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

func (il *InlineView) ShouldRemoveWithLine(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (int, bool) {
	return -1, false
}

func (il *InlineView) ShouldRemoveWithSpecifiedColumnAfter(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (model.Position, bool) {
	return model.Position{}, false
}

func (il *InlineView) ShouldRemoveWithSpecifiedRangeLines(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (model.Range, bool) {
	return model.Range{}, false
}

func (il *InlineView) ShouldRemoveWithSpecifiedRangeColumns(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (model.Range, bool) {
	panic("ShouldRemoveWithSpecifiedRangeColumns is not implemented")
}

func (il *InlineView) ShouldRemoveLastCharacter(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (model.Element, bool) {
	return nil, false
}
