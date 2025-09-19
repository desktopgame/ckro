package base

type Runtime interface {
	Stackable

	BeginBackground()
	EndBackground()
	Repaint()
}
