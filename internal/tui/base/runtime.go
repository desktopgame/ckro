package base

// Runtime is provide interface to for driving application.
// in this, contain supports stack of controls.
// controls able to erase self from stack through this Runtime object.
type Runtime interface {
	Stackable

	// BeginBackground is increment number of background tasks.
	// Runtime is ignore events while zero than bigger of number of background tasks.
	BeginBackground()

	// EndBackground is decrement number of background tasks.
	EndBackground()

	// Repaint is do request rerender terminal.
	Repaint()
}
