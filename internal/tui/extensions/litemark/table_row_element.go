package litemark

import "github.com/desktopgame/ckro/internal/tui/model"

type TableRowElement struct {
	Range    model.Range
	Children []model.Element
}

func (t *TableRowElement) GetRange(index int) model.Range {
	return t.Range
}

func (t *TableRowElement) GetRangeCount() int {
	return 1
}

func (t *TableRowElement) GetElement(index int) model.Element {
	return t.Children[index]
}

func (t *TableRowElement) GetElementCount() int {
	return len(t.Children)
}
