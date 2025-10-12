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
			cellElement := tableElement.Children[i*tableElement.Columns+j]
			cellView := ctx.Resolver.Resolve(cellElement)

			cellView.Layout(ctx, textLayout.Children[i*tableElement.Columns+j], offsetX, yy, widthTable[j], hh)
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

func (tv *TableView) MoveLength(ctx view.Context, textLayout *view.TextLayout) int {
	return 1
}

func (tv *TableView) MoveUp(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	return -1
}

func (tv *TableView) MoveDown(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	return -1
}

func (tv *TableView) MoveLeft(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	return -1
}

func (tv *TableView) MoveRight(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	return -1
}

func (tv *TableView) ConvertPos(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (ViewLocalX int, ViewLocalY int) {
	return 0, 0
}

func (tv *TableView) ConvertModel(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) view.CharacterReference {
	e := textLayout.Element
	r := e.GetRange(0)
	st := r.StartPosition
	return view.CharacterReference{
		StartPosition: model.Position{
			Row:    st.Row,
			Column: st.Column,
		},
		Bytes: 0,
	}
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

func (tv *TableView) ConvertViewLocalPos(ctx view.Context, textLayout *view.TextLayout, bytePos model.Position) int {
	return 0
}
