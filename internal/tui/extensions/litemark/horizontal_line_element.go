package litemark

import "github.com/desktopgame/ckro/internal/tui/model"

type HorizontalLineElement struct {
	Range model.Range
}

func (h *HorizontalLineElement) GetRange(index int) model.Range {
	return h.Range
}

func (h *HorizontalLineElement) GetRangeCount() int {
	return 1
}

func (h *HorizontalLineElement) GetElement(index int) model.Element {
	return nil
}

func (h *HorizontalLineElement) GetElementCount() int {
	return 0
}
