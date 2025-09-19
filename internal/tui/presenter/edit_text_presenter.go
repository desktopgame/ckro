package presenter

import (
	"github.com/desktopgame/ckro/internal/text"
	"github.com/gdamore/tcell/v2"
)

type EditTextPresenter struct {
	inputBuffer []rune
	OnModified  func()
}

func (edit *EditTextPresenter) modify() {
	if edit.OnModified != nil {
		edit.OnModified()
	}
}

func (edit *EditTextPresenter) Present(view View) {
}

func (edit *EditTextPresenter) Handle(view View, ev tcell.Event) {
	doc := view.GetDocument()
	switch e := ev.(type) {
	case *tcell.EventKey:
		switch e.Key() {
		case tcell.KeyUp:
			doc.MoveUp()
		case tcell.KeyDown:
			doc.MoveDown()
		case tcell.KeyLeft:
			doc.MoveLeft()
		case tcell.KeyRight:
			doc.MoveRight()
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			doc.RemoveChar()
		case tcell.KeyEnter:
			doc.InsertLine()
			edit.modify()
		case tcell.KeyTAB:
			doc.InsertString("\t")
			edit.modify()
		case tcell.KeyRune:
			edit.inputBuffer = append(edit.inputBuffer, e.Rune())
			inputString := string(edit.inputBuffer)
			if text.GraphemeLength(inputString) == 1 {
				doc.InsertString(inputString)
				edit.inputBuffer = []rune{}
				edit.modify()
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
