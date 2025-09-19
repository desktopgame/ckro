package tui

import (
	"github.com/desktopgame/ckro/internal/tui/presenter"
	"github.com/gdamore/tcell/v2"
)

// TextPresenter is show fancy content through TextBox object.
// in example, show frame by text, show tree by text, etc.
// in addition, TextPresenter is define event process.
// input a string when event, or change selected item, or something.
type TextPresenter interface {
	// Present is set shown content to View.
	Present(view presenter.View)

	// Handle is process a event.
	Handle(view presenter.View, ev tcell.Event)

	// ShowCursor returns true if should be show cursor.
	ShowCursor() bool

	// IsFocusable returns true if should be capture the focus.
	IsFocusable() bool
}
