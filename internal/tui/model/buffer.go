package model

import (
	"errors"
	"slices"
	"strings"
)

type Line struct {
	content string
}

func (l *Line) PrependString(s string) {
	l.content = s + l.content
}

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

func (l *Line) AppendString(s string) {
	l.content += s
}

func (l *Line) Remove(offset int, length int) {
	l.content = l.content[:offset] + l.content[offset+length:]
}

func (l *Line) GetContent() string {
	return l.content
}

type Buffer struct {
	lines []*Line
}

func (buf *Buffer) Init() {
	buf.lines = []*Line{
		new(Line),
	}
}

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

func (buf *Buffer) PrependString(row int, s string) (Position, error) {
	return buf.InsertString(row, 0, s)
}

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

func (buf *Buffer) AppendString(row int, s string) (Position, error) {
	if row >= 0 && row < len(buf.lines) {
		line := buf.lines[row]
		return buf.InsertString(row, len(line.GetContent()), s)
	}
	return Position{}, errors.New("out of range")
}

func (buf *Buffer) RemoveLine(row int) {
	if row >= 0 && row < len(buf.lines) {
		buf.lines = slices.Delete(buf.lines, row, row+1)
	}
}

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

func (buf *Buffer) GetLineAt(row int) *Line {
	return buf.lines[row]
}

func (buf *Buffer) GetLineCount() int {
	return len(buf.lines)
}
