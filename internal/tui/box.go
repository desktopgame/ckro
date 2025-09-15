package tui

type Box struct {
	Controls    []Control
	orientation Orientation
	x           int
	y           int
}

func (b *Box) Init(orientation Orientation) {
	b.Controls = []Control{}
	b.orientation = orientation
}

func (b *Box) Layout(width int, height int) {
	switch b.orientation {
	case Horizontal:
		offsetX := b.x

		staticWidth := 0
		flexibleControls := 0
		for _, ctrl := range b.Controls {
			if ctrl.IsFlexibleWidth() {
				flexibleControls++
			} else {
				w, _ := ctrl.MinimumSize()
				staticWidth += w
			}
		}

		for _, ctrl := range b.Controls {
			ctrl.Move(offsetX, b.y)

			if ctrl.IsFlexibleWidth() {
				fw := (width - staticWidth) / flexibleControls
				ctrl.Layout(fw, height)

				offsetX += fw
			} else {
				w, _ := ctrl.MinimumSize()
				ctrl.Layout(w, height)

				offsetX += w
			}
		}
	case Vertical:

	}
}

func (b *Box) GetOrientation() Orientation {
	return b.orientation
}
