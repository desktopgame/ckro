package view

import (
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/gdamore/tcell/v2"
)

type FoldBlockView struct {
}

func (fv *FoldBlockView) Layout(ctx Context, textLayout *TextLayout, x, y, w, h int) {
	headerHeight := 1
	if ctx.FoldManager.IsFolded(ctx.Document, textLayout.Element) {
		childElement := textLayout.Children[0].Element
		childView := ctx.Resolver.Resolve(childElement)

		mw := textLayout.Children[0].MinimumWidth
		childView.Layout(ctx, textLayout.Children[0], 2, 1+headerHeight, mw, 1)
	} else {
		offsetY := 1 + headerHeight
		for i := 0; i < len(textLayout.Children); i++ {
			childElement := textLayout.Children[i].Element
			childView := ctx.Resolver.Resolve(childElement)

			mw := textLayout.Children[i].MinimumWidth
			mh := textLayout.Children[i].MinimumHeight
			childView.Layout(ctx, textLayout.Children[i], 2, offsetY, mw, mh)
			offsetY += textLayout.Children[i].Height
		}
	}
	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
}

func (fv *FoldBlockView) Draw(ctx Context, textLayout *TextLayout, renderer Renderer) {
	foldFrameStyle := tcell.StyleDefault.Foreground(tcell.ColorYellow)
	foldLabel := "[-]"
	if ctx.FoldManager.IsFolded(ctx.Document, textLayout.Element) {
		foldLabel = "[+]"
	}

	for i, r := range foldLabel {
		renderer.SetContent(i, 0, r, nil, foldFrameStyle)
	}
	for i := 2; i < textLayout.Height-1; i++ {
		renderer.SetContent(1, i, '|', nil, foldFrameStyle)
	}
	renderer.SetContent(1, 1, '*', nil, foldFrameStyle)
	renderer.SetContent(1, textLayout.Height-1, '*', nil, foldFrameStyle)

	for _, child := range textLayout.Children {
		childView := ctx.Resolver.Resolve(child.Element)
		childView.Draw(ctx, child, renderer.Translate(child.RelativeX, child.RelativeY).Region(child.Width, child.Height))
	}
}

func (fv *FoldBlockView) Measure(ctx Context, e model.Element, width int, height int) *TextLayout {
	children := []*TextLayout{}

	var minimumHeight int
	if ctx.FoldManager.IsFolded(ctx.Document, e) {
		minimumHeight = 2

		childElement := e.GetElement(0)
		childView := ctx.Resolver.Resolve(childElement)
		child := childView.Measure(ctx, childElement, width-2, 1)
		if child.MinimumHeight > 1 {
			r := childElement.GetRange(0)
			if r.StartPosition.Row != r.EndPosition.Row {
				r = model.Range{
					StartPosition: r.StartPosition,
					EndPosition: model.Position{
						Row:    r.StartPosition.Row,
						Column: ctx.Document.GetLineBytes(r.StartPosition.Row),
					},
				}
			}
			childElement = &model.PlainElement{
				Range: r,
			}
			childView = &PlainTextView{}
			child = childView.Measure(ctx, childElement, width-2, 9999)
			minimumHeight += child.MinimumHeight
			children = append(children, child)
		} else if child.MinimumWidth+2 > width {
			r := childElement.GetRange(0)
			if r.StartPosition.Row == r.EndPosition.Row {
				childElement = &model.PlainElement{
					Range: r,
				}
				childView = &PlainTextView{}
				child = childView.Measure(ctx, childElement, width-2, 9999)
				minimumHeight += child.MinimumHeight
				children = append(children, child)
			} else {
				for j := r.StartPosition.Row; j <= r.EndPosition.Row; j++ {
					r2 := model.Range{
						StartPosition: model.Position{
							Row:    j,
							Column: 0,
						},
						EndPosition: model.Position{
							Row:    j,
							Column: ctx.Document.GetLineBytes(j),
						},
					}
					childElement = &model.PlainElement{
						Range: r2,
					}
					childView = &PlainTextView{}
					child = childView.Measure(ctx, childElement, width-2, 9999)
					minimumHeight += child.MinimumHeight
					children = append(children, child)
				}
			}
		} else {
			minimumHeight++
			children = append(children, child)
		}

	} else {
		minimumHeight = 2

		for i := 0; i < e.GetElementCount(); i++ {
			childElement := e.GetElement(i)
			childView := ctx.Resolver.Resolve(childElement)

			child := childView.Measure(ctx, childElement, width-2, 9999)

			if child.MinimumWidth+2 > width {

				r := childElement.GetRange(0)
				if r.StartPosition.Row == r.EndPosition.Row {
					childElement = &model.PlainElement{
						Range: r,
					}
					childView = &PlainTextView{}
					child = childView.Measure(ctx, childElement, width-2, 9999)
					minimumHeight += child.MinimumHeight
					children = append(children, child)
				} else {

					for j := r.StartPosition.Row; j <= r.EndPosition.Row; j++ {
						r2 := model.Range{
							StartPosition: model.Position{
								Row:    j,
								Column: 0,
							},
							EndPosition: model.Position{
								Row:    j,
								Column: ctx.Document.GetLineBytes(j),
							},
						}
						childElement = &model.PlainElement{
							Range: r2,
						}
						childView = &PlainTextView{}
						child = childView.Measure(ctx, childElement, width-2, 9999)
						minimumHeight += child.MinimumHeight
						children = append(children, child)
					}
				}
			} else {
				minimumHeight += child.MinimumHeight
				children = append(children, child)
			}

		}
	}
	return &TextLayout{
		Element:       e,
		MinimumWidth:  width,
		MinimumHeight: minimumHeight + 1,
		Children:      children,
	}
}

func (fv *FoldBlockView) MoveLength(ctx Context, textLayout *TextLayout) int {
	if ctx.FoldManager.IsFolded(ctx.Document, textLayout.Element) {
		childElement := textLayout.Children[0].Element
		childView := ctx.Resolver.Resolve(childElement)
		return childView.MoveLength(ctx, textLayout.Children[0]) + 1
	} else {
		_, ttl := CompositeViewLengthTable(ctx, textLayout)
		return ttl + 1
	}
}

func (fv *FoldBlockView) MoveUp(ctx Context, textLayout *TextLayout, viewLocalPos int) int {
	if viewLocalPos == 0 {
		return -1
	}
	n := CompositeMoveUp(ctx, textLayout, viewLocalPos-1)
	if n == -1 {
		return 0
	}
	return n + 1
}

func (fv *FoldBlockView) MoveDown(ctx Context, textLayout *TextLayout, viewLocalPos int) int {
	if viewLocalPos == 0 {
		return 1
	}
	n := CompositeMoveDown(ctx, textLayout, viewLocalPos-1)
	if n >= 0 {
		n++
	}
	return n
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

	// skip header
	if viewLocalPos == 0 {
		return 1, 0
	}
	viewLocalPos--

	if ctx.FoldManager.IsFolded(ctx.Document, e) {
		childElement := textLayout.Children[0].Element
		childView := ctx.Resolver.Resolve(childElement)
		lx, ly := childView.ConvertPos(ctx, textLayout.Children[0], viewLocalPos)
		return 2 + lx, 1 + ly + 1
	}
	table, _ := CompositeViewLengthTable(ctx, textLayout)
	index, col := CompositeViewIndex(table, viewLocalPos)

	h := 0
	for i := 0; i < index; i++ {
		h += textLayout.Children[i].Height
	}

	child := textLayout.Children[index]
	v := ctx.Resolver.Resolve(child.Element)
	lx, ly := v.ConvertPos(ctx, child, col)
	return lx + 2, h + ly + 1 + 1
}

func (fv *FoldBlockView) ConvertModel(ctx Context, textLayout *TextLayout, viewLocalPos int) CharacterReference {
	// skip header
	if viewLocalPos == 0 {
		r := textLayout.Element.GetRange(0)
		fb := textLayout.Element.(*model.FoldBlockElement)
		return CharacterReference{
			StartPosition: model.Position{
				Row:    r.StartPosition.Row,
				Column: r.StartPosition.Column + fb.Level,
			},
			Bytes: 0,
		}
	}
	viewLocalPos--

	if ctx.FoldManager.IsFolded(ctx.Document, textLayout.Element) {
		childElement := textLayout.Children[0].Element
		childView := ctx.Resolver.Resolve(childElement)
		return childView.ConvertModel(ctx, textLayout.Children[0], viewLocalPos)
	} else {

		table, _ := CompositeViewLengthTable(ctx, textLayout)
		index, col := CompositeViewIndex(table, viewLocalPos)
		//if viewLocalPos == fv.sumTableValue(table, len(table)-1) {
		//	child := textLayout.Children[len(textLayout.Children)-1]
		//	childView := ctx.Resolver.Resolve(child.Element)
		//	childLen := childView.MoveLength(ctx, child)
		//	return childView.ConvertModel(ctx, child, childLen-1)
		//}
		child := textLayout.Children[index]

		v := ctx.Resolver.Resolve(child.Element)
		return v.ConvertModel(ctx, child, col)
	}
}

func (fv *FoldBlockView) ConvertViewLocalPos(ctx Context, textLayout *TextLayout, bytePos model.Position) int {
	e := textLayout.Element
	r := e.GetRange(0)
	if bytePos.Row == r.StartPosition.Row {
		return 0
	}
	if ctx.FoldManager.IsFolded(ctx.Document, textLayout.Element) {
		r := e.GetRange(0)
		child := textLayout.Children[0]
		childElement := textLayout.Children[0].Element
		childView := ctx.Resolver.Resolve(childElement)
		if bytePos.Row == r.StartPosition.Row {
			return childView.ConvertViewLocalPos(ctx, textLayout.Children[0], bytePos) + 1
		}
		return childView.MoveLength(ctx, child) - 1

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
						return viewOffset + childView.ConvertViewLocalPos(ctx, child, bytePos) + 1
					}
				}
				// inclusive line end, because of FoldBlockView is only contain line orientated view
				if bytePos.Column >= st.Column && (bytePos.Column <= ed.Column || ed.Row > st.Row) {
					return viewOffset + childView.ConvertViewLocalPos(ctx, child, bytePos) + 1
				}
			}
			viewOffset += childView.MoveLength(ctx, child)
		}

		return fv.MoveLength(ctx, textLayout) - 1
	}
}

func (fv *FoldBlockView) FindFoldElementAt(ctx Context, textLayout *TextLayout, viewLocalPos int) (model.Element, int, bool) {
	if ctx.FoldManager.IsFolded(ctx.Document, textLayout.Element) {
		return textLayout.Element, viewLocalPos, true
	}

	table, _ := CompositeViewLengthTable(ctx, textLayout)
	index, col := CompositeViewIndex(table, viewLocalPos-1)
	if index >= 0 {
		child := textLayout.Children[index]
		childView := ctx.Resolver.Resolve(child.Element)
		e, pos, ok := childView.FindFoldElementAt(ctx, child, col)
		if ok {
			return e, pos, true
		}
	}
	if viewLocalPos == 0 {
		return textLayout.Element, 0, true
	}
	return nil, -1, false
}

func (fv *FoldBlockView) ShouldBeforeInsertionNewLineOnLineBegin(ctx Context, textLayout *TextLayout, viewLocalPos int) bool {
	// skip header
	if viewLocalPos == 0 {
		return false
	}
	viewLocalPos--

	if ctx.FoldManager.IsFolded(ctx.Document, textLayout.Element) {
		return false
	} else {
		table, _ := CompositeViewLengthTable(ctx, textLayout)
		index, col := CompositeViewIndex(table, viewLocalPos)
		child := textLayout.Children[index]

		v := ctx.Resolver.Resolve(child.Element)
		return v.ShouldBeforeInsertionNewLineOnLineBegin(ctx, child, col)
	}
}

func (fv *FoldBlockView) ShouldRemoveWithLine(ctx Context, textLayout *TextLayout, viewLocalPos int) (int, bool) {
	// skip header
	if viewLocalPos == 0 {
		return -1, false
	}
	viewLocalPos--

	if ctx.FoldManager.IsFolded(ctx.Document, textLayout.Element) {
		return -1, false
	} else {
		table, _ := CompositeViewLengthTable(ctx, textLayout)
		index, col := CompositeViewIndex(table, viewLocalPos)
		child := textLayout.Children[index]

		v := ctx.Resolver.Resolve(child.Element)
		return v.ShouldRemoveWithLine(ctx, child, col)
	}
}

func (fv *FoldBlockView) ShouldRemoveWithSpecifiedColumnAfter(ctx Context, textLayout *TextLayout, viewLocalPos int) (model.Position, bool) {
	// skip header
	if viewLocalPos == 0 {
		r2 := textLayout.Element.GetRange(1)
		return model.Position{
			Row:    r2.EndPosition.Row,
			Column: r2.EndPosition.Column - 1,
		}, true
	}
	viewLocalPos--

	if ctx.FoldManager.IsFolded(ctx.Document, textLayout.Element) {
		return model.Position{}, false
	} else {
		table, _ := CompositeViewLengthTable(ctx, textLayout)
		index, col := CompositeViewIndex(table, viewLocalPos)
		child := textLayout.Children[index]

		v := ctx.Resolver.Resolve(child.Element)
		p, ok := v.ShouldRemoveWithSpecifiedColumnAfter(ctx, child, col)
		if ok {
			return p, ok
		}
	}
	return model.Position{}, false
}

func (fv *FoldBlockView) ShouldRemoveWithSpecifiedRangeLines(ctx Context, textLayout *TextLayout, viewLocalPos int) (model.Range, bool) {
	// skip header
	if viewLocalPos == 0 {
		return model.Range{}, false
	}
	viewLocalPos--

	if ctx.FoldManager.IsFolded(ctx.Document, textLayout.Element) {
		return model.Range{}, false
	} else {
		table, _ := CompositeViewLengthTable(ctx, textLayout)
		index, col := CompositeViewIndex(table, viewLocalPos)
		child := textLayout.Children[index]

		v := ctx.Resolver.Resolve(child.Element)
		return v.ShouldRemoveWithSpecifiedRangeLines(ctx, child, col)
	}
}

func (fv *FoldBlockView) ShouldRemoveWithSpecifiedRangeColumns(ctx Context, textLayout *TextLayout, viewLocalPos int) (model.Range, bool) {
	// skip header
	if viewLocalPos == 0 {
		return model.Range{}, false
	}
	viewLocalPos--

	if ctx.FoldManager.IsFolded(ctx.Document, textLayout.Element) {
		return model.Range{}, false
	} else {
		table, _ := CompositeViewLengthTable(ctx, textLayout)
		index, col := CompositeViewIndex(table, viewLocalPos)
		child := textLayout.Children[index]

		v := ctx.Resolver.Resolve(child.Element)
		return v.ShouldRemoveWithSpecifiedRangeColumns(ctx, child, col)
	}
}

func (fv *FoldBlockView) ShouldRemoveLastCharacter(ctx Context, textLayout *TextLayout, viewLocalPos int) (model.Element, bool) {
	if viewLocalPos == fv.MoveLength(ctx, textLayout)-1 {
		return textLayout.Element, true
	}

	// skip header
	if viewLocalPos == 0 {
		return nil, false
	}
	viewLocalPos--

	if ctx.FoldManager.IsFolded(ctx.Document, textLayout.Element) {
		return nil, false
	} else {
		table, _ := CompositeViewLengthTable(ctx, textLayout)
		index, col := CompositeViewIndex(table, viewLocalPos)
		child := textLayout.Children[index]

		v := ctx.Resolver.Resolve(child.Element)
		return v.ShouldRemoveLastCharacter(ctx, child, col)
	}
}
