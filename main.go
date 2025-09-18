package main

import (
	"log"
	"os"

	"github.com/desktopgame/ckro/internal/tui"
	"github.com/desktopgame/ckro/internal/tui/controls"
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

	scrollBar := tui.NewTile(&presenter.ScrollBarTextPresenter{
		TargetView: textArea.TextBox,
	})
	scrollBar.MinimumWidth = 1
	scrollBar.FlexibleHeight = true

	// 水平レイアウトで行番号とテキストエリアを並べる
	editorBox := tui.Box{}
	editorBox.Init(tui.Horizontal)
	editorBox.Controls = append(editorBox.Controls, &lineNumbers)
	editorBox.Controls = append(editorBox.Controls, &gutter)
	editorBox.Controls = append(editorBox.Controls, &textArea)
	editorBox.Controls = append(editorBox.Controls, scrollBar)

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

	ctrlLine := tui.Box{}
	ctrlLine.Init(tui.Horizontal)

	ctrl1 := tui.Tile{}
	ctrl1.Init()
	ctrl1.MinimumWidth = 10
	ctrl1.MinimumHeight = 10
	ctrl1.TextPresenter = &presenter.LabelTextPresenter{
		Text:        "Label",
		AlignCenter: true,
	}

	ctrl2 := tui.Tile{}
	ctrl2.Init()
	ctrl2.FlexibleWidth = true
	ctrl2.MinimumHeight = 10
	ctrl2.TextPresenter = &presenter.EditTextPresenter{}

	frame := tui.Frame{}
	frame.Control = &ctrl2

	ctrlLine.Controls = append(ctrlLine.Controls, &ctrl1)
	ctrlLine.Controls = append(ctrlLine.Controls, &frame)

	// QuickCommandPaletteのテスト
	commands := []string{
		"File: Open",
		"File: Save",
		"File: Save As",
		"Edit: Copy",
		"Edit: Paste",
		"Edit: Find",
		"View: Toggle Sidebar",
		"Help: About",
	}
	commandPalette := controls.NewCommandPalette(commands, func(command string) {
		log.Printf("Executed command: %s", command)
	})

	commandPaletteUI := tui.WithCenter(tui.WithFrame(commandPalette), 80, 20)

	window := tui.Window{}
	window.Push(&vbox)

	w, h := s.Size()
	window.Resize(w, h)

	s.Show()

	for {
		s.Clear()

		window.Frame(s, w, h)

		s.Show()

		ev := s.PollEvent()
		switch e := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
			w, h = e.Size()
			window.Resize(w, h)
		case *tcell.EventKey:
			if e.Rune() == 'p' && (e.Modifiers()&tcell.ModAlt != 0) {
				if window.Top() == commandPaletteUI {
					window.Pop()
				} else if window.GetLayerCount() == 1 {
					window.Push(commandPaletteUI)
				}
				continue
			}
			if e.Key() == tcell.KeyBacktab {
				window.FocusNext()
				continue
			}
			switch e.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				return
			}
			window.Handle(ev)
		}
	}
}
