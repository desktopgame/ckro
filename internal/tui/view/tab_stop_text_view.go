package view

import (
	"github.com/desktopgame/ckro/internal/tui/model"
)

type TabStopTextView interface {
	DrawWithTabStop(ctx Context, textLayout *TextLayout, renderer Renderer, column int) int
	WidthWithTabStop(ctx Context, e model.Element, column int) int
}
