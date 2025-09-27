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
			elements = append(elements, &model.PlainElement{
				Text: GetText(doc, block.LineIndex, Span{
					StartColumn: 0,
					EndColumn:   len(doc.GetLineAt(block.LineIndex)),
				}),
				StartPosition: model.Position{
					Row:    block.LineIndex,
					Column: 0,
				},
				EndPosition: model.Position{
					Row:    block.LineIndex,
					Column: text.GraphemeLength(doc.GetLineAt(block.LineIndex)) - 1,
				},
			})
		}
	}
	return elements
}
