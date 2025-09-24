package view

import (
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/gdamore/tcell/v2"
)

type TableContainerView struct {
}

func (t *TableContainerView) Layout(textViewResolver TextViewResolver, e model.Element, width int) *TextLayout {
	if table, ok := e.(*model.TableContainerElement); ok {
		children := []*TextLayout{}

		for i := 0; i < table.GetElementCount(); i++ {
			row := &TextLayout{}

			for j := 0; j < table.GetElement(i).GetElementCount(); j++ {
				elem := table.GetElement(i).GetElement(j)
				view := textViewResolver.Resolve(elem)

				row.Children = append(row.Children, view.Layout(textViewResolver, elem, width))
			}

			children = append(children, row)
		}

		return &TextLayout{
			Element:  e,
			Children: children,
		}
	}
	panic("TableContainerView requires TableContainerElement")
}

func (t *TableContainerView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer, x int, y int, localViewLine int) {
	// x := 0
	// y :=

	if _, ok := textLayout.Element.(*model.TableContainerElement); ok {
		// yBorders := table.GetElementCount() + 1
		heightTable := t.HeightTable(textViewResolver, textLayout)
		seek := localViewLine
		drawX := x + 1
		drawY := y + 1
		for i, h := range heightTable {

			if seek < h {
				widthTable := t.WidthTable(textViewResolver, textLayout, i)

				columns := len(textLayout.Children[i].Children)
				for j := 0; j < columns; j++ {
					cellElem := textLayout.Children[i].Children[j].Element
					cellView := textViewResolver.Resolve(cellElem)

					cellView.Draw(textViewResolver, textLayout.Children[i].Children[j], renderer, drawX, drawY, seek)
					drawX += widthTable[j] + 1
				}
				break
			}
			drawY += h
			seek -= h
		}

		height := t.Height(textViewResolver, textLayout)
		for i := 0; i < height; i++ {
			if i != localViewLine {
				continue
			}
			width := t.Width(textViewResolver, textLayout, i)
			if i == 0 || i == height-1 {
				for j := 0; j < width; j++ {
					if j == 0 || j == width-1 {
						renderer.SetContent(j, i, '*', nil, tcell.StyleDefault)
					} else {
						renderer.SetContent(j, i, '-', nil, tcell.StyleDefault)
					}
				}
			} else {
				renderer.SetContent(0, i, '|', nil, tcell.StyleDefault)
				renderer.SetContent(width-1, i, '|', nil, tcell.StyleDefault)
			}
		}
		return
	}
	panic("TableContainerView requires TableContainerElement")
}

func (t *TableContainerView) Width(textViewResolver TextViewResolver, textLayout *TextLayout, row int) int {
	if table, ok := textLayout.Element.(*model.TableContainerElement); ok {
		columns := table.GetElement(0).GetElementCount()
		xBorders := columns + 1
		widthTable := t.WidthTable(textViewResolver, textLayout, row)
		width := 0
		for _, w := range widthTable {
			width += w
		}
		return width + xBorders
	}
	panic("TableContainerView requires TableContainerElement")
}

func (t *TableContainerView) WidthTable(textViewResolver TextViewResolver, textLayout *TextLayout, row int) []int {
	if table, ok := textLayout.Element.(*model.TableContainerElement); ok {
		columns := table.GetElement(0).GetElementCount()
		var widthTable []int

		for j := 0; j < columns; j++ {
			cMaxWidth := -1
			for i := 0; i < table.GetElementCount(); i++ {
				elem := table.GetElement(i).GetElement(j)
				view := textViewResolver.Resolve(elem)

				rows := view.Height(textViewResolver, textLayout.Children[i].Children[j])
				cvMaxWidth := -1
				for r := 0; r < rows; r++ {
					cvWidth := view.Width(textViewResolver, textLayout.Children[i].Children[j], r)
					if cvWidth > cvMaxWidth {
						cvMaxWidth = cvWidth
					}
				}
				if cvMaxWidth > cMaxWidth {
					cMaxWidth = cvMaxWidth
				}
			}
			widthTable = append(widthTable, cMaxWidth)
		}
		return widthTable
	}
	panic("TableContainerView requires TableContainerElement")
}

func (t *TableContainerView) Height(textViewResolver TextViewResolver, textLayout *TextLayout) int {
	if table, ok := textLayout.Element.(*model.TableContainerElement); ok {
		yBorders := table.GetElementCount() + 1
		heightTable := t.HeightTable(textViewResolver, textLayout)
		height := 0
		for _, h := range heightTable {
			height += h
		}
		return height + yBorders
	}
	panic("TableContainerView requires TableContainerElement")
}

func (t *TableContainerView) HeightTable(textViewResolver TextViewResolver, textLayout *TextLayout) []int {
	if table, ok := textLayout.Element.(*model.TableContainerElement); ok {
		var heightTable []int
		for i := 0; i < table.GetElementCount(); i++ {
			rowElement := table.GetElement(i)
			maxHeight := -1

			for j := 0; j < rowElement.GetElementCount(); j++ {
				columnElement := rowElement.GetElement(j)
				columnView := textViewResolver.Resolve(columnElement)

				cvHeight := columnView.Height(textViewResolver, textLayout.Children[i].Children[j])
				if cvHeight > maxHeight {
					maxHeight = cvHeight
				}
			}

			heightTable = append(heightTable, maxHeight)
		}
		return heightTable
	}
	panic("TableContainerView requires TableContainerElement")
}
