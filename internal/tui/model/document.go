package model

type Document interface {
	Read(r Range) Segment
	Render() []Element

	InsertString(row int, bytePos int, s string)
	Remove(row int, bytePos int, byteLen int)
	Clear()

	GetLineBytes(index int) int
	GetLineCount() int
	GetVersion() uint
}
