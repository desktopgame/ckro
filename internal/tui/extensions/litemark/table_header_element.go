package litemark

import "github.com/desktopgame/ckro/internal/tui/model"

type TableHeaderElement struct {
	Range model.Range
	Level int
}

func (t *TableHeaderElement) GetRange(index int) model.Range {
	return t.Range
}

func (t *TableHeaderElement) GetRangeCount() int {
	return 1
}

func (t *TableHeaderElement) GetElement(index int) model.Element {
	return nil
}

func (t *TableHeaderElement) GetElementCount() int {
	return 0
}
