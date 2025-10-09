package model

import "io"

// Document is model for TextBox, and contain string.
// most important method of Document is Render.
// TextBox get Element through Render method for rendering.
// and TextViewResolver is mapping Element to TextView.
// This results in TextBox is got TextView from structured Element.
// This is it rendering core.
// more details is see TextView, TextViewResolver
type Document interface {
	// Read returns Segment on specified range.
	Read(r Range) Segment

	// Render is presentation a contain string as structured element.
	Render() []Element

	// InsertString is insert string into specified position.
	InsertString(row int, bytePos int, s string)

	// Remove is remove string in specified range.
	Remove(row int, bytePos int, byteLen int)

	// ReplaceAll is replace a all content by Reader content
	ReplaceAll(r io.Reader)

	// Clear is remove all strings.
	Clear()

	// CreateTrack creates a new position tracker.
	CreateTrack(row int, bytePos int) *Track

	// GetLineBytes returns total byte count at specified line.
	GetLineBytes(index int) int

	// GetLineCount returns total line count.
	GetLineCount() int

	// GetVersion returns version. version is changed on modify string.
	GetVersion() uint
}
