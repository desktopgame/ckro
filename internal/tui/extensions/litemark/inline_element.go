package litemark

import "github.com/desktopgame/ckro/internal/tui/model"

type InlineElement struct {
	Ranges   []model.Range
	Children []model.Element
}

func (il *InlineElement) GetRange(index int) model.Range {
	return il.Ranges[index]
}

func (il *InlineElement) GetRangeCount() int {
	return len(il.Ranges)
}

func (il *InlineElement) GetElement(index int) model.Element {
	return nil
}

func (il *InlineElement) GetElementCount() int {
	return 0
}
