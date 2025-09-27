package litemark

import "github.com/desktopgame/ckro/internal/tui/model"

type TextElement struct {
	Range    model.Range
	Children []model.Element
}

func (t *TextElement) GetRange(index int) model.Range {
	return t.Range
}

func (t *TextElement) GetRangeCount() int {
	return 1
}

func (t *TextElement) GetElement(index int) model.Element {
	return t.Children[index]
}

func (t *TextElement) GetElementCount() int {
	return len(t.Children)
}
