package presenter

import "github.com/desktopgame/ckro/internal/tui/model"

type View interface {
	TextFrame()
	TextVertical()
	TextHorizontal()
	TextClear()
	GetDocument() *model.Document
}
