package tui

import "github.com/gdamore/tcell/v2"

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

func (b *Box) Update() {
	for _, ctrl := range b.Controls {
		ctrl.Update()
	}
}

func (b *Box) Handle(ev tcell.Event) {
	for _, ctrl := range b.Controls {
		ctrl.Handle(ev)
	}
}

func (b *Box) Draw(s tcell.Screen) {
	for _, ctrl := range b.Controls {
		ctrl.Draw(s)
	}
}

func (b *Box) MinimumSize(width int, height int) (Width int, Height int) {
	mw := width
	mh := height

	switch b.orientation {
	case Horizontal:
		mh = -1
		totalMW := 0

		for _, ctrl := range b.Controls {
			w, h := ctrl.MinimumSize(mw, height)

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

		for _, ctrl := range b.Controls {
			w, h := ctrl.MinimumSize(width, mh)

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

func (b *Box) Move(x int, y int) {
	b.x = x
	b.y = y
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
				w, _ := ctrl.MinimumSize(width, height)
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
				ctrl.Layout(width, fh)

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

func (b *Box) IsFlexibleWidth() bool {
	for _, ctrl := range b.Controls {
		if ctrl.IsFlexibleWidth() {
			return true
		}
	}
	return false
}

func (b *Box) IsFlexibleHeight() bool {
	for _, ctrl := range b.Controls {
		if ctrl.IsFlexibleHeight() {
			return true
		}
	}
	return false
}

func (b *Box) GetOrientation() Orientation {
	return b.orientation
}
