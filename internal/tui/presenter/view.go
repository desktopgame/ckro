package presenter

import (
	"iter"

	"github.com/desktopgame/ckro/internal/tui/model"
)

type Segment struct {
	Text      string
	ModelLine int
	ViewLine  int
}

type View interface {
	TextFrame()
	TextVertical()
	TextHorizontal()
	TextClear()
	CursorUpdate()
	BreakIter() iter.Seq[Segment]
	GetDocument() *model.Document
}
