package presenter

import "github.com/desktopgame/ckro/internal/tui/model"

type View interface {
	GetDocument() *model.Document
}
