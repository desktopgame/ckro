package model

// Span is range of single line.
// end column is exclusive.
type Span struct {
	StartColumn int
	EndColumn   int
}

func (sp Span) IsZero() bool {
	return sp.StartColumn == sp.EndColumn
}
