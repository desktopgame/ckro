package view

import (
	"github.com/desktopgame/ckro/internal/tui/model"
)

type GhostView struct {
}

func (g *GhostView) Layout(ctx Context, textLayout *TextLayout, x, y, w, h int) {
	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
}

func (g *GhostView) Draw(ctx Context, textLayout *TextLayout, renderer Renderer) {
}

func (g *GhostView) MinimumSize(ctx Context, e model.Element, width int, height int) *TextLayout {
	return &TextLayout{
		Element:       e,
		MinimumWidth:  width,
		MinimumHeight: 1,
	}
}

func (g *GhostView) MoveLength(ctx Context, e model.Element) int {
	return 1
}

func (g *GhostView) MoveUp(ctx Context, e model.Element, viewLocalPos int) int {
	return -1
}

func (g *GhostView) MoveDown(ctx Context, e model.Element, viewLocalPos int) int {
	return -1
}

func (g *GhostView) MoveLeft(ctx Context, e model.Element, viewLocalPos int) int {
	return -1
}

func (g *GhostView) MoveRight(ctx Context, e model.Element, viewLocalPos int) int {
	return -1
}

func (g *GhostView) ConvertPos(ctx Context, textLayout *TextLayout, viewLocalPos int) (ViewLocalX int, ViewLocalY int) {
	return 0, 0
}

func (g *GhostView) ConvertModel(ctx Context, textLayout *TextLayout, viewLocalPos int) CharacterReference {
	e := textLayout.Element
	r := e.GetRange(0)
	ed := r.EndPosition
	return CharacterReference{
		StartPosition: model.Position{
			Row:    ed.Row,
			Column: ed.Column,
		},
		Bytes: 0,
	}
}

func (g *GhostView) ConvertViewLocalPos(ctx Context, textLayout *TextLayout, bytePos model.Position) int {
	return 0
}
