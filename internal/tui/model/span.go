package model

type Span struct {
	StartColumn int
	EndColumn   int
}

func (sp Span) IsZero() bool {
	return sp.StartColumn == sp.EndColumn
}
