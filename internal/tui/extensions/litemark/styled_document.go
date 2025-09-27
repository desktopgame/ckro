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
		case *Heading:
			elements = append(elements, &HeadingElement{
				StartPosition: model.Position{
					Row:    block.LineIndex,
					Column: 0,
				},
				EndPosition: model.Position{
					Row:    block.LineIndex,
					Column: text.GraphemeLength(doc.GetLineAt(block.LineIndex)) - 1,
				},
				Text:  GetText(doc, block.LineIndex, block.Span),
				Level: block.Level,
			})
		case *CodeBlock:
			codeLines := []model.Element{}
			for i := 0; i < block.LineCount-2; i++ {
				lineIndex := block.LineIndex + i + 1
				codeLines = append(codeLines, &TextElement{
					Children: []model.Element{
						&InlineElement{
							Text: doc.GetLineAt(lineIndex),
						},
					},
				})
			}
			elements = append(elements, &CodeBlockElement{
				StartPosition: model.Position{
					Row:    block.LineIndex,
					Column: 0,
				},
				EndPosition: model.Position{
					Row:    block.LineIndex,
					Column: text.GraphemeLength(doc.GetLineAt(block.LineIndex)) - 1,
				},
				Children: codeLines,
			})
		case *Text:
			texts := []model.Element{}
			for _, aInline := range block.Inlines {
				var text string
				var style *model.Style
				switch inline := aInline.(type) {
				case *Italic:
					text = GetText(doc, block.LineIndex, inline.Spans[1])
					style = &model.Style{
						IsItalic: true,
					}
				case *Bold:
					text = GetText(doc, block.LineIndex, inline.Spans[1])
					style = &model.Style{
						IsBold: true,
					}
				case *Code:
					text = GetText(doc, block.LineIndex, inline.Spans[1])
					style = &model.Style{
						Foreground: model.Black,
						Background: model.White,
					}
				case *Link:
					text = GetText(doc, block.LineIndex, inline.Spans[1])
					style = &model.Style{
						Foreground:  model.Blue,
						IsUnderline: true,
					}
				case *Image:
					text = GetText(doc, block.LineIndex, inline.Spans[1])
					style = &model.Style{}
				case *PlainText:
					text = GetText(doc, block.LineIndex, inline.Spans[0])
					style = &model.Style{}
				}

				texts = append(texts, &InlineElement{
					Text: text,
					StartPosition: model.Position{
						Row:    block.LineIndex,
						Column: aInline.BaseInline().Spans[0].StartColumn,
					},
					EndPosition: model.Position{
						Row:    block.LineIndex,
						Column: aInline.BaseInline().Spans[0].EndColumn - 1,
					},
					Style: style,
				})
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
