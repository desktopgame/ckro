package litemark

import "github.com/desktopgame/ckro/internal/tui/model"

type TableElement struct {
	Ranges   []model.Range
	Children []model.Element
	Columns  int
	Aligns   []int
}

func (t *TableElement) GetRange(index int) model.Range {
	return t.Ranges[index]
}

func (t *TableElement) GetRangeCount() int {
	return len(t.Ranges)
}

func (t *TableElement) GetElement(index int) model.Element {
	return t.Children[index]
}

func (t *TableElement) GetElementCount() int {
	return len(t.Children)
}
