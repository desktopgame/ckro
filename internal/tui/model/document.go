package model

type Document interface {
	Read(r Range) Segment
	Render() []Element

	WriteString(row int, bytePos int, s string)
	Remove(row int, bytePos int, byteLen int)

	GetVersion() uint

	InsertLine()
	InsertString(s string)
	RemoveChar()
	Clear()

	FindPrev(searchStr string) bool
	FindNext(searchStr string) bool
	Replace(deleteCount int, replaceStr string) bool

	MoveLeft()
	MoveRight()
	MoveUp()
	MoveDown()
	MoveReset()
	GetBuffer() *Buffer
	GetCursorRow() int
	GetCursorColumn() int
}
