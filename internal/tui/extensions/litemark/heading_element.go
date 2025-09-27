package litemark

import "github.com/desktopgame/ckro/internal/tui/model"

type HeadingElement struct {
	Ranges []model.Range
	Level  int
}

func (h *HeadingElement) GetRange(index int) model.Range {
	return h.Ranges[index]
}

func (h *HeadingElement) GetRangeCount() int {
	return len(h.Ranges)
}

func (h *HeadingElement) GetElement(index int) model.Element {
	return nil
}

func (h *HeadingElement) GetElementCount() int {
	return 0
}
