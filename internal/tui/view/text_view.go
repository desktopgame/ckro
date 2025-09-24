package view

import "github.com/desktopgame/ckro/internal/tui/model"

type TextView interface {
	Layout(textViewResolver TextViewResolver, e model.Element, width int) *TextLayout
	Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer)

	Width(textViewResolver TextViewResolver, textLayout *TextLayout, row int) int
	Height(textViewResolver TextViewResolver, textLayout *TextLayout) int
}
