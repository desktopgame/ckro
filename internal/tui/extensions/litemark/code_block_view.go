package litemark

import (
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/desktopgame/ckro/internal/tui/view"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

type CodeBlockView struct {
}

func (c *CodeBlockView) Layout(ctx view.Context, textLayout *view.TextLayout, x, y, w, h int) {
	cbe := textLayout.Element.(*CodeBlockElement)
	additionalOffset := 0
	if len(cbe.Lang) > 0 {
		additionalOffset = 2
	}

	offsetY := 1 + additionalOffset
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
	cbe := textLayout.Element.(*CodeBlockElement)

	subLines := 0
	if len(cbe.Lang) > 0 {
		w := runewidth.StringWidth(cbe.Lang)
		renderer.SetContent(0, 0, '*', nil, tcell.StyleDefault)
		for i := 1; i < w+2; i++ {
			renderer.SetContent(i, 0, '-', nil, tcell.StyleDefault)
		}
		renderer.SetContent(w+2, 0, '*', nil, tcell.StyleDefault)

		for i, r := range cbe.Lang {
			renderer.SetContent(i+1, 1, r, nil, tcell.StyleDefault)
		}
		renderer.SetContent(0, 1, '|', nil, tcell.StyleDefault)
		renderer.SetContent(w+2, 1, '|', nil, tcell.StyleDefault)
		renderer = renderer.Translate(0, 2)
		subLines = 2
	}

	for i := 1; i < textLayout.Width-1; i++ {
		renderer.SetContent(i, 0, '-', nil, tcell.StyleDefault)
		renderer.SetContent(i, textLayout.Height-1-subLines, '-', nil, tcell.StyleDefault)
	}
	for i := 1; i < textLayout.Height-1-subLines; i++ {
		renderer.SetContent(0, i, '|', nil, tcell.StyleDefault)
		renderer.SetContent(textLayout.Width-1, i, '|', nil, tcell.StyleDefault)
	}
	renderer.SetContent(0, 0, '*', nil, tcell.StyleDefault)
	renderer.SetContent(textLayout.Width-1, 0, '*', nil, tcell.StyleDefault)
	renderer.SetContent(0, textLayout.Height-1-subLines, '*', nil, tcell.StyleDefault)
	renderer.SetContent(textLayout.Width-1, textLayout.Height-1-subLines, '*', nil, tcell.StyleDefault)
	for _, child := range textLayout.Children {
		childView := ctx.Resolver.Resolve(child.Element)
		childView.Draw(ctx, child, renderer.Translate(child.RelativeX, child.RelativeY-subLines))
	}
}

func (c *CodeBlockView) MinimumSize(ctx view.Context, e model.Element, width int, height int) *view.TextLayout {
	cbe := e.(*CodeBlockElement)

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
	additionalHeight := 0
	requireWidth := 0
	if len(cbe.Lang) > 0 {
		additionalHeight = 2
		requireWidth = runewidth.StringWidth(cbe.Lang) + 2 + 1
	}
	return &view.TextLayout{
		Element:       e,
		MinimumWidth:  max(maxWidth+2, requireWidth),
		MinimumHeight: e.GetElementCount() + 2 + additionalHeight,
		Children:      children,
	}
}

func (c *CodeBlockView) ViewLengthTable(ctx view.Context, textLayout *view.TextLayout) ([]int, int) {
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

func (c *CodeBlockView) findTableIndex(textLayout *view.TextLayout, table []int, viewLocalPos int) (Row int, Column int) {
	e := textLayout.Element
	cbe := e.(*CodeBlockElement)
	if len(cbe.Lang) > 0 {
		viewLocalPos -= runewidth.StringWidth(cbe.Lang) + 1
	}

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

func (c *CodeBlockView) sumTableValue(textLayout *view.TextLayout, table []int, index int) int {
	e := textLayout.Element
	cbe := e.(*CodeBlockElement)
	moves := 0
	if len(cbe.Lang) > 0 {
		moves = runewidth.StringWidth(cbe.Lang) + 1
	}

	v := 0
	for i := 0; i <= index; i++ {
		v += table[i]
	}
	return v + moves
}

func (c *CodeBlockView) MoveLength(ctx view.Context, textLayout *view.TextLayout) int {
	e := textLayout.Element
	cbe := e.(*CodeBlockElement)
	additionalMoves := 0
	if len(cbe.Lang) > 0 {
		additionalMoves = runewidth.StringWidth(cbe.Lang) + 1
	}

	_, ttl := c.ViewLengthTable(ctx, textLayout)
	return ttl + additionalMoves
}

func (c *CodeBlockView) MoveUp(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	e := textLayout.Element
	cbe := e.(*CodeBlockElement)
	additionalMoves := 0
	if len(cbe.Lang) > 0 {
		additionalMoves = runewidth.StringWidth(cbe.Lang) + 1
	}

	table, _ := c.ViewLengthTable(ctx, textLayout)
	index, _ := c.findTableIndex(textLayout, table, viewLocalPos)
	if index == -1 {
		return -1
	}
	if index <= 0 {
		if len(cbe.Lang) > 0 {
			return additionalMoves - 1
		}
		return -1
	}
	return c.sumTableValue(textLayout, table, index-2)
}

func (c *CodeBlockView) MoveDown(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	table, _ := c.ViewLengthTable(ctx, textLayout)
	index, _ := c.findTableIndex(textLayout, table, viewLocalPos)
	if index == len(table)-1 {
		return -1
	}
	return c.sumTableValue(textLayout, table, index)
}

func (c *CodeBlockView) MoveLeft(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	if viewLocalPos <= 0 {
		return -1
	}
	return viewLocalPos - 1
}

func (c *CodeBlockView) MoveRight(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	if viewLocalPos >= c.MoveLength(ctx, textLayout)-1 {
		return -1
	}
	return viewLocalPos + 1
}

func (c *CodeBlockView) ConvertPos(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (ViewLocalX int, ViewLocalY int) {
	e := textLayout.Element
	cbe := e.(*CodeBlockElement)

	additionalHeight := 0
	if len(cbe.Lang) > 0 {
		w := runewidth.StringWidth(cbe.Lang) + 1
		if viewLocalPos < w {
			return 1 + viewLocalPos, 1
		}
		additionalHeight = 2
	}

	table, _ := c.ViewLengthTable(ctx, textLayout)
	index, col := c.findTableIndex(textLayout, table, viewLocalPos)
	if viewLocalPos == c.sumTableValue(textLayout, table, len(table)-1) {
		child := textLayout.Children[len(textLayout.Children)-1]
		childView := ctx.Resolver.Resolve(child.Element)
		childLen := childView.MoveLength(ctx, child)
		lx, ly := childView.ConvertPos(ctx, child, childLen-1)
		return lx + 1, ly + 1 + additionalHeight
	}

	child := textLayout.Children[index]
	v := ctx.Resolver.Resolve(child.Element)
	lx, _ := v.ConvertPos(ctx, child, col)
	return lx + 1, index + 1 + additionalHeight
}

func (c *CodeBlockView) ConvertModel(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) view.CharacterReference {
	e := textLayout.Element
	cbe := e.(*CodeBlockElement)

	if len(cbe.Lang) > 0 {
		w := runewidth.StringWidth(cbe.Lang) + 1
		if viewLocalPos < w {
			r := e.GetRange(1)
			bytes := 1
			if viewLocalPos == w-1 {
				bytes = 0
			}
			return view.CharacterReference{
				StartPosition: model.Position{
					Row:    r.StartPosition.Row,
					Column: r.StartPosition.Column + viewLocalPos,
				},
				Bytes: bytes,
			}
		}
	}

	table, _ := c.ViewLengthTable(ctx, textLayout)
	index, col := c.findTableIndex(textLayout, table, viewLocalPos)
	if viewLocalPos == c.sumTableValue(textLayout, table, len(table)-1) {
		child := textLayout.Children[len(textLayout.Children)-1]
		childView := ctx.Resolver.Resolve(child.Element)
		childLen := childView.MoveLength(ctx, child)
		return childView.ConvertModel(ctx, child, childLen-1)
	}
	child := textLayout.Children[index]

	v := ctx.Resolver.Resolve(child.Element)
	return v.ConvertModel(ctx, child, col)
}

func (c *CodeBlockView) ConvertViewLocalPos(ctx view.Context, textLayout *view.TextLayout, bytePos model.Position) int {
	e := textLayout.Element
	cbe := e.(*CodeBlockElement)

	moves := 0
	if len(cbe.Lang) > 0 {
	}
	if len(cbe.Lang) > 0 {
		r := e.GetRange(1)

		if bytePos.Row == r.StartPosition.Row {
			if bytePos.Column >= r.StartPosition.Column && bytePos.Column < r.EndPosition.Column+1 {
				return bytePos.Column - r.StartPosition.Column
			}
		}
		moves = runewidth.StringWidth(cbe.Lang) + 1
	}

	viewOffset := moves
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

	return c.MoveLength(ctx, textLayout) - 1
}
