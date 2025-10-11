package litemark

import "github.com/desktopgame/ckro/internal/tui/model"

type TableElement struct {
	Range model.Range
	Level int
}

func (t *TableElement) GetRange(index int) model.Range {
	return t.Range
}

func (t *TableElement) GetRangeCount() int {
	return 1
}

func (t *TableElement) GetElement(index int) model.Element {
	return nil
}

func (t *TableElement) GetElementCount() int {
	return 0
}
