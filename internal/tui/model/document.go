package model

import "github.com/desktopgame/ckro/internal/text"

// Document is wrapper of Buffer.
// track a current cursor position.
type Document struct {
	buffer       Buffer
	cursorRow    int
	cursorColumn int
}

// Init is initialize Buffer.
func (doc *Document) Init() {
	doc.buffer = Buffer{}
	doc.buffer.Init()
	doc.cursorRow = 0
	doc.cursorColumn = 0
}

// InsertLine is break line at current cursor position.
func (doc *Document) InsertLine() {
	currLine := doc.buffer.GetLineAt(doc.cursorRow)
	codepointColumn := text.GraphemeToCodepointPos(currLine.GetContent(), doc.cursorColumn)

	doc.buffer.InsertLine(doc.cursorRow, codepointColumn)
	doc.cursorRow++
	doc.cursorColumn = 0
}

// InsertString is insert string at current cursor position.
func (doc *Document) InsertString(s string) {
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
func (doc *Document) RemoveChar() {
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

// MoveLeft is move cursor to left.
func (doc *Document) MoveLeft() {
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
func (doc *Document) MoveRight() {
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
func (doc *Document) MoveUp() {
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
func (doc *Document) MoveDown() {
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
func (doc *Document) MoveReset() {
	doc.cursorRow = 0
	doc.cursorColumn = 0
}

// GetBuffer returns Buffer.
func (doc *Document) GetBuffer() *Buffer {
	return &doc.buffer
}

// GetCursorRow returns row of cursor.
func (doc *Document) GetCursorRow() int {
	return doc.cursorRow
}

// GetCursorColumn returns column of cursor.
func (doc *Document) GetCursorColumn() int {
	return doc.cursorColumn
}
