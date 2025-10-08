package presenter

import (
	"github.com/desktopgame/ckro/internal/text"
	"github.com/gdamore/tcell/v2"
)

type EditTextPresenter struct {
	inputBuffer []rune
	OnModified  func()
	ReadOnly    bool
	Selection   bool
}

func (edit *EditTextPresenter) modify() {
	if edit.OnModified != nil {
		edit.OnModified()
	}
}

func (edit *EditTextPresenter) Present(view View) {
}

func (edit *EditTextPresenter) selectionEnd(view View) {
	if edit.Selection {
		edit.Selection = false
		view.SelectionEnd()
	}
}

func (edit *EditTextPresenter) Handle(view View, ev tcell.Event) {
	// doc := view.GetDocument()
	switch e := ev.(type) {
	case *tcell.EventKey:
		switch e.Key() {
		case tcell.KeyUp, tcell.KeyCtrlP:
			if tcell.ModShift&e.Modifiers() > 0 {
				if !edit.Selection {
					edit.Selection = true
					view.SelectionStart()
				}
				view.MoveUp()
			} else {
				edit.selectionEnd(view)
				view.MoveUp()
			}
		case tcell.KeyDown, tcell.KeyCtrlN:
			if tcell.ModShift&e.Modifiers() > 0 {
				if !edit.Selection {
					edit.Selection = true
					view.SelectionStart()
				}
				view.MoveDown()
			} else {
				edit.selectionEnd(view)
				view.MoveDown()
			}
		case tcell.KeyLeft, tcell.KeyCtrlB:
			if tcell.ModShift&e.Modifiers() > 0 {
				if !edit.Selection {
					edit.Selection = true
					view.SelectionStart()
				}
				view.MoveLeft()
			} else {
				edit.selectionEnd(view)
				view.MoveLeft()
			}
		case tcell.KeyRight, tcell.KeyCtrlF:
			if tcell.ModShift&e.Modifiers() > 0 {
				if !edit.Selection {
					edit.Selection = true
					view.SelectionStart()
				}
				view.MoveRight()
			} else {
				edit.selectionEnd(view)
				view.MoveRight()
			}
		case tcell.KeyCtrlA:
			edit.selectionEnd(view)
			view.MoveLineStart()
		case tcell.KeyCtrlE:
			edit.selectionEnd(view)
			view.MoveLineEnd()
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			edit.selectionEnd(view)
			view.RemoveChar()
		case tcell.KeyEnter:
			edit.selectionEnd(view)
			if !view.Submit() {
				view.InsertString("\n")
			}
			edit.modify()
		case tcell.KeyTAB:
			edit.selectionEnd(view)
			view.InsertString("\t")
			edit.modify()
		case tcell.KeyRune:
			if !edit.ReadOnly {
				edit.inputBuffer = append(edit.inputBuffer, e.Rune())
				inputString := string(edit.inputBuffer)
				if text.GraphemeLength(inputString) == 1 {
					edit.selectionEnd(view)
					view.InsertString(inputString)
					edit.inputBuffer = []rune{}
					edit.modify()
				}
			}
		}
	}

	view.CursorUpdate()
}

func (edit *EditTextPresenter) ShowCursor() bool {
	return true
}

func (edit *EditTextPresenter) IsFocusable() bool {
	return true
}
