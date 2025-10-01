package presenter

import (
	"github.com/gdamore/tcell/v2"
)

type ListTextPresenter struct {
	Items         []string
	SelectedIndex int
	CursorChar    rune
	Prefix        string
}

func (lp *ListTextPresenter) Present(view View) {
	// TODO: impl
	//view.TextClear()
	//doc := view.GetDocument()
	//
	//cursorChar := lp.CursorChar
	//if cursorChar == 0 {
	//	cursorChar = '>'
	//}
	//
	//prefix := lp.Prefix
	//if prefix == "" {
	//	prefix = " "
	//}
	//
	//// show list items.
	//for i, item := range lp.Items {
	//	if i == lp.SelectedIndex {
	//		doc.InsertString(string(cursorChar))
	//	} else {
	//		doc.InsertString(" ")
	//	}
	//
	//	doc.InsertString(prefix + item)
	//
	//	if i < len(lp.Items)-1 {
	//		doc.InsertLine()
	//	}
	//}

	lp.setCursorToSelectedItem(view)
}

// setCursorToSelectedItem is sets cursor to the selected item
func (lp *ListTextPresenter) setCursorToSelectedItem(view View) {
	// TODO: impl
	//doc := view.GetDocument()
	//doc.MoveReset()
	//
	//for i := 0; i < lp.SelectedIndex && i < len(lp.Items)-1; i++ {
	//	doc.MoveDown()
	//}
	//
	//for doc.GetCursorColumn() > 0 {
	//	doc.MoveLeft()
	//}
	//
	//view.CursorUpdate()
}

func (lp *ListTextPresenter) Handle(view View, ev tcell.Event) {
	if keyEvent, ok := ev.(*tcell.EventKey); ok {
		switch keyEvent.Key() {
		case tcell.KeyUp:
			if lp.SelectedIndex > 0 {
				lp.SelectedIndex--
				lp.Present(view)
			}
		case tcell.KeyDown:
			if lp.SelectedIndex < len(lp.Items)-1 {
				lp.SelectedIndex++
				lp.Present(view)
			}
		case tcell.KeyHome:
			if len(lp.Items) > 0 {
				lp.SelectedIndex = 0
				lp.Present(view)
			}
		case tcell.KeyEnd:
			if len(lp.Items) > 0 {
				lp.SelectedIndex = len(lp.Items) - 1
				lp.Present(view)
			}
		}
	}

	view.CursorUpdate()
}

func (lp *ListTextPresenter) ShowCursor() bool {
	return false
}

func (lp *ListTextPresenter) IsFocusable() bool {
	return true
}

// GetSelectedItem returns the currently selected item
func (lp *ListTextPresenter) GetSelectedItem() string {
	if lp.SelectedIndex >= 0 && lp.SelectedIndex < len(lp.Items) {
		return lp.Items[lp.SelectedIndex]
	}
	return ""
}

// GetSelectedIndex returns the currently selected index
func (lp *ListTextPresenter) GetSelectedIndex() int {
	return lp.SelectedIndex
}

// SetSelectedIndex is sets the selected index
func (lp *ListTextPresenter) SetSelectedIndex(index int) {
	if index >= 0 && index < len(lp.Items) {
		lp.SelectedIndex = index
	}
}

// AddItem is adds a new item to the list
func (lp *ListTextPresenter) AddItem(item string) {
	lp.Items = append(lp.Items, item)
}

// RemoveItem is removes an item at the specified index
func (lp *ListTextPresenter) RemoveItem(index int) {
	if index >= 0 && index < len(lp.Items) {
		lp.Items = append(lp.Items[:index], lp.Items[index+1:]...)

		if lp.SelectedIndex >= len(lp.Items) && len(lp.Items) > 0 {
			lp.SelectedIndex = len(lp.Items) - 1
		} else if len(lp.Items) == 0 {
			lp.SelectedIndex = 0
		}
	}
}

// Clear is removes all items from the list
func (lp *ListTextPresenter) Clear() {
	lp.Items = nil
	lp.SelectedIndex = 0
}
