package tui

func WithFrame(ctrl Control) *Frame {
	fr := Frame{}
	fr.Control = ctrl
	return &fr
}

func WithCenter(ctrl Control, width int, height int) *Center {
	c := Center{}
	c.Control = ctrl
	c.Width = width
	c.Height = height
	return &c
}
