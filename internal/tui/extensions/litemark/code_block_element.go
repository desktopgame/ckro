package litemark

import "github.com/desktopgame/ckro/internal/tui/model"

type CodeBlockElement struct {
	Range    model.Range
	Children []model.Element
}

func (c *CodeBlockElement) GetRange(index int) model.Range {
	return c.Range
}

func (c *CodeBlockElement) GetRangeCount() int {
	return 1
}

func (c *CodeBlockElement) GetElement(index int) model.Element {
	return c.Children[index]
}

func (c *CodeBlockElement) GetElementCount() int {
	return len(c.Children)
}
