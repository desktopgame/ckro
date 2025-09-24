package view

import (
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/gdamore/tcell/v2"
)

type ButtonTextView struct {
}

func (b *ButtonTextView) Layout(textViewResolver TextViewResolver, e model.Element, width int) *TextLayout {
	return &TextLayout{
		Element:  e,
		Children: nil,
	}
}

func (b *ButtonTextView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	renderer.SetContent(0, 0, '*', nil, tcell.StyleDefault)
	renderer.SetContent(0, 0, '-', nil, tcell.StyleDefault)
	renderer.SetContent(0, 0, '-', nil, tcell.StyleDefault)
	renderer.SetContent(0, 0, '-', nil, tcell.StyleDefault)
	renderer.SetContent(0, 0, '*', nil, tcell.StyleDefault)
	renderer.SetContent(0, 1, '|', nil, tcell.StyleDefault)
	renderer.SetContent(0, 1, 'A', nil, tcell.StyleDefault)
	renderer.SetContent(0, 1, 'B', nil, tcell.StyleDefault)
	renderer.SetContent(0, 1, 'C', nil, tcell.StyleDefault)
	renderer.SetContent(0, 1, '|', nil, tcell.StyleDefault)
	renderer.SetContent(0, 2, '*', nil, tcell.StyleDefault)
	renderer.SetContent(0, 2, '-', nil, tcell.StyleDefault)
	renderer.SetContent(0, 2, '-', nil, tcell.StyleDefault)
	renderer.SetContent(0, 2, '-', nil, tcell.StyleDefault)
	renderer.SetContent(0, 2, '*', nil, tcell.StyleDefault)
}

func (b *ButtonTextView) Width(textViewResolver TextViewResolver, textLayout *TextLayout, row int) int {
	return 5
}

func (b *ButtonTextView) Height(textViewResolver TextViewResolver, textLayout *TextLayout) int {
	return 3
}
