package presenter

import (
	"github.com/gdamore/tcell/v2"
)

type ScrollBarTextPresenter struct {
	TargetView View
}

func (sb *ScrollBarTextPresenter) Present(view View) {
	if sb.TargetView == nil {
		return
	}

	view.TextClear()
	// doc := view.GetDocument()
	// targetDoc := sb.TargetView.GetDocument()

	scrollY := sb.TargetView.GetScrollY()
	viewHeight := sb.TargetView.GetHeight()

	// get the count of lines.
	totalLines := sb.TargetView.GetViewHeight()
	if totalLines == 0 {
		totalLines = 1
	}

	scrollBarHeight := viewHeight
	maxScrollY := max(0, totalLines-viewHeight)

	// show scrollbar
	for i := 0; i < scrollBarHeight; i++ {
		var char rune

		if totalLines <= viewHeight {
			// when scroll is not needed
			char = '│'
		} else {
			normalizedScrollY := max(0, min(scrollY, maxScrollY))
			thumbSize := max(1, (viewHeight*scrollBarHeight)/totalLines)

			thumbStart := 0
			if maxScrollY > 0 {
				thumbStart = (normalizedScrollY * (scrollBarHeight - thumbSize)) / maxScrollY
			}
			thumbEnd := thumbStart + thumbSize

			if i >= thumbStart && i < thumbEnd {
				char = '█'
			} else {
				char = '░'
			}
		}

		view.InsertString(string(char))

		if i < scrollBarHeight-1 {
			view.InsertString("\n")
		}
	}

	//sb.syncCursorPosition(view, targetDoc)
}

// syncCursorPosition is synchronizes cursor position with target view
func (sb *ScrollBarTextPresenter) syncCursorPosition(view View, targetDoc interface{}) {
	// TODO: impl
	//if doc, ok := targetDoc.(interface{ GetCursorRow() int }); ok {
	//	cursorRow := doc.GetCursorRow()
	//
	//	scrollBarDoc := view.GetDocument()
	//	scrollBarDoc.MoveReset()
	//
	//	for i := 0; i < cursorRow && i < scrollBarDoc.GetBuffer().GetLineCount()-1; i++ {
	//		scrollBarDoc.MoveDown()
	//	}
	//
	//	for scrollBarDoc.GetCursorColumn() > 0 {
	//		scrollBarDoc.MoveLeft()
	//	}
	//}
}

func (sb *ScrollBarTextPresenter) Handle(view View, ev tcell.Event) {
	view.CursorUpdate()
}

func (sb *ScrollBarTextPresenter) ShowCursor() bool {
	return false
}

func (sb *ScrollBarTextPresenter) IsFocusable() bool {
	return false
}
