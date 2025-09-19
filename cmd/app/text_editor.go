package main

import (
	"github.com/desktopgame/ckro/internal/tui"
	"github.com/desktopgame/ckro/internal/tui/presenter"
)

type TextEditor struct {
	TextArea  tui.Tile
	container tui.Box
}

func (t *TextEditor) Init(onModified func()) {
	t.TextArea = tui.Tile{}
	t.TextArea.Init()
	t.TextArea.MinimumWidth = 3
	t.TextArea.FlexibleWidth = true
	t.TextArea.FlexibleHeight = true
	t.TextArea.TextBox.ShowCursor = true
	t.TextArea.TextPresenter = &presenter.EditTextPresenter{
		OnModified: onModified,
	}
	// 行番号エリア
	lineNumbers := tui.Tile{}
	lineNumbers.Init()
	lineNumbers.MinimumWidth = 4 // 行番号の幅（桁数に応じて調整）
	lineNumbers.FlexibleHeight = true
	lineNumbers.TextPresenter = &presenter.LineNumberTextPresenter{
		TargetView: t.TextArea.TextBox, // テキストエリアを対象に設定
	}
	// lineNumbers.TextPresenter = &presenter.FrameTextPresenter{}

	gutter := tui.Tile{}
	gutter.Init()
	gutter.MinimumWidth = 1 // 行番号の幅（桁数に応じて調整）
	gutter.FlexibleHeight = true
	gutter.TextPresenter = &presenter.VerticalSeparatorTextPresenter{}

	scrollBar := tui.NewTile(&presenter.ScrollBarTextPresenter{
		TargetView: t.TextArea.TextBox,
	})
	scrollBar.MinimumWidth = 1
	scrollBar.FlexibleHeight = true

	// 水平レイアウトで行番号とテキストエリアを並べる
	t.container = tui.Box{}
	t.container.Init(tui.Horizontal)
	t.container.Controls = append(t.container.Controls, &lineNumbers)
	t.container.Controls = append(t.container.Controls, &gutter)
	t.container.Controls = append(t.container.Controls, &t.TextArea)
	t.container.Controls = append(t.container.Controls, scrollBar)
}

func (t *TextEditor) Traverse(fm *tui.FocusManager) {
	t.container.Traverse(fm)
}

func (t *TextEditor) Update() {
	t.container.Update()
}

func (t *TextEditor) Draw(g *tui.Graphics) {
	t.container.Draw(g)
}

func (t *TextEditor) MinimumSize(width int, height int) (Width int, Height int) {
	return t.container.MinimumSize(width, height)
}

func (t *TextEditor) Move(x int, y int) {
	t.container.Move(x, y)
}

func (t *TextEditor) Layout(w int, h int) {
	t.container.Layout(w, h)
}

func (t *TextEditor) IsFlexibleWidth() bool {
	return t.container.IsFlexibleWidth()
}

func (t *TextEditor) IsFlexibleHeight() bool {
	return t.container.IsFlexibleHeight()
}
