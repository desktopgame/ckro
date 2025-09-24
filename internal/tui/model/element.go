package model

type Element interface {
	GetStyle() *Style
	GetText() string
	GetStartPosition() Position
	GetEndPosition() Position
	GetElement(index int) Element
	GetElementCount() int
}
