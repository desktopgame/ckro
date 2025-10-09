package model

// Track is position of moves depend on Document changes.
type Track struct {
	Position Position
	Lost     bool
}
