package view

import "github.com/desktopgame/ckro/internal/tui/model"

type TextLayout struct {
	Element  model.Element
	Children []*TextLayout

	// この要素を描画するとき、親はこのオフセット分TranslateしたRendererを子に渡して描画を移譲する
	RelativeX int
	RelativeY int

	Width         int
	Height        int
	MinimumWidth  int
	MinimumHeight int
	WidthTable    []int
	Indent        int
}
