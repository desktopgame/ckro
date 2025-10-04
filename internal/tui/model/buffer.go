package model

import (
	"bufio"
	"errors"
	"io"
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
	lines  []*Line
	tracks []*Track
}

// Init is initialize Buffer.
func (buf *Buffer) Init() {
	buf.lines = []*Line{
		new(Line),
	}
}

func (buf *Buffer) insertLine(row int, column int) *Line {
	if column == 0 {
		newLine := &Line{}
		if column == len(buf.lines[row].GetContent()) {
			buf.lines = slices.Insert(buf.lines, row+1, newLine)
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

// InsertString is insert string into specified position.
func (buf *Buffer) InsertString(row int, column int, s string) (Position, error) {
	if row >= 0 && row < len(buf.lines) {
		line := buf.lines[row]

		if s == "\n" {
			buf.insertLine(row, column)
			return Position{}, nil
		}

		insertLines := strings.Split(s, "\n")
		position := Position{
			Row:    row,
			Column: column,
		}

		if len(insertLines) == 1 {
			line.InsertString(column, s)
			position.Column += len(s)
		} else {
			line.InsertString(column, insertLines[0])
			breakAt := column + len(insertLines[0])
			position.Column = breakAt

			for i := 1; i < len(insertLines); i++ {
				nextLine := buf.insertLine(row, breakAt)
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

// RemoveLine is remove specified line.
func (buf *Buffer) removeLine(row int) {
	if row >= 0 && row < len(buf.lines) {
		buf.lines = slices.Delete(buf.lines, row, row+1)
	}
}

// RemoveString is remove string specified range.
func (buf *Buffer) RemoveString(row int, column int, length int) {
	if row >= 0 && row < len(buf.lines) {
		line := buf.lines[row]

		for length > 0 {
			lineLen := len(line.GetContent())
			if lineLen == 0 {
				buf.removeLine(row)
				length--
			} else if column == lineLen {
				nextLine := buf.GetLineAt(row + 1)
				buf.GetLineAt(row).AppendString(nextLine.GetContent())
				buf.removeLine(row + 1)
				length--
			} else {
				removeChars := min(lineLen-column, length)

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
}

func (buf *Buffer) ReplaceAll(r io.Reader) {
	lines := []*Line{}
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 64<<10), 16<<20)
	for sc.Scan() {
		line := sc.Text()
		lines = append(lines, &Line{
			content: line,
		})
	}
	if len(lines) == 0 {
		lines = []*Line{
			new(Line),
		}
	}
	buf.lines = lines

	for _, t := range buf.tracks {
		t.Lost = true
	}
	buf.tracks = nil
}

func (buf *Buffer) CreateTrack(row int, bytePos int) *Track {
	t := &Track{
		Position: Position{
			Row:    row,
			Column: bytePos,
		},
		Lost: false,
	}
	buf.tracks = append(buf.tracks, t)
	return t
}

// GetLineAt returns specified line.
func (buf *Buffer) GetLineAt(row int) *Line {
	return buf.lines[row]
}

// GetLineCount returns count of lines.
func (buf *Buffer) GetLineCount() int {
	return len(buf.lines)
}
