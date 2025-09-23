package model

type Element interface {
	GetText() string
	GetElement(index int) Element
	GetElementCount() int
}
