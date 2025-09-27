package litemark

import (
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
