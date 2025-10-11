package litemark

import (
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/desktopgame/ckro/internal/tui/view"
	"github.com/gdamore/tcell/v2"
)

type TableHeaderView struct {
}

func (thv *TableHeaderView) Layout(ctx view.Context, textLayout *view.TextLayout, x, y, w, h int) {
	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
}

func (thv *TableHeaderView) Draw(ctx view.Context, textLayout *view.TextLayout, renderer view.Renderer) {
	r := textLayout.Element.GetRange(0)
	at := r.StartPosition
	sel := ctx.TextSelection

	if sel.Contain(at) {
		selectStyle := tcell.StyleDefault.Reverse(true)
		renderer.SetContent(0, 0, ' ', nil, selectStyle)
	}
}

func (thv *TableHeaderView) Measure(ctx view.Context, e model.Element, width int, height int) *view.TextLayout {
	return &view.TextLayout{
		Element:       e,
		MinimumWidth:  1,
		MinimumHeight: 1,
	}
}

func (thv *TableHeaderView) MoveLength(ctx view.Context, textLayout *view.TextLayout) int {
	return 1
}

func (thv *TableHeaderView) MoveUp(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	return -1
}

func (thv *TableHeaderView) MoveDown(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	return -1
}

func (thv *TableHeaderView) MoveLeft(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	return -1
}

func (thv *TableHeaderView) MoveRight(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	return -1
}

func (thv *TableHeaderView) ConvertPos(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (ViewLocalX int, ViewLocalY int) {
	return 0, 0
}

func (thv *TableHeaderView) ConvertModel(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) view.CharacterReference {
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

func (thv *TableHeaderView) FindFoldElementAt(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (model.Element, int, bool) {
	return nil, 0, false
}

func (thv *TableHeaderView) ShouldBeforeInsertionNewLineOnLineBegin(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) bool {
	return false
}

func (thv *TableHeaderView) ShouldRemoveWithLine(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (int, bool) {
	return -1, false
}

func (thv *TableHeaderView) ShouldRemoveWithSpecifiedColumnAfter(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (model.Position, bool) {
	return model.Position{}, false
}

func (thv *TableHeaderView) ShouldRemoveWithSpecifiedRangeLines(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (model.Range, bool) {
	return model.Range{}, false
}

func (thv *TableHeaderView) ShouldRemoveWithSpecifiedRangeColumns(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (model.Range, bool) {
	return model.Range{}, false
}

func (thv *TableHeaderView) ShouldRemoveLastCharacter(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (model.Element, bool) {
	return nil, false
}

func (thv *TableHeaderView) ConvertViewLocalPos(ctx view.Context, textLayout *view.TextLayout, bytePos model.Position) int {
	return 0
}
