package litemark

import (
	"github.com/desktopgame/ckro/internal/optional"
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/gdamore/tcell/v2"
)

type InlineElement struct {
	Ranges      []model.Range
	Children    []model.Element
	Pad         int
	IsBold      bool
	IsItalic    bool
	IsUnderline bool
	Foreground  optional.Optional[tcell.Color]
	Background  optional.Optional[tcell.Color]
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
