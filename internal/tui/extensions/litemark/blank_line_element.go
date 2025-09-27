package litemark

import "github.com/desktopgame/ckro/internal/tui/model"

type BlankLineElement struct {
	Ranges []model.Range
	Level  int
}

func (b *BlankLineElement) GetRange(index int) model.Range {
	return model.Range{}
}

func (b *BlankLineElement) GetRangeCount() int {
	return 0
}

func (b *BlankLineElement) GetElement(index int) model.Element {
	return nil
}

func (b *BlankLineElement) GetElementCount() int {
	return 0
}
