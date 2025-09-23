package tui

import "github.com/desktopgame/ckro/internal/tui/model"

type TextEngine interface {
	ProvideView(e model.Element) TextView
	ProvideInputHandler(e model.Element) TextInputHandler
}
