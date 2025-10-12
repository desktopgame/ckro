package litemark

import (
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/desktopgame/ckro/internal/tui/view"
	"github.com/gdamore/tcell/v2"
)

type TableView struct {
}

func (tv *TableView) calculateTable(e model.Element, children []*view.TextLayout) (WidthTable []int, HeightTable []int, TotalWidth int, TotalHeight int) {
	tableElement := e.(*TableElement)
	rowCount := tableElement.GetElementCount() / tableElement.Columns
	columnCount := tableElement.Columns

	var heightTable []int
	for i := 0; i < rowCount; i++ {
		maxHeight := -1
		for j := 0; j < columnCount; j++ {
			mh := children[i*columnCount+j].MinimumHeight

			if mh > maxHeight {
				maxHeight = mh
			}
		}
		heightTable = append(heightTable, maxHeight)
	}

	var widthTable []int
	for j := 0; j < columnCount; j++ {
		maxWidth := -1
		for i := 0; i < rowCount; i++ {
			mw := children[i*columnCount+j].MinimumWidth

			if mw > maxWidth {
				maxWidth = mw
			}
		}
		widthTable = append(widthTable, maxWidth)
	}

	totalWidth := 0
	for _, w := range widthTable {
		totalWidth += w
	}

	totalHeight := 0
	for _, h := range heightTable {
		totalHeight += h
	}

	return widthTable, heightTable, totalWidth, totalHeight
}

func (tv *TableView) Layout(ctx view.Context, textLayout *view.TextLayout, x, y, w, h int) {
	tableElement := textLayout.Element.(*TableElement)
	widthTable, heightTable, _, _ := tv.calculateTable(textLayout.Element, textLayout.Children)

	totalHeight := 0
	yy := 1
	for i, hh := range heightTable {
		offsetX := 1
		for j := 0; j < tableElement.Columns; j++ {
			cellChild := textLayout.Children[i*tableElement.Columns+j]
			cellMw := cellChild.MinimumWidth
			cellMh := cellChild.MinimumHeight

			cellElement := tableElement.Children[i*tableElement.Columns+j]
			cellView := ctx.Resolver.Resolve(cellElement)
			cellAlign := tableElement.Aligns[j]

			switch cellAlign {
			case TABLE_ALIGN_LEFT:
				cellView.Layout(ctx, textLayout.Children[i*tableElement.Columns+j], offsetX, yy, cellMw, cellMh)
			case TABLE_ALIGN_CENTER:
				pad := widthTable[j] - cellMw
				if pad > 0 {
					cellView.Layout(ctx, textLayout.Children[i*tableElement.Columns+j], offsetX+(pad/2), yy, cellMw, cellMh)
				} else {
					cellView.Layout(ctx, textLayout.Children[i*tableElement.Columns+j], offsetX, yy, cellMw, cellMh)
				}
			case TABLE_ALIGN_RIGHT:
				pad := widthTable[j] - cellMw
				if pad > 0 {
					cellView.Layout(ctx, textLayout.Children[i*tableElement.Columns+j], offsetX+pad, yy, cellMw, cellMh)
				} else {
					cellView.Layout(ctx, textLayout.Children[i*tableElement.Columns+j], offsetX, yy, cellMw, cellMh)
				}
			}
			offsetX += widthTable[j] + 1
		}
		if i == 0 {
			yy++ // header
		}
		totalHeight += hh
		yy += hh
	}

	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
}

func (tv *TableView) Draw(ctx view.Context, textLayout *view.TextLayout, renderer view.Renderer) {
	tableWidth := textLayout.Width
	tableHeight := textLayout.Height

	// Draw top border
	renderer.SetContent(0, 0, '┌', nil, tcell.StyleDefault)
	for x := 1; x < tableWidth-1; x++ {
		renderer.SetContent(x, 0, '─', nil, tcell.StyleDefault)
	}
	renderer.SetContent(tableWidth-1, 0, '┐', nil, tcell.StyleDefault)

	// Draw bottom border
	renderer.SetContent(0, tableHeight-1, '└', nil, tcell.StyleDefault)
	for x := 1; x < tableWidth-1; x++ {
		renderer.SetContent(x, 2, '─', nil, tcell.StyleDefault)
		renderer.SetContent(x, tableHeight-1, '─', nil, tcell.StyleDefault)
	}
	renderer.SetContent(tableWidth-1, tableHeight-1, '┘', nil, tcell.StyleDefault)

	// Draw left and right borders
	for y := 1; y < tableHeight-1; y++ {
		renderer.SetContent(0, y, '│', nil, tcell.StyleDefault)
		renderer.SetContent(tableWidth-1, y, '│', nil, tcell.StyleDefault)
	}

	widthTable, heightTable, _, _ := tv.calculateTable(textLayout.Element, textLayout.Children)
	borderY := 1
	for i, hh := range heightTable {
		borderX := 1
		for j, ww := range widthTable {
			if j == len(widthTable)-1 {
				continue
			}
			borderX += ww
			renderer.SetContent(borderX, borderY, '│', nil, tcell.StyleDefault)
			borderX++
		}
		if i == 0 {
			borderY++
		}
		borderY += hh
	}

	// Draw table content and header separator
	for i, child := range textLayout.Children {
		childElement := textLayout.Element.GetElement(i)
		childView := ctx.Resolver.Resolve(childElement)
		childView.Draw(ctx, child, renderer.Translate(child.RelativeX, child.RelativeY)) // Offset by left border
	}
}

func (tv *TableView) Measure(ctx view.Context, e model.Element, width int, height int) *view.TextLayout {
	children := []*view.TextLayout{}
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := ctx.Resolver.Resolve(childElement)
		child := childView.Measure(ctx, childElement, width, height)
		children = append(children, child)
	}

	widthTable, _, totalWidth, totalHeight := tv.calculateTable(e, children)

	return &view.TextLayout{
		Element:       e,
		MinimumWidth:  totalWidth + len(widthTable) + 1, // borders
		MinimumHeight: totalHeight + 3,
		Children:      children,
	}
}

func (tv *TableView) moveTable(ctx view.Context, textLayout *view.TextLayout) (Table [][]int, Total int) {
	totalMoves := 0
	moveTable := [][]int{}
	tableElement := textLayout.Element.(*TableElement)
	rowCount := len(textLayout.Children) / tableElement.Columns

	for i := 0; i < rowCount; i++ {
		moveLine := []int{}
		for j := 0; j < tableElement.Columns; j++ {
			index := i*tableElement.Columns + j
			child := textLayout.Children[index]
			childElement := child.Element
			childView := ctx.Resolver.Resolve(childElement)

			moves := childView.MoveLength(ctx, child)
			moveLine = append(moveLine, moves)
			totalMoves += moves
		}
		moveTable = append(moveTable, moveLine)
	}
	return moveTable, totalMoves
}

func (tv *TableView) moveGridPos(table [][]int, viewLocalPos int) (int, int, int) {
	rowCount := len(table)
	n := 0
	for i := 0; i < rowCount; i++ {
		for j := 0; j < len(table[i]); j++ {
			start := n
			l := table[i][j]
			if viewLocalPos >= start && viewLocalPos < start+l {
				return i, j, viewLocalPos - start
			}
			n += l
		}
	}
	return -1, -1, -1
}

func (tv *TableView) moveByGridPos(table [][]int, row int, column int) (int, int) {
	total := 0
	for i := 0; i < len(table); i++ {
		for j := 0; j < len(table[i]); j++ {
			l := table[i][j]
			if i == row && j == column {
				return total, l
			}
			total += l
		}
	}
	return -1, -1
}

func (tv *TableView) moveOffset(table [][]int, row int, column int) int {
	moves := 0
	for i := 0; i < row; i++ {
		for j := 0; j < column; j++ {
			moves += table[i][j]
		}
	}
	for j := 0; j < column; j++ {
		moves += table[row][j]
	}
	return moves
}

func (tv *TableView) MoveLength(ctx view.Context, textLayout *view.TextLayout) int {
	_, ttl := tv.moveTable(ctx, textLayout)
	return ttl
}

func (tv *TableView) MoveUp(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	table, _ := tv.moveTable(ctx, textLayout)
	row, col, _ := tv.moveGridPos(table, viewLocalPos)
	if row == 0 {
		return -1
	}
	vl, _ := tv.moveByGridPos(table, row-1, col)
	return vl
}

func (tv *TableView) MoveDown(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	table, _ := tv.moveTable(ctx, textLayout)
	row, col, _ := tv.moveGridPos(table, viewLocalPos)
	if row == len(table)-1 {
		return -1
	}
	vl, _ := tv.moveByGridPos(table, row+1, col)
	return vl
}

func (tv *TableView) MoveLeft(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	table, _ := tv.moveTable(ctx, textLayout)
	row, col, offset := tv.moveGridPos(table, viewLocalPos)
	if offset == 0 {
		if col == 0 {
			if row == 0 {
				return -1
			}
			vl, l := tv.moveByGridPos(table, row-1, len(table[0])-1)
			return vl + (l - 1)
		}
		vl, l := tv.moveByGridPos(table, row, col-1)
		return vl + (l - 1)
	}
	return viewLocalPos - 1
}

func (tv *TableView) MoveRight(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	table, _ := tv.moveTable(ctx, textLayout)
	row, col, offset := tv.moveGridPos(table, viewLocalPos)
	_, l := tv.moveByGridPos(table, row, col)
	if offset == l-1 {
		if col == len(table[0])-1 {
			if row == len(table)-1 {
				return -1
			}
			vl, _ := tv.moveByGridPos(table, row+1, 0)
			return vl
		}
		vl, _ := tv.moveByGridPos(table, row, col+1)
		return vl
	}
	return viewLocalPos + 1
}

func (tv *TableView) ConvertPos(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (ViewLocalX int, ViewLocalY int) {
	tableElement := textLayout.Element.(*TableElement)
	wt, _, _, _ := tv.calculateTable(textLayout.Element, textLayout.Children)
	table, _ := tv.moveTable(ctx, textLayout)
	row, col, offset := tv.moveGridPos(table, viewLocalPos)
	child := textLayout.Children[row*tableElement.Columns+col]
	childElement := child.Element
	childView := ctx.Resolver.Resolve(childElement)
	vlx, vly := childView.ConvertPos(ctx, child, offset)
	if row >= 1 {
		row++
	}
	offsetX := 1
	for j := 0; j < col; j++ {
		offsetX += wt[j] + 1
	}

	cw := child.Width
	switch tableElement.Aligns[col] {
	case TABLE_ALIGN_LEFT:
	case TABLE_ALIGN_CENTER:
		pad := wt[col] - cw
		if pad > 0 {
			offsetX += (pad / 2)
		}
	case TABLE_ALIGN_RIGHT:
		pad := wt[col] - cw
		if pad > 0 {
			offsetX += pad
		}
	}
	return offsetX + vlx, row + 1 + vly
}

func (tv *TableView) ConvertModel(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) view.CharacterReference {
	tableElement := textLayout.Element.(*TableElement)
	table, _ := tv.moveTable(ctx, textLayout)
	row, col, offset := tv.moveGridPos(table, viewLocalPos)
	child := textLayout.Children[row*tableElement.Columns+col]
	childElement := child.Element
	childView := ctx.Resolver.Resolve(childElement)
	return childView.ConvertModel(ctx, child, offset)
}

func (tv *TableView) ConvertViewLocalPos(ctx view.Context, textLayout *view.TextLayout, bytePos model.Position) int {
	viewOffset := 0
	for i := 0; i < len(textLayout.Children); i++ {
		child := textLayout.Children[i]
		r := child.Element.GetRange(0)
		st := r.StartPosition
		ed := r.EndPosition
		childView := ctx.Resolver.Resolve(child.Element)

		if bytePos.Row >= st.Row && bytePos.Row <= ed.Row {
			if bytePos.Column >= st.Column && (bytePos.Column <= ed.Column || ed.Row > st.Row) {
				return viewOffset + childView.ConvertViewLocalPos(ctx, child, bytePos)
			}
		}
		viewOffset += childView.MoveLength(ctx, child)
	}

	return tv.MoveLength(ctx, textLayout) - 1
}

func (tv *TableView) FindFoldElementAt(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (model.Element, int, bool) {
	return nil, 0, false
}

func (tv *TableView) ShouldBeforeInsertionNewLineOnLineBegin(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) bool {
	return false
}

func (tv *TableView) ShouldRemoveWithLine(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (int, bool) {
	return -1, false
}

func (tv *TableView) ShouldRemoveWithSpecifiedColumnAfter(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (model.Position, bool) {
	if viewLocalPos == 0 {
		return model.Position{}, false
	}
	table, _ := tv.moveTable(ctx, textLayout)
	ttl := 0
	for i, v := range table {
		if viewLocalPos == ttl {
			index := (i-1)*len(table[0]) + (len(table[0]) - 1)
			return textLayout.Children[index].Element.GetRange(0).EndPosition, true
		}
		for _, vv := range v {
			ttl += vv
		}
	}
	return model.Position{}, false
}

func (tv *TableView) ShouldRemoveWithSpecifiedRangeLines(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (model.Range, bool) {
	return model.Range{}, false
}

func (tv *TableView) ShouldRemoveWithSpecifiedRangeColumns(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (model.Range, bool) {
	return model.Range{}, false
}

func (tv *TableView) ShouldRemoveLastCharacter(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (model.Element, bool) {
	return nil, false
}
