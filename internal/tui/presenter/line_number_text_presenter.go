package presenter

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

type LineNumberTextPresenter struct {
	TargetView View
}

func (ln *LineNumberTextPresenter) Present(view View) {
	if ln.TargetView == nil {
		return
	}

	view.TextClear()
	doc := view.GetDocument()
	targetDoc := ln.TargetView.GetDocument()

	// get the count of lines.
	lineCount := targetDoc.GetBuffer().GetLineCount()
	if lineCount == 0 {
		lineCount = 1
	}

	// calculate digits of line number.
	maxDigits := len(fmt.Sprintf("%d", lineCount))

	// collect lines.
	scrollY := ln.TargetView.GetScrollY()
	segments := make([]Segment, 0)

	for segment := range ln.TargetView.BreakIter() {
		segments = append(segments, segment)
	}

	// show line number.
	for i, segment := range segments {
		if segment.ViewLine >= scrollY {
			lineNumber := fmt.Sprintf("%*d", maxDigits, segment.ModelLine+1)
			if segment.IsGhostLine {
				lineNumber = "~"
			}
			doc.InsertString(lineNumber)

			if i < len(segments)-1 {
				doc.InsertLine()
			}
		}
	}

	ln.syncCursorPosition(view, targetDoc)
}

// syncCursorPosition is synchronizes cursor position with target view
func (ln *LineNumberTextPresenter) syncCursorPosition(view View, targetDoc interface{}) {
	if doc, ok := targetDoc.(interface{ GetCursorRow() int }); ok {
		cursorRow := doc.GetCursorRow()

		lineNumDoc := view.GetDocument()
		lineNumDoc.MoveReset()

		for i := 0; i < cursorRow && i < lineNumDoc.GetBuffer().GetLineCount()-1; i++ {
			lineNumDoc.MoveDown()
		}

		for lineNumDoc.GetCursorColumn() > 0 {
			lineNumDoc.MoveLeft()
		}
	}
}

func (ln *LineNumberTextPresenter) Handle(view View, ev tcell.Event) {
	view.CursorUpdate()
}

func (ln *LineNumberTextPresenter) ShowCursor() bool {
	return false
}

func (ln *LineNumberTextPresenter) IsFocusable() bool {
	return false
}
