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
			buf.updateTracksAfterInsert(row, column, 1, 0)
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
			buf.updateTracksAfterInsert(row, column, 0, len(s))
		} else {
			line.InsertString(column, insertLines[0])
			breakAt := column + len(insertLines[0])
			position.Column = breakAt

			insertedLines := 0
			lastLineLen := 0
			for i := 1; i < len(insertLines); i++ {
				nextLine := buf.insertLine(row, breakAt)
				nextLine.PrependString(insertLines[i])

				position.Row += 1
				position.Column = len(insertLines[i])

				breakAt = len(insertLines[i])
				lastLineLen = len(insertLines[i])
				row += 1
				insertedLines++
			}
			buf.updateTracksAfterMultiLineInsert(row-insertedLines, column, insertedLines, lastLineLen)
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
		startRow := row
		startColumn := column
		endRow := row
		endColumn := column
		line := buf.lines[row]

		for length > 0 {
			lineLen := len(line.GetContent())
			if lineLen == 0 {
				buf.removeLine(row)
				endRow = row
				endColumn = 0
				length--
			} else if column == lineLen {
				nextLine := buf.GetLineAt(row + 1)
				buf.GetLineAt(row).AppendString(nextLine.GetContent())
				buf.removeLine(row + 1)
				endRow = row + 1
				endColumn = 0
				length--
			} else {
				removeChars := min(lineLen-column, length)

				line.Remove(column, removeChars)
				endRow = row
				endColumn = column + removeChars
				length -= removeChars
				length -= 1

				if length > 0 {
					row += 1
					column = 0
					line = buf.lines[row]
				}
			}
		}

		buf.updateTracksAfterRemove(startRow, startColumn, endRow, endColumn)
	}
}

// ReplaceAll is replace a all content by Reader content
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

// CreateTrack creates a new position tracker.
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

// RemoveTrack removes a position tracker from the buffer.
func (buf *Buffer) RemoveTrack(track *Track) {
	for i, t := range buf.tracks {
		if t == track {
			buf.tracks = slices.Delete(buf.tracks, i, i+1)
			return
		}
	}
}

// updateTracksAfterInsert updates all tracked positions after a single-line insertion.
// Parameters:
//   - row: the row where insertion occurred
//   - column: the column where insertion occurred
//   - insertedLines: number of lines inserted (1 for newline, 0 for single-line)
//   - insertedColumns: number of columns inserted (for single-line insertion)
func (buf *Buffer) updateTracksAfterInsert(row int, column int, insertedLines int, insertedColumns int) {
	for _, track := range buf.tracks {
		if track.Lost {
			continue
		}

		if insertedLines > 0 {
			// Newline insertion
			if track.Position.Row > row {
				// Track is after insertion row → shift down
				track.Position.Row += insertedLines
			} else if track.Position.Row == row && track.Position.Column > column {
				// Track is on same row, after insertion point → move to next line
				track.Position.Row += insertedLines
				track.Position.Column -= column
			}
		} else {
			// Single-line insertion
			if track.Position.Row == row && track.Position.Column >= column {
				// Track is on same row, at or after insertion point → shift right
				track.Position.Column += insertedColumns
			}
		}
	}
}

// updateTracksAfterMultiLineInsert updates all tracked positions after a multi-line insertion.
// Parameters:
//   - row: the row where insertion occurred
//   - column: the column where insertion occurred
//   - insertedLines: number of lines inserted
//   - lastLineLen: length of the last inserted line
func (buf *Buffer) updateTracksAfterMultiLineInsert(row int, column int, insertedLines int, lastLineLen int) {
	for _, track := range buf.tracks {
		if track.Lost {
			continue
		}

		if track.Position.Row > row {
			// Track is after insertion row → shift down
			track.Position.Row += insertedLines
		} else if track.Position.Row == row && track.Position.Column >= column {
			// Track is on same row, at or after insertion point
			// Move to new line with adjusted column
			track.Position.Row += insertedLines
			track.Position.Column = lastLineLen + (track.Position.Column - column)
		}
	}
}

// updateTracksAfterRemove updates all tracked positions after a removal.
// Parameters:
//   - startRow: the row where removal started
//   - startColumn: the column where removal started
//   - endRow: the row where removal ended
//   - endColumn: the column where removal ended
func (buf *Buffer) updateTracksAfterRemove(startRow int, startColumn int, endRow int, endColumn int) {
	removedLines := endRow - startRow

	for _, track := range buf.tracks {
		if track.Lost {
			continue
		}

		trackPos := track.Position

		if removedLines == 0 {
			// Single-line removal
			if trackPos.Row == startRow {
				if trackPos.Column >= startColumn && trackPos.Column < endColumn {
					// Track is within removed range → mark as lost
					track.Lost = true
				} else if trackPos.Column >= endColumn {
					// Track is after removed range → shift left
					track.Position.Column -= (endColumn - startColumn)
				}
			}
		} else {
			// Multi-line removal
			if trackPos.Row >= startRow && trackPos.Row <= endRow {
				// Track is within removed line range
				if trackPos.Row == startRow && trackPos.Column < startColumn {
					// Track is before removal start on start line → keep as is
					continue
				} else if trackPos.Row == endRow && trackPos.Column >= endColumn {
					// Track is after removal end on end line → move to start line
					track.Position.Row = startRow
					track.Position.Column = startColumn + (trackPos.Column - endColumn)
				} else {
					// Track is within removed range → mark as lost
					track.Lost = true
				}
			} else if trackPos.Row > endRow {
				// Track is after removed range → shift up
				track.Position.Row -= removedLines
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
