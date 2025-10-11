package litemark

import (
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/desktopgame/ckro/internal/tui/view"
	"github.com/gdamore/tcell/v2"
)

type TableView struct {
}

func (tv *TableView) Layout(ctx view.Context, textLayout *view.TextLayout, x, y, w, h int) {
	var heightTable []int
	for i := 0; i < len(textLayout.Children); i++ {
		row := textLayout.Children[i].Element
		maxHeight := -1
		for j := 0; j < row.GetElementCount(); j++ {
			mh := textLayout.Children[i].MinimumHeight

			if mh > maxHeight {
				maxHeight = mh
			}
		}
		heightTable = append(heightTable, maxHeight)
	}

	var widthTable []int
	for j := 0; j < len(textLayout.Children[0].Children); j++ {
		maxWidth := -1
		for i := 0; i < len(textLayout.Children); i++ {
			mw := textLayout.Children[i].Children[j].MinimumWidth

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
	yy := 1
	for i, h := range heightTable {
		row := textLayout.Children[i].Element
		rowView := ctx.Resolver.Resolve(row)

		textLayout.Children[i].WidthTable = widthTable
		rowView.Layout(ctx, textLayout.Children[i], 1, yy, w, h)
		if i == 0 {
			yy++
		}
		totalHeight += h
		yy += h
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
		renderer.SetContent(x, tableHeight-1, '─', nil, tcell.StyleDefault)
	}
	renderer.SetContent(tableWidth-1, tableHeight-1, '┘', nil, tcell.StyleDefault)

	// Draw left and right borders
	for y := 1; y < tableHeight-1; y++ {
		renderer.SetContent(0, y, '│', nil, tcell.StyleDefault)
		renderer.SetContent(tableWidth-1, y, '│', nil, tcell.StyleDefault)
	}

	// Draw table content and header separator
	y := 1 // Start after top border
	for i, child := range textLayout.Children {
		childElement := textLayout.Element.GetElement(i)
		childView := ctx.Resolver.Resolve(childElement)
		childView.Draw(ctx, child, renderer.Translate(child.RelativeX, child.RelativeY)) // Offset by left border
		y += child.Height

		// Draw horizontal separator after header
		if _, isHeader := childElement.(*TableHeaderElement); isHeader {
			// Draw header separator line
			renderer.SetContent(0, y, '├', nil, tcell.StyleDefault)
			for x := 1; x < tableWidth-1; x++ {
				renderer.SetContent(x, y, '─', nil, tcell.StyleDefault)
			}
			renderer.SetContent(tableWidth-1, y, '┤', nil, tcell.StyleDefault)
			y++ // Move to next line after separator
		}
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

	var heightTable []int
	for i := 0; i < e.GetElementCount(); i++ {
		row := e.GetElement(i)
		maxHeight := -1
		for j := 0; j < row.GetElementCount(); j++ {
			mh := children[i].Children[j].MinimumHeight

			if mh > maxHeight {
				maxHeight = mh
			}
		}
		heightTable = append(heightTable, maxHeight)
	}

	var widthTable []int
	for j := 0; j < e.GetElement(0).GetElementCount(); j++ {
		maxWidth := -1
		for i := 0; i < e.GetElementCount(); i++ {
			mw := children[i].Children[j].MinimumWidth

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

	return &view.TextLayout{
		Element:       e,
		MinimumWidth:  totalWidth + (e.GetElement(0).GetElementCount() + 1),
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
