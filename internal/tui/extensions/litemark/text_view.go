package litemark

import (
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/desktopgame/ckro/internal/tui/view"
)

type TextView struct {
}

func (t *TextView) Layout(textViewResolver view.TextViewResolver, textLayout *view.TextLayout, x, y, w, h int) {
	offsetX := 0
	for i := 0; i < len(textLayout.Children); i++ {
		childElement := textLayout.Children[i].Element
		childView := textViewResolver.Resolve(childElement)

		mw := textLayout.Children[i].MinimumWidth
		childView.Layout(textViewResolver, textLayout.Children[i], offsetX, 0, mw, 1)
		offsetX += textLayout.Children[i].Width
	}
	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
}

func (t *TextView) Draw(textViewResolver view.TextViewResolver, textLayout *view.TextLayout, renderer view.Renderer) {
	for _, child := range textLayout.Children {
		childView := textViewResolver.Resolve(child.Element)
		childView.Draw(textViewResolver, child, renderer.Translate(child.RelativeX, child.RelativeY))
	}
}

func (t *TextView) MinimumSize(textViewResolver view.TextViewResolver, e model.Element, width int, height int) *view.TextLayout {
	totalWidth := 0
	children := []*view.TextLayout{}
	column := 0
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)

		var child *view.TextLayout
		if tsView, ok := childView.(view.TabStopTextView); ok {
			width := tsView.WidthWithTabStop(textViewResolver, childElement, column)
			child = &view.TextLayout{
				Element:       childElement,
				MinimumWidth:  width,
				MinimumHeight: 1,
			}
		} else {
			child = childView.MinimumSize(textViewResolver, childElement, width, 1)
		}
		children = append(children, child)

		column += child.MinimumWidth
		totalWidth += child.MinimumWidth
	}
	return &view.TextLayout{
		Element:       e,
		MinimumWidth:  totalWidth,
		MinimumHeight: 1,
		Children:      children,
	}
}
