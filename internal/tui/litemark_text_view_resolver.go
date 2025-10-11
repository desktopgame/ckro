package tui

import (
	"github.com/desktopgame/ckro/internal/tui/extensions/litemark"
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/desktopgame/ckro/internal/tui/view"
)

// LitemarkTextViewResolver implement TextViewResolver for litemark language.
type LitemarkTextViewResolver struct {
}

// Resolve returns TextView, supports litemark language.
func (l *LitemarkTextViewResolver) Resolve(e model.Element) view.TextView {
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
	case *litemark.TableElement:
		return &litemark.TableView{}
	case *litemark.TableHeaderElement:
		return &litemark.TableHeaderView{}
	case *litemark.TableRowElement:
		return &litemark.TableRowView{}
	// Legacy elements
	case *model.FoldBlockElement:
		return &view.FoldBlockView{}
	case *model.GhostElement:
		return &view.GhostView{}
	case *model.PlainElement:
		return &view.PlainTextView{}
	}
	return &view.PlainTextView{}
}
