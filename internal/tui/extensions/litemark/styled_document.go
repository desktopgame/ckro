package litemark

import (
	"github.com/desktopgame/ckro/internal/optional"
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/gdamore/tcell/v2"
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
		case *BlankLine:
			elements = append(elements, &BlankLineElement{
				Range: model.Range{
					StartPosition: model.Position{
						Row:    block.LineIndex,
						Column: 0,
					},
					EndPosition: model.Position{
						Row:    block.LineIndex,
						Column: 0,
					},
				},
			})
		case *CodeBlock:
			codeLines := []model.Element{}
			for i := 0; i < block.LineCount-2; i++ {
				lineIndex := block.LineIndex + i + 1
				codeLines = append(codeLines, &TextElement{
					Range: model.Range{
						StartPosition: model.Position{
							Row:    lineIndex,
							Column: 0,
						},
						EndPosition: model.Position{
							Row:    lineIndex,
							Column: len(doc.GetLineAt(lineIndex)),
						},
					},
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
				Ranges: []model.Range{
					{
						StartPosition: model.Position{
							Row:    block.LineIndex,
							Column: 0,
						},
						EndPosition: model.Position{
							Row:    block.LineIndex + block.LineCount - 1,
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

				isBold := false
				isItalic := false
				isUnderline := false
				fg := optional.None[tcell.Color]()
				bg := optional.None[tcell.Color]()

				switch aInline.(type) {
				case *Bold:
					isBold = true
				case *Italic:
					isItalic = true
				case *Link:
					isUnderline = true
					fg = optional.Some(tcell.ColorBlue)
				case *Code:
					fg = optional.Some(tcell.ColorBlack)
					bg = optional.Some(tcell.ColorWhite)
				}

				texts = append(texts, &InlineElement{
					Ranges:      ranges,
					IsBold:      isBold,
					IsItalic:    isItalic,
					IsUnderline: isUnderline,
					Foreground:  fg,
					Background:  bg,
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
