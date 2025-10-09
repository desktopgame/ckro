package view

import "github.com/desktopgame/ckro/internal/tui/model"

// TextViewResolver resolve TextView from Element.
// in normally, determine by element type.
type TextViewResolver interface {
	Resolve(e model.Element) TextView
}
