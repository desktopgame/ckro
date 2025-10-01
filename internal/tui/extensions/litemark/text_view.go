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

func (t *TextView) MoveLength(ctx view.Context, e model.Element) int {
	totalLength := 0
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := ctx.Resolver.Resolve(childElement)

		totalLength += childView.MoveLength(ctx, childElement)
	}
	return totalLength + 1
}

func (t *TextView) MoveUp(ctx view.Context, e model.Element, viewLocalPos int) int {
	return -1
}

func (t *TextView) MoveDown(ctx view.Context, e model.Element, viewLocalPos int) int {
	return -1
}

func (t *TextView) MoveLeft(ctx view.Context, e model.Element, viewLocalPos int) int {
	if viewLocalPos <= 0 {
		return -1
	}
	return viewLocalPos - 1
}

func (t *TextView) MoveRight(ctx view.Context, e model.Element, viewLocalPos int) int {
	if viewLocalPos >= t.MoveLength(ctx, e)-1 {
		return -1
	}
	return viewLocalPos + 1
}

func (t *TextView) ConvertPos(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (ViewLocalX int, ViewLocalY int) {
	e := textLayout.Element
	return text.DisplayPos(ctx.GetText(e), viewLocalPos), 0
}

func (t *TextView) ConvertViewLocalPos(ctx view.Context, textLayout *view.TextLayout, bytePos model.Position) int {

	vls := 0
	for i := 0; i < len(textLayout.Children); i++ {
		child := textLayout.Children[i]
		childElement := child.Element

		r := childElement.GetRange(0)
		st := r.StartPosition
		ed := r.EndPosition
		childView := ctx.Resolver.Resolve(childElement)
		if bytePos.Column >= st.Column && (bytePos.Column < ed.Column || ed.Row > st.Row) {
			return vls + childView.ConvertViewLocalPos(ctx, child, bytePos)
		}
		vls += childView.MoveLength(ctx, childElement)
	}
	return t.MoveLength(ctx, textLayout.Element) - 1
}

func (t *TextView) ConvertModel(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) view.CharacterReference {
	e := textLayout.Element
	r := e.GetRange(0)
	str := ctx.Document.Read(r).GetLine(0)
	st := r.StartPosition
	ed := r.EndPosition

	graphemes := text.GraphemeLength(str)
	if viewLocalPos >= graphemes {
		return view.CharacterReference{
			StartPosition: model.Position{
				Row:    st.Row,
				Column: st.Column + len(str),
			},
			Bytes: 0,
		}
	}

	totalLen := 0
	for i := 0; i < len(textLayout.Children); i++ {
		child := textLayout.Children[i]
		childElement := child.Element
		childView := ctx.Resolver.Resolve(childElement)
		childViewLen := childView.MoveLength(ctx, childElement)
		start := totalLen
		end := start + childViewLen

		if viewLocalPos >= start && viewLocalPos < end {
			return childView.ConvertModel(ctx, child, viewLocalPos-start)
		}

		totalLen += childViewLen
	}
	//lastElement := textLayout.Children[len(textLayout.Children)-1].Element
	//lastView := ctx.Resolver.Resolve(lastElement)
	//panic("")

	return view.CharacterReference{
		StartPosition: model.Position{
			Row:    ed.Row,
			Column: ed.Column,
		},
		Bytes: 0,
	}
}

func (t *TextView) ConvertRelativeX(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	return viewLocalPos
}

func (t *TextView) MoveFirstLine(ctx view.Context, e model.Element, relX int) int {
	l := t.MoveLength(ctx, e)
	return min(relX, l-1)
}

func (t *TextView) MoveLastLine(ctx view.Context, e model.Element, relX int) int {
	l := t.MoveLength(ctx, e)
	return min(relX, l-1)
}
