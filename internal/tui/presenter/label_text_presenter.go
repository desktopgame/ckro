package presenter

import (
	"github.com/gdamore/tcell/v2"
)

type LabelTextPresenter struct {
	Text string
}

func (label *LabelTextPresenter) Present(view View) {
	doc := view.GetDocument()
	doc.Init()
	doc.InsertString(label.Text)
}

func (label *LabelTextPresenter) Handle(view View, ev tcell.Event) {
}
