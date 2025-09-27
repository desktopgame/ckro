package litemark

import (
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
				Ranges: []model.Range{
					{
						StartPosition: model.Position{
							Row:    block.LineIndex,
							Column: 0,
						},
						EndPosition: model.Position{
							Row:    block.LineIndex,
							Column: len(doc.GetLineAt(block.LineIndex)),
						},
					},
					{
						StartPosition: model.Position{
							Row:    block.LineIndex,
							Column: block.Span.StartColumn,
						},
						EndPosition: model.Position{
							Row:    block.LineIndex,
							Column: block.Span.EndColumn,
						},
					},
				},
				Level: block.Level,
			})
		case *CodeBlock:
			codeLines := []model.Element{}
			for i := 0; i < block.LineCount-2; i++ {
				lineIndex := block.LineIndex + i + 1
				codeLines = append(codeLines, &TextElement{
					Children: []model.Element{
						&InlineElement{
							Ranges: []model.Range{
								{
									StartPosition: model.Position{
										Row:    lineIndex,
										Column: 0,
									},
									EndPosition: model.Position{
										Row:    lineIndex,
										Column: len(doc.GetLineAt(lineIndex)),
									},
								},
								{
									StartPosition: model.Position{
										Row:    lineIndex,
										Column: 0,
									},
									EndPosition: model.Position{
										Row:    lineIndex,
										Column: len(doc.GetLineAt(lineIndex)),
									},
								},
							},
						},
					},
				})
			}
			elements = append(elements, &CodeBlockElement{
				Range: model.Range{
					StartPosition: model.Position{
						Row:    block.LineIndex,
						Column: 0,
					},
					EndPosition: model.Position{
						Row:    block.LineIndex,
						Column: len(doc.GetLineAt(block.LineIndex)),
					},
				},
				Children: codeLines,
			})
		case *Text:
			texts := []model.Element{}
			for _, aInline := range block.Inlines {

				ranges := []model.Range{
					{
						StartPosition: model.Position{
							Row:    block.LineIndex,
							Column: aInline.BaseInline().Spans[0].StartColumn,
						},
						EndPosition: model.Position{
							Row:    block.LineIndex,
							Column: aInline.BaseInline().Spans[0].EndColumn,
						},
					},
				}

				spanIndex := 1
				if _, ok := aInline.(*PlainText); ok {
					spanIndex = 0
				}

				ranges = append(ranges,
					model.Range{
						StartPosition: model.Position{
							Row:    block.LineIndex,
							Column: aInline.BaseInline().Spans[spanIndex].StartColumn,
						},
						EndPosition: model.Position{
							Row:    block.LineIndex,
							Column: aInline.BaseInline().Spans[spanIndex].EndColumn,
						},
					},
				)

				texts = append(texts, &InlineElement{
					Ranges: ranges,
				})
			}

			elements = append(elements, &TextElement{
				Range: model.Range{
					StartPosition: model.Position{
						Row:    block.LineIndex,
						Column: 0,
					},
					EndPosition: model.Position{
						Row:    block.LineIndex,
						Column: len(doc.GetLineAt(block.LineIndex)),
					},
				},
				Children: texts,
			})
		}
	}
	return elements
}
