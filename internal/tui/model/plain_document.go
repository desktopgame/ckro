package model

import (
	"strings"

	"github.com/desktopgame/ckro/internal/text"
)

// PlainDocument is wrapper of Buffer.
// track a current cursor position.
type PlainDocument struct {
	buffer       Buffer
	cursorRow    int
	cursorColumn int
	Styled       bool
}

// Init is initialize Buffer.
func (doc *PlainDocument) Init() {
	doc.Clear()
}

func (doc *PlainDocument) Render() []Element {
	if doc.Styled {
		// Get all text content from buffer
		var content strings.Builder
		for i := 0; i < doc.buffer.GetLineCount(); i++ {
			line := doc.buffer.GetLineAt(i)
			content.WriteString(line.GetContent())
			if i < doc.buffer.GetLineCount()-1 {
				content.WriteString("\n")
			}
		}

		// Convert markdown to elements using goldmark
		return ConvertMarkdownToElements(content.String())
	}

	// Fallback to plain text rendering
	elements := []Element{}
	for i := 0; i < doc.buffer.GetLineCount(); i++ {
		line := doc.buffer.GetLineAt(i)
		lineStr := line.GetContent()

		elements = append(elements, &PlainElement{
			Text: line.GetContent(),
			StartPosition: Position{
				Row:    i,
				Column: 0,
			},
			EndPosition: Position{
				Row:    i,
				Column: text.GraphemeLength(lineStr) - 1,
			},
		})
	}
	return elements
}

// InsertLine is break line at current cursor position.
func (doc *PlainDocument) InsertLine() {
	currLine := doc.buffer.GetLineAt(doc.cursorRow)
	codepointColumn := text.GraphemeToCodepointPos(currLine.GetContent(), doc.cursorColumn)

	doc.buffer.InsertLine(doc.cursorRow, codepointColumn)
	doc.cursorRow++
	doc.cursorColumn = 0
}

// InsertString is insert string at current cursor position.
func (doc *PlainDocument) InsertString(s string) {
	currLine := doc.buffer.GetLineAt(doc.cursorRow)
	codepointColumn := text.GraphemeToCodepointPos(currLine.GetContent(), doc.cursorColumn)

	at, err := doc.buffer.InsertString(doc.cursorRow, codepointColumn, s)
	if err != nil {
		panic("illegal state")
	}
	if at.Row < doc.buffer.GetLineCount() {
		resultLine := doc.buffer.GetLineAt(at.Row)
		doc.cursorRow = at.Row
		doc.cursorColumn = text.CodepointToGraphemePos(resultLine.GetContent(), at.Column)
	}
}

// RemoveChar is remove character at current cursor position.
func (doc *PlainDocument) RemoveChar() {
	currLine := doc.buffer.GetLineAt(doc.cursorRow)
	if doc.cursorRow > 0 {
		if len(currLine.GetContent()) == 0 || doc.cursorColumn == 0 {
			aboveLine := doc.buffer.GetLineAt(doc.cursorRow - 1)
			cursorAt := text.GraphemeLength(aboveLine.GetContent())

			doc.buffer.AppendString(doc.cursorRow-1, currLine.GetContent())
			doc.buffer.RemoveLine(doc.cursorRow)
			doc.cursorRow--
			doc.cursorColumn = cursorAt
		} else if doc.cursorColumn > 0 {
			newContent := text.GraphemeRemove(currLine.GetContent(), doc.cursorColumn-1, 1)

			currLine.Remove(0, len(currLine.GetContent()))
			currLine.InsertString(0, newContent)
			doc.cursorColumn--
		}
	} else {
		if doc.cursorColumn > 0 {
			newContent := text.GraphemeRemove(currLine.GetContent(), doc.cursorColumn-1, 1)

			currLine.Remove(0, len(currLine.GetContent()))
			currLine.InsertString(0, newContent)
			doc.cursorColumn--
		}
	}
}

// Clear is initialize Buffer.
func (doc *PlainDocument) Clear() {
	doc.buffer = Buffer{}
	doc.buffer.Init()
	doc.cursorRow = 0
	doc.cursorColumn = 0
}

// FindPrev is searches for the specified string backward from the current cursor position
// and moves the cursor to the beginning of the found string.
// Returns true if found, false otherwise.
func (doc *PlainDocument) FindPrev(searchStr string) bool {
	if searchStr == "" {
		return false
	}

	// Start searching from current position
	startRow := doc.cursorRow
	startCol := doc.cursorColumn

	// Search in current line from beginning to current position
	if startRow < doc.buffer.GetLineCount() && startCol > 0 {
		currentLine := doc.buffer.GetLineAt(startRow).GetContent()
		searchText := currentLine[:startCol]
		if index := strings.LastIndex(searchText, searchStr); index != -1 {
			doc.cursorColumn = index
			return true
		}
	}

	// Search in previous lines (from bottom to top)
	for row := startRow - 1; row >= 0; row-- {
		line := doc.buffer.GetLineAt(row).GetContent()
		if index := strings.LastIndex(line, searchStr); index != -1 {
			doc.cursorRow = row
			doc.cursorColumn = index
			return true
		}
	}

	// Wrap around: search from end to current position
	for row := doc.buffer.GetLineCount() - 1; row >= startRow; row-- {
		line := doc.buffer.GetLineAt(row).GetContent()
		var searchText string
		if row == startRow {
			if startCol < len(line) {
				searchText = line[startCol:]
			}
		} else {
			searchText = line
		}

		if len(searchText) > 0 {
			if index := strings.LastIndex(searchText, searchStr); index != -1 {
				if row == startRow {
					doc.cursorColumn = startCol + index
				} else {
					doc.cursorRow = row
					doc.cursorColumn = index
				}
				return true
			}
		}
	}

	return false
}

// FindNext is searches for the specified string forward from the current cursor position
// and moves the cursor to the beginning of the found string.
// Returns true if found, false otherwise.
func (doc *PlainDocument) FindNext(searchStr string) bool {
	if searchStr == "" {
		return false
	}

	// Start searching from current position
	startRow := doc.cursorRow
	startCol := doc.cursorColumn

	// Search in current line from current position
	if startRow < doc.buffer.GetLineCount() {
		currentLine := doc.buffer.GetLineAt(startRow).GetContent()
		// Search from current column + 0 to avoid finding the same occurrence
		searchStart := startCol + 0
		if searchStart < len(currentLine) {
			if index := strings.Index(currentLine[searchStart:], searchStr); index != -1 {
				doc.cursorColumn = searchStart + index
				return true
			}
		}
	}

	// Search in subsequent lines
	for row := startRow + 1; row < doc.buffer.GetLineCount(); row++ {
		line := doc.buffer.GetLineAt(row).GetContent()
		if index := strings.Index(line, searchStr); index != -1 {
			doc.cursorRow = row
			doc.cursorColumn = index
			return true
		}
	}

	// Wrap around: search from beginning to current position
	for row := 0; row <= startRow; row++ {
		line := doc.buffer.GetLineAt(row).GetContent()
		var searchEnd int
		if row == startRow {
			searchEnd = startCol
		} else {
			searchEnd = len(line)
		}

		if searchEnd > 0 {
			searchText := line[:searchEnd]
			if index := strings.Index(searchText, searchStr); index != -1 {
				doc.cursorRow = row
				doc.cursorColumn = index
				return true
			}
		}
	}

	return false
}

// Replace deletes the specified number of characters from the current cursor position
// and inserts the replacement string at that position.
// After replacement, the cursor is moved to the end of the inserted text.
// Returns true if the operation was successful, false otherwise.
func (doc *PlainDocument) Replace(deleteCount int, replaceStr string) bool {
	if deleteCount < 0 {
		return false
	}

	// Save the starting position for cursor positioning after insertion
	startRow := doc.cursorRow
	startCol := doc.cursorColumn

	// If deleteCount is 0, just insert the string
	if deleteCount == 0 {
		doc.InsertString(replaceStr)
		return true
	}

	// Delete characters
	for i := 0; i < deleteCount; i++ {
		// Check if we're at the end of document
		if doc.cursorRow >= doc.buffer.GetLineCount() {
			break
		}

		currLine := doc.buffer.GetLineAt(doc.cursorRow)
		lineLength := text.GraphemeLength(currLine.GetContent())

		// If we're at the end of current line
		if doc.cursorColumn >= lineLength {
			// If there's a next line, move to it and continue deleting
			if doc.cursorRow < doc.buffer.GetLineCount()-1 {
				// Delete the line break (merge with next line)
				nextLine := doc.buffer.GetLineAt(doc.cursorRow + 1)
				doc.buffer.AppendString(doc.cursorRow, nextLine.GetContent())
				doc.buffer.RemoveLine(doc.cursorRow + 1)
				// Stay at current position to continue deleting
			} else {
				// We're at the end of document, stop deleting
				break
			}
		} else {
			// Delete character at current position
			newContent := text.GraphemeRemove(currLine.GetContent(), doc.cursorColumn, 1)
			currLine.Remove(0, len(currLine.GetContent()))
			currLine.InsertString(0, newContent)
			// Don't move cursor position as we're deleting forward
		}
	}

	// Reset cursor to the start position before insertion
	doc.cursorRow = startRow
	doc.cursorColumn = startCol

	// Insert replacement string at the current position
	// InsertString will automatically move the cursor to the end of inserted text
	doc.InsertString(replaceStr)

	return true
}

// MoveLeft is move cursor to left.
func (doc *PlainDocument) MoveLeft() {
	if doc.cursorColumn > 0 {
		doc.cursorColumn--
	} else {
		if doc.cursorRow > 0 {
			doc.cursorRow--
			doc.cursorColumn = text.GraphemeLength(doc.buffer.GetLineAt(doc.cursorRow).GetContent())
		}
	}
}

// MoveRight is move cursor to right.
func (doc *PlainDocument) MoveRight() {
	currLine := doc.buffer.GetLineAt(doc.cursorRow)
	if doc.cursorColumn < text.GraphemeLength(currLine.GetContent()) {
		doc.cursorColumn++
	} else {
		if doc.cursorRow < doc.buffer.GetLineCount()-1 {
			doc.cursorRow++
			doc.cursorColumn = 0
		}
	}
}

// MoveUp is move cursor to up.
func (doc *PlainDocument) MoveUp() {
	if doc.cursorRow > 0 {
		doc.cursorRow--

		currLine := doc.buffer.GetLineAt(doc.cursorRow)
		lineLength := text.GraphemeLength(currLine.GetContent())
		if doc.cursorColumn > lineLength {
			doc.cursorColumn = lineLength
		}
	}
}

// MoveDown is move cursor to down.
func (doc *PlainDocument) MoveDown() {
	if doc.cursorRow < doc.buffer.GetLineCount()-1 {
		doc.cursorRow++

		currLine := doc.buffer.GetLineAt(doc.cursorRow)
		lineLength := text.GraphemeLength(currLine.GetContent())
		if doc.cursorColumn > lineLength {
			doc.cursorColumn = lineLength
		}
	}
}

// MoveReset is cursor position reset to zero.
func (doc *PlainDocument) MoveReset() {
	doc.cursorRow = 0
	doc.cursorColumn = 0
}

// GetBuffer returns Buffer.
func (doc *PlainDocument) GetBuffer() *Buffer {
	return &doc.buffer
}

// GetCursorRow returns row of cursor.
func (doc *PlainDocument) GetCursorRow() int {
	return doc.cursorRow
}

// GetCursorColumn returns column of cursor.
func (doc *PlainDocument) GetCursorColumn() int {
	return doc.cursorColumn
}
