package tui

import (
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/desktopgame/ckro/internal/tui/view"
)

type MarkdownTextEngine struct {
}

func (m *MarkdownTextEngine) Resolve(e model.Element) view.TextView {
	switch e.(type) {
	// Goldmark Elements
	case *model.DocumentElement:
		return &view.DocumentTextView{}
	case *model.ParagraphElementGM:
		return &view.ParagraphTextViewGM{}
	case *model.HeadingElement:
		return &view.HeadingTextView{}
	case *model.CodeBlockElement:
		return &view.CodeBlockTextView{}
	case *model.BlockquoteElement:
		return &view.BlockquoteTextView{}
	case *model.ListElement:
		return &view.ListTextView{}
	case *model.ListItemElement:
		return &view.ListItemTextView{}
	case *model.ThematicBreakElement:
		return &view.ThematicBreakTextView{}
	case *model.TextElement:
		return &view.TextElementView{}
	case *model.EmphasisElement:
		return &view.EmphasisTextView{}
	case *model.StrongElement:
		return &view.StrongTextView{}
	case *model.CodeElement:
		return &view.CodeElementView{}
	case *model.LinkElement:
		return &view.LinkTextView{}
	case *model.ImageElement:
		return &view.ImageTextView{}
	case *model.SoftBreakElement:
		return &view.SoftBreakTextView{}
	case *model.HardBreakElement:
		return &view.HardBreakTextView{}
	case *model.TableElement:
		return &view.TableTextView{}
	case *model.TableHeaderElement:
		return &view.TableHeaderTextView{}
	case *model.TableRowElementGM:
		return &view.TableRowTextViewGM{}
	case *model.TableCellElement:
		return &view.TableCellTextView{}
	// Legacy elements
	case *model.LineContainerElement:
		return &view.LineContainerView{}
	case *model.InlineElement:
		return &view.InlineTextView{}
	case *model.ParagraphElement:
		return &view.InlineTextView{}
	}
	return nil
}

func (m *MarkdownTextEngine) ProvideInputHandler(e model.Element) TextInputHandler {
	return nil
}
