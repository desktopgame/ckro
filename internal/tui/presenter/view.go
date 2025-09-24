package presenter

import (
	"iter"

	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/desktopgame/ckro/internal/tui/view"
)

type Segment struct {
	TextLayout    *view.TextLayout
	ModelLine     int
	ViewLine      int
	LocalViewLine int
}

type View interface {
	TextFrame()
	TextVertical()
	TextHorizontal()
	TextClear()
	CursorUpdate()
	BreakIter() iter.Seq[Segment]
	GetDocument() model.Document
	GetWidth() int
	GetHeight() int
	GetScrollX() int
	GetScrollY() int
}
