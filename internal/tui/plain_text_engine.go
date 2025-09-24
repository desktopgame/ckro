package tui

import (
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/desktopgame/ckro/internal/tui/view"
)

type PlainTextEngine struct {
}

func (p *PlainTextEngine) ProvideView(e model.Element) view.TextView {
	return &view.PlainTextView{}
}

func (p *PlainTextEngine) ProvideInputHandler(e model.Element) TextInputHandler {
	return nil
}
