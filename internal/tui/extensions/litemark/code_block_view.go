package litemark

import (
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/desktopgame/ckro/internal/tui/view"
)

type CodeBlockView struct {
}

func (c *CodeBlockView) Layout(textViewResolver view.TextViewResolver, textLayout *view.TextLayout, x, y, w, h int) {
	offsetY := 0
	for i := 0; i < len(textLayout.Children); i++ {
		childElement := textLayout.Children[i].Element
		childView := textViewResolver.Resolve(childElement)

		mw := textLayout.Children[i].MinimumWidth
		childView.Layout(textViewResolver, textLayout.Children[i], 0, offsetY, mw, 1)
		offsetY++
	}
	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
}

func (c *CodeBlockView) Draw(textViewResolver view.TextViewResolver, textLayout *view.TextLayout, renderer view.Renderer) {
	for _, child := range textLayout.Children {
		childView := textViewResolver.Resolve(child.Element)
		childView.Draw(textViewResolver, child, renderer.Translate(child.RelativeX, child.RelativeY))
	}
}

func (c *CodeBlockView) MinimumSize(textViewResolver view.TextViewResolver, e model.Element, width int, height int) *view.TextLayout {
	totalWidth := 0
	children := []*view.TextLayout{}
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)

		child := childView.MinimumSize(textViewResolver, childElement, width, 1)
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
