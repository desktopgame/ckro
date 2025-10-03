package view

import (
	"github.com/desktopgame/ckro/internal/tui/model"
)

type FoldBlockView struct {
}

func (fv *FoldBlockView) Layout(ctx Context, textLayout *TextLayout, x, y, w, h int) {
	offsetY := 0
	for i := 0; i < len(textLayout.Children); i++ {
		childElement := textLayout.Children[i].Element
		childView := ctx.Resolver.Resolve(childElement)

		mw := textLayout.Children[i].MinimumWidth
		mh := textLayout.Children[i].MinimumHeight
		childView.Layout(ctx, textLayout.Children[i], 1, offsetY, mw, mh)
		offsetY += mh
	}
	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
}

func (fv *FoldBlockView) Draw(ctx Context, textLayout *TextLayout, renderer Renderer) {
}

func (fv *FoldBlockView) MinimumSize(ctx Context, e model.Element, width int, height int) *TextLayout {
	maxWidth := -1
	children := []*TextLayout{}
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := ctx.Resolver.Resolve(childElement)

		child := childView.MinimumSize(ctx, childElement, width, 1)
		children = append(children, child)

		if child.MinimumWidth > maxWidth {
			maxWidth = child.MinimumWidth
		}
	}
	return &TextLayout{
		Element:       e,
		MinimumWidth:  min(width, maxWidth),
		MinimumHeight: 1,
		Children:      children,
	}
}

func (fv *FoldBlockView) MoveLength(ctx Context, e model.Element) int {
	return 1
}

func (fv *FoldBlockView) MoveUp(ctx Context, e model.Element, viewLocalPos int) int {
	return -1
}

func (fv *FoldBlockView) MoveDown(ctx Context, e model.Element, viewLocalPos int) int {
	return -1
}

func (fv *FoldBlockView) MoveLeft(ctx Context, e model.Element, viewLocalPos int) int {
	return -1
}

func (fv *FoldBlockView) MoveRight(ctx Context, e model.Element, viewLocalPos int) int {
	return -1
}

func (fv *FoldBlockView) ConvertPos(ctx Context, textLayout *TextLayout, viewLocalPos int) (ViewLocalX int, ViewLocalY int) {
	return 0, 0
}

func (fv *FoldBlockView) ConvertModel(ctx Context, textLayout *TextLayout, viewLocalPos int) CharacterReference {
	e := textLayout.Element
	r := e.GetRange(0)
	st := r.StartPosition
	return CharacterReference{
		StartPosition: model.Position{
			Row:    st.Row,
			Column: st.Column,
		},
		Bytes: 0,
	}
}

func (fv *FoldBlockView) ConvertViewLocalPos(ctx Context, textLayout *TextLayout, bytePos model.Position) int {
	return 0
}
