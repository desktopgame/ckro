package model

// Element is part of structured document.
// in example, Heading, Bold, Italic...
type Element interface {
	// GetRange returns range at specified index.
	// zero index is always inclusive all range.
	GetRange(index int) Range

	// GetRangeCount returns count of ranges.
	GetRangeCount() int

	// GetElement returns subelement at specified index.
	GetElement(index int) Element

	// GetElementCount returns count of elements.
	GetElementCount() int
}
