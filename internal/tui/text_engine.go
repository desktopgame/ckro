package tui

import (
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/desktopgame/ckro/internal/tui/view"
)

type TextEngine interface {
	Resolve(e model.Element) view.TextView
	ProvideInputHandler(e model.Element) TextInputHandler
}
