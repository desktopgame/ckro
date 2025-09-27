package litemark

import (
	"github.com/desktopgame/ckro/internal/text"
	"github.com/desktopgame/ckro/internal/tui/model"
)

type StyledDocument struct {
	model.PlainDocument
}

func (doc *StyledDocument) Render() []model.Element {
	elements := []model.Element{}

	blocks := Parse(doc)
	for _, aBlock := range blocks {
		switch block := aBlock.(type) {
		case *Text:
			texts := []model.Element{}
			for _, aInline := range block.Inlines {
				switch inline := aInline.(type) {
				case *Italic:
					texts = append(texts, &InlineElement{
						Text: GetText(doc, block.LineIndex, inline.Spans[1]),
						StartPosition: model.Position{
							Row:    block.LineIndex,
							Column: inline.Spans[0].StartColumn,
						},
						EndPosition: model.Position{
							Row:    block.LineIndex,
							Column: inline.Spans[0].EndColumn - 1,
						},
						Style: &model.Style{
							IsItalic: true,
						},
					})
				case *Bold:
					texts = append(texts, &InlineElement{
						Text: GetText(doc, block.LineIndex, inline.Spans[1]),
						StartPosition: model.Position{
							Row:    block.LineIndex,
							Column: inline.Spans[0].StartColumn,
						},
						EndPosition: model.Position{
							Row:    block.LineIndex,
							Column: inline.Spans[0].EndColumn - 1,
						},
						Style: &model.Style{
							IsBold: true,
						},
					})
				case *PlainText:
					texts = append(texts, &InlineElement{
						Text: GetText(doc, block.LineIndex, inline.Spans[0]),
						StartPosition: model.Position{
							Row:    block.LineIndex,
							Column: inline.Spans[0].StartColumn,
						},
						EndPosition: model.Position{
							Row:    block.LineIndex,
							Column: inline.Spans[0].EndColumn - 1,
						},
						Style: &model.Style{},
					})
				}
			}
			elements = append(elements, &TextElement{
				StartPosition: model.Position{
					Row:    block.LineIndex,
					Column: 0,
				},
				EndPosition: model.Position{
					Row:    block.LineIndex,
					Column: text.GraphemeLength(doc.GetLineAt(block.LineIndex)) - 1,
				},
				Children: texts,
			})
		}
	}
	return elements
}
