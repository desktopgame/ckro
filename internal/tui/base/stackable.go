package base

type Layer struct {
	Control         Control
	ClearBackground bool
	OnPop           func(returnCode int)
}

type Stackable interface {
	Push(layer Layer)
	Pop(returnCode int)
}
