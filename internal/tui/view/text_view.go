package view

import "github.com/desktopgame/ckro/internal/tui/model"

type TextView interface {
	// x,y はこのビューを置くべき左上の座標
	// w,h はこのビューに与えられたサイズ
	Layout(textViewResolver TextViewResolver, e model.Element, x, y, w, h int) *TextLayout
	Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer)

	// width,height はこのビューに与えられたサイズ
	MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) (Width int, Height int)
}
