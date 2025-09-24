package view

import "github.com/desktopgame/ckro/internal/tui/model"

type TextViewResolver interface {
	Resolve(e model.Element) TextView
}
