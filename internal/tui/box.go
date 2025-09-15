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
		offsetY := b.y

		staticHeight := 0
		flexibleControls := 0
		for _, ctrl := range b.Controls {
			if ctrl.IsFlexibleWidth() {
				flexibleControls++
			} else {
				_, h := ctrl.MinimumSize()
				staticHeight += h
			}
		}

		for _, ctrl := range b.Controls {
			ctrl.Move(b.x, offsetY)

			if ctrl.IsFlexibleWidth() {
				fh := (height - staticHeight) / flexibleControls
				ctrl.Layout(width, fh)

				offsetY += fh
			} else {
				_, h := ctrl.MinimumSize()
				ctrl.Layout(width, h)

				offsetY += h
			}
		}
	}
}

func (b *Box) GetOrientation() Orientation {
	return b.orientation
}
