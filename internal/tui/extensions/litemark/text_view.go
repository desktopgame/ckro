package litemark

import (
	"github.com/desktopgame/ckro/internal/text"
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/desktopgame/ckro/internal/tui/view"
)

type TextView struct {
}

func (t *TextView) Layout(ctx view.Context, textLayout *view.TextLayout, x, y, w, h int) {
	offsetX := 0
	for i := 0; i < len(textLayout.Children); i++ {
		childElement := textLayout.Children[i].Element
		childView := ctx.Resolver.Resolve(childElement)

		mw := textLayout.Children[i].MinimumWidth
		childView.Layout(ctx, textLayout.Children[i], offsetX, 0, mw, 1)
		offsetX += textLayout.Children[i].Width
	}
	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
}

func (t *TextView) Draw(ctx view.Context, textLayout *view.TextLayout, renderer view.Renderer) {
	for _, child := range textLayout.Children {
		childView := ctx.Resolver.Resolve(child.Element)
		childView.Draw(ctx, child, renderer.Translate(child.RelativeX, child.RelativeY))
	}
}

func (t *TextView) MinimumSize(ctx view.Context, e model.Element, width int, height int) *view.TextLayout {
	totalWidth := 0
	children := []*view.TextLayout{}
	column := 0
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := ctx.Resolver.Resolve(childElement)

		var child *view.TextLayout
		if tsView, ok := childView.(view.TabStopTextView); ok {
			width := tsView.WidthWithTabStop(ctx, childElement, column)
			child = &view.TextLayout{
				Element:       childElement,
				MinimumWidth:  width,
				MinimumHeight: 1,
			}
		} else {
			child = childView.MinimumSize(ctx, childElement, width, 1)
		}
		children = append(children, child)

		column += child.MinimumWidth
		totalWidth += child.MinimumWidth
	}
	return &view.TextLayout{
		Element:       e,
		MinimumWidth:  totalWidth,
		MinimumHeight: 1,
		Children:      children,
	}
}

func (t *TextView) MoveLength(ctx view.Context, textLayout *view.TextLayout) int {
	totalLength := 0
	for i := 0; i < len(textLayout.Children); i++ {
		childElement := textLayout.Children[i].Element
		childView := ctx.Resolver.Resolve(childElement)

		totalLength += childView.MoveLength(ctx, textLayout.Children[i])
	}
	return totalLength + 1
}

func (t *TextView) MoveUp(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	return -1
}

func (t *TextView) MoveDown(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	return -1
}

func (t *TextView) MoveLeft(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	if viewLocalPos <= 0 {
		return -1
	}
	return viewLocalPos - 1
}

func (t *TextView) MoveRight(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	if viewLocalPos >= t.MoveLength(ctx, textLayout)-1 {
		return -1
	}
	return viewLocalPos + 1
}

func (t *TextView) RemoveCombine(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (bool, view.CharacterReference) {
	vls := 0
	for i := 0; i < len(textLayout.Children); i++ {
		child := textLayout.Children[i]
		childElement := child.Element

		childView := ctx.Resolver.Resolve(childElement)
		childViewLen := childView.MoveLength(ctx, textLayout.Children[i])
		if childViewLen == 1 && vls+1 == viewLocalPos {
			r := childElement.GetRange(0)
			bPos := view.CharacterReference{
				StartPosition: model.Position{
					Row:    r.StartPosition.Row,
					Column: r.StartPosition.Column,
				},
				Bytes: (r.EndPosition.Column - r.StartPosition.Column),
			}
			return true, bPos
		}
		vls += childViewLen
	}
	if len(textLayout.Children) == 1 {

	}
	return false, view.CharacterReference{}
}

func (t *TextView) ConvertPos(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (ViewLocalX int, ViewLocalY int) {
	e := textLayout.Element
	return text.DisplayPos(ctx.GetText(e), viewLocalPos), 0
}

func (t *TextView) ConvertViewLocalPos(ctx view.Context, textLayout *view.TextLayout, bytePos model.Position) int {
	if bytePos.Column == 0 {
		return 0
	}
	vls := 0
	for i := 0; i < len(textLayout.Children); i++ {
		child := textLayout.Children[i]
		childElement := child.Element

		r := childElement.GetRange(0)
		st := r.StartPosition
		ed := r.EndPosition
		if il, ok := childElement.(*InlineElement); ok {
			if il.IsBold || il.IsItalic || il.IsUnderline || il.Foreground.IsSome() || il.Background.IsSome() {
				ed.Column++
			}
		}
		childView := ctx.Resolver.Resolve(childElement)
		if bytePos.Column >= st.Column && (bytePos.Column < ed.Column || ed.Row > st.Row) {
			return vls + childView.ConvertViewLocalPos(ctx, child, bytePos) // + 1
		}
		vls += childView.MoveLength(ctx, textLayout.Children[i])
	}
	return t.MoveLength(ctx, textLayout) - 1
}

func (t *TextView) ConvertModel(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) view.CharacterReference {
	e := textLayout.Element
	r := e.GetRange(0)
	ed := r.EndPosition

	totalLen := 0
	for i := 0; i < len(textLayout.Children); i++ {
		child := textLayout.Children[i]
		childElement := child.Element
		childView := ctx.Resolver.Resolve(childElement)
		childViewLen := childView.MoveLength(ctx, child)
		start := totalLen
		end := start + childViewLen

		if viewLocalPos >= start && viewLocalPos < end {
			return childView.ConvertModel(ctx, child, viewLocalPos-start)
		}

		totalLen += childViewLen
	}

	return view.CharacterReference{
		StartPosition: model.Position{
			Row:    ed.Row,
			Column: ed.Column,
		},
		Bytes: 0,
	}
}

func (t *TextView) ShouldBeforeInsertionNewLineOnLineBegin(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) bool {
	if viewLocalPos == 0 {
		if len(textLayout.Children) == 0 {
			return false
		}
		child := textLayout.Children[0]
		childElement := child.Element
		childView := ctx.Resolver.Resolve(childElement)
		return childView.ShouldBeforeInsertionNewLineOnLineBegin(ctx, child, 0)
	}
	return false
}

func (t *TextView) ShouldRemoveWithLine(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) bool {
	return false
}

func (t *TextView) ShouldRemoveWithSpecifiedColumnAfter(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (model.Position, bool) {
	return model.Position{}, false
}

func (t *TextView) ShouldRemoveWithSpecifiedRangeLines(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (model.Range, bool) {
	return model.Range{}, false
}

func (t *TextView) ShouldRemoveWithSpecifiedRangeColumns(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (model.Range, bool) {
	vls := 0
	for i := 0; i < len(textLayout.Children); i++ {
		child := textLayout.Children[i]
		childElement := child.Element

		childView := ctx.Resolver.Resolve(childElement)
		childViewLen := childView.MoveLength(ctx, textLayout.Children[i])
		if childViewLen == 1 && vls+1 == viewLocalPos {
			r := childElement.GetRange(0)
			return r, true
		}
		vls += childViewLen
	}
	return model.Range{}, false
}

func (t *TextView) ShouldRemoveLastCharacter(ctx view.Context, textLayout *view.TextLayout) bool {
	return false
}

func (t *TextView) ConvertRelativeX(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	return viewLocalPos
}

func (t *TextView) MoveFirstLine(ctx view.Context, textLayout *view.TextLayout, relX int) int {
	l := t.MoveLength(ctx, textLayout)
	return min(relX, l-1)
}

func (t *TextView) MoveLastLine(ctx view.Context, textLayout *view.TextLayout, relX int) int {
	l := t.MoveLength(ctx, textLayout)
	return min(relX, l-1)
}
