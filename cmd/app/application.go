package main

import (
	"log"
	"os"

	"github.com/desktopgame/ckro/internal/tui"
	"github.com/desktopgame/ckro/internal/tui/controls"
	"github.com/desktopgame/ckro/internal/tui/presenter"
	"github.com/gdamore/tcell/v2"
)

type Application struct {
	screen         tcell.Screen
	commandPalette tui.Control
	window         tui.Window
	width          int
	height         int
}

func (app *Application) Init() {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}
	if err = s.Init(); err != nil {
		log.Fatal(err)
	}

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
	commands := []controls.Command{
		&controls.DelegateCommand{
			Label: "File: Open",
			Func: func(cp *controls.CommandPalette, stackable tui.Stackable) {
				// 現在のディレクトリを取得
				currentDir, err := os.Getwd()
				if err != nil {
					currentDir = "."
				}

				// ファイルチューザーを作成
				fileChooser := controls.NewFileChooser(
					currentDir,
					func(stackable tui.Stackable, selectedFile string) {
						// ファイルが選択された時の処理
						content, err := os.ReadFile(selectedFile)
						if err != nil {
							log.Printf("Error reading file: %v", err)
							return
						}

						// テキストエリアにファイル内容を表示
						doc := textArea.TextBox.GetDocument()
						doc.Init()
						doc.InsertString(string(content))
						textArea.TextBox.CursorReset()

						// ファイルチューザーを閉じる
						stackable.Pop(0)
					},
					func(stackable tui.Stackable) {
						// キャンセル時の処理
						stackable.Pop(1)
					},
				)
				fileChooserUI := tui.WithCenter(tui.WithFrame(fileChooser), 80, 20)

				stackable.Push(tui.Layer{
					Control: fileChooserUI,
					OnPop: func(returnCode int) {
						if returnCode == 0 {
							stackable.Pop(-1)
						}
					},
				})
			},
		},
	}
	commandPalette := controls.NewCommandPalette(commands)
	app.commandPalette = tui.WithCenter(tui.WithFrame(commandPalette), 80, 20)

	w, h := s.Size()
	app.screen = s
	app.window.Init(s, w, h)
	app.window.Push(tui.Layer{
		Control: &vbox,
	})
	app.width = w
	app.height = h

}

func (app *Application) Run() {
	defer app.screen.Fini()
	app.screen.Show()

	for {
		app.screen.Clear()
		app.window.Blit(app.width, app.height)
		app.screen.Show()

		ev := app.screen.PollEvent()

		switch e := ev.(type) {
		case *tcell.EventResize:
			app.screen.Sync()
			w, h := e.Size()
			app.window.Resize(w, h)
			app.width = w
			app.height = h
		case *tcell.EventKey:
			if e.Rune() == 'p' && (e.Modifiers()&tcell.ModAlt != 0) {
				if app.window.GetLayerCount() == 1 {
					app.window.Push(tui.Layer{
						Control: app.commandPalette,
					})
				}
				continue
			}
			if e.Key() == tcell.KeyBacktab {
				app.window.FocusNext()
				continue
			}
			switch e.Key() {
			case tcell.KeyCtrlC:
				return
			}
			tuiEvent := tui.Event{}
			tuiEvent.Init(&app.window, ev)
			app.window.Handle(tuiEvent)
		}
	}
}
