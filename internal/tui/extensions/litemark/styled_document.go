package litemark

import (
	"github.com/desktopgame/ckro/internal/optional"
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/gdamore/tcell/v2"
)

type StyledDocument struct {
	model.PlainDocument

	cache        []model.Element
	cacheVersion uint
}

func (doc *StyledDocument) text2Element(block *Text) model.Element {
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
		if _, ok := aInline.(*Link); ok {
			spanIndex = 0
		}
		if _, ok := aInline.(*Image); ok {
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
		case *Image:
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

	inl1 := block.Inlines[0].BaseInline().Spans[0]
	inl2 := block.Inlines[len(block.Inlines)-1].BaseInline().Spans[0]

	return &TextElement{
		Range: model.Range{
			StartPosition: model.Position{
				Row:    block.LineIndex,
				Column: inl1.StartColumn,
			},
			EndPosition: model.Position{
				Row:    block.LineIndex,
				Column: inl2.EndColumn,
			},
		},
		Children: texts,
	}
}

func (doc *StyledDocument) renderElement(blocks []AbstractBlock) []model.Element {
	elements := []model.Element{}

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
			if block.LineCount > 2 {

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

				langRange := model.Range{
					StartPosition: model.Position{
						Row:    block.LineIndex,
						Column: block.Span.StartColumn,
					},
					EndPosition: model.Position{
						Row:    block.LineIndex,
						Column: block.Span.EndColumn,
					},
				}
				lang := ""
				if !langRange.IsZero() {
					sg := doc.Read(langRange)
					lang = sg.GetLine(0)
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
					Lang:     lang,
					Children: codeLines,
				})
			} else {
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
					Children: []model.Element{
						&InlineElement{
							Ranges: []model.Range{
								{
									StartPosition: model.Position{
										Row:    block.LineIndex,
										Column: 0,
									},
									EndPosition: model.Position{
										Row:    block.LineIndex,
										Column: doc.GetLineBytes(block.LineIndex),
									},
								},
								{
									StartPosition: model.Position{
										Row:    block.LineIndex,
										Column: 0,
									},
									EndPosition: model.Position{
										Row:    block.LineIndex,
										Column: doc.GetLineBytes(block.LineIndex),
									},
								},
							},
						},
					},
				})
				elements = append(elements, &TextElement{
					Range: model.Range{
						StartPosition: model.Position{
							Row:    block.LineIndex + 1,
							Column: 0,
						},
						EndPosition: model.Position{
							Row:    block.LineIndex + 1,
							Column: len(doc.GetLineString(block.LineIndex + 1)),
						},
					},
					Children: []model.Element{
						&InlineElement{
							Ranges: []model.Range{
								{
									StartPosition: model.Position{
										Row:    block.LineIndex + 1,
										Column: 0,
									},
									EndPosition: model.Position{
										Row:    block.LineIndex + 1,
										Column: doc.GetLineBytes(block.LineIndex + 1),
									},
								},
								{
									StartPosition: model.Position{
										Row:    block.LineIndex + 1,
										Column: 0,
									},
									EndPosition: model.Position{
										Row:    block.LineIndex + 1,
										Column: doc.GetLineBytes(block.LineIndex + 1),
									},
								},
							},
						},
					},
				})
			}
		case *FoldBlock:
			if block.LineCount > 2 {
				r := model.Range{
					StartPosition: model.Position{
						Row:    block.LineIndex + 1,
						Column: 0,
					},
					EndPosition: model.Position{
						Row:    block.LineIndex + (block.LineCount - 2),
						Column: doc.GetLineBytes(block.LineIndex + (block.LineCount - 2)),
					},
				}
				sg := doc.Read(r)
				lines := []string{}
				for i := 0; i < sg.GetLineCount(); i++ {
					lines = append(lines, sg.GetLine(i))
				}
				sr := &StringReader{Source: lines}
				aBlocks := Parse(sr)
				for _, aBlock := range aBlocks {
					if table, ok := aBlock.(*Table); ok {
						for _, h := range table.Headers {
							h.Block.LineIndex += block.LineIndex + 1
						}
						for _, r := range table.Rows {
							for _, c := range r.Columns {
								c.LineIndex += block.LineIndex + 1
							}
						}
					}
					aBlock.BaseBlock().LineIndex += block.LineIndex + 1
				}

				elements = append(elements, &model.FoldBlockElement{
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
					Children: doc.renderElement(aBlocks),
					Level:    block.Level,
				})
			} else {
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
					Children: []model.Element{
						&InlineElement{
							Ranges: []model.Range{
								{
									StartPosition: model.Position{
										Row:    block.LineIndex,
										Column: 0,
									},
									EndPosition: model.Position{
										Row:    block.LineIndex,
										Column: doc.GetLineBytes(block.LineIndex),
									},
								},
								{
									StartPosition: model.Position{
										Row:    block.LineIndex,
										Column: 0,
									},
									EndPosition: model.Position{
										Row:    block.LineIndex,
										Column: doc.GetLineBytes(block.LineIndex),
									},
								},
							},
						},
					},
				})
				elements = append(elements, &TextElement{
					Range: model.Range{
						StartPosition: model.Position{
							Row:    block.LineIndex + 1,
							Column: 0,
						},
						EndPosition: model.Position{
							Row:    block.LineIndex + 1,
							Column: len(doc.GetLineString(block.LineIndex + 1)),
						},
					},
					Children: []model.Element{
						&InlineElement{
							Ranges: []model.Range{
								{
									StartPosition: model.Position{
										Row:    block.LineIndex + 1,
										Column: 0,
									},
									EndPosition: model.Position{
										Row:    block.LineIndex + 1,
										Column: doc.GetLineBytes(block.LineIndex + 1),
									},
								},
								{
									StartPosition: model.Position{
										Row:    block.LineIndex + 1,
										Column: 0,
									},
									EndPosition: model.Position{
										Row:    block.LineIndex + 1,
										Column: doc.GetLineBytes(block.LineIndex + 1),
									},
								},
							},
						},
					},
				})
			}
		case *Table:
			tableCells := []model.Element{}
			tableColumns := 0
			for _, h := range block.Headers {
				tableCells = append(tableCells, doc.text2Element(h))
				tableColumns++
			}

			for _, t := range block.Rows {
				for _, c := range t.Columns {
					tableCells = append(tableCells, doc.text2Element(c))
				}
			}
			elements = append(elements, &TableElement{
				Ranges: []model.Range{
					{
						StartPosition: model.Position{
							Row:    block.LineIndex,
							Column: 0,
						},
						EndPosition: model.Position{
							Row:    block.LineIndex + (block.LineCount - 1),
							Column: doc.GetLineBytes(block.LineIndex + (block.LineCount - 1)),
						},
					},
					{
						StartPosition: model.Position{
							Row:    block.LineIndex + 1,
							Column: 0,
						},
						EndPosition: model.Position{
							Row:    block.LineIndex + 1,
							Column: doc.GetLineBytes(block.LineIndex + 1),
						},
					},
				},
				Children: tableCells,
				Columns:  tableColumns,
				Aligns:   block.Aligns,
			})
		case *Text:
			elements = append(elements, doc.text2Element(block))
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
	return elements
}

func (doc *StyledDocument) doRender() []model.Element {
	blocks := Parse(doc)
	elements := doc.renderElement(blocks)

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

func (doc *StyledDocument) Render() []model.Element {
	if doc.cacheVersion == 0 || (doc.cacheVersion != doc.GetVersion()) {
		doc.cache = doc.doRender()
	}
	doc.cacheVersion = doc.GetVersion()
	return doc.cache
}
