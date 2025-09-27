package litemark

import "github.com/desktopgame/ckro/internal/tui/model"

type TextElement struct {
	Text          string
	Style         *model.Style
	StartPosition model.Position
	EndPosition   model.Position
	Children      []model.Element
}

func (t *TextElement) GetStyle() *model.Style {
	return nil
}

func (t *TextElement) GetText() string {
	return t.Text
}

func (t *TextElement) GetStartPosition() model.Position {
	return t.StartPosition
}

func (t *TextElement) GetEndPosition() model.Position {
	return t.EndPosition
}

func (t *TextElement) GetElement(index int) model.Element {
	return t.Children[index]
}

func (t *TextElement) GetElementCount() int {
	return len(t.Children)
}
