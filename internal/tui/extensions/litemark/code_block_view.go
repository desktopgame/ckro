package litemark

import (
	"github.com/desktopgame/ckro/internal/text"
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/desktopgame/ckro/internal/tui/view"
)

type CodeBlockView struct {
}

func (c *CodeBlockView) Layout(ctx view.Context, textLayout *view.TextLayout, x, y, w, h int) {
	offsetY := 0
	for i := 0; i < len(textLayout.Children); i++ {
		childElement := textLayout.Children[i].Element
		childView := ctx.Resolver.Resolve(childElement)

		mw := textLayout.Children[i].MinimumWidth
		childView.Layout(ctx, textLayout.Children[i], 0, offsetY, mw, 1)
		offsetY++
	}
	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
}

func (c *CodeBlockView) Draw(ctx view.Context, textLayout *view.TextLayout, renderer view.Renderer) {
	for _, child := range textLayout.Children {
		childView := ctx.Resolver.Resolve(child.Element)
		childView.Draw(ctx, child, renderer.Translate(child.RelativeX, child.RelativeY))
	}
}

func (c *CodeBlockView) MinimumSize(ctx view.Context, e model.Element, width int, height int) *view.TextLayout {
	totalWidth := 0
	children := []*view.TextLayout{}
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := ctx.Resolver.Resolve(childElement)

		child := childView.MinimumSize(ctx, childElement, width, 1)
		children = append(children, child)
		totalWidth += child.MinimumWidth
	}
	return &view.TextLayout{
		Element:       e,
		MinimumWidth:  totalWidth,
		MinimumHeight: e.GetElementCount(),
		Children:      children,
	}
}

func (c *CodeBlockView) MoveLength(ctx view.Context, e model.Element) int {
	return text.GraphemeLength(ctx.GetText(e)) + 1 // include newline
}

func (c *CodeBlockView) MoveUp(ctx view.Context, e model.Element, viewLocalPos int) int {
	return -1
}

func (c *CodeBlockView) MoveDown(ctx view.Context, e model.Element, viewLocalPos int) int {
	return -1
}

func (c *CodeBlockView) MoveLeft(ctx view.Context, e model.Element, viewLocalPos int) int {
	if viewLocalPos <= 0 {
		return -1
	}
	return viewLocalPos - 1
}

func (c *CodeBlockView) MoveRight(ctx view.Context, e model.Element, viewLocalPos int) int {
	if viewLocalPos > c.MoveLength(ctx, e) {
		return -1
	}
	return viewLocalPos + 1
}

func (c *CodeBlockView) ConvertPos(ctx view.Context, e model.Element, viewLocalPos int) (ViewLocalX int, ViewLocalY int) {
	return text.DisplayPos(ctx.GetText(e), viewLocalPos), 0
}

func (c *CodeBlockView) ConvertModel(ctx view.Context, e model.Element, viewLocalPos int) model.Position {
	st := e.GetRange(0).StartPosition
	return model.Position{
		Row:    st.Row,
		Column: st.Column + viewLocalPos,
	}
}
