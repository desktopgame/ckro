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
							Column: len(doc.GetLineString(block.LineIndex)),
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
				if len(doc.GetLineString(lineIndex)) == 0 {
					codeLines = append(codeLines, &BlankLineElement{
						Range: model.Range{
							StartPosition: model.Position{
								Row:    lineIndex,
								Column: 0,
							},
							EndPosition: model.Position{
								Row:    lineIndex,
								Column: 0,
							},
						},
					})
				} else {
					codeLines = append(codeLines, &TextElement{
						Range: model.Range{
							StartPosition: model.Position{
								Row:    lineIndex,
								Column: 0,
							},
							EndPosition: model.Position{
								Row:    lineIndex,
								Column: len(doc.GetLineString(lineIndex)),
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
											Column: len(doc.GetLineString(lineIndex)),
										},
									},
									{
										StartPosition: model.Position{
											Row:    lineIndex,
											Column: 0,
										},
										EndPosition: model.Position{
											Row:    lineIndex,
											Column: len(doc.GetLineString(lineIndex)),
										},
									},
								},
							},
						},
					})
				}
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
							Column: len(doc.GetLineString(block.LineIndex + block.LineCount - 1)),
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

				pad := 0
				isBold := false
				isItalic := false
				isUnderline := false
				fg := optional.None[tcell.Color]()
				bg := optional.None[tcell.Color]()

				switch aInline.(type) {
				case *Bold:
					pad = 2
					isBold = true
				case *Italic:
					pad = 1
					isItalic = true
				case *Strike:
					pad = 2
				case *Link:
					isUnderline = true
					fg = optional.Some(tcell.ColorBlue)
				case *Code:
					pad = 1
					fg = optional.Some(tcell.ColorBlack)
					bg = optional.Some(tcell.ColorWhite)
				}

				texts = append(texts, &InlineElement{
					Ranges:      ranges,
					Pad:         pad,
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
						Column: len(doc.GetLineString(block.LineIndex)),
					},
				},
				Children: texts,
			})
		case *HorizontalLine:
			elements = append(elements, &HorizontalLineElement{
				Range: model.Range{
					StartPosition: model.Position{
						Row:    block.LineIndex,
						Column: 0,
					},
					EndPosition: model.Position{
						Row:    block.LineIndex,
						Column: len(doc.GetLineString(block.LineIndex)),
					},
				},
			})
		}
	}
	for i := 0; i < 10; i++ {
		bytes := doc.GetLineBytes(doc.GetLineCount() - 1)
		elements = append(elements, &model.GhostElement{
			Index: i,
			Range: model.Range{
				StartPosition: model.Position{
					Row:    doc.GetLineCount() - 1,
					Column: 0,
				},
				EndPosition: model.Position{
					Row:    doc.GetLineCount() - 1,
					Column: bytes,
				},
			},
		})
	}
	return elements
}
