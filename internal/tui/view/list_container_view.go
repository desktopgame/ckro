package view

import (
	"github.com/desktopgame/ckro/internal/tui/model"
)

type ListContainerView struct {
}

func (l *ListContainerView) Layout(textViewResolver TextViewResolver, e model.Element, width int) *TextLayout {
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

func (l *ListContainerView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer, x int, y int, localViewLine int) {
	if listContainer, ok := textLayout.Element.(*model.ListContainerElement); ok {
		for i := 0; i < listContainer.GetElementCount(); i++ {
			childElement := listContainer.GetElement(i)
			childView := textViewResolver.Resolve(childElement)

			height := childView.Height(textViewResolver, textLayout.Children[i])
			if localViewLine < height {
				childView.Draw(textViewResolver, textLayout.Children[i], renderer, x, y, localViewLine)
				break
			}
			localViewLine -= height
		}
		return
	}
	panic("ListContainerView requires ListContainerElement")
}

func (l *ListContainerView) Width(textViewResolver TextViewResolver, textLayout *TextLayout, row int) int {
	if listContainer, ok := textLayout.Element.(*model.ListContainerElement); ok {
		for i := 0; i < listContainer.GetElementCount(); i++ {
			childElement := listContainer.GetElement(i)
			childView := textViewResolver.Resolve(childElement)

			height := childView.Height(textViewResolver, textLayout.Children[i])
			if row < height {
				return childView.Width(textViewResolver, textLayout.Children[i], row)
			}
			row -= height
		}
		return -1
	}
	panic("ListContainerView requires ListContainerElement")
}

func (l *ListContainerView) Height(textViewResolver TextViewResolver, textLayout *TextLayout) int {
	if listContainer, ok := textLayout.Element.(*model.ListContainerElement); ok {
		height := 0
		for i := 0; i < listContainer.GetElementCount(); i++ {
			childElement := listContainer.GetElement(i)
			childView := textViewResolver.Resolve(childElement)

			height += childView.Height(textViewResolver, textLayout.Children[i])
		}
		return height
	}
	panic("ListContainerView requires ListContainerElement")
}
