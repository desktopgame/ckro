package view

import "github.com/desktopgame/ckro/internal/tui/model"

type TextView interface {
	Layout(textViewResolver TextViewResolver, e model.Element, width int) *TextLayout
	Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer)
}
