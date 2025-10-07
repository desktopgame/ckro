package litemark

import (
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/desktopgame/ckro/internal/tui/view"
	"github.com/gdamore/tcell/v2"
)

type HorizontalLineView struct {
}

func (hl *HorizontalLineView) Layout(ctx view.Context, textLayout *view.TextLayout, x, y, w, h int) {
	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
}

func (hl *HorizontalLineView) Draw(ctx view.Context, textLayout *view.TextLayout, renderer view.Renderer) {
	for i := 0; i < textLayout.Width; i++ {
		renderer.SetContent(i, 0, '-', nil, tcell.StyleDefault)
	}
}

func (hl *HorizontalLineView) MinimumSize(ctx view.Context, e model.Element, width int, height int) *view.TextLayout {
	return &view.TextLayout{
		Element:       e,
		MinimumWidth:  width,
		MinimumHeight: 1,
	}
}

func (hl *HorizontalLineView) MoveLength(ctx view.Context, textLayout *view.TextLayout) int {
	return 1
}

func (hl *HorizontalLineView) MoveUp(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	return -1
}

func (hl *HorizontalLineView) MoveDown(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	return -1
}

func (hl *HorizontalLineView) MoveLeft(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	return -1
}

func (hl *HorizontalLineView) MoveRight(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) int {
	return -1
}

func (hl *HorizontalLineView) ConvertPos(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (ViewLocalX int, ViewLocalY int) {
	return 0, 0
}

func (hl *HorizontalLineView) ConvertModel(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) view.CharacterReference {
	e := textLayout.Element
	r := e.GetRange(0)
	ed := r.EndPosition
	return view.CharacterReference{
		StartPosition: model.Position{
			Row:    ed.Row,
			Column: ed.Column,
		},
		Bytes: 0,
	}
}

func (hl *HorizontalLineView) ConvertViewLocalPos(ctx view.Context, textLayout *view.TextLayout, bytePos model.Position) int {
	return 0
}

func (hl *HorizontalLineView) ShouldBeforeInsertionNewLineOnLineBegin(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) bool {
	return false
}

func (hl *HorizontalLineView) ShouldRemoveWithLine(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (int, bool) {
	return -1, false
}

func (hl *HorizontalLineView) ShouldRemoveWithSpecifiedColumnAfter(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (model.Position, bool) {
	return model.Position{}, false
}

func (hl *HorizontalLineView) ShouldRemoveWithSpecifiedRangeLines(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (model.Range, bool) {
	return model.Range{}, false
}

func (hl *HorizontalLineView) ShouldRemoveWithSpecifiedRangeColumns(ctx view.Context, textLayout *view.TextLayout, viewLocalPos int) (model.Range, bool) {
	return model.Range{}, false
}

func (hl *HorizontalLineView) ShouldRemoveLastCharacter(ctx view.Context, textLayout *view.TextLayout) bool {
	return true
}
