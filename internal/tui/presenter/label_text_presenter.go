package presenter

import (
	"github.com/desktopgame/ckro/internal/text"
	"github.com/gdamore/tcell/v2"
)

type LabelTextPresenter struct {
	Text        string
	AlignCenter bool
}

func (label *LabelTextPresenter) Present(view View) {
	view.TextClear()

	if label.AlignCenter {
		width := view.GetWidth()
		lines := view.GetHeight()
		if lines > 3 {
			for i := 0; i < lines/2; i++ {
				view.InsertString("\n")
			}

			length := text.DisplayWidth(label.Text)
			if length >= width {
				view.InsertString(label.Text)
			} else {
				for i := 0; i < (width-length)/2; i++ {
					view.InsertString(" ")
				}
				view.InsertString(label.Text)
			}

		} else {
			view.InsertString(label.Text)
		}
	} else {
		view.InsertString(label.Text)
	}
}

func (label *LabelTextPresenter) Handle(view View, ev tcell.Event) {
}

func (label *LabelTextPresenter) ShowCursor() bool {
	return false
}

func (label *LabelTextPresenter) IsFocusable() bool {
	return false
}
