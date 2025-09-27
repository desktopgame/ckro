package view

import "github.com/desktopgame/ckro/internal/tui/model"

type TabStopTextView interface {
	WidthWithTabStop(textViewResolver TextViewResolver, e model.Element, column int) int
}
