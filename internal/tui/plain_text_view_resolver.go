package tui

import (
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/desktopgame/ckro/internal/tui/view"
)

// PlainTextViewResolver implement TextViewResolver for plain text.
type PlainTextViewResolver struct {
}

// Resolve returns TextView, supports plain text.
func (p *PlainTextViewResolver) Resolve(e model.Element) view.TextView {
	switch e.(type) {
	case *model.PlainElement:
		return &view.PlainTextView{}
	case *model.GhostElement:
		return &view.GhostView{}
	}
	return &view.PlainTextView{}
}
