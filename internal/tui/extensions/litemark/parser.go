package litemark

import (
	"regexp"
	"strings"
)

const TABLE_ALIGN_LEFT = 0
const TABLE_ALIGN_CENTER = 1
const TABLE_ALIGN_RIGHT = 2

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

// MarkerInfo represents a fold marker position
type MarkerInfo struct {
	LineIndex int
	Level     int
	IsStart   bool
}

// countFoldMarkers counts start and end markers and returns if they match
func countFoldMarkers(reader Reader) (startMarkers []MarkerInfo, endMarkers []MarkerInfo, balanced bool) {
	startMarkers = []MarkerInfo{}
	endMarkers = []MarkerInfo{}

	for i := 0; i < reader.GetLineCount(); i++ {
		line := reader.GetLineString(i)
		if len(line) == 0 {
			continue
		}

		// Check for fold start marker
		if line[0] == '{' {
			column := 0
			for column < len(line) && line[column] == '{' {
				column++
			}
			if column >= 3 && line == strings.Repeat("{", column) {
				startMarkers = append(startMarkers, MarkerInfo{
					LineIndex: i,
					Level:     column,
					IsStart:   true,
				})
			}
		}

		// Check for fold end marker
		if line[0] == '}' {
			column := 0
			for column < len(line) && line[column] == '}' {
				column++
			}
			if column >= 3 && line == strings.Repeat("}", column) {
				endMarkers = append(endMarkers, MarkerInfo{
					LineIndex: i,
					Level:     column,
					IsStart:   false,
				})
			}
		}
	}

	balanced = len(startMarkers) == len(endMarkers)
	return
}

// countCodeBlockMarkers counts start and end markers for code blocks
func countCodeBlockMarkers(reader Reader) (startMarkers []MarkerInfo, endMarkers []MarkerInfo, balanced bool) {
	startMarkers = []MarkerInfo{}
	endMarkers = []MarkerInfo{}
	inCodeBlock := false
	currentLevel := 0

	for i := 0; i < reader.GetLineCount(); i++ {
		line := reader.GetLineString(i)
		if len(line) == 0 {
			continue
		}

		if line[0] == '`' {
			column := 0
			for column < len(line) && line[column] == '`' {
				column++
			}

			if column >= 3 {
				if !inCodeBlock {
					// Start marker
					startMarkers = append(startMarkers, MarkerInfo{
						LineIndex: i,
						Level:     column,
						IsStart:   true,
					})
					inCodeBlock = true
					currentLevel = column
				} else {
					// Potential end marker - check if level matches
					if column == currentLevel {
						endMarkers = append(endMarkers, MarkerInfo{
							LineIndex: i,
							Level:     column,
							IsStart:   false,
						})
						inCodeBlock = false
						currentLevel = 0
					}
				}
			}
		}
	}

	balanced = len(startMarkers) == len(endMarkers)
	return
}

// matchFoldMarkers matches start and end markers using stack
func matchFoldMarkers(startMarkers []MarkerInfo, endMarkers []MarkerInfo) map[int]int {
	// Returns map: startLineIndex -> endLineIndex
	pairs := make(map[int]int)
	stack := []MarkerInfo{}

	// Merge and sort markers by line index
	allMarkers := make([]MarkerInfo, 0, len(startMarkers)+len(endMarkers))
	allMarkers = append(allMarkers, startMarkers...)
	allMarkers = append(allMarkers, endMarkers...)

	// Simple sort by line index
	for i := 0; i < len(allMarkers); i++ {
		for j := i + 1; j < len(allMarkers); j++ {
			if allMarkers[i].LineIndex > allMarkers[j].LineIndex {
				allMarkers[i], allMarkers[j] = allMarkers[j], allMarkers[i]
			}
		}
	}

	for _, marker := range allMarkers {
		if marker.IsStart {
			stack = append(stack, marker)
		} else {
			// End marker - try to match with stack top
			if len(stack) > 0 {
				top := stack[len(stack)-1]
				if top.Level == marker.Level {
					// Match found
					pairs[top.LineIndex] = marker.LineIndex
					stack = stack[:len(stack)-1] // pop
				}
			}
		}
	}

	return pairs
}

func parseBlock(sc *Scanner, line string, lineIndex int, re *regexp.Regexp) *Table {
	if !strings.HasPrefix(line, "|") || !strings.HasSuffix(line, "|") {
		sc.lineIndex = lineIndex + 1
		return nil
	}
	if !sc.Ready() {
		sc.lineIndex = lineIndex + 1
		return nil
	}
	headers := strings.Split(line, "|")
	headers = headers[1 : len(headers)-1]

	aligns := sc.Next()
	if !strings.HasPrefix(aligns, "|") || !strings.HasSuffix(aligns, "|") {
		sc.lineIndex = lineIndex + 1
		return nil
	}
	if !re.MatchString(aligns) {
		sc.lineIndex = lineIndex + 1
		return nil
	}

	alignsSplit := strings.Split(aligns, "|")
	alignsSplit = alignsSplit[1 : len(alignsSplit)-1]
	for _, align := range alignsSplit {
		for _, b := range align {
			if b == ':' {
				continue
			}
			if b == '-' {
				continue
			}
			if b == ' ' {
				continue
			}
			sc.lineIndex = lineIndex + 1
			return nil
		}
	}

	if len(headers) != len(alignsSplit) {
		sc.lineIndex = lineIndex + 1
		return nil
	}

	alingsParsed := make([]int, len(alignsSplit))
	for i := 0; i < len(alignsSplit); i++ {
		lColon := strings.HasPrefix(alignsSplit[i], ":")
		rColon := strings.HasSuffix(alignsSplit[i], ":")

		if lColon && rColon {
			alingsParsed[i] = TABLE_ALIGN_CENTER
		} else if lColon {
			alingsParsed[i] = TABLE_ALIGN_LEFT
		} else if rColon {
			alingsParsed[i] = TABLE_ALIGN_RIGHT
		}
	}

	tableHeaders := []*Text{}
	tableHeaderOffset := 1
	for i := 0; i < len(headers); i++ {
		inlines := ParseInline(headers[i])
		for _, il := range inlines {
			bil := il.BaseInline()
			for j := 0; j < len(bil.Spans); j++ {
				bil.Spans[j] = Span{
					StartColumn: bil.Spans[j].StartColumn + tableHeaderOffset,
					EndColumn:   bil.Spans[j].EndColumn + tableHeaderOffset,
				}
			}
		}

		tableHeaders = append(tableHeaders, &Text{
			Block: Block{
				LineIndex: lineIndex,
				LineCount: 1,
			},
			Inlines: inlines,
		})
		tableHeaderOffset = byteIndexOf(line, '|', tableHeaderOffset+1) + 1
	}

	tableRows := []TableRow{}
	columnCount := -1
	for sc.Ready() {
		row := sc.Next()
		if re.MatchString(row) && strings.HasPrefix(aligns, "|") && strings.HasSuffix(aligns, "|") {
			texts := strings.Split(row, "|")
			texts = texts[1 : len(texts)-1]
			columns := []*Text{}

			columnOffset := 1
			for _, tex := range texts {
				if len(tex) == 0 {
					continue
				}
				inlines := ParseInline(tex)
				for _, il := range inlines {
					bil := il.BaseInline()
					for i := 0; i < len(bil.Spans); i++ {
						bil.Spans[i] = Span{
							StartColumn: bil.Spans[i].StartColumn + columnOffset,
							EndColumn:   bil.Spans[i].EndColumn + columnOffset,
						}
					}
				}

				columns = append(columns, &Text{
					Block: Block{
						LineIndex: lineIndex + 2 + len(tableRows),
						LineCount: 1,
					},
					Inlines: inlines,
				})
				columnOffset = byteIndexOf(row, '|', columnOffset+1) + 1
			}
			if columnCount == -1 {
				columnCount = len(columns)
			}
			if columnCount != len(columns) {
				sc.lineIndex--
				break
			}
			tableRows = append(tableRows, TableRow{
				LineIndex: lineIndex + 2 + len(tableRows),
				Columns:   columns,
			})
		} else {
			sc.lineIndex--
			break
		}
	}
	if len(tableRows) > 0 && len(tableHeaders) == len(tableRows[0].Columns) {
		table := Table{
			Block: Block{
				LineIndex: lineIndex,
				LineCount: 2 + len(tableRows),
			},
			Headers: tableHeaders,
			Aligns:  alingsParsed,
			Rows:    tableRows,
		}
		return &table
	} else {
		sc.lineIndex = lineIndex + 1
		return nil
	}
}

func Parse(reader Reader) []AbstractBlock {
	// Phase 1: Count and match fold markers
	foldStartMarkers, foldEndMarkers, foldBalanced := countFoldMarkers(reader)
	var foldPairs map[int]int
	if foldBalanced {
		foldPairs = matchFoldMarkers(foldStartMarkers, foldEndMarkers)
	} else {
		foldPairs = make(map[int]int) // Empty map - no folds will be created
	}

	// Phase 2: Count code block markers
	codeStartMarkers, codeEndMarkers, codeBalanced := countCodeBlockMarkers(reader)
	codeBlockPairs := make(map[int]int) // startLine -> endLine
	if codeBalanced {
		for i := range codeStartMarkers {
			if i < len(codeEndMarkers) {
				codeBlockPairs[codeStartMarkers[i].LineIndex] = codeEndMarkers[i].LineIndex
			}
		}
	}

	table_re := regexp.MustCompile(`^(\|.+)+\|$`)

	sc := Scanner{Reader: reader}
	blocks := []AbstractBlock{}

	for sc.Ready() {
		lineIndex := sc.lineIndex
		line := sc.Next()

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

		// Table
		if table_re.MatchString(line) {
			table := parseBlock(&sc, line, lineIndex, table_re)
			if table != nil {
				blocks = append(blocks, table)
				continue
			}
		}

		// CodeBlock - using pre-computed pairs
		if endLine, exists := codeBlockPairs[lineIndex]; exists {
			column := 0
			for column < len(line) && line[column] == '`' {
				column++
			}

			codeBlock := &CodeBlock{
				Block: Block{
					LineIndex: lineIndex,
					LineCount: endLine - lineIndex + 1,
				},
				Span: Span{
					StartColumn: column,
					EndColumn:   len(line),
				},
				Level: column,
			}
			blocks = append(blocks, codeBlock)

			// Skip to after end marker
			sc.lineIndex = endLine + 1
			continue
		}

		// Unmatched code block marker - treat as plain text
		if line[0] == '`' {
			column := 0
			for column < len(line) && line[column] == '`' {
				column++
			}
			if column >= 3 {
				blocks = append(blocks, &Text{
					Block: Block{
						LineIndex: lineIndex,
						LineCount: 1,
					},
					Inlines: ParseInline(line),
				})
				continue
			}
		}

		// Fold - using pre-computed pairs from count-based matching
		if endLine, exists := foldPairs[lineIndex]; exists {
			column := 0
			for column < len(line) && line[column] == '{' {
				column++
			}

			foldBlock := &FoldBlock{
				Block: Block{
					LineIndex: lineIndex,
					LineCount: endLine - lineIndex + 1,
				},
				Span: Span{
					StartColumn: column,
					EndColumn:   len(line),
				},
				Level: column,
			}
			blocks = append(blocks, foldBlock)

			// Skip to after end marker
			sc.lineIndex = endLine + 1
			continue
		}

		// Unmatched fold marker (when not balanced) - treat as plain text
		if line[0] == '{' {
			column := 0
			for column < len(line) && line[column] == '{' {
				column++
			}
			if column >= 3 && line == strings.Repeat("{", column) {
				// This is a fold start marker but not matched - treat as text
				blocks = append(blocks, &Text{
					Block: Block{
						LineIndex: lineIndex,
						LineCount: 1,
					},
					Inlines: ParseInline(line),
				})
				continue
			}
		}
		if line[0] == '}' {
			column := 0
			for column < len(line) && line[column] == '}' {
				column++
			}
			if column >= 3 && line == strings.Repeat("}", column) {
				// This is a fold end marker but not matched - treat as text
				blocks = append(blocks, &Text{
					Block: Block{
						LineIndex: lineIndex,
						LineCount: 1,
					},
					Inlines: ParseInline(line),
				})
				continue
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
