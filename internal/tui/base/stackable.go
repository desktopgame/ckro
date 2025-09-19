package base

type Layer struct {
	Control Control
	OnPop   func(stackable Stackable)
}

type Stackable interface {
	Push(layer Layer)
	Pop()
}
