package litemark

import "github.com/desktopgame/ckro/internal/tui/model"

type InlineElement struct {
	Text          string
	Style         *model.Style
	StartPosition model.Position
	EndPosition   model.Position
	Children      []model.Element
}

func (il *InlineElement) GetStyle() *model.Style {
	return il.Style
}

func (il *InlineElement) GetText() string {
	return il.Text
}

func (il *InlineElement) GetStartPosition() model.Position {
	return il.StartPosition
}

func (il *InlineElement) GetEndPosition() model.Position {
	return il.EndPosition
}

func (il *InlineElement) GetElement(index int) model.Element {
	return nil
}

func (il *InlineElement) GetElementCount() int {
	return 0
}
