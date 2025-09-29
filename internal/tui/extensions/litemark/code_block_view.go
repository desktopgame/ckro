package litemark

import (
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/desktopgame/ckro/internal/tui/view"
	"github.com/gdamore/tcell/v2"
)

type CodeBlockView struct {
}

func (c *CodeBlockView) Layout(ctx view.Context, textLayout *view.TextLayout, x, y, w, h int) {
	offsetY := 1
	for i := 0; i < len(textLayout.Children); i++ {
		childElement := textLayout.Children[i].Element
		childView := ctx.Resolver.Resolve(childElement)

		mw := textLayout.Children[i].MinimumWidth
		childView.Layout(ctx, textLayout.Children[i], 1, offsetY, mw, 1)
		offsetY++
	}
	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
}

func (c *CodeBlockView) Draw(ctx view.Context, textLayout *view.TextLayout, renderer view.Renderer) {
	for i := 1; i < textLayout.Width-1; i++ {
		renderer.SetContent(i, 0, '-', nil, tcell.StyleDefault)
		renderer.SetContent(i, textLayout.Height-1, '-', nil, tcell.StyleDefault)
	}
	for i := 1; i < textLayout.Height-1; i++ {
		renderer.SetContent(0, i, '|', nil, tcell.StyleDefault)
		renderer.SetContent(textLayout.Width-1, i, '|', nil, tcell.StyleDefault)
	}
	renderer.SetContent(0, 0, '*', nil, tcell.StyleDefault)
	renderer.SetContent(textLayout.Width-1, 0, '*', nil, tcell.StyleDefault)
	renderer.SetContent(0, textLayout.Height-1, '*', nil, tcell.StyleDefault)
	renderer.SetContent(textLayout.Width-1, textLayout.Height-1, '*', nil, tcell.StyleDefault)
	for _, child := range textLayout.Children {
		childView := ctx.Resolver.Resolve(child.Element)
		childView.Draw(ctx, child, renderer.Translate(child.RelativeX, child.RelativeY))
	}
}

func (c *CodeBlockView) MinimumSize(ctx view.Context, e model.Element, width int, height int) *view.TextLayout {
	maxWidth := -1
	children := []*view.TextLayout{}
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := ctx.Resolver.Resolve(childElement)

		child := childView.MinimumSize(ctx, childElement, width, 1)
		children = append(children, child)

		if child.MinimumWidth > maxWidth {
			maxWidth = child.MinimumWidth
		}
	}
	return &view.TextLayout{
		Element:       e,
		MinimumWidth:  maxWidth + 2,
		MinimumHeight: e.GetElementCount() + 2,
		Children:      children,
	}
}

func (c *CodeBlockView) ViewLengthTable(ctx view.Context, e model.Element) ([]int, int) {
	var table []int
	total := 0
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := ctx.Resolver.Resolve(childElement)

		l := childView.MoveLength(ctx, childElement)
		table = append(table, l)
		total += l
	}
	return table, total
}

func (c *CodeBlockView) findTableIndex(table []int, viewLocalPos int) (Row int, Column int) {
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

func (c *CodeBlockView) sumTableValue(table []int, index int) int {
	v := 0
	for i := 0; i <= index; i++ {
		v += table[i]
	}
	return v
}

func (c *CodeBlockView) MoveLength(ctx view.Context, e model.Element) int {
	_, ttl := c.ViewLengthTable(ctx, e)
	return ttl
}

func (c *CodeBlockView) MoveUp(ctx view.Context, e model.Element, viewLocalPos int) int {
	table, _ := c.ViewLengthTable(ctx, e)
	index, _ := c.findTableIndex(table, viewLocalPos)
	if index <= 0 {
		return -1
	}
	if index == 1 {
		return 0
	}
	return c.sumTableValue(table, index-2)
}

func (c *CodeBlockView) MoveDown(ctx view.Context, e model.Element, viewLocalPos int) int {
	table, _ := c.ViewLengthTable(ctx, e)
	index, _ := c.findTableIndex(table, viewLocalPos)
	if index == len(table)-1 {
		return -1
	}
	return c.sumTableValue(table, index)
}

func (c *CodeBlockView) MoveLeft(ctx view.Context, e model.Element, viewLocalPos int) int {
	if viewLocalPos <= 0 {
		return -1
	}
	return viewLocalPos - 1
}

func (c *CodeBlockView) MoveRight(ctx view.Context, e model.Element, viewLocalPos int) int {
	if viewLocalPos >= c.MoveLength(ctx, e) {
		return -1
	}
	return viewLocalPos + 1
}

func (c *CodeBlockView) ConvertPos(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (ViewLocalX int, ViewLocalY int) {
	e := textLayout.Element
	table, _ := c.ViewLengthTable(ctx, e)
	index, col := c.findTableIndex(table, viewLocalPos)

	child := textLayout.Children[index]
	v := ctx.Resolver.Resolve(child.Element)
	lx, _ := v.ConvertPos(ctx, child, col)
	return lx + 1, index + 1
}

func (c *CodeBlockView) ConvertModel(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) view.CharacterReference {
	table, _ := c.ViewLengthTable(ctx, textLayout.Element)
	index, col := c.findTableIndex(table, viewLocalPos)
	if viewLocalPos == c.sumTableValue(table, len(table)-1) {
		child := textLayout.Children[len(textLayout.Children)-1]
		childView := ctx.Resolver.Resolve(child.Element)
		childLen := childView.MoveLength(ctx, child.Element)
		return childView.ConvertModel(ctx, child, childLen-1)
	}
	child := textLayout.Children[index]

	v := ctx.Resolver.Resolve(child.Element)
	return v.ConvertModel(ctx, child, col)
}

func (c *CodeBlockView) ConvertViewLocalPos(ctx view.Context, textLayout *view.TextLayout, bytePos model.Position) int {
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
		viewOffset += childView.MoveLength(ctx, child.Element)
	}

	return c.MoveLength(ctx, textLayout.Element) - 1
}
