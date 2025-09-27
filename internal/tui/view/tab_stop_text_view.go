package view

import (
	"github.com/desktopgame/ckro/internal/tui/model"
)

type TabStopTextView interface {
	WidthWithTabStop(ctx Context, e model.Element, column int) int
}
