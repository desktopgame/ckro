package litemark

import (
	"github.com/desktopgame/ckro/internal/text"
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/desktopgame/ckro/internal/tui/view"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

type HeadingView struct {
}

func (hv *HeadingView) Layout(ctx view.Context, textLayout *view.TextLayout, x, y, w, h int) {
	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
}

func (hv *HeadingView) Draw(ctx view.Context, textLayout *view.TextLayout, renderer view.Renderer) {
	// x := 0
	// y := 0
	headingElement := textLayout.Element.(*HeadingElement)

	var color tcell.Color
	switch headingElement.Level {
	case 1:
		color = tcell.ColorRed // H1: Red
	case 2:
		color = tcell.ColorBlue // H2: Blue
	case 3:
		color = tcell.ColorGreen // H3: Green
	case 4:
		color = tcell.ColorYellow // H4: Yellow
	case 5:
		color = tcell.ColorPurple // H5: Purple
	case 6:
		color = tcell.ColorTeal // H6: Teal
	default:
		color = tcell.ColorWhite // Default: White
	}

	// def := tcell.StyleDefault
	style := tcell.StyleDefault.Bold(true).Foreground(color)
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

func (hv *HeadingView) MinimumSize(ctx view.Context, e model.Element, width int, height int) *view.TextLayout {
	return &view.TextLayout{
		Element:       e,
		MinimumWidth:  text.DisplayWidth(ctx.GetSegment(e, 1).GetLine(0)),
		MinimumHeight: 1,
	}
}

func (hv *HeadingView) MoveLength(ctx view.Context, textLayout *view.TextLayout) int {
	line := ctx.GetSegment(textLayout.Element, 1).GetLine(0)
	return text.GraphemeLength(line) + 1 // include newline
}

func (hv *HeadingView) MoveUp(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	return -1
}

func (hv *HeadingView) MoveDown(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	return -1
}

func (hv *HeadingView) MoveLeft(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	if viewLocalPos <= 0 {
		return -1
	}
	return viewLocalPos - 1
}

func (hv *HeadingView) MoveRight(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	if viewLocalPos >= hv.MoveLength(ctx, textLayout)-1 {
		return -1
	}
	return viewLocalPos + 1
}

func (hv *HeadingView) ConvertPos(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (ViewLocalX int, ViewLocalY int) {
	e := textLayout.Element
	return text.DisplayPos(ctx.GetSegment(e, 1).GetLine(0), viewLocalPos), 0
}

func (hv *HeadingView) ConvertModel(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) view.CharacterReference {
	e := textLayout.Element
	rs := e.GetRange(0)
	r := e.GetRange(1)
	st := rs.StartPosition
	str := ctx.Document.Read(r).GetLine(0)

	if viewLocalPos >= hv.MoveLength(ctx, textLayout)-1 {
		return view.CharacterReference{
			StartPosition: model.Position{
				Row:    st.Row,
				Column: st.Column + len(str) + (e.(*HeadingElement).Level + 1),
			},
			Bytes: 0,
		}
	}

	bPos, bLen := text.GraphemeToByteRange(str, viewLocalPos)
	return view.CharacterReference{
		StartPosition: model.Position{
			Row:    st.Row,
			Column: st.Column + bPos + (e.(*HeadingElement).Level + 1),
		},
		Bytes: bLen,
	}
}

func (hv *HeadingView) ConvertViewLocalPos(ctx view.Context, textLayout *view.TextLayout, bytePos model.Position) int {
	e := textLayout.Element
	r := e.GetRange(1)

	if bytePos.Row < r.StartPosition.Row {
		return 0
	}
	if bytePos.Row > r.EndPosition.Row {
		return hv.MoveLength(ctx, textLayout) - 1
	}

	if bytePos.Column < r.StartPosition.Column+1 {
		return 0
	}
	if bytePos.Column >= r.EndPosition.Column {
		return hv.MoveLength(ctx, textLayout) - 1
	}

	r = e.GetRange(1)
	str := ctx.Document.Read(r).GetLine(bytePos.Row - r.StartPosition.Row)

	bytes := 0
	clusters := text.GraphemeClusters(str)
	for i, cluster := range clusters {
		if bytes == bytePos.Column {
			ofs := e.(*HeadingElement).Level + 1
			return max(ofs, i) - ofs
		}
		bytes += len(cluster)
	}
	return hv.MoveLength(ctx, textLayout) - 1
	//panic("")
}

func (hv *HeadingView) ShouldBeforeInsertionNewLineOnLineBegin(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) bool {
	return viewLocalPos == 0
}

func (hv *HeadingView) ShouldRemoveWithLine(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (int, bool) {
	r := textLayout.Element.GetRange(0)
	return r.StartPosition.Row, viewLocalPos == 1 && hv.MoveLength(ctx, textLayout) == 2
}

func (hv *HeadingView) ShouldRemoveWithSpecifiedColumnAfter(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (model.Position, bool) {
	return model.Position{}, false
}

func (hv *HeadingView) ShouldRemoveWithSpecifiedRangeLines(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (model.Range, bool) {
	return model.Range{}, false
}

func (hv *HeadingView) ShouldRemoveWithSpecifiedRangeColumns(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (model.Range, bool) {
	return model.Range{}, false
}

func (hv *HeadingView) ShouldRemoveLastCharacter(ctx view.Context, textLayout *view.TextLayout) bool {
	return false
}
