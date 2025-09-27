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
			blocks = append(blocks, &Heading{
				Block: Block{
					LineIndex: lineIndex,
					LineCount: 1,
				},
				Level: column,
			})
			continue
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
						StartColumn: at,
						EndColumn:   end + 1,
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
						StartColumn: at,
						EndColumn:   end + 2,
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
						StartColumn: at,
						EndColumn:   end + 1,
					},
				})
				column = end + 1
				continue
			}
		}

		// Strike
		if c == '~' && column+1 < n && line[column+1] == '~' {
			end := strIndexOf(line, "~~", column+2)
			if end != -1 {
				inlines = append(inlines, &Strike{
					Inline: Inline{
						StartColumn: at,
						EndColumn:   end + 2,
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
				StartColumn: column,
				EndColumn:   next,
			},
		})
		column = next
	}
	return inlines
}
