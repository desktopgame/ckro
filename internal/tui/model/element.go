package model

type Element interface {
	GetRange(index int) Range
	GetRangeCount() int
	GetElement(index int) Element
	GetElementCount() int
}
