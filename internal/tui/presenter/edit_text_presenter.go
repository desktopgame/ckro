package presenter

import (
	"github.com/desktopgame/ckro/internal/text"
	"github.com/gdamore/tcell/v2"
)

type EditTextPresenter struct {
	inputBuffer []rune
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
		case tcell.KeyRune:
			edit.inputBuffer = append(edit.inputBuffer, e.Rune())
			inputString := string(edit.inputBuffer)
			if text.GraphemeLength(inputString) == 1 {
				doc.InsertString(inputString)
				edit.inputBuffer = []rune{}
			}
		}
	}

	view.CursorUpdate()
}
