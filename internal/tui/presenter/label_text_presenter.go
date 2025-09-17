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
	doc := view.GetDocument()
	doc.Init()

	if label.AlignCenter {
		width := view.GetWidth()
		lines := view.GetHeight()
		if lines > 3 {
			for i := 0; i < lines/2; i++ {
				doc.InsertLine()
			}

			length := text.DisplayWidth(label.Text)
			if length >= width {
				doc.InsertString(label.Text)
			} else {
				for i := 0; i < (width-length)/2; i++ {
					doc.InsertString(" ")
				}
				doc.InsertString(label.Text)
			}

		} else {
			doc.InsertString(label.Text)
		}
	} else {
		doc.InsertString(label.Text)
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
