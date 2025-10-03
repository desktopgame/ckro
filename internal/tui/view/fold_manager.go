package view

import "github.com/desktopgame/ckro/internal/tui/model"

type FoldManager interface {
	IsFolded(doc model.Document, e model.Element) bool
}
