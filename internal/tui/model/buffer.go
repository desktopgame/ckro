package model

import (
	"errors"
	"slices"
	"strings"
)

// Line is parts of buffer.
type Line struct {
	content string
}

// PrependString is insert string into line ahead.
func (l *Line) PrependString(s string) {
	l.content = s + l.content
}

// InsertString is insert string into specified column.
func (l *Line) InsertString(column int, s string) {
	if len(s) == 0 {
		return
	}
	if column == 0 {
		l.content = s + l.content
	} else if column == len(l.content) {
		l.content = l.content + s
	} else {
		l.content = l.content[:column] + s + l.content[column:]
	}
}

// AppendString is appending string into tail.
func (l *Line) AppendString(s string) {
	l.content += s
}

// Remove is remove string specified range.
func (l *Line) Remove(offset int, length int) {
	l.content = l.content[:offset] + l.content[offset+length:]
}

// GetContent is returns string of content.
func (l *Line) GetContent() string {
	return l.content
}

// Buffer is array of line strings.
// Buffer API's require and returns unit of codepoint positions.
type Buffer struct {
	lines []*Line
}

// Init is initialize Buffer.
func (buf *Buffer) Init() {
	buf.lines = []*Line{
		new(Line),
	}
}

// InsertLine is break line at specified position.
func (buf *Buffer) InsertLine(row int, column int) *Line {
	if column == 0 {
		newLine := &Line{}
		if len(buf.lines) == 1 {
			buf.lines = append(buf.lines, newLine)
		} else {
			buf.lines = slices.Insert(buf.lines, row, newLine)
		}
		return newLine
	} else if column == len(buf.lines[row].GetContent()) {
		newLine := &Line{}
		buf.lines = slices.Insert(buf.lines, row+1, newLine)
		return newLine
	} else {
		line := buf.lines[row]
		br := line.GetContent()[column:]
		newLine := &Line{}
		newLine.InsertString(0, br)
		line.Remove(column, len(line.GetContent())-column)
		buf.lines = slices.Insert(buf.lines, row+1, newLine)
		return newLine
	}
}

// PrependString is insert string into specified line ahead.
func (buf *Buffer) PrependString(row int, s string) (Position, error) {
	return buf.InsertString(row, 0, s)
}

// InsertString is insert string into specified position.
func (buf *Buffer) InsertString(row int, column int, s string) (Position, error) {
	if row >= 0 && row < len(buf.lines) {
		line := buf.lines[row]

		insertLines := strings.Split(s, "\n")
		position := Position{
			Row:    row,
			Column: column,
		}

		if len(insertLines) == 0 {
			line.InsertString(column, s)
			position.Column += len(s)
		} else {
			line.InsertString(column, insertLines[0])
			breakAt := column + len(insertLines[0])
			position.Column = breakAt

			for i := 1; i < len(insertLines); i++ {
				nextLine := buf.InsertLine(row, breakAt)
				nextLine.PrependString(insertLines[i])

				position.Row += 1
				position.Column = len(insertLines[i])

				breakAt = len(insertLines[i])
				row += 1
			}
		}
		return position, nil
	}
	return Position{}, errors.New("out of range")
}

// AppendString is appending string into tail.
func (buf *Buffer) AppendString(row int, s string) (Position, error) {
	if row >= 0 && row < len(buf.lines) {
		line := buf.lines[row]
		return buf.InsertString(row, len(line.GetContent()), s)
	}
	return Position{}, errors.New("out of range")
}

// RemoveLine is remove specified line.
func (buf *Buffer) RemoveLine(row int) {
	if row >= 0 && row < len(buf.lines) {
		buf.lines = slices.Delete(buf.lines, row, row+1)
	}
}

// RemoveString is remove string specified range.
func (buf *Buffer) RemoveString(row int, column int, length int) {
	if row >= 0 && row < len(buf.lines) {
		line := buf.lines[row]

		for length > 0 {
			removeChars := min(len(line.GetContent())-column, length)
			line.Remove(column, removeChars)
			length -= removeChars
			length -= 1

			if length > 0 {
				row += 1
				column = 0
				line = buf.lines[row]
			}
		}
	}
}

// GetLineAt returns specified line.
func (buf *Buffer) GetLineAt(row int) *Line {
	return buf.lines[row]
}

// GetLineCount returns count of lines.
func (buf *Buffer) GetLineCount() int {
	return len(buf.lines)
}
