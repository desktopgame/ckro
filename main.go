package main

import (
	"log"
	"os"

	"github.com/desktopgame/ckro/internal/tui"
	"github.com/desktopgame/ckro/internal/tui/presenter"
	"github.com/gdamore/tcell/v2"
)

func main() {

	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}
	if err = s.Init(); err != nil {
		log.Fatal(err)
	}
	defer s.Fini()

	def := tcell.StyleDefault
	s.SetStyle(def)
	s.Clear()

	hbox := tui.Box{}
	hbox.Init(tui.Horizontal)

	vbox := tui.Box{}
	vbox.Init(tui.Vertical)

	textArea := tui.Tile{}
	textArea.Init()
	textArea.MinimumWidth = 3
	textArea.FlexibleWidth = true
	textArea.FlexibleHeight = true
	textArea.TextBox.ShowCursor = true
	textArea.TextPresenter = &presenter.EditTextPresenter{}
	// 行番号エリア
	lineNumbers := tui.Tile{}
	lineNumbers.Init()
	lineNumbers.MinimumWidth = 4 // 行番号の幅（桁数に応じて調整）
	lineNumbers.FlexibleHeight = true
	lineNumbers.TextPresenter = &presenter.LineNumberTextPresenter{
		TargetView: textArea.TextBox, // テキストエリアを対象に設定
	}
	// lineNumbers.TextPresenter = &presenter.FrameTextPresenter{}

	gutter := tui.Tile{}
	gutter.Init()
	gutter.MinimumWidth = 1 // 行番号の幅（桁数に応じて調整）
	gutter.FlexibleHeight = true
	gutter.TextPresenter = &presenter.VerticalSeparatorTextPresenter{}

	// 水平レイアウトで行番号とテキストエリアを並べる
	editorBox := tui.Box{}
	editorBox.Init(tui.Horizontal)
	editorBox.Controls = append(editorBox.Controls, &lineNumbers)
	editorBox.Controls = append(editorBox.Controls, &gutter)
	editorBox.Controls = append(editorBox.Controls, &textArea)

	tree := tui.Tile{}
	tree.Init()
	tree.MinimumWidth = 50
	tree.FlexibleHeight = true
	tree.TextPresenter = &presenter.TreeTextPresenter{
		RootDirectory: ".",
		OnFileOpen: func(filePath string) {
			// ファイルを読み込んでテキストエリアに表示
			content, err := os.ReadFile(filePath)
			if err != nil {
				return
			}

			// テキストエリアのドキュメントをクリアして新しい内容を設定
			doc := textArea.TextBox.GetDocument()
			doc.Init()
			doc.InsertString(string(content))
			textArea.TextBox.CursorReset()
		},
	}

	treeSeparator := tui.Tile{}
	treeSeparator.Init()
	treeSeparator.MinimumWidth = 1
	treeSeparator.FlexibleHeight = true
	treeSeparator.TextPresenter = &presenter.VerticalSeparatorTextPresenter{}

	textAreaSeparator := tui.Tile{}
	textAreaSeparator.Init()
	textAreaSeparator.MinimumHeight = 1
	textAreaSeparator.FlexibleWidth = true
	textAreaSeparator.TextPresenter = &presenter.HorizontalSeparatorTextPresenter{}

	modeline := tui.Tile{}
	modeline.Init()
	modeline.FlexibleWidth = true
	modeline.MinimumHeight = 1

	modelineSeparator := tui.Tile{}
	modelineSeparator.Init()
	modelineSeparator.MinimumHeight = 1
	modelineSeparator.FlexibleWidth = true
	modelineSeparator.TextPresenter = &presenter.HorizontalSeparatorTextPresenter{}

	minibuffer := tui.Tile{}
	minibuffer.Init()
	minibuffer.FlexibleWidth = true
	minibuffer.MinimumHeight = 1
	minibuffer.TextPresenter = &presenter.EditTextPresenter{}

	hbox.Controls = append(hbox.Controls, &tree)
	hbox.Controls = append(hbox.Controls, &treeSeparator)
	hbox.Controls = append(hbox.Controls, &editorBox)

	vbox.Controls = append(vbox.Controls, &hbox)
	vbox.Controls = append(vbox.Controls, &textAreaSeparator)
	vbox.Controls = append(vbox.Controls, &modeline)
	vbox.Controls = append(vbox.Controls, &modelineSeparator)
	vbox.Controls = append(vbox.Controls, &minibuffer)

	gl := tui.GridLayout{}
	gl.Init(2, 2)
	gl.Set(0, 0, 9, 9, false, false, &presenter.EditTextPresenter{})
	gl.Set(0, 1, 9, 9, false, false, &presenter.EditTextPresenter{})
	gl.Set(1, 0, 9, 9, false, false, &presenter.EditTextPresenter{})
	gl.Set(1, 1, 9, 9, false, false, &presenter.EditTextPresenter{})

	stack := tui.Stack{}
	stack.Init()
	stack.Layers = append(stack.Layers, &vbox)
	stack.Layers = append(stack.Layers, gl.Build())
	stack.Top = 0

	focusManager := tui.FocusManager{}
	stack.Traverse(&focusManager)
	focusManager.Grab()

	w, h := s.Size()
	stack.Layout(w, h)

	s.Show()

	for {
		s.Clear()

		mw, mh := stack.MinimumSize(w, h)

		if mw <= w && mh <= h {
			stack.Update()
			stack.Draw(s)
		}

		s.Show()

		ev := s.PollEvent()
		switch e := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
			w, h = e.Size()
			stack.Layout(w, h)
		case *tcell.EventKey:
			switch e.Key() {
			case tcell.KeyTAB:
				focusManager.FocusNext()
			case tcell.KeyEscape, tcell.KeyCtrlC:
				return
			}
			focusManager.Handle(ev)
		}
	}
}
