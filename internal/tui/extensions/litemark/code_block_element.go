package litemark

import "github.com/desktopgame/ckro/internal/tui/model"

type CodeBlockElement struct {
	Text          string
	Style         *model.Style
	StartPosition model.Position
	EndPosition   model.Position
	Children      []model.Element
}

func (c *CodeBlockElement) GetStyle() *model.Style {
	return nil
}

func (c *CodeBlockElement) GetText() string {
	return c.Text
}

func (c *CodeBlockElement) GetStartPosition() model.Position {
	return c.StartPosition
}

func (c *CodeBlockElement) GetEndPosition() model.Position {
	return c.EndPosition
}

func (c *CodeBlockElement) GetElement(index int) model.Element {
	return c.Children[index]
}

func (c *CodeBlockElement) GetElementCount() int {
	return len(c.Children)
}
