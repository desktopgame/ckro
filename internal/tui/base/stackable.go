package base

type Layer struct {
	Control         Control
	ClearBackground bool
	OnPop           func()
}

type Stackable interface {
	Push(layer Layer)
	Pop()
}
