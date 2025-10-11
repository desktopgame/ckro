package litemark

import (
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/desktopgame/ckro/internal/tui/view"
	"github.com/gdamore/tcell/v2"
)

type TableRowView struct {
}

func (trv *TableRowView) Layout(ctx view.Context, textLayout *view.TextLayout, x, y, w, h int) {
	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
}

func (trv *TableRowView) Draw(ctx view.Context, textLayout *view.TextLayout, renderer view.Renderer) {
	r := textLayout.Element.GetRange(0)
	at := r.StartPosition
	sel := ctx.TextSelection

	if sel.Contain(at) {
		selectStyle := tcell.StyleDefault.Reverse(true)
		renderer.SetContent(0, 0, ' ', nil, selectStyle)
	}
}

func (trv *TableRowView) Measure(ctx view.Context, e model.Element, width int, height int) *view.TextLayout {
	return &view.TextLayout{
		Element:       e,
		MinimumWidth:  1,
		MinimumHeight: 1,
	}
}

func (trv *TableRowView) MoveLength(ctx view.Context, textLayout *view.TextLayout) int {
	return 1
}

func (trv *TableRowView) MoveUp(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	return -1
}

func (trv *TableRowView) MoveDown(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	return -1
}

func (trv *TableRowView) MoveLeft(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	return -1
}

func (trv *TableRowView) MoveRight(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	return -1
}

func (trv *TableRowView) ConvertPos(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (ViewLocalX int, ViewLocalY int) {
	return 0, 0
}

func (trv *TableRowView) ConvertModel(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) view.CharacterReference {
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

func (trv *TableRowView) FindFoldElementAt(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (model.Element, int, bool) {
	return nil, 0, false
}

func (trv *TableRowView) ShouldBeforeInsertionNewLineOnLineBegin(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) bool {
	return false
}

func (trv *TableRowView) ShouldRemoveWithLine(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (int, bool) {
	return -1, false
}

func (trv *TableRowView) ShouldRemoveWithSpecifiedColumnAfter(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (model.Position, bool) {
	return model.Position{}, false
}

func (trv *TableRowView) ShouldRemoveWithSpecifiedRangeLines(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (model.Range, bool) {
	return model.Range{}, false
}

func (trv *TableRowView) ShouldRemoveWithSpecifiedRangeColumns(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (model.Range, bool) {
	return model.Range{}, false
}

func (trv *TableRowView) ShouldRemoveLastCharacter(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (model.Element, bool) {
	return nil, false
}

func (trv *TableRowView) ConvertViewLocalPos(ctx view.Context, textLayout *view.TextLayout, bytePos model.Position) int {
	return 0
}
