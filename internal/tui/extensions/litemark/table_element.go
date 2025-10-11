package litemark

import "github.com/desktopgame/ckro/internal/tui/model"

type TableElement struct {
	Range    model.Range
	Children []model.Element
}

func (t *TableElement) GetRange(index int) model.Range {
	return t.Range
}

func (t *TableElement) GetRangeCount() int {
	return 1
}

func (t *TableElement) GetElement(index int) model.Element {
	return t.Children[index]
}

func (t *TableElement) GetElementCount() int {
	return len(t.Children)
}
