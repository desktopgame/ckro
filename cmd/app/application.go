package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/desktopgame/ckro/internal/llm"
	"github.com/desktopgame/ckro/internal/tui"
	"github.com/desktopgame/ckro/internal/tui/base"
	"github.com/desktopgame/ckro/internal/tui/controls"
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/desktopgame/ckro/internal/tui/presenter"
	"github.com/gdamore/tcell/v2"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

type Application struct {
	screen         tcell.Screen
	card           tui.Card
	textEdior      TextEditor
	treePresenter  *presenter.TreeTextPresenter
	modeLine       ModeLine
	miniBuffer     MiniBuffer
	commandPalette tui.Control

	window tui.Window
	width  int
	height int

	filePath string
	modified bool

	chatManager    llm.ChatManager
	chatResponseId int

	vaultManager VaultManager
	config       Config
}

func (app *Application) newFile() {
	app.filePath = ""
	app.modified = false
	doc := app.textEdior.TextArea.TextBox.GetDocument()
	doc.Clear()
	app.textEdior.TextArea.TextBox.CursorReset()
}

func (app *Application) openFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	doc := app.textEdior.TextArea.TextBox.Document
	doc.ReplaceAll(file)

	app.filePath = filePath
	app.modified = false
	app.textEdior.TextArea.TextBox.MoveReset()
	return nil
}

func (app *Application) saveFile() error {
	if app.filePath == "" {
		return errors.New("filePath is empty")
	}
	sb := strings.Builder{}
	doc := app.textEdior.TextArea.TextBox.GetDocument()

	for i := 0; i < doc.GetLineCount(); i++ {
		sg := doc.Read(model.Range{
			StartPosition: model.Position{
				Row:    i,
				Column: 0,
			},
			EndPosition: model.Position{
				Row:    i,
				Column: doc.GetLineBytes(i),
			},
		})
		sb.WriteString(sg.GetLine(0))

		if i < doc.GetLineCount()-1 {
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
	doc := app.textEdior.TextArea.TextBox.GetDocument()

	for i := 0; i < doc.GetLineCount(); i++ {
		sg := doc.Read(model.Range{
			StartPosition: model.Position{
				Row:    i,
				Column: 0,
			},
			EndPosition: model.Position{
				Row:    i,
				Column: doc.GetLineBytes(i),
			},
		})
		sb.WriteString(sg.GetLine(0))

		if i < doc.GetLineCount()-1 {
			sb.WriteRune('\n')
		}
	}

	err := os.WriteFile(filePath, []byte(sb.String()), 0644)
	if err == nil {
		app.filePath = filePath
		app.modified = false
		app.treePresenter.Reload()
	}
	return err
}

func showSaveAsDialogAndThenForTree(app *Application, runtime base.Runtime, callback func()) {
	// デフォルトのファイル名を設定
	defaultFileName := "untitled.txt"
	if app.filePath != "" {
		// 既存のファイルパスがある場合はそのファイル名を使用
		defaultFileName = app.filePath
	}

	// 入力ダイアログを作成
	inputDialog := controls.NewInputDialog(
		"Save As",
		"Enter filename:",
		defaultFileName,
		func(runtime base.Runtime, filename string) {
			// OKが選択された場合
			if filename != "" {
				err := app.saveFileAs(filename)
				if err != nil {
					log.Printf("Failed to save file: %v", err)
				}
				runtime.Pop(0) // ダイアログを閉じる
				callback()     // コールバック実行
			}
		},
		func(runtime base.Runtime) {
			// キャンセルが選択された場合
			runtime.Pop(1) // ダイアログを閉じる
		},
	)
	inputDialogUI := tui.WithCenter(tui.WithFrame(inputDialog), 70, 12)

	runtime.Push(tui.Layer{
		Control: inputDialogUI,
	})
}

func (app *Application) doLayout() {
	app.window.Resize(app.width, app.height)
}

func (app *Application) initView() {
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
	app.miniBuffer.Init()

	tree := tui.Tile{}
	tree.Init()
	tree.MinimumWidth = 50
	tree.FlexibleHeight = true
	app.treePresenter = &presenter.TreeTextPresenter{
		RootDirectory: ".",
		OnFileOpen: func(filePath string) {
			if app.modified {
				// 確認ダイアログを表示
				confirmDialog := controls.NewConfirmationDialog(
					"Unsaved Changes",
					"You have unsaved changes.\nDo you want to save before opening a new file?",
					func(runtime base.Runtime) {
						// Yesが選択された場合 - 保存してからファイルを開く
						runtime.Pop(-1) // ダイアログを閉じる
						if app.filePath == "" {
							// 名前のないファイルの場合は「名前を付けて保存」
							showSaveAsDialogAndThenForTree(app, runtime, func() {
								app.openFile(filePath)
							})
						} else {
							// 既存ファイルの場合は直接保存
							err := app.saveFile()
							if err != nil {
								log.Printf("Failed to save file: %v", err)
							}
							app.openFile(filePath)
						}
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
	tree.TextPresenter = app.treePresenter

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
	vbox.Controls = append(vbox.Controls, &app.miniBuffer)

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
			Label: "File: Delete",
			Func:  FileDeleteCommand(app),
		},
		&controls.DelegateCommand{
			Label: "File: Open folder",
			Func:  FileOpenFolderCommand(app),
		},
		&controls.DelegateCommand{
			Label: "Vault: Init",
			Func:  VaultInitCommand(app),
		},
		&controls.DelegateCommand{
			Label: "Vault: Open",
			Func:  VaultOpenCommand(app),
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
}

func (app *Application) initSystem() {
	// Vaultの初期化
	app.vaultManager.Init()
	app.config.Update(&app.vaultManager)
	// MCP関連の初期化
	servers := map[string]*mcp.Server{}
	servers["time"] = NewTimeMcp()
	servers["file_system"] = NewFileSystemMcp()

	mcpClients := map[string]*llm.McpClient{}
	waitGroup := sync.WaitGroup{}
	for name, server := range servers {
		mcpClient := llm.McpClient{}
		mcpClient.Init()

		waitGroup.Add(1)
		go func() {
			mcpClient.ConnectLocal(context.Background(), server)
			mcpClients[name] = &mcpClient
			waitGroup.Done()
		}()
	}

	waitGroup.Wait()

	client := openai.NewClient(
		option.WithAPIKey(app.config.ApiKey),
		option.WithBaseURL(app.config.BaseUrl),
	)
	app.chatManager.Init(&client, app.config.Model, app.config.SystemPrompt, mcpClients)
	app.chatManager.Setup()
}

func (app *Application) Init() {
	app.initView()
	app.initSystem()
}

func (app *Application) loopMiniBuffer() {
	for {
		s, ok := <-app.miniBuffer.ch

		if !ok {
			break
		}

		log.Println(s)

		tb := app.textEdior.TextArea.TextBox
		sb := strings.Builder{}
		sb.WriteString("\n")
		sb.WriteString("{{{\n")
		sb.WriteString("USER:\n")
		sb.WriteString(s)
		sb.WriteString("\n")
		sb.WriteString("}}}\n")

		tb.Document.InsertString(
			tb.Document.GetLineCount()-1,
			tb.Document.GetLineBytes(tb.Document.GetLineCount()-1),
			sb.String(),
		)
		app.modified = true
		app.window.Repaint()

		app.miniBuffer.ReadOnly()
		pipe := make(chan llm.Event)
		done := make(chan struct{})
		go app.chatManager.Post(context.Background(), s, pipe)
		go func() {
			defer close(pipe)
			defer close(done)

			for {
				ev := <-pipe

				if conf, ok := ev.(*llm.ConfirmEvent); ok {
					message := fmt.Sprintf("want to use tool of `%s`, are you ok? [y/n]", conf.ToolName)
					app.modeLine.Text(message)

					app.miniBuffer.Editable()
					app.window.Repaint()

					s = <-app.miniBuffer.ch

					app.modeLine.Text("")
					app.miniBuffer.ReadOnly()
					app.window.Repaint()
					conf.Approve = s == "y" || s == "Y"
				}

				ev.Consume(context.Background())

				tb.MoveTextEnd()
				{
					if log, ok := ev.(*llm.LogEvent); ok {
						if log.Body.OfAssistant != nil {
							sb = strings.Builder{}
							sb.WriteString("{{{\n")
							sb.WriteString("BOT:\n")
							sb.WriteString(log.Body.OfAssistant.Content.OfString.Value)
							sb.WriteString("\n")
							sb.WriteString("}}}\n")

							tb.InsertString(sb.String())
						} else if log.Body.OfTool != nil {
							sb = strings.Builder{}
							sb.WriteString("{{{\n")
							sb.WriteString("TOOL:\n")
							sb.WriteString(log.Body.OfTool.Content.OfString.Value)
							sb.WriteString("\n")
							sb.WriteString("}}}\n")

							tb.InsertString(sb.String())
						} else if log.Body.OfSystem != nil {
							sb = strings.Builder{}
							sb.WriteString("{{{\n")
							sb.WriteString("SYSTEM:\n")
							sb.WriteString(log.Body.OfSystem.Content.OfString.Value)
							sb.WriteString("\n")
							sb.WriteString("}}}\n")

							tb.InsertString(sb.String())
						}
					}
					tb.MoveRight()
					tb.CursorUpdate()
				}

				if _, ok := ev.(*llm.MessageEvent); ok {
					break
				}

				if e, ok := ev.(*llm.ErrorEvent); ok {
					tb.MoveTextEnd()
					tb.InsertString(e.GetError().Error())
					tb.MoveRight()
					tb.CursorUpdate()
					break
				}
			}

			app.miniBuffer.Editable()
			app.window.Repaint()
		}()

		<-done
		app.chatResponseId++
	}
}

func (app *Application) Run() {
	defer app.screen.Fini()
	app.screen.Show()

	go app.loopMiniBuffer()

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

func (app *Application) Close() {
	app.miniBuffer.Close()
}
