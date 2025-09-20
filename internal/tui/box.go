package tui

// Box is layout sub controls by horizontal or vertical.
type Box struct {
	Controls    []Control
	orientation Orientation
	x           int
	y           int
}

// Init is initialize Box.
func (b *Box) Init(orientation Orientation) {
	b.Controls = []Control{}
	b.orientation = orientation
}

// Traverse is register sub controls into FocusManager by the array order.
func (b *Box) Traverse(fm *FocusManager) {
	for _, ctrl := range b.Controls {
		ctrl.Traverse(fm)
	}
}

// Update is delegate to sub controls.
func (b *Box) Update() {
	for _, ctrl := range b.Controls {
		ctrl.Update()
	}
}

// Draw is delegate to sub controls.
func (b *Box) Draw(g *Graphics) {
	for _, ctrl := range b.Controls {
		ctrl.Draw(g)
	}
}

// MinimumSize is calculate minimum size in consideration flex.
// flexible controls assigned evenly divided size of subtract static size from total size.
func (b *Box) MinimumSize(width int, height int) (Width int, Height int) {
	mw := width
	mh := height

	switch b.orientation {
	case Horizontal:
		mh = -1
		totalMW := 0

		staticWidth := 0
		flexibleControls := 0
		for _, ctrl := range b.Controls {
			if ctrl.IsFlexibleWidth() {
				flexibleControls++
			} else {
				w, _ := ctrl.MinimumSize(width, height)
				staticWidth += w
			}
		}

		for _, ctrl := range b.Controls {
			w, h := ctrl.MinimumSize(mw, height)
			if ctrl.IsFlexibleWidth() {
				// w = (mw - staticWidth) / flexibleControls
				w, h = ctrl.MinimumSize((mw-staticWidth)/flexibleControls, height)
			}

			if h > mh {
				mh = h
			}

			totalMW += w
		}

		if mh == -1 {
			mh = height
		}
		mw = max(mw, totalMW)
	case Vertical:
		mw = -1
		totalMH := 0

		staticHeight := 0
		flexibleControls := 0
		for _, ctrl := range b.Controls {
			if ctrl.IsFlexibleHeight() {
				flexibleControls++
			} else {
				_, h := ctrl.MinimumSize(width, height)
				staticHeight += h
			}
		}

		for _, ctrl := range b.Controls {
			w, h := ctrl.MinimumSize(width, mh)
			if ctrl.IsFlexibleHeight() {
				//h = (mh - staticHeight) / flexibleControls
				w, h = ctrl.MinimumSize(width, (mh-staticHeight)/flexibleControls)
			}

			if w > mw {
				mw = w
			}

			totalMH += h
		}

		if mw == -1 {
			mw = width
		}
		mh = max(mh, totalMH)
	}

	return mw, mh
}

// Move is move the base point of Box.
func (b *Box) Move(x int, y int) {
	b.x = x
	b.y = y
}

// Layout is arrangement sub controls into line by horizontal or vertical.
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
				w, _ := ctrl.MinimumSize(width, height)
				staticWidth += w
			}
		}

		for _, ctrl := range b.Controls {
			ctrl.Move(offsetX, b.y)

			if ctrl.IsFlexibleWidth() {
				fw := (width - staticWidth) / flexibleControls
				_, h := ctrl.MinimumSize(width, height)
				if ctrl.IsFlexibleHeight() {
					h = height
				}
				ctrl.Layout(fw, h)

				offsetX += fw
			} else {
				w, h := ctrl.MinimumSize(width, height)
				if ctrl.IsFlexibleHeight() {
					h = height
				}
				ctrl.Layout(w, h)

				offsetX += w
			}
		}
	case Vertical:
		offsetY := b.y

		staticHeight := 0
		flexibleControls := 0
		for _, ctrl := range b.Controls {
			if ctrl.IsFlexibleHeight() {
				flexibleControls++
			} else {
				_, h := ctrl.MinimumSize(width, height)
				staticHeight += h
			}
		}

		for _, ctrl := range b.Controls {
			ctrl.Move(b.x, offsetY)

			if ctrl.IsFlexibleHeight() {
				fh := (height - staticHeight) / flexibleControls
				w, _ := ctrl.MinimumSize(width, height)
				if ctrl.IsFlexibleWidth() {
					w = width
				}
				ctrl.Layout(w, fh)

				offsetY += fh
			} else {
				w, h := ctrl.MinimumSize(width, height)
				if ctrl.IsFlexibleWidth() {
					w = width
				}
				ctrl.Layout(w, h)

				offsetY += h
			}
		}
	}
}

// IsFlexibleWidth returns true if have to any flexible width controls.
func (b *Box) IsFlexibleWidth() bool {
	for _, ctrl := range b.Controls {
		if ctrl.IsFlexibleWidth() {
			return true
		}
	}
	return false
}

// IsFlexibleHeight returns true if have to any flexible height controls.
func (b *Box) IsFlexibleHeight() bool {
	for _, ctrl := range b.Controls {
		if ctrl.IsFlexibleHeight() {
			return true
		}
	}
	return false
}

// GetOrientation returns orientation.
func (b *Box) GetOrientation() Orientation {
	return b.orientation
}
