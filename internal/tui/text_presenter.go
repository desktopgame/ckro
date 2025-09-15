package tui

import (
	"github.com/desktopgame/ckro/internal/tui/presenter"
	"github.com/gdamore/tcell/v2"
)

type TextPresenter interface {
	Present(view presenter.View)
	Handle(view presenter.View, ev tcell.Event)
	ShowCursor() bool
	IsFocusable() bool
}
