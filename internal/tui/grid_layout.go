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

func (gl *GridLayout) Set(row int, column int, minimumWidth int, minimumHeight int, flexibleWidth bool, flexibleHeight bool, presenter TextPresenter) *Tile {
	if row < 0 || row >= gl.rowCount || column < 0 || column >= gl.columnCount {
		return nil
	}

	t := &Tile{}
	t.Init()
	t.MinimumWidth = minimumWidth
	t.MinimumHeight = minimumHeight
	t.FlexibleWidth = flexibleWidth
	t.FlexibleHeight = flexibleHeight
	if presenter != nil {
		t.TextPresenter = presenter
	}
	gl.controls[row][column] = t
	return t
}

func (gl *GridLayout) SetStatic(row int, column int, minimunWidth int, minimumHeight int) {
	gl.Set(row, column, minimunWidth, minimumHeight, false, false, nil)
}

func (gl *GridLayout) SetFlex(row int, column int) {
	gl.Set(row, column, 0, 0, true, true, nil)
}

func (gl *GridLayout) Build() *Box {
	vbox := Box{}
	vbox.Init(Vertical)

	// Top spacing
	header := &Tile{}
	header.Init()
	header.FlexibleWidth = true
	header.FlexibleHeight = true
	vbox.Controls = append(vbox.Controls, header)

	// Top border row with proper spacing
	topBorderRow := Box{}
	topBorderRow.Init(Horizontal)

	// Left spacing for top border
	topLeftSpace := &Tile{}
	topLeftSpace.Init()
	topLeftSpace.FlexibleWidth = true
	topLeftSpace.FlexibleHeight = true
	topBorderRow.Controls = append(topBorderRow.Controls, topLeftSpace)

	// Left border for top
	topLeftBorder := &Tile{}
	topLeftBorder.Init()
	topLeftBorder.MinimumWidth = 1
	topLeftBorder.MinimumHeight = 1
	topLeftBorder.TextPresenter = &presenter.HorizontalSeparatorTextPresenter{}
	topBorderRow.Controls = append(topBorderRow.Controls, topLeftBorder)

	// Content area borders - match the actual content width
	for j := 0; j < gl.columnCount; j++ {
		contentBorder := &Tile{}
		contentBorder.Init()
		contentBorder.MinimumHeight = 1
		contentBorder.FlexibleWidth = false

		// Set minimum width based on the content in this column
		if gl.controls[0][j] != nil {
			contentBorder.MinimumWidth = gl.controls[0][j].MinimumWidth
			if gl.controls[0][j].FlexibleWidth {
				contentBorder.FlexibleWidth = true
			}
		} else {
			contentBorder.MinimumWidth = 6 // Default width for empty cells
		}

		contentBorder.TextPresenter = &presenter.HorizontalSeparatorTextPresenter{}
		topBorderRow.Controls = append(topBorderRow.Controls, contentBorder)

		// Add separator between columns (except after last column)
		if j < gl.columnCount-1 {
			separator := &Tile{}
			separator.Init()
			separator.MinimumWidth = 1
			separator.MinimumHeight = 1
			separator.TextPresenter = &presenter.HorizontalSeparatorTextPresenter{}
			topBorderRow.Controls = append(topBorderRow.Controls, separator)
		}
	}

	// Right border for top
	topRightBorder := &Tile{}
	topRightBorder.Init()
	topRightBorder.MinimumWidth = 1
	topRightBorder.MinimumHeight = 1
	topRightBorder.TextPresenter = &presenter.HorizontalSeparatorTextPresenter{}
	topBorderRow.Controls = append(topBorderRow.Controls, topRightBorder)

	// Right spacing for top border
	topRightSpace := &Tile{}
	topRightSpace.Init()
	topRightSpace.FlexibleWidth = true
	topRightSpace.FlexibleHeight = true
	topBorderRow.Controls = append(topBorderRow.Controls, topRightSpace)

	vbox.Controls = append(vbox.Controls, &topBorderRow)

	// Grid rows with separators
	for i := 0; i < gl.rowCount; i++ {
		// Create row
		hbox := Box{}
		hbox.Init(Horizontal)

		// Left spacing
		leftSpace := &Tile{}
		leftSpace.Init()
		leftSpace.FlexibleWidth = true
		leftSpace.FlexibleHeight = true
		hbox.Controls = append(hbox.Controls, leftSpace)

		// Left border
		leftBorder := &Tile{}
		leftBorder.Init()
		leftBorder.MinimumWidth = 1
		leftBorder.FlexibleHeight = true
		leftBorder.TextPresenter = &presenter.VerticalSeparatorTextPresenter{}
		hbox.Controls = append(hbox.Controls, leftBorder)

		// Add controls with vertical separators between them
		for j := 0; j < gl.columnCount; j++ {
			// Add the control (or empty tile if nil)
			if gl.controls[i][j] != nil {
				hbox.Controls = append(hbox.Controls, gl.controls[i][j])
			} else {
				// Create empty flexible tile if no control is set
				emptyTile := &Tile{}
				emptyTile.Init()
				emptyTile.FlexibleWidth = true
				emptyTile.FlexibleHeight = true
				hbox.Controls = append(hbox.Controls, emptyTile)
			}

			// Add vertical separator between columns (except after last column)
			if j < gl.columnCount-1 {
				separator := &Tile{}
				separator.Init()
				separator.MinimumWidth = 1
				separator.FlexibleHeight = true
				separator.TextPresenter = &presenter.VerticalSeparatorTextPresenter{}
				hbox.Controls = append(hbox.Controls, separator)
			}
		}

		// Right border
		rightBorder := &Tile{}
		rightBorder.Init()
		rightBorder.MinimumWidth = 1
		rightBorder.FlexibleHeight = true
		rightBorder.TextPresenter = &presenter.VerticalSeparatorTextPresenter{}
		hbox.Controls = append(hbox.Controls, rightBorder)

		// Right spacing
		rightSpace := &Tile{}
		rightSpace.Init()
		rightSpace.FlexibleWidth = true
		rightSpace.FlexibleHeight = true
		hbox.Controls = append(hbox.Controls, rightSpace)

		vbox.Controls = append(vbox.Controls, &hbox)

		// Add horizontal separator between rows (except after last row)
		if i < gl.rowCount-1 {
			rowSeparatorRow := Box{}
			rowSeparatorRow.Init(Horizontal)

			// Left spacing for row separator
			rowSepLeftSpace := &Tile{}
			rowSepLeftSpace.Init()
			rowSepLeftSpace.FlexibleWidth = true
			rowSepLeftSpace.FlexibleHeight = true
			rowSeparatorRow.Controls = append(rowSeparatorRow.Controls, rowSepLeftSpace)

			// Left border for row separator
			rowSepLeftBorder := &Tile{}
			rowSepLeftBorder.Init()
			rowSepLeftBorder.MinimumWidth = 1
			rowSepLeftBorder.MinimumHeight = 1
			rowSepLeftBorder.TextPresenter = &presenter.HorizontalSeparatorTextPresenter{}
			rowSeparatorRow.Controls = append(rowSeparatorRow.Controls, rowSepLeftBorder)

			// Content area separators - match the actual content width
			for j := 0; j < gl.columnCount; j++ {
				contentSeparator := &Tile{}
				contentSeparator.Init()
				contentSeparator.MinimumHeight = 1
				contentSeparator.FlexibleWidth = false

				// Set minimum width based on the content in this column
				if gl.controls[0][j] != nil {
					contentSeparator.MinimumWidth = gl.controls[0][j].MinimumWidth
					if gl.controls[0][j].FlexibleWidth {
						contentSeparator.FlexibleWidth = true
					}
				} else {
					contentSeparator.MinimumWidth = 6 // Default width for empty cells
				}

				contentSeparator.TextPresenter = &presenter.HorizontalSeparatorTextPresenter{}
				rowSeparatorRow.Controls = append(rowSeparatorRow.Controls, contentSeparator)

				// Add separator between columns (except after last column)
				if j < gl.columnCount-1 {
					separator := &Tile{}
					separator.Init()
					separator.MinimumWidth = 1
					separator.MinimumHeight = 1
					separator.TextPresenter = &presenter.HorizontalSeparatorTextPresenter{}
					rowSeparatorRow.Controls = append(rowSeparatorRow.Controls, separator)
				}
			}

			// Right border for row separator
			rowSepRightBorder := &Tile{}
			rowSepRightBorder.Init()
			rowSepRightBorder.MinimumWidth = 1
			rowSepRightBorder.MinimumHeight = 1
			rowSepRightBorder.TextPresenter = &presenter.HorizontalSeparatorTextPresenter{}
			rowSeparatorRow.Controls = append(rowSeparatorRow.Controls, rowSepRightBorder)

			// Right spacing for row separator
			rowSepRightSpace := &Tile{}
			rowSepRightSpace.Init()
			rowSepRightSpace.FlexibleWidth = true
			rowSepRightSpace.FlexibleHeight = true
			rowSeparatorRow.Controls = append(rowSeparatorRow.Controls, rowSepRightSpace)

			vbox.Controls = append(vbox.Controls, &rowSeparatorRow)
		}
	}

	// Bottom border row with proper spacing
	bottomBorderRow := Box{}
	bottomBorderRow.Init(Horizontal)

	// Left spacing for bottom border
	bottomLeftSpace := &Tile{}
	bottomLeftSpace.Init()
	bottomLeftSpace.FlexibleWidth = true
	bottomLeftSpace.FlexibleHeight = true
	bottomBorderRow.Controls = append(bottomBorderRow.Controls, bottomLeftSpace)

	// Left border for bottom
	bottomLeftBorder := &Tile{}
	bottomLeftBorder.Init()
	bottomLeftBorder.MinimumWidth = 1
	bottomLeftBorder.MinimumHeight = 1
	bottomLeftBorder.TextPresenter = &presenter.HorizontalSeparatorTextPresenter{}
	bottomBorderRow.Controls = append(bottomBorderRow.Controls, bottomLeftBorder)

	// Content area borders for bottom - match the actual content width
	for j := 0; j < gl.columnCount; j++ {
		contentBorder := &Tile{}
		contentBorder.Init()
		contentBorder.MinimumHeight = 1
		contentBorder.FlexibleWidth = false

		// Set minimum width based on the content in this column
		if gl.controls[0][j] != nil {
			contentBorder.MinimumWidth = gl.controls[0][j].MinimumWidth
			if gl.controls[0][j].FlexibleWidth {
				contentBorder.FlexibleWidth = true
			}
		} else {
			contentBorder.MinimumWidth = 6 // Default width for empty cells
		}

		contentBorder.TextPresenter = &presenter.HorizontalSeparatorTextPresenter{}
		bottomBorderRow.Controls = append(bottomBorderRow.Controls, contentBorder)

		// Add separator between columns (except after last column)
		if j < gl.columnCount-1 {
			separator := &Tile{}
			separator.Init()
			separator.MinimumWidth = 1
			separator.MinimumHeight = 1
			separator.TextPresenter = &presenter.HorizontalSeparatorTextPresenter{}
			bottomBorderRow.Controls = append(bottomBorderRow.Controls, separator)
		}
	}

	// Right border for bottom
	bottomRightBorder := &Tile{}
	bottomRightBorder.Init()
	bottomRightBorder.MinimumWidth = 1
	bottomRightBorder.MinimumHeight = 1
	bottomRightBorder.TextPresenter = &presenter.HorizontalSeparatorTextPresenter{}
	bottomBorderRow.Controls = append(bottomBorderRow.Controls, bottomRightBorder)

	// Right spacing for bottom border
	bottomRightSpace := &Tile{}
	bottomRightSpace.Init()
	bottomRightSpace.FlexibleWidth = true
	bottomRightSpace.FlexibleHeight = true
	bottomBorderRow.Controls = append(bottomBorderRow.Controls, bottomRightSpace)

	vbox.Controls = append(vbox.Controls, &bottomBorderRow)

	// Bottom spacing
	footer := &Tile{}
	footer.Init()
	footer.FlexibleWidth = true
	footer.FlexibleHeight = true
	vbox.Controls = append(vbox.Controls, footer)

	return &vbox
}
