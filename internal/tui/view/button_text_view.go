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

func (b *ButtonTextView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer, x int, y int, localViewLine int) {
	// x := 0
	// y :=

	s := ""
	if localViewLine == 0 {
		s = "*---*"
	} else if localViewLine == 1 {
		s = "|ABC|"
	} else if localViewLine == 2 {
		s = "*---*"
	}
	runes := []rune(s)
	for i, r := range runes {
		renderer.SetContent(x+i, y, r, nil, tcell.StyleDefault)
	}
}

func (b *ButtonTextView) Width(textViewResolver TextViewResolver, textLayout *TextLayout, row int) int {
	return 5
}

func (b *ButtonTextView) Height(textViewResolver TextViewResolver, textLayout *TextLayout) int {
	return 3
}
