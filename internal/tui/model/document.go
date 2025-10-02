package model

import "io"

type Document interface {
	Read(r Range) Segment
	Render() []Element

	InsertString(row int, bytePos int, s string)
	Remove(row int, bytePos int, byteLen int)
	ReplaceAll(r io.Reader)
	Clear()

	GetLineBytes(index int) int
	GetLineCount() int
	GetVersion() uint
}
