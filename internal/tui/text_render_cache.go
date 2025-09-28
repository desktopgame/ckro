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
}

func (trc *TextRenderCache) Update(ctx view.Context, textBoxWidth int) {
	newDocVersion := ctx.Document.GetVersion()
	if trc.documentVersion > 0 && trc.documentVersion == newDocVersion {
		return
	}

	elements := ctx.Document.Render()

	entries := []*view.TextLayout{}
	for i := 0; i < len(elements); i++ {
		element := elements[i]
		textView := ctx.Resolver.Resolve(element)
		tl := textView.MinimumSize(ctx, element, textBoxWidth, 9999)
		entries = append(entries, tl)
	}

	totalViewLen := 0
	viewLenTable := []int{}
	for _, elem := range elements {
		view := ctx.Resolver.Resolve(elem)
		viewLen := view.MoveLength(ctx, elem)

		viewLenTable = append(viewLenTable, viewLen)
		totalViewLen += viewLen
	}

	trc.documentVersion = newDocVersion
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
		for i := 0; i < len(trc.viewLenTable); i++ {
			ttl += trc.viewLenTable[i]
		}
		elementIndex = len(trc.elements) - 1
		viewLocalPosition = trc.viewLenTable[len(trc.viewLenTable)-1]
		elementStart = ttl
	}
	return totalViewLen, elementIndex, elementStart, viewLocalPosition
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
