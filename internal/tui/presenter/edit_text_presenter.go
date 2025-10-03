package presenter

import (
	"github.com/desktopgame/ckro/internal/text"
	"github.com/gdamore/tcell/v2"
)

type EditTextPresenter struct {
	inputBuffer []rune
	OnModified  func()
	ReadOnly    bool
}

func (edit *EditTextPresenter) modify() {
	if edit.OnModified != nil {
		edit.OnModified()
	}
}

func (edit *EditTextPresenter) Present(view View) {
}

func (edit *EditTextPresenter) Handle(view View, ev tcell.Event) {
	// doc := view.GetDocument()
	switch e := ev.(type) {
	case *tcell.EventKey:
		switch e.Key() {
		case tcell.KeyUp, tcell.KeyCtrlP:
			view.MoveUp()
		case tcell.KeyDown, tcell.KeyCtrlN:
			view.MoveDown()
		case tcell.KeyLeft, tcell.KeyCtrlB:
			view.MoveLeft()
		case tcell.KeyRight, tcell.KeyCtrlF:
			view.MoveRight()
		case tcell.KeyCtrlA:
			view.MoveLineStart()
		case tcell.KeyCtrlE:
			view.MoveLineEnd()
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			view.RemoveChar()
		case tcell.KeyEnter:
			if !view.Submit() {
				view.InsertString("\n")
			}
			edit.modify()
		case tcell.KeyTAB:
			view.InsertString("\t")
			edit.modify()
		case tcell.KeyRune:
			if !edit.ReadOnly {
				edit.inputBuffer = append(edit.inputBuffer, e.Rune())
				inputString := string(edit.inputBuffer)
				if text.GraphemeLength(inputString) == 1 {
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
