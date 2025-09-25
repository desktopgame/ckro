package tui

import (
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/desktopgame/ckro/internal/tui/view"
)

type PlainTextEngine struct {
}

func (p *PlainTextEngine) Resolve(e model.Element) view.TextView {
	switch e.(type) {
	case *model.LineContainerElement:
		return &view.LineContainerView{}
	case *model.InlineElement:
		return &view.InlineTextView{}
	case *model.ParagraphElement:
		return &view.InlineTextView{}
	}
	return nil
}

func (p *PlainTextEngine) ProvideInputHandler(e model.Element) TextInputHandler {
	return nil
}
