package base

// Layer is stackable object.
type Layer struct {
	Control         Control
	ClearBackground bool
	OnPop           func(returnCode int)
}

// Stackable is stack of layer.
// by doing layers draw by index order, to represent controls overlay.
type Stackable interface {
	// Push is stack a layer.
	Push(layer Layer)

	// Pop is remove layer from stack top.
	Pop(returnCode int)
}
