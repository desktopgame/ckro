package presenter

import (
	"github.com/gdamore/tcell/v2"
)

type HorizontalSeparatorTextPresenter struct {
}

func (h *HorizontalSeparatorTextPresenter) Present(view View) {
	view.TextHorizontal()
}

func (h *HorizontalSeparatorTextPresenter) Handle(view View, ev tcell.Event) {
}

func (h *HorizontalSeparatorTextPresenter) ShowCursor() bool {
	return false
}

func (h *HorizontalSeparatorTextPresenter) IsFocusable() bool {
	return false
}

type VerticalSeparatorTextPresenter struct {
}

func (v *VerticalSeparatorTextPresenter) Present(view View) {
	view.TextVertical()
}

func (v *VerticalSeparatorTextPresenter) Handle(view View, ev tcell.Event) {
}

func (v *VerticalSeparatorTextPresenter) ShowCursor() bool {
	return false
}

func (v *VerticalSeparatorTextPresenter) IsFocusable() bool {
	return false
}
