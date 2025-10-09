package model

// Range is range presentation by two Position.
// end column is exclusive.
type Range struct {
	StartPosition Position
	EndPosition   Position
}

func (r Range) IsZero() bool {
	return r.StartPosition.Row == r.EndPosition.Row && r.StartPosition.Column == r.EndPosition.Column
}
