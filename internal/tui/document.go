package tui

type Document struct {
	buffer       Buffer
	cursorRow    int
	cursorColumn int
}

func (doc *Document) Init() {
	doc.buffer = Buffer{}
	doc.buffer.Init()
	doc.cursorRow = 0
	doc.cursorColumn = 0
}

func (doc *Document) InsertLine() {
	doc.buffer.InsertLine(doc.cursorRow, doc.cursorColumn)
	doc.cursorRow++
	doc.cursorColumn = 0
}

func (doc *Document) InsertString(s string) {
	at, err := doc.buffer.InsertString(doc.cursorRow, doc.cursorColumn, s)
	if err != nil {
		panic("illegal state")
	}
	doc.cursorRow = at.Row
	doc.cursorColumn = at.Column
}

func (doc *Document) RemoveChar() {
	currLine := doc.buffer.GetLineAt(doc.cursorRow)
	if doc.cursorRow > 0 {
		if len(currLine.GetContent()) == 0 || doc.cursorColumn == 0 {
			aboveLine := doc.buffer.GetLineAt(doc.cursorRow - 1)
			cursorAt := len(aboveLine.GetContent())

			doc.buffer.AppendString(doc.cursorRow-1, currLine.GetContent())
			doc.buffer.RemoveLine(doc.cursorRow)
			doc.cursorRow--
			doc.cursorColumn = cursorAt
		} else if doc.cursorColumn > 0 {
			doc.buffer.RemoveString(doc.cursorRow, doc.cursorColumn-1, 1)
			doc.cursorColumn--
		}
	} else {
		if doc.cursorColumn > 0 {
			doc.buffer.RemoveString(doc.cursorRow, doc.cursorColumn-1, 1)
			doc.cursorColumn--
		}
	}
}

func (doc *Document) MoveLeft() {
	if doc.cursorColumn > 0 {
		doc.cursorColumn--
	} else {
		if doc.cursorRow > 0 {
			doc.cursorRow--
			doc.cursorColumn = len(doc.buffer.GetLineAt(doc.cursorRow).GetContent())
		}
	}
}

func (doc *Document) MoveRight() {
	currLine := doc.buffer.GetLineAt(doc.cursorRow)
	if doc.cursorColumn < len(currLine.GetContent()) {
		doc.cursorColumn++
	} else {
		if doc.cursorRow < doc.buffer.GetLineCount()-1 {
			doc.cursorRow++
			doc.cursorColumn = 0
		}
	}
}

func (doc *Document) MoveUp() {
	if doc.cursorRow > 0 {
		doc.cursorRow--

		currLine := doc.buffer.GetLineAt(doc.cursorRow)
		if doc.cursorColumn > len(currLine.GetContent()) {
			doc.cursorColumn = len(currLine.GetContent())
		}
	}
}

func (doc *Document) MoveDown() {
	if doc.cursorRow < doc.buffer.GetLineCount()-1 {
		doc.cursorRow++

		currLine := doc.buffer.GetLineAt(doc.cursorRow)
		if doc.cursorColumn > len(currLine.GetContent()) {
			doc.cursorColumn = len(currLine.GetContent())
		}
	}
}

func (doc *Document) GetCursorRow() int {
	return doc.cursorRow
}

func (doc *Document) GetCursorColumn() int {
	return doc.cursorColumn
}
