package tui

import (
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/desktopgame/ckro/internal/tui/view"
)

type TextRenderCache struct {
	documentVersion uint
	elements        []model.Element
	layoutCache     []*view.TextLayout
	totalViewLen    int
	viewLenTable    []int
	textBoxWidth    int
}

func (trc *TextRenderCache) ForceUpdate(ctx view.Context, textBoxWidth int) {
	trc.updateImpl(ctx, textBoxWidth, true)
}

func (trc *TextRenderCache) Update(ctx view.Context, textBoxWidth int) {
	trc.updateImpl(ctx, textBoxWidth, false)
}

func (trc *TextRenderCache) updateImpl(ctx view.Context, textBoxWidth int, forceUpdate bool) {
	layoutChanged := trc.textBoxWidth != textBoxWidth
	newDocVersion := ctx.Document.GetVersion()
	if (trc.documentVersion > 0 && trc.documentVersion == newDocVersion) && !layoutChanged && !forceUpdate {
		return
	}
	trc.textBoxWidth = textBoxWidth
	trc.documentVersion = newDocVersion

	elements := ctx.Document.Render()

	entries := []*view.TextLayout{}
	newElements := []model.Element{}
	viewLine := 0
	for i := 0; i < len(elements); i++ {
		element := elements[i]
		textView := ctx.Resolver.Resolve(element)
		tl := textView.MinimumSize(ctx, element, textBoxWidth, 9999)
		textView.Layout(ctx, tl, 0, viewLine, tl.MinimumWidth, tl.MinimumHeight)

		if tl.Width > textBoxWidth {
			sg := ctx.Document.Read(element.GetRange(0))
			if sg.GetLineCount() > 1 {
				for j := 0; j < sg.GetLineCount(); j++ {
					row := element.GetRange(0).StartPosition.Row + j
					pElement := &model.PlainElement{
						Range: model.Range{
							StartPosition: model.Position{
								Row:    row,
								Column: sg.GetSpan(j).StartColumn,
							},
							EndPosition: model.Position{
								Row:    row,
								Column: sg.GetSpan(j).EndColumn,
							},
						},
					}

					textView = ctx.Resolver.Resolve(pElement)
					tl = textView.MinimumSize(ctx, pElement, textBoxWidth, 9999)
					textView.Layout(ctx, tl, 0, viewLine, tl.MinimumWidth, tl.MinimumHeight)

					entries = append(entries, tl)
					newElements = append(newElements, pElement)
					viewLine += tl.Height
				}
			} else {
				pElement := &model.PlainElement{
					Range: element.GetRange(0),
				}

				textView = ctx.Resolver.Resolve(pElement)
				tl = textView.MinimumSize(ctx, pElement, textBoxWidth, 9999)
				textView.Layout(ctx, tl, 0, viewLine, tl.MinimumWidth, tl.MinimumHeight)

				entries = append(entries, tl)
				newElements = append(newElements, pElement)
				viewLine += tl.Height
			}
		} else {
			entries = append(entries, tl)
			newElements = append(newElements, element)
			viewLine += tl.Height
		}
	}

	elements = newElements

	totalViewLen := 0
	viewLenTable := []int{}
	for _, l := range entries {
		view := ctx.Resolver.Resolve(l.Element)
		viewLen := view.MoveLength(ctx, l)

		viewLenTable = append(viewLenTable, viewLen)
		totalViewLen += viewLen
	}

	trc.elements = elements
	trc.layoutCache = entries
	trc.totalViewLen = totalViewLen
	trc.viewLenTable = viewLenTable
}

func (trc *TextRenderCache) Stats(viewPosition int) (TotalViewLen int, ElementIndex int, ViewStart int, ViewLocalPosition int) {
	totalViewLen := 0
	elementIndex := -1
	elementStart := -1
	viewLocalPosition := 0
	for i, viewLen := range trc.viewLenTable {
		viewStart := totalViewLen
		viewEnd := viewStart + viewLen

		if viewPosition >= viewStart && viewPosition < viewEnd {
			elementIndex = i
			viewLocalPosition = viewPosition - viewStart
			elementStart = viewStart
		}
		totalViewLen += viewLen
	}
	if elementIndex == -1 {
		ttl := 0
		for i := 0; i < len(trc.viewLenTable)-1; i++ {
			ttl += trc.viewLenTable[i]
		}
		elementIndex = len(trc.elements) - 1
		viewLocalPosition = trc.viewLenTable[len(trc.viewLenTable)-1]
		elementStart = ttl
	}
	return totalViewLen, elementIndex, elementStart, viewLocalPosition
}

func (trc *TextRenderCache) Total() int {
	return trc.totalViewLen
}

func (trc *TextRenderCache) GetItem(index int) (Element model.Element, Layout *view.TextLayout) {
	return trc.GetElement(index), trc.GetLayout(index)
}

func (trc *TextRenderCache) GetElement(index int) model.Element {
	return trc.elements[index]
}

func (trc *TextRenderCache) GetLayout(index int) *view.TextLayout {
	return trc.layoutCache[index]
}

func (trc *TextRenderCache) GetItemCount() int {
	return len(trc.elements)
}
