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

func (hv *HeadingView) MoveLength(ctx view.Context, e model.Element) int {
	line := ctx.GetSegment(e, 1).GetLine(0)
	return text.GraphemeLength(line) + 1 // include newline
}

func (hv *HeadingView) MoveUp(ctx view.Context, e model.Element, viewLocalPos int) int {
	return -1
}

func (hv *HeadingView) MoveDown(ctx view.Context, e model.Element, viewLocalPos int) int {
	return -1
}

func (hv *HeadingView) MoveLeft(ctx view.Context, e model.Element, viewLocalPos int) int {
	if viewLocalPos <= 0 {
		return -1
	}
	return viewLocalPos - 1
}

func (hv *HeadingView) MoveRight(ctx view.Context, e model.Element, viewLocalPos int) int {
	if viewLocalPos >= hv.MoveLength(ctx, e) {
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

	graphemes := text.GraphemeLength(str)
	if viewLocalPos >= graphemes {
		return view.CharacterReference{
			StartPosition: model.Position{
				Row:    st.Row,
				Column: st.Column + graphemes + (e.(*HeadingElement).Level + 1),
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
