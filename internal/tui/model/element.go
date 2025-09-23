package model

type Element interface {
	GetElement(index int) Element
	GetElementCount() int
}
