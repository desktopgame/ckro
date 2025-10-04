package view

import (
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/gdamore/tcell/v2"
)

type FoldBlockView struct {
}

func (fv *FoldBlockView) Layout(ctx Context, textLayout *TextLayout, x, y, w, h int) {
	if ctx.FoldManager.IsFolded(ctx.Document, textLayout.Element) {
		childElement := textLayout.Children[0].Element
		childView := ctx.Resolver.Resolve(childElement)

		mw := textLayout.Children[0].MinimumWidth
		mh := textLayout.Children[0].MinimumHeight
		childView.Layout(ctx, textLayout.Children[0], 1+2, 1, mw, mh)
	} else {
		offsetY := 1
		for i := 0; i < len(textLayout.Children); i++ {
			childElement := textLayout.Children[i].Element
			childView := ctx.Resolver.Resolve(childElement)

			mw := textLayout.Children[i].MinimumWidth
			mh := textLayout.Children[i].MinimumHeight
			childView.Layout(ctx, textLayout.Children[i], 1, offsetY, mw, mh)
			offsetY += mh
		}
	}
	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
}

func (fv *FoldBlockView) Draw(ctx Context, textLayout *TextLayout, renderer Renderer) {
	for i := 1; i < textLayout.Width-1; i++ {
		renderer.SetContent(i, 0, '-', nil, tcell.StyleDefault)
		renderer.SetContent(i, textLayout.Height-1, '-', nil, tcell.StyleDefault)
	}
	for i := 1; i < textLayout.Height-1; i++ {
		renderer.SetContent(0, i, '|', nil, tcell.StyleDefault)
		renderer.SetContent(textLayout.Width-1, i, '|', nil, tcell.StyleDefault)
	}
	renderer.SetContent(0, 0, '+', nil, tcell.StyleDefault)
	renderer.SetContent(textLayout.Width-1, 0, '+', nil, tcell.StyleDefault)
	renderer.SetContent(0, textLayout.Height-1, '+', nil, tcell.StyleDefault)
	renderer.SetContent(textLayout.Width-1, textLayout.Height-1, '+', nil, tcell.StyleDefault)

	if ctx.FoldManager.IsFolded(ctx.Document, textLayout.Element) {
		renderer.SetContent(1, 1, '>', nil, tcell.StyleDefault)
		renderer.SetContent(1, 2, ' ', nil, tcell.StyleDefault)
		for _, child := range textLayout.Children {
			childView := ctx.Resolver.Resolve(child.Element)
			childView.Draw(ctx, child, renderer.Translate(child.RelativeX, child.RelativeY))
		}
	} else {
		for _, child := range textLayout.Children {
			childView := ctx.Resolver.Resolve(child.Element)
			childView.Draw(ctx, child, renderer.Translate(child.RelativeX, child.RelativeY))
		}
	}
}

func (fv *FoldBlockView) MinimumSize(ctx Context, e model.Element, width int, height int) *TextLayout {
	maxWidth := -1
	children := []*TextLayout{}

	var minimumHeight int
	if ctx.FoldManager.IsFolded(ctx.Document, e) {
		minimumHeight = 2

		childElement := e.GetElement(0)
		childView := ctx.Resolver.Resolve(childElement)
		child := childView.MinimumSize(ctx, childElement, width-4, 1)

		if child.MinimumWidth+4 > width {
			childElement = &model.PlainElement{
				Range: childElement.GetRange(0),
			}
			childView = &PlainTextView{}
			child = childView.MinimumSize(ctx, childElement, width-4, 9999)
			minimumHeight += child.MinimumHeight
		} else {
			minimumHeight++
		}

		children = append(children, child)

		maxWidth = child.MinimumWidth + 2
	} else {
		minimumHeight = 2

		for i := 0; i < e.GetElementCount(); i++ {
			childElement := e.GetElement(i)
			childView := ctx.Resolver.Resolve(childElement)

			child := childView.MinimumSize(ctx, childElement, width-2, 1)

			if child.MinimumWidth+2 > width {
				childElement = &model.PlainElement{
					Range: childElement.GetRange(0),
				}
				childView = &PlainTextView{}
				child = childView.MinimumSize(ctx, childElement, width-2, 9999)
				minimumHeight += child.MinimumHeight
			} else {
				minimumHeight++
			}

			children = append(children, child)

			if child.MinimumWidth > maxWidth {
				maxWidth = child.MinimumWidth
			}
		}
	}
	return &TextLayout{
		Element:       e,
		MinimumWidth:  width,
		MinimumHeight: minimumHeight,
		Children:      children,
	}
}

func (fv *FoldBlockView) ViewLengthTable(ctx Context, textLayout *TextLayout) ([]int, int) {
	var table []int
	total := 0
	for i := 0; i < len(textLayout.Children); i++ {
		childElement := textLayout.Children[i].Element
		childView := ctx.Resolver.Resolve(childElement)

		l := childView.MoveLength(ctx, textLayout.Children[i])
		table = append(table, l)
		total += l
	}
	return table, total
}

func (fv *FoldBlockView) findTableIndex(table []int, viewLocalPos int) (Row int, Column int) {
	n := 0
	index := -1
	col := -1
	for i, l := range table {
		start := n
		if viewLocalPos >= start && viewLocalPos < n+l {
			index = i
			col = viewLocalPos - start
			break
		}
		n += l
	}
	return index, col
}

func (fv *FoldBlockView) sumTableValue(table []int, index int) int {
	v := 0
	for i := 0; i <= index; i++ {
		v += table[i]
	}
	return v
}

func (fv *FoldBlockView) MoveLength(ctx Context, textLayout *TextLayout) int {
	if ctx.FoldManager.IsFolded(ctx.Document, textLayout.Element) {
		childElement := textLayout.Children[0].Element
		childView := ctx.Resolver.Resolve(childElement)
		return childView.MoveLength(ctx, textLayout.Children[0])
	} else {
		_, ttl := fv.ViewLengthTable(ctx, textLayout)
		return ttl
	}
}

func (fv *FoldBlockView) MoveUp(ctx Context, textLayout *TextLayout, viewLocalPos int) int {
	table, _ := fv.ViewLengthTable(ctx, textLayout)
	index, _ := fv.findTableIndex(table, viewLocalPos)
	if index <= 0 {
		return -1
	}
	if index == 1 {
		return 0
	}
	return fv.sumTableValue(table, index-2)
}

func (fv *FoldBlockView) MoveDown(ctx Context, textLayout *TextLayout, viewLocalPos int) int {
	table, _ := fv.ViewLengthTable(ctx, textLayout)
	index, _ := fv.findTableIndex(table, viewLocalPos)
	if index == len(table)-1 {
		return -1
	}
	return fv.sumTableValue(table, index)
}

func (fv *FoldBlockView) MoveLeft(ctx Context, textLayout *TextLayout, viewLocalPos int) int {
	if viewLocalPos <= 0 {
		return -1
	}
	return viewLocalPos - 1
}

func (fv *FoldBlockView) MoveRight(ctx Context, textLayout *TextLayout, viewLocalPos int) int {
	if viewLocalPos >= fv.MoveLength(ctx, textLayout)-1 {
		return -1
	}
	return viewLocalPos + 1
}

func (fv *FoldBlockView) ConvertPos(ctx Context, textLayout *TextLayout, viewLocalPos int) (ViewLocalX int, ViewLocalY int) {
	e := textLayout.Element

	if ctx.FoldManager.IsFolded(ctx.Document, e) {
		childElement := e.GetElement(0)
		childView := ctx.Resolver.Resolve(childElement)
		lx, ly := childView.ConvertPos(ctx, textLayout.Children[0], viewLocalPos)
		return 1 + lx + 2, 1 + ly
	}
	table, _ := fv.ViewLengthTable(ctx, textLayout)
	index, col := fv.findTableIndex(table, viewLocalPos)
	if viewLocalPos == fv.sumTableValue(table, len(table)-1) {
		child := textLayout.Children[len(textLayout.Children)-1]
		childView := ctx.Resolver.Resolve(child.Element)
		childLen := childView.MoveLength(ctx, child)
		lx, ly := childView.ConvertPos(ctx, child, childLen-1)
		return lx + 1, ly + 1
	}

	child := textLayout.Children[index]
	v := ctx.Resolver.Resolve(child.Element)
	lx, ly := v.ConvertPos(ctx, child, col)
	return lx + 1, ly + 1
}

func (fv *FoldBlockView) ConvertModel(ctx Context, textLayout *TextLayout, viewLocalPos int) CharacterReference {
	e := textLayout.Element
	if ctx.FoldManager.IsFolded(ctx.Document, textLayout.Element) {
		childElement := e.GetElement(0)
		childView := ctx.Resolver.Resolve(childElement)
		return childView.ConvertModel(ctx, textLayout.Children[0], viewLocalPos)
	} else {

		table, _ := fv.ViewLengthTable(ctx, textLayout)
		index, col := fv.findTableIndex(table, viewLocalPos)
		if viewLocalPos == fv.sumTableValue(table, len(table)-1) {
			child := textLayout.Children[len(textLayout.Children)-1]
			childView := ctx.Resolver.Resolve(child.Element)
			childLen := childView.MoveLength(ctx, child)
			return childView.ConvertModel(ctx, child, childLen-1)
		}
		child := textLayout.Children[index]

		v := ctx.Resolver.Resolve(child.Element)
		return v.ConvertModel(ctx, child, col)
	}
}

func (fv *FoldBlockView) ConvertViewLocalPos(ctx Context, textLayout *TextLayout, bytePos model.Position) int {
	e := textLayout.Element
	if ctx.FoldManager.IsFolded(ctx.Document, textLayout.Element) {
		r := e.GetRange(0)
		if bytePos.Row == r.StartPosition.Row {
			return 0
		} else if bytePos.Row == r.EndPosition.Row {
			return fv.MoveLength(ctx, textLayout) - 1
		}

		childElement := textLayout.Children[0].Element
		childView := ctx.Resolver.Resolve(childElement)
		return childView.ConvertViewLocalPos(ctx, textLayout.Children[0], bytePos)
	} else {
		viewOffset := 0
		for i := 0; i < len(textLayout.Children); i++ {
			child := textLayout.Children[i]
			r := child.Element.GetRange(0)
			st := r.StartPosition
			ed := r.EndPosition
			childView := ctx.Resolver.Resolve(child.Element)

			if bytePos.Row >= st.Row && bytePos.Row <= ed.Row {

				if st.Row == ed.Row && st.Column == ed.Column {
					if bytePos.Row == st.Row && bytePos.Column == st.Column {
						return viewOffset + childView.ConvertViewLocalPos(ctx, child, bytePos)
					}
				}
				if bytePos.Column >= st.Column && (bytePos.Column <= ed.Column || ed.Row > st.Row) {
					return viewOffset + childView.ConvertViewLocalPos(ctx, child, bytePos)
				}
			}
			viewOffset += childView.MoveLength(ctx, child)
		}

		return fv.MoveLength(ctx, textLayout) - 1
	}
}
