package litemark

import (
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/desktopgame/ckro/internal/tui/view"
)

type BlankLineView struct {
}

func (b *BlankLineView) Layout(ctx view.Context, textLayout *view.TextLayout, x, y, w, h int) {
	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
}

func (b *BlankLineView) Draw(ctx view.Context, textLayout *view.TextLayout, renderer view.Renderer) {
}

func (b *BlankLineView) MinimumSize(ctx view.Context, e model.Element, width int, height int) *view.TextLayout {
	return &view.TextLayout{
		Element:       e,
		MinimumWidth:  1,
		MinimumHeight: 1,
	}
}

func (b *BlankLineView) MoveLength(ctx view.Context, e model.Element) int {
	return 1
}

func (b *BlankLineView) MoveUp(ctx view.Context, e model.Element, viewLocalPos int) int {
	return -1
}

func (b *BlankLineView) MoveDown(ctx view.Context, e model.Element, viewLocalPos int) int {
	return -1
}

func (b *BlankLineView) MoveLeft(ctx view.Context, e model.Element, viewLocalPos int) int {
	return -1
}

func (b *BlankLineView) MoveRight(ctx view.Context, e model.Element, viewLocalPos int) int {
	return -1
}

func (b *BlankLineView) ConvertPos(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (ViewLocalX int, ViewLocalY int) {
	return 0, 0
}

func (b *BlankLineView) ConvertModel(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) view.CharacterReference {
	e := textLayout.Element
	r := e.GetRange(0)
	st := r.StartPosition
	return view.CharacterReference{
		StartPosition: model.Position{
			Row:    st.Row,
			Column: st.Column,
		},
		Bytes: 0,
	}
}

func (b *BlankLineView) ConvertViewLocalPos(ctx view.Context, textLayout *view.TextLayout, bytePos model.Position) int {
	return 0
}
