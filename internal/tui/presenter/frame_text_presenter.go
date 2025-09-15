package presenter

import "github.com/gdamore/tcell/v2"

type FrameTextPresenter struct {
}

func (f *FrameTextPresenter) Present(view View) {
	view.TextFrame()
}

func (f *FrameTextPresenter) Handle(view View, ev tcell.Event) {
}

func (f *FrameTextPresenter) ShowCursor() bool {
	return false
}

func (f *FrameTextPresenter) IsFocusable() bool {
	return false
}
