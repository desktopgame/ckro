package tui

import "github.com/desktopgame/ckro/internal/tui/model"

type PlainTextEngine struct {
}

func (p *PlainTextEngine) ProvideView(e model.Element) TextView {
	return &PlainTextView{}
}

func (p *PlainTextEngine) ProvideInputHandler(e model.Element) TextInputHandler {
	return nil
}
