package litemark

import "strings"

func strIndexOf(s, sub string, at int) int {
	i := strings.Index(s[at:], sub)
	if i < 0 {
		return -1
	}
	return at + i
}

func byteIndexOf(s string, sub byte, at int) int {
	i := strings.IndexByte(s[at:], sub)
	if i < 0 {
		return -1
	}
	return at + i
}

func GetText(reader Reader, lineIndex int, span Span) string {
	return reader.GetLineString(lineIndex)[span.StartColumn:span.EndColumn]
}

func Parse(reader Reader) []AbstractBlock {
	sc := Scanner{Reader: reader}
	blocks := []AbstractBlock{}

	codeBlockScope := false
	codeBlockMarkerLen := 0
	codeBlockCurrent := &CodeBlock{}

	for sc.Ready() {
		lineIndex := sc.lineIndex
		line := sc.Next()

		// CodeBlock
		if codeBlockScope {
			codeBlockEnded := false
			if len(line) > 0 && line[0] == '`' {
				codeBlockScope = true

				column := 0
				for column < len(line) && line[column] == '`' {
					column++
				}
				if codeBlockMarkerLen == column-1 {
					codeBlockEnded = true
				}
			}
			if codeBlockEnded {
				codeBlockScope = false
				codeBlockMarkerLen = 0
				codeBlockCurrent.LineCount++
			} else {
				codeBlockCurrent.LineCount++
			}
			continue
		}

		// Soft break
		if len(line) == 0 {
			blocks = append(blocks, &BlankLine{
				Block: Block{
					LineIndex: lineIndex,
					LineCount: 1,
				},
			})
			continue
		}

		// Heading
		if line[0] == '#' {
			column := 0
			for column < len(line) && line[column] == '#' {
				column++
			}
			if column < len(line) && line[column] == ' ' {
				blocks = append(blocks, &Heading{
					Block: Block{
						LineIndex: lineIndex,
						LineCount: 1,
					},
					Span: Span{
						StartColumn: column + 1,
						EndColumn:   len(line),
					},
					Level: column,
				})
				continue
			}
		}

		// CodeBlock
		if line[0] == '`' {
			column := 0
			for column < len(line) {
				column++
			}

			if column >= 3 && line == strings.Repeat("`", column) {
				codeBlockScope = true

				codeBlockMarkerLen = column - 1

				codeBlockCurrent = &CodeBlock{
					Block: Block{
						LineIndex: lineIndex,
						LineCount: 1,
					},
					Span: Span{
						StartColumn: column,
						EndColumn:   len(line),
					},
					Level: column,
				}
				blocks = append(blocks, codeBlockCurrent)
				continue
			}
		}

		// Fold
		if line[0] == '{' {
			column := 0
			for column < len(line) && line[column] == '{' {
				column++
			}

			if column >= 3 && line == strings.Repeat("{", column) {
				lineCount := 1
				foldBlock := &FoldBlock{
					Block: Block{
						LineIndex: lineIndex,
					},
					Span: Span{
						StartColumn: column,
						EndColumn:   len(line),
					},
					Level: column,
				}

				foundClose := false
				for sc.Ready() {
					innerLine := sc.Next()

					lineCount++

					if innerLine == strings.Repeat("}", column) {
						foundClose = true
						break
					}
				}
				if !foundClose {
					sc.lineIndex = lineIndex + 1
					blocks = append(blocks, &Text{
						Block: Block{
							LineIndex: lineIndex,
							LineCount: 1,
						},
						Inlines: ParseInline(line),
					})
					continue
				} else {
					foldBlock.LineCount = lineCount
					blocks = append(blocks, foldBlock)
					continue
				}
			}
		}

		if line == "***" || line == "---" {
			blocks = append(blocks, &HorizontalLine{
				Block: Block{
					LineIndex: lineIndex,
					LineCount: 1,
				},
			})
			continue
		}

		// Inline text
		blocks = append(blocks, &Text{
			Block: Block{
				LineIndex: lineIndex,
				LineCount: 1,
			},
			Inlines: ParseInline(line),
		})
	}

	// remaining lines
	if codeBlockScope {
		blocks = blocks[:len(blocks)-1]

		for i := 0; i < codeBlockCurrent.LineCount; i++ {
			lineIndex := codeBlockCurrent.LineIndex + i
			blocks = append(blocks, &Text{
				Block: Block{
					LineIndex: lineIndex,
					LineCount: 1,
				},
				Inlines: ParseInline(reader.GetLineString(lineIndex)),
			})
		}
		codeBlockScope = false
	}
	return blocks
}

func ParseInline(line string) []AbstractInline {
	inlines := []AbstractInline{}
	column := 0
	n := len(line)
	SPECIAL_CHARS := "`!*[$~#\n"
	for column < len(line) {
		at := column
		c := line[column]

		// Code
		if c == '`' {
			end := byteIndexOf(line, '`', column+1)
			if end != -1 && end != column+1 {
				inlines = append(inlines, &Code{
					Inline: Inline{
						Spans: []Span{
							{
								StartColumn: at,
								EndColumn:   end + 1,
							},
							{
								StartColumn: at + 1,
								EndColumn:   end,
							},
						},
					},
				})
				column = end + 1
				continue
			}
		}

		// Bold
		if c == '*' && column+1 < n && line[column+1] == '*' {
			end := strIndexOf(line, "**", column+2)
			if end != -1 {
				inlines = append(inlines, &Bold{
					Inline: Inline{
						Spans: []Span{
							{
								StartColumn: at,
								EndColumn:   end + 2,
							},
							{
								StartColumn: at + 2,
								EndColumn:   end,
							},
						},
					},
				})
				column = end + 2
				continue
			}
		}
		if c == '_' && column+1 < n && line[column+1] == '_' {
			end := strIndexOf(line, "__", column+2)
			if end != -1 {
				inlines = append(inlines, &Bold{
					Inline: Inline{
						Spans: []Span{
							{
								StartColumn: at,
								EndColumn:   end + 2,
							},
							{
								StartColumn: at + 2,
								EndColumn:   end,
							},
						},
					},
				})
				column = end + 2
				continue
			}
		}

		// Italic
		if c == '*' {
			end := byteIndexOf(line, '*', column+1)
			if end != -1 && end != column+1 {
				inlines = append(inlines, &Italic{
					Inline: Inline{
						Spans: []Span{
							{
								StartColumn: at,
								EndColumn:   end + 1,
							},
							{
								StartColumn: at + 1,
								EndColumn:   end,
							},
						},
					},
				})
				column = end + 1
				continue
			}
		}
		if c == '_' {
			end := byteIndexOf(line, '_', column+1)
			if end != -1 && end != column+1 {
				inlines = append(inlines, &Italic{
					Inline: Inline{
						Spans: []Span{
							{
								StartColumn: at,
								EndColumn:   end + 1,
							},
							{
								StartColumn: at + 1,
								EndColumn:   end,
							},
						},
					},
				})
				column = end + 1
				continue
			}
		}

		// Link / Image
		if c == '[' || (c == '!' && column+1 < n && line[column+1] == '[') {
			bang := c == '!'
			start := column + 1
			if bang {
				start = column + 2
			}
			if start < n && start+1 < n {
				end := byteIndexOf(line, ']', start+1)
				if end != -1 {
					altStart := start
					altEnd := end
					start = end
					if start+1 < n && line[start+1] == '(' {
						end = byteIndexOf(line, ')', start+2)
						linkStart := start + 2
						linkEnd := end
						if bang {
							inlines = append(inlines, &Image{
								Inline: Inline{
									Spans: []Span{
										{
											StartColumn: column,
											EndColumn:   end + 1,
										},
										{
											StartColumn: altStart,
											EndColumn:   altEnd,
										},
										{
											StartColumn: linkStart,
											EndColumn:   linkEnd,
										},
									},
								},
							})
						} else {
							inlines = append(inlines, &Link{
								Inline: Inline{
									Spans: []Span{
										{
											StartColumn: column,
											EndColumn:   end + 1,
										},
										{
											StartColumn: altStart,
											EndColumn:   altEnd,
										},
										{
											StartColumn: linkStart,
											EndColumn:   linkEnd,
										},
									},
								},
							})
						}
						column = end + 1
						continue
					}
				}
			}
		}

		// Strike
		if c == '~' && column+1 < n && line[column+1] == '~' {
			end := strIndexOf(line, "~~", column+2)
			if end != -1 {
				inlines = append(inlines, &Strike{
					Inline: Inline{
						Spans: []Span{
							{
								StartColumn: at,
								EndColumn:   end + 2,
							},
							{
								StartColumn: at + 2,
								EndColumn:   end,
							},
						},
					},
				})
				column = end + 2
				continue
			}
		}

		// Plain text
		next := column + 1
		for next < n && !strings.ContainsRune(SPECIAL_CHARS, rune(line[next])) {
			next++
		}
		inlines = append(inlines, &PlainText{
			Inline: Inline{
				Spans: []Span{
					{
						StartColumn: column,
						EndColumn:   next,
					},
				},
			},
		})
		column = next
	}
	return inlines
}
