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
	targetDoc := ln.TargetView.GetDocument()

	// get the count of lines.
	lineCount := targetDoc.GetLineCount()
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
	count := 0
	for _, segment := range segments {
		if segment.ViewLine >= scrollY {
			lineNumber := fmt.Sprintf("%*d", maxDigits, segment.ModelLine+1)
			if segment.IsGhostLine {
				lineNumber = "~"
			}
			view.InsertString(lineNumber)

			if count < view.GetHeight() {
				view.InsertString("\n")
			}
			count++
			if count >= view.GetHeight() {
				break
			}
		}
	}

	//view.CursorUpdate()
	ln.syncCursorPosition(view, ln.TargetView)
}

// syncCursorPosition is synchronizes cursor position with target view
func (ln *LineNumberTextPresenter) syncCursorPosition(view View, targetView View) {
	// _, y, _, _ := targetView.CursorPosition()

	// view.MoveReset()
	// view.CursorTo(y)
	//	view.CursorUpdate()
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
