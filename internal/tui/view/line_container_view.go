package view

import (
	"github.com/desktopgame/ckro/internal/tui/model"
)

type LineContainerView struct {
}

func (l *LineContainerView) Layout(textViewResolver TextViewResolver, e model.Element, width int) *TextLayout {
	children := []*TextLayout{}
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		children = append(children, childView.Layout(textViewResolver, childElement, width))
	}

	return &TextLayout{
		Element:  e,
		Children: children,
	}
}

func (l *LineContainerView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	if lineContainer, ok := textLayout.Element.(*model.LineContainerElement); ok {
		x := 0
		for i := 0; i < lineContainer.GetElementCount(); i++ {
			childElement := lineContainer.GetElement(i)
			childView := textViewResolver.Resolve(childElement)

			childView.Draw(textViewResolver, textLayout.Children[i], renderer.Translate(x, 0))

			width := 0
			if tChildView, ok := childView.(TabStopTextView); ok {
				width = tChildView.WidthWithTabStop(textViewResolver, textLayout.Children[i], width)
			} else {
				width = childView.Width(textViewResolver, textLayout.Children[i], 0)
			}
			x += width
		}
		return
	}
	panic("LineContainerView requires LineContainerElement")
}

func (l *LineContainerView) Width(textViewResolver TextViewResolver, textLayout *TextLayout, row int) int {
	if lineContainer, ok := textLayout.Element.(*model.LineContainerElement); ok {
		width := 0
		for i := 0; i < lineContainer.GetElementCount(); i++ {
			childElement := lineContainer.GetElement(i)
			childView := textViewResolver.Resolve(childElement)

			if tChildView, ok := childView.(TabStopTextView); ok {
				width += tChildView.WidthWithTabStop(textViewResolver, textLayout.Children[i], width)
			} else {
				width += childView.Width(textViewResolver, textLayout.Children[i], 0)
			}
		}
		return width
	}
	panic("LineContainerView requires LineContainerElement")
}

func (l *LineContainerView) Height(textViewResolver TextViewResolver, textLayout *TextLayout) int {
	return 1
}
