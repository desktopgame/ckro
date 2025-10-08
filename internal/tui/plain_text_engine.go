package tui

import (
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/desktopgame/ckro/internal/tui/view"
)

type PlainTextEngine struct {
}

func (p *PlainTextEngine) Resolve(e model.Element) view.TextView {
	switch e.(type) {
	case *model.PlainElement:
		return &view.PlainTextView{}
	case *model.GhostElement:
		return &view.GhostView{}
	}
	return &view.PlainTextView{}
}
