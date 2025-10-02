package litemark

import "github.com/desktopgame/ckro/internal/tui/model"

type CodeBlockElement struct {
	Ranges   []model.Range
	Children []model.Element
	Lang     string
}

func (c *CodeBlockElement) GetRange(index int) model.Range {
	return c.Ranges[index]
}

func (c *CodeBlockElement) GetRangeCount() int {
	return len(c.Ranges)
}

func (c *CodeBlockElement) GetElement(index int) model.Element {
	return c.Children[index]
}

func (c *CodeBlockElement) GetElementCount() int {
	return len(c.Children)
}
