package model

// Segment is reference to multiline string.
type Segment interface {
	// GetSpan returns span at specified line.
	GetSpan(lineIndex int) Span

	// GetLine returns string of line at specified line.
	GetLine(lineIndex int) string

	// GetLineCount returns line count.
	GetLineCount() int
}
