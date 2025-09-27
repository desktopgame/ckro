package model

type Segment interface {
	GetLine(lineIndex int) string
	GetLineCount() int
}
