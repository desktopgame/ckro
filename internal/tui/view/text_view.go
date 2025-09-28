package view

import "github.com/desktopgame/ckro/internal/tui/model"

type TextView interface {
	// x,y はこのビューを置くべき左上の座標
	// w,h はこのビューに与えられたサイズ
	Layout(ctx Context, textLayout *TextLayout, x, y, w, h int)

	// rendererは0,0をビューの左上として描画を開始する
	Draw(ctx Context, textLayout *TextLayout, renderer Renderer)

	// width,height はこのビューに与えられたサイズ
	MinimumSize(ctx Context, e model.Element, width int, height int) *TextLayout

	MoveLength(ctx Context, e model.Element) int
	MoveUp(ctx Context, e model.Element, viewLocalPos int) int
	MoveDown(ctx Context, e model.Element, viewLocalPos int) int
	MoveLeft(ctx Context, e model.Element, viewLocalPos int) int
	MoveRight(ctx Context, e model.Element, viewLocalPos int) int
	ConvertPos(ctx Context, e model.Element, viewLocalPos int) (ViewLocalX int, ViewLocalY int)
	ConvertModel(ctx Context, e model.Element, viewLocalPos int) CharacterReference
}
