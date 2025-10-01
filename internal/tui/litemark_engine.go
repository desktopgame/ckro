package tui

import (
	"github.com/desktopgame/ckro/internal/tui/extensions/litemark"
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/desktopgame/ckro/internal/tui/view"
)

type LitemarkEngine struct {
}

func (l *LitemarkEngine) Resolve(e model.Element) view.TextView {
	switch e.(type) {
	case *litemark.HeadingElement:
		return &litemark.HeadingView{}
	case *litemark.TextElement:
		return &litemark.TextView{}
	case *litemark.InlineElement:
		return &litemark.InlineView{}
	case *litemark.CodeBlockElement:
		return &litemark.CodeBlockView{}
	case *litemark.BlankLineElement:
		return &litemark.BlankLineView{}
	case *litemark.HorizontalLineElement:
		return &litemark.HorizontalLineView{}
	// Legacy elements
	case *model.GhostElement:
		return &view.GhostView{}
	case *model.PlainElement:
		return &view.PlainTextView{}
	}
	return nil
}

func (l *LitemarkEngine) ProvideInputHandler(e model.Element) TextInputHandler {
	return nil
}
