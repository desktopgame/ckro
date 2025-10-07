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

	MoveLength(ctx Context, textLayout *TextLayout) int
	MoveUp(ctx Context, textLayout *TextLayout, viewLocalPos int) int
	MoveDown(ctx Context, textLayout *TextLayout, viewLocalPos int) int
	MoveLeft(ctx Context, textLayout *TextLayout, viewLocalPos int) int
	MoveRight(ctx Context, textLayout *TextLayout, viewLocalPos int) int
	ConvertPos(ctx Context, textLayout *TextLayout, viewLocalPos int) (ViewLocalX int, ViewLocalY int)
	ConvertModel(ctx Context, textLayout *TextLayout, viewLocalPos int) CharacterReference

	// バイト座標をローカルビュー座標に変換する
	// 飾りと重なるときは内側に寄せる
	ConvertViewLocalPos(ctx Context, textLayout *TextLayout, bytePos model.Position) int

	// 指定のローカルビュー位置における改行を行の前に移動するなら true を返す
	ShouldBeforeInsertionNewLineOnLineBegin(ctx Context, textLayout *TextLayout, viewLocalPos int) bool

	// 指定のローカルビュー位置における削除で行全体を削除するべきなら true を返す
	ShouldRemoveWithLine(ctx Context, textLayout *TextLayout, viewLocalPos int) (int, bool)

	// 指定のローカルビュー位置における削除で指定位置以降の列をすべて削除するべきなら true を返す
	ShouldRemoveWithSpecifiedColumnAfter(ctx Context, textLayout *TextLayout, viewLocalPos int) (model.Position, bool)

	// 指定のローカルビュー位置における削除で指定範囲をすべて削除するべきなら true を返す
	ShouldRemoveWithSpecifiedRangeLines(ctx Context, textLayout *TextLayout, viewLocalPos int) (model.Range, bool)

	// 指定のローカルビュー位置における削除で指定範囲をすべて削除するべきなら true を返す
	ShouldRemoveWithSpecifiedRangeColumns(ctx Context, textLayout *TextLayout, viewLocalPos int) (model.Range, bool)

	// このビューの直後の削除が直前のこのビューに食い込むなら true を返す
	ShouldRemoveLastCharacter(ctx Context, textLayout *TextLayout, viewLocalPos int) (model.Element, bool)
}
