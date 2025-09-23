package model

type Element interface {
	GetText() string
	GetStartPosition() Position
	GetEndPosition() Position
	GetElement(index int) Element
	GetElementCount() int
}
