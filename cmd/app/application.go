package main

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/desktopgame/ckro/internal/llm"
	"github.com/desktopgame/ckro/internal/tui"
	"github.com/desktopgame/ckro/internal/tui/base"
	"github.com/desktopgame/ckro/internal/tui/controls"
	"github.com/desktopgame/ckro/internal/tui/presenter"
	"github.com/gdamore/tcell/v2"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

type Application struct {
	screen         tcell.Screen
	chatManager    llm.ChatManager
	card           tui.Card
	textEdior      TextEditor
	modeLine       ModeLine
	filePath       string
	modified       bool
	commandPalette tui.Control
	window         tui.Window
	width          int
	height         int
}

func (app *Application) newFile() {
	app.filePath = ""
	app.modified = false
	doc := app.textEdior.TextArea.TextBox.GetDocument()
	doc.Init()
	app.textEdior.TextArea.TextBox.CursorReset()
}

func (app *Application) openFile(filePath string) error {
	// ファイルを読み込んでテキストエリアに表示
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	app.filePath = filePath
	app.modified = false

	// テキストエリアのドキュメントをクリアして新しい内容を設定
	doc := app.textEdior.TextArea.TextBox.GetDocument()
	doc.Init()
	doc.InsertString(string(content))
	app.textEdior.TextArea.TextBox.CursorReset()
	return nil
}

func (app *Application) saveFile() error {
	if app.filePath == "" {
		return errors.New("filePath is empty")
	}
	sb := strings.Builder{}
	buf := app.textEdior.TextArea.TextBox.GetDocument().GetBuffer()

	for i := 0; i < buf.GetLineCount(); i++ {
		sb.WriteString(buf.GetLineAt(i).GetContent())

		if i < buf.GetLineCount()-1 {
			sb.WriteRune('\n')
		}
	}

	err := os.WriteFile(app.filePath, []byte(sb.String()), 0644)
	if err == nil {
		app.modified = false
	}
	return err
}

func (app *Application) saveFileAs(filePath string) error {
	sb := strings.Builder{}
	buf := app.textEdior.TextArea.TextBox.GetDocument().GetBuffer()

	for i := 0; i < buf.GetLineCount(); i++ {
		sb.WriteString(buf.GetLineAt(i).GetContent())

		if i < buf.GetLineCount()-1 {
			sb.WriteRune('\n')
		}
	}

	err := os.WriteFile(filePath, []byte(sb.String()), 0644)
	if err == nil {
		app.filePath = filePath
		app.modified = false
	}
	return err
}

func (app *Application) doLayout() {
	app.window.Resize(app.width, app.height)
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

	app.textEdior.Init(func() {
		app.modified = true
	})
	app.modeLine.Init()

	tree := tui.Tile{}
	tree.Init()
	tree.MinimumWidth = 50
	tree.FlexibleHeight = true
	tree.TextPresenter = &presenter.TreeTextPresenter{
		RootDirectory: ".",
		OnFileOpen: func(filePath string) {
			if app.modified {
				// 確認ダイアログを表示
				confirmDialog := controls.NewConfirmationDialog(
					"Unsaved Changes",
					"You have unsaved changes.\nDo you want to save before opening a new file?",
					func(runtime base.Runtime) {
						// Yesが選択された場合 - 保存してからファイルを開く
						// TODO: 保存処理を実装
						log.Println("Save and open new file")
						runtime.Pop(-1) // ダイアログを閉じる
						app.openFile(filePath)
					},
					func(runtime base.Runtime) {
						// Noが選択された場合 - 保存せずにファイルを開く
						runtime.Pop(-1) // ダイアログを閉じる
						app.openFile(filePath)
					},
				)
				confirmDialogUI := tui.WithCenter(tui.WithFrame(confirmDialog), 60, 15)

				app.window.Push(tui.Layer{
					Control: confirmDialogUI,
				})
				return
			}
			app.openFile(filePath)
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

	editorWithSide := tui.Box{}
	editorWithSide.Init(tui.Horizontal)

	app.card.Controls = append(app.card.Controls, &tui.Blank{})

	editorWithSide.Controls = append(editorWithSide.Controls, &app.textEdior)
	editorWithSide.Controls = append(editorWithSide.Controls, &app.card)

	hbox.Controls = append(hbox.Controls, &tree)
	hbox.Controls = append(hbox.Controls, &treeSeparator)
	hbox.Controls = append(hbox.Controls, &editorWithSide)

	vbox.Controls = append(vbox.Controls, &hbox)
	vbox.Controls = append(vbox.Controls, &textAreaSeparator)
	vbox.Controls = append(vbox.Controls, &app.modeLine)
	vbox.Controls = append(vbox.Controls, &modelineSeparator)
	vbox.Controls = append(vbox.Controls, &minibuffer)

	// QuickCommandPaletteのテスト
	commands := []controls.Command{
		&controls.DelegateCommand{
			Label: "File: New",
			Func:  FileNewCommand(app),
		},
		&controls.DelegateCommand{
			Label: "File: Open",
			Func:  FileOpenCommand(app),
		},
		&controls.DelegateCommand{
			Label: "File: Save",
			Func:  FileSaveCommand(app),
		},
		&controls.DelegateCommand{
			Label: "File: Save As",
			Func:  FileSaveAsCommand(app),
		},
		&controls.DelegateCommand{
			Label: "Chat: Message",
			Func:  ChatMessage(app),
		},
		&controls.DelegateCommand{
			Label: "Debug: Card prev",
			Func:  DebugCardPrev(app),
		},
		&controls.DelegateCommand{
			Label: "Debug: Card next",
			Func:  DebugCardNext(app),
		},
	}
	commandPalette := controls.NewCommandPalette(commands)
	app.commandPalette = tui.WithCenter(tui.WithFrame(commandPalette), 80, 20)

	w, h := s.Size()
	app.screen = s
	app.filePath = ""
	app.modified = false
	app.window.Init(s, w, h)
	app.window.Push(tui.Layer{
		Control: &vbox,
	})
	app.width = w
	app.height = h

	client := openai.NewClient(
		option.WithAPIKey("lmstudio"),
		option.WithBaseURL("http://localhost:1234/v1"),
	)
	app.chatManager.Init(&client, "openai/gpt-oss-20b", "あなたは親切なアシスタントです。")
	app.chatManager.Setup()
}

func (app *Application) Run() {
	defer app.screen.Fini()
	app.screen.Show()

	for {
		app.screen.Clear()
		app.window.Blit(app.width, app.height)
		app.screen.Show()

		ev := app.screen.PollEvent()

		if intr, ok := ev.(*tcell.EventInterrupt); ok {
			if _, ok := intr.Data().(tui.RepaintMessage); ok {
				continue
			}
		}

		if app.window.DoInBackground() {
			continue
		}

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
			if e.Key() == tcell.KeyCtrlC {
				if app.window.GetLayerCount() == 1 {
					return
				}
			}
			tuiEvent := tui.Event{}
			tuiEvent.Init(&app.window, ev)
			app.window.Handle(tuiEvent)
		}
	}
}
