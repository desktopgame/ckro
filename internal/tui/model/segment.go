package model

type Segment interface {
	GetSpan(lineIndex int) Span
	GetLine(lineIndex int) string
	GetLineCount() int
}
