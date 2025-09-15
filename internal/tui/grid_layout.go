package tui

import "github.com/desktopgame/ckro/internal/tui/presenter"

type GridLayout struct {
	rowCount    int
	columnCount int
	controls    [][]*Tile
}

func (gl *GridLayout) Init(rowCount int, columnCount int) {
	gl.rowCount = rowCount
	gl.columnCount = columnCount

	for i := 0; i < rowCount; i++ {
		line := []*Tile{}
		for j := 0; j < columnCount; j++ {
			line = append(line, nil)
		}
		gl.controls = append(gl.controls, line)
	}
}

func (gl *GridLayout) Set(row int, column int, minimunWidth int, minimumHeight int, flexibleWidth bool, flexibleHeight bool) {
	t := &Tile{}
	t.Init()
	t.MinimumWidth = minimunWidth
	t.MinimumHeight = minimumHeight
	t.FlexibleWidth = flexibleWidth
	t.FlexibleHeight = flexibleHeight
	t.TextPresenter = &presenter.FrameTextPresenter{}
	gl.controls[row][column] = t
}

func (gl *GridLayout) SetStatic(row int, column int, minimunWidth int, minimumHeight int) {
	gl.Set(row, column, minimunWidth, minimumHeight, false, false)
}

func (gl *GridLayout) SetFlex(row int, column int) {
	gl.Set(row, column, 0, 0, true, true)
}

func (gl *GridLayout) Build() *Box {
	vbox := Box{}
	vbox.Init(Vertical)

	header := &Tile{}
	header.Init()
	header.FlexibleWidth = true
	header.FlexibleHeight = true
	vbox.Controls = append(vbox.Controls, header)

	/*
		line1 := &Tile{}
		line1.Init()
		line1.MinimumHeight = 1
		line1.FlexibleWidth = true
		line1.TextPresenter = &presenter.HorizontalSeparatorTextPresenter{}
		vbox.Controls = append(vbox.Controls, line1)
	*/

	for i := 0; i < gl.rowCount; i++ {
		hbox := Box{}
		hbox.Init(Horizontal)

		space1 := &Tile{}
		space1.Init()
		space1.FlexibleWidth = true
		space1.FlexibleHeight = true
		hbox.Controls = append(hbox.Controls, space1)

		t1 := &Tile{}
		t1.Init()
		t1.MinimumWidth = 1
		t1.FlexibleHeight = true
		t1.TextPresenter = &presenter.VerticalSeparatorTextPresenter{}
		hbox.Controls = append(hbox.Controls, t1)

		for j := 0; j < gl.columnCount; j++ {
			t := gl.controls[i][j]
			hbox.Controls = append(hbox.Controls, t)
		}

		t2 := &Tile{}
		t2.Init()
		t2.MinimumWidth = 1
		t2.FlexibleHeight = true
		t2.TextPresenter = &presenter.VerticalSeparatorTextPresenter{}
		hbox.Controls = append(hbox.Controls, t2)

		space2 := &Tile{}
		space2.Init()
		space2.FlexibleWidth = true
		space2.FlexibleHeight = true
		hbox.Controls = append(hbox.Controls, space2)

		vbox.Controls = append(vbox.Controls, &hbox)
	}
	/*
		line2 := &Tile{}
		line2.Init()
		line2.MinimumHeight = 1
		line2.FlexibleWidth = true
		line2.TextPresenter = &presenter.HorizontalSeparatorTextPresenter{}
		vbox.Controls = append(vbox.Controls, line2)
	*/

	footer := &Tile{}
	footer.Init()
	footer.FlexibleWidth = true
	footer.FlexibleHeight = true
	vbox.Controls = append(vbox.Controls, footer)
	return &vbox
}
