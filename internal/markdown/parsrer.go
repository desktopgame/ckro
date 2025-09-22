package markdown

import (
	"regexp"
	"strings"

	"github.com/desktopgame/ckro/internal/optional"
)

var FENCE_RE = regexp.MustCompile("^```(\\S+)?\\s*$")
var MATH_FENCE_RE = regexp.MustCompile(`^\$\$.+$`)
var HEADER_RE = regexp.MustCompile(`^(#{1,6})\s+(.*)$`)
var QUOTE_RE = regexp.MustCompile(`^>\s?(.*)$`)
var UL_RE = regexp.MustCompile(`^([*-])\s+(.*)$`)
var OL_RE = regexp.MustCompile(`^(\d+)[\.)]\s+(.*)$`)
var CHECK_OFF_RE = regexp.MustCompile(`^\[\s\]\s+(.*)$`)
var CHECK_ON_RE = regexp.MustCompile(`^\[x\]\s+(.*)$`)
var TABLE_HEADER_RE = regexp.MustCompile(`^(\|.+)+\|$`)
var TABLE_LAYOUT_RE = regexp.MustCompile(`^(\|.+)+\|$`)
var TABLE_CONTENT_RE = regexp.MustCompile(`^(\|.+)+\|$`)
var HR_RE = regexp.MustCompile(`^\s{0,3}(\*\s*\*\s*\*|-\s*-\s*-|_\s*_\s*_)\s*$`)
var CALLOUT_TAG_RE = regexp.MustCompile(`^\[!([A-Za-z0-9_-]+)\]\s*(.*)$`)

var SPECIAL_INLINE_CHARS = "`!*[$~#\n"
var TRIM_CHARS = " 　\t"

func parseLineStarts(text string) []int {
	starts := []int{0}
	pos := 0
	for _, r := range text {
		pos += 1
		if r == '\n' {
			starts = append(starts, pos)
		}
	}
	return starts
}

func locationRange(text string, lineStarts []int, si int, ei int) (int, int, int, int) {
	startRune := lineStarts[si]
	endRune := len(text)
	if ei < len(lineStarts) {
		endRune = lineStarts[ei]
	}
	return startRune, endRune, si, ei
}

func strIndexOf(s, sub string, at int) int {
	i := strings.Index(s[at:], sub)
	if i < 0 {
		return -1
	}
	return at + i
}

func runeIndexOf(s string, sub rune, at int) int {
	i := strings.IndexRune(s[at:], sub)
	if i < 0 {
		return -1
	}
	return at + i
}

func parseWikilinkInner(inner string) map[string]string {
	filePart := inner
	slug := optional.None[string]()
	alias := optional.None[string]()
	if strings.Contains(inner, "|") {
		a := strings.Split(inner, "|")
		filePart = a[0]
		alias = optional.Some(a[1])
	}
	if strings.Contains(filePart, "#") {
		a := strings.Split(inner, "|")
		filePart = a[0]
		slug = optional.Some(a[1])
	}
	return map[string]string{
		"file":  strings.Trim(filePart, TRIM_CHARS),
		"slug":  strings.Trim(slug.OrElse(""), TRIM_CHARS),
		"alias": strings.Trim(alias.OrElse(""), TRIM_CHARS),
	}
}

func parseInline(text string) []AbstractInline {
	out := []AbstractInline{}
	i := 0
	n := len(text)
	taggableAt := true
	for i < n {
		ch := text[i]
		// Obsidian like newline
		if ch == '\n' {
			out = append(out, &HardBreak{})
			i++
			continue
		}
		// Inline code `code`
		if ch == '`' {
			j := runeIndexOf(text, '`', i+1)
			if j != -1 {
				out = append(out, &Code{
					Text: text[i+1 : j],
				})
				i = j + 1
				continue
			}
		}
		// Strong **bold** takes precedence over *em*
		if ch == '*' && i+1 < n && text[i+1] == '*' {
			j := strIndexOf(text, "**", i+2)
			if j != -1 {
				children := parseInline(text[i+2 : j])
				out = append(out, &Strong{
					Children: children,
				})
				i = j + 2
				continue
			}
		}
		// Emphasis *em*
		if ch == '*' {
			j := runeIndexOf(text, '*', i+1)
			if j != -1 {
				children := parseInline(text[i+1 : j])
				out = append(out, &Emph{
					Children: children,
				})
				i = j + 1
				continue
			}
		}
		// Wikilink / Embed: [[...]] or ![[...]]
		if ch == '[' || (ch == '!' && i+1 < n && text[i+1] == '[') {
			bang := ch == '!'
			start := i + 1
			if bang {
				start = i + 2
			}
			if start < n && start+1 < n && text[start] == '[' {
				end := strIndexOf(text, "]]", start+1)
				if end != -1 {
					inner := text[start+1 : end]
					data := parseWikilinkInner(inner)
					name := "wikilink"
					if bang {
						name = "embed"
					}
					out = append(out, &ExtInline{
						Name: name,
						Data: data,
					})
					i = end + 2
					continue
				}
			}
		}
		// Link / Image: [...](...) or ![...](...)
		if ch == '[' || (ch == '!' && i+1 < n && text[i+1] == '[') {
			bang := ch == '!'
			start := i + 1
			if bang {
				start = i + 2
			}
			if start < n && start+1 < n {
				end := runeIndexOf(text, ']', start+1)
				if end != -1 {
					alt := text[start:end]
					start = end
					if start+1 < n && text[start+1] == '(' {
						end = runeIndexOf(text, ')', start+2)
						link := text[start+2 : end]
						if bang {
							out = append(out, &Image{
								Alt: alt,
								Url: link,
							})
						} else {
							out = append(out, &Link{
								Title: optional.Some(alt),
								Url:   link,
							})
						}
						i = end + 1
						continue
					}
				}
			}
		}
		// Math inline (like Obsidian)
		if ch == '$' && (i+1 < n && text[i+1] == '$') {
			j := i + 1
			end := strIndexOf(text, "$$", j+1)
			if end != -1 {
				start := i + 2
				i = end + 2
				end--
				out = append(out, &Math{
					Tex: text[start : end+1],
				})
				continue
			}
		}
		// Math inline $...$ (naive; avoids $$)
		if ch == '$' && !(i+1 < n && text[i+1] == '$') {
			j := runeIndexOf(text, '$', i+1)
			if j != -1 {
				// Quick guard to avoid $100$ captureing; require non-digit boundaries
				leftOk := i == 0 || (text[i-1] < '0' || text[i-1] > '9')
				rightOk := (j+1 >= n) || (text[j+1] < '0' || text[j+1] > '9')
				if leftOk && rightOk {
					out = append(out, &Math{
						Tex: text[i+1 : j],
					})
					i = j + 1
					continue
				}
			}
		}
		// Strike
		if ch == '~' && (i+1 < n && text[i+1] == '~') {
			j := i + 1
			end := strIndexOf(text, "~~", j+1)
			if end != -1 {
				start := i + 2
				i = end + 2
				end--
				out = append(out, &Strike{
					Children: parseInline(text[start : end+1]),
				})
				continue
			}
		}
		// Tag
		if ch == '#' && (i+1 < n && text[i+1] != ' ') && taggableAt {
			j := i + 1
			for j < n && (text[j] != ' ' && text[j] != '\n') {
				j++
			}
			out = append(out, &Tag{
				Text: text[i+1 : j],
			})
			i = j
			continue
		}
		// Fallback: plain text until next special char
		j := i + 1
		for j < n && !strings.ContainsRune(SPECIAL_INLINE_CHARS, rune(text[j])) {
			j++
		}
		out = append(out, &Text{
			Text: text[i:j],
		})
		// separate check
		taggableAt = false
		if j > 0 && (text[j-1] == ' ' || text[j-1] == '\n') {
			taggableAt = true
		}
		i = j
	}
	return out
}

func hasIndent(line string, il int) bool {
	spaces := strings.Repeat("  ", il)
	tabs := strings.Repeat("\t", il)
	if il > 0 {
		if strings.HasPrefix(line, spaces) {
			return true
		} else if strings.HasPrefix(line, tabs) {
			return true
		} else {
			return false
		}
	}
	return true
}

func cutIndent(line string, il int) string {
	spaces := strings.Repeat("  ", il)
	tabs := strings.Repeat("\t", il)
	offset := 0
	if il > 0 {
		if strings.HasPrefix(line, spaces) {
			offset = len(spaces)
		} else if strings.HasPrefix(line, tabs) {
			offset = len(tabs)
		}
	}
	line = line[offset:]
	return line
}

func parseList(lines []string, block *ListBlock, line int, il int) int {
	currentBlock := optional.None[*ListBlock]()
	for line < len(lines) {
		if !hasIndent(lines[line], il) {
			break
		}
		scanLine := cutIndent(lines[line], il)
		mmUl := UL_RE.FindStringSubmatch(scanLine)
		mmOl := OL_RE.FindStringSubmatch(scanLine)
		mmOrdered := len(mmOl) > 0
		if len(mmUl) == 0 && len(mmOl) == 0 {
			if currentBlock.IsSome() {
				if hasIndent(lines[line], il+1) {
					aheadLine := cutIndent(lines[line], il+1)
					if UL_RE.MatchString(aheadLine) || OL_RE.MatchString(aheadLine) {
						line = parseList(lines, currentBlock.Value(), line, il+1)
						continue
					} else {
						line++
						continue
					}
				}
			}
			break
		}
		var itemText string
		checked := optional.None[bool]()
		if len(mmOl) > 0 {
			itemText = mmOl[2]

			matchOn := CHECK_ON_RE.FindStringSubmatch(itemText)
			if len(matchOn) > 0 {
				checked = optional.Some(true)
				itemText = matchOn[1]
			}

			matchOff := CHECK_OFF_RE.FindStringSubmatch(itemText)
			if len(matchOff) > 0 {
				checked = optional.Some(false)
				itemText = matchOff[1]
			}
		} else if len(mmUl) > 0 {
			itemText = mmUl[2]

			matchOn := CHECK_ON_RE.FindStringSubmatch(itemText)
			if len(matchOn) > 0 {
				checked = optional.Some(true)
				itemText = matchOn[1]
			}

			matchOff := CHECK_OFF_RE.FindStringSubmatch(itemText)
			if len(matchOff) > 0 {
				checked = optional.Some(false)
				itemText = matchOff[1]
			}
		} else {
			break
		}
		currentBlock = optional.Some(&ListBlock{
			Ordered: mmOrdered,
		})
		item := ListItem{
			Paragraph: optional.Some(Paragraph{
				Inlines: parseInline(itemText),
			}),
			Checked: checked,
			Blocks: []AbstractBlock{
				currentBlock.Value(),
			},
		}
		block.Items = append(block.Items, &item)
		line += 1
	}
	return line
}

func parseAlign(align string) string {
	padLeft := false
	padRight := false
	if strings.HasPrefix(align, ":") {
		padRight = true
	}
	if strings.HasSuffix(align, ":") {
		padLeft = true
	}
	if padLeft && padRight {
		return "center"
	} else if padLeft {
		return "right"
	} else if padRight {
		return "left"
	}
	return "left"
}

func parseTableContents(match []string) TableRow {
	columns := strings.Split(match[0], "|")

	var paragraphs []Paragraph
	for _, column := range columns {
		inlines := parseInline(column)
		paragraphs = append(paragraphs, Paragraph{
			Inlines: inlines,
		})
	}
	return TableRow{Columns: paragraphs}
}

func Parse(text string) *Document {
	lineStarts := parseLineStarts(text)
	lines := strings.Split(text, "\n")
	doc := Document{
		Block: Block{
			Node: Node{
				StartRune: 0,
				EndRune:   len(text),
				StartLine: 0,
				EndLine:   len(lines),
			},
		},
	}
	i := 0
	for i < len(lines) {
		ln := lines[i]

		// blank line
		if len(strings.Trim(ln, TRIM_CHARS)) == 0 {
			i++
			continue
		}

		// Thematic break
		if HR_RE.MatchString(ln) {
			startRune, endRune, sl, el := locationRange(text, lineStarts, i, i+1)
			doc.Blocks = append(doc.Blocks, &ThematicBreak{
				Block: Block{
					Node: Node{
						StartRune: startRune,
						EndRune:   endRune,
						StartLine: sl,
						EndLine:   el,
					},
				},
			})
			i += 1
			continue
		}

		// Fence ```
		m := FENCE_RE.FindStringSubmatch(ln)
		if len(m) > 0 {
			info := ""
			if len(m) > 1 {
				info = m[1]
			}
			info = strings.Trim(info, TRIM_CHARS)
			j := i + 1

			bodyLines := []string{}
			for j < len(lines) && !FENCE_RE.MatchString(lines[j]) {
				bodyLines = append(bodyLines, lines[j])
				j++
			}
			// consume closing fence if present
			end := j
			if j < len(lines) && FENCE_RE.MatchString(lines[j]) {
				end = j + 1
			}
			startRune, endRune, sl, el := locationRange(text, lineStarts, i, end)
			doc.Blocks = append(doc.Blocks, &CodeBlock{
				Info: info,
				Code: strings.Join(bodyLines, "\n"),
				Block: Block{
					Node{
						StartRune: startRune,
						EndRune:   endRune,
						StartLine: sl,
						EndLine:   el,
					},
				},
			})
			i = end
			continue
		}

		// Math fence $$
		if MATH_FENCE_RE.MatchString(ln) && !strings.Contains(ln[2:], "$$") {
			j := i + 1
			bodyLines := []string{}
			for j < len(lines) && !MATH_FENCE_RE.MatchString(lines[j]) {
				bodyLines = append(bodyLines, lines[j])
				j++
			}
			end := j
			if j < len(lines) && MATH_FENCE_RE.MatchString(lines[j]) {
				end = j + 1
			}
			startRune, endRune, sl, el := locationRange(text, lineStarts, i, end)
			doc.Blocks = append(doc.Blocks, &MathBlock{
				Tex: strings.Join(bodyLines, "\n"),
				Block: Block{
					Node: Node{
						StartRune: startRune,
						EndRune:   endRune,
						StartLine: sl,
						EndLine:   el,
					},
				},
			})
			i = end
			continue
		}

		// Heading
		hm := HEADER_RE.FindStringSubmatch(ln)
		if len(hm) > 0 {
			level := len(hm[1])
			content := hm[2]
			startRune, endRune, sl, el := locationRange(text, lineStarts, i, i+1)
			inlines := parseInline(content)
			doc.Blocks = append(doc.Blocks, &Heading{
				Level:   level,
				Inlines: inlines,
				Block: Block{
					Node{
						StartRune: startRune,
						EndRune:   endRune,
						StartLine: sl,
						EndLine:   el,
					},
				},
			})
			i++
			continue
		}

		// Blockquote (with optional callout detection on first line)
		if QUOTE_RE.MatchString(ln) {
			j := i
			innerLines := []string{}
			for j < len(lines) && QUOTE_RE.MatchString(lines[j]) {
				innerLines = append(innerLines, QUOTE_RE.FindStringSubmatch(lines[j])[1])
				j++
			}
			// Callout?
			calloutKind := ""
			if len(innerLines) > 0 {
				cm := CALLOUT_TAG_RE.FindStringSubmatch(innerLines[0])
				if len(cm) > 0 {
					calloutKind = strings.ToLower(cm[1])
					innerLines[0] = cm[2]
				}
			}
			innerText := strings.Join(innerLines, "\n")
			childDoc := Parse(innerText)
			startRune, endRune, sl, el := locationRange(text, lineStarts, i, j)
			var node AbstractBlock
			if len(calloutKind) > 0 {
				node = &ExtBlock{
					Name:   "callout",
					Blocks: childDoc.Blocks,
					Block: Block{
						Node: Node{
							StartRune: startRune,
							EndRune:   endRune,
							StartLine: sl,
							EndLine:   el,
							Attrs: map[string]string{
								"kind": calloutKind,
							},
						},
					},
				}
			} else {
				node = &Blockquote{
					Blocks: childDoc.Blocks,
					Block: Block{
						Node: Node{
							StartRune: startRune,
							EndRune:   endRune,
							StartLine: sl,
							EndLine:   el,
							Attrs: map[string]string{
								"kind": calloutKind,
							},
						},
					},
				}
			}
			doc.Blocks = append(doc.Blocks, node)
			i = j
			continue
		}

		// List
		mUl := UL_RE.FindStringSubmatch(ln)
		mOl := OL_RE.FindStringSubmatch(ln)
		if len(mUl) > 0 || len(mOl) > 0 {
			ordered := len(mOl) > 0
			j := i

			root := ListBlock{}
			doc.Blocks = append(doc.Blocks, &root)
			j = parseList(lines, &root, j, 0)

			startRune, endRune, sl, el := locationRange(text, lineStarts, i, j)
			root.Ordered = ordered
			root.StartRune = startRune
			root.EndRune = endRune
			root.StartLine = sl
			root.EndLine = el
			i = j
			continue
		}

		// Table
		tableHeader := TABLE_HEADER_RE.FindStringSubmatch(ln)
		if len(tableHeader) > 0 {
			j := i + 1
			tableAligns := TABLE_LAYOUT_RE.FindStringSubmatch(lines[j])
			if len(tableAligns) > 0 {
				j++
				tableContents := [][]string{}
				for j < len(lines) {
					content := TABLE_CONTENT_RE.FindStringSubmatch(lines[j])
					if len(content) == 0 {
						break
					}
					tableContents = append(tableContents, content)
					j++
				}
				if len(tableContents) > 0 {
					headers := strings.Split(tableHeader[0], "|")
					aligns := strings.Split(tableAligns[0], "|")

					var dAligns []string
					for _, align := range aligns {
						dAligns = append(dAligns, parseAlign(align))
					}

					var dContents []TableRow
					for _, content := range tableContents {
						dContents = append(dContents, parseTableContents(content))
					}

					startRune, endRune, sl, el := locationRange(text, lineStarts, i, j)
					doc.Blocks = append(doc.Blocks, &Table{
						Headers: headers[1 : len(headers)-1],
						Aligns:  dAligns,
						Rows:    dContents,
						Block: Block{
							Node: Node{
								StartRune: startRune,
								EndRune:   endRune,
								StartLine: sl,
								EndLine:   el,
							},
						},
					})
					i = j
					continue
				}
			}
		}

		// Paragraph
		j := i
		buf := []string{}
		for j < len(lines) {
			cur := lines[j]
			if len(strings.Trim(cur, TRIM_CHARS)) == 0 {
				break
			}
			mFenceRe := FENCE_RE.MatchString(cur)
			mMathFence := MATH_FENCE_RE.MatchString(cur) && !strings.Contains(cur[2:], "$$")
			mHeader := HEADER_RE.MatchString(cur)
			mQuote := QUOTE_RE.MatchString(cur)
			mUL := UL_RE.MatchString(cur) || OL_RE.MatchString(cur)
			mHR := HR_RE.MatchString(cur)
			anyMatch := mFenceRe || mMathFence || mHeader || mQuote || mUL || mHR

			if anyMatch {
				break
			}
			buf = append(buf, cur)
			j++
		}
		paraText := strings.Join(buf, "\n")
		startRune, endRune, sl, el := locationRange(text, lineStarts, i, j)
		inlines := parseInline(paraText)
		doc.Blocks = append(doc.Blocks,
			&Paragraph{
				Inlines: inlines,
				Block: Block{
					Node{
						StartRune: startRune,
						EndRune:   endRune,
						StartLine: sl,
						EndLine:   el,
					},
				},
			})
		i = j
	}
	return &doc
}
