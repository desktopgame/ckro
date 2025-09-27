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
	return reader.GetLine(lineIndex)[span.StartColumn:span.EndColumn]
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
			if line[0] == '`' {
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

		// Soft break
		if len(line) == 0 {
			blocks = append(blocks, &SoftBreak{
				Block: Block{
					LineIndex: lineIndex,
					LineCount: 1,
				},
			})
			continue
		}

		// CodeBlock
		if line[0] == '`' {
			codeBlockScope = true

			column := 0
			for column < len(line) && line[column] == '`' {
				column++
			}
			codeBlockMarkerLen = column - 1

			lang := ""
			if column < len(line) {
				lang = line[column:]
			}

			codeBlockCurrent = &CodeBlock{
				Block: Block{
					LineIndex: lineIndex,
					LineCount: 1,
				},
				Lang: lang,
			}
			blocks = append(blocks, codeBlockCurrent)
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
			if end != -1 {
				inlines = append(inlines, &Code{
					Inline: Inline{
						Spans: []Span{
							{
								StartColumn: at,
								EndColumn:   end + 1,
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
			if end != -1 {
				inlines = append(inlines, &Italic{
					Inline: Inline{
						Spans: []Span{
							{
								StartColumn: at,
								EndColumn:   end + 1,
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
					// alt := line[start:end]
					start = end
					if start+1 < n && line[start+1] == '(' {
						end = byteIndexOf(line, ')', start+2)
						// link := line[start+2 : end]
						if bang {
							inlines = append(inlines, &Image{
								Inline: Inline{
									Spans: []Span{
										{
											StartColumn: column,
											EndColumn:   end + 1,
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
					Span{
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
