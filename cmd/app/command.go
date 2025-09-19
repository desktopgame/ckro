package main

import (
	"context"
	"log"
	"os"

	"github.com/desktopgame/ckro/internal/tui"
	"github.com/desktopgame/ckro/internal/tui/base"
	"github.com/desktopgame/ckro/internal/tui/controls"
)

func FileNewCommand(app *Application) func(runtime base.Runtime, cp *controls.CommandPalette) {
	return func(runtime base.Runtime, cp *controls.CommandPalette) {
		if app.modified {
			// 確認ダイアログを表示
			confirmDialog := controls.NewConfirmationDialog(
				"Unsaved Changes",
				"You have unsaved changes.\nDo you want to save before creating a new file?",
				func(runtime base.Runtime) {
					// Yesが選択された場合 - 保存してからファイルを開く
					// TODO: 保存処理を実装
					log.Println("Save and open new file")
					runtime.Pop(-1) // ダイアログを閉じる
					app.newFile()
				},
				func(runtime base.Runtime) {
					// Noが選択された場合 - 保存せずにファイルを開く
					runtime.Pop(-1) // ダイアログを閉じる
					app.newFile()
				},
			)
			confirmDialogUI := tui.WithCenter(tui.WithFrame(confirmDialog), 60, 15)

			runtime.Push(tui.Layer{
				Control: confirmDialogUI,
				OnPop: func(returnCode int) {
					runtime.Pop(-1)
				},
			})
			return
		}
		app.newFile()
		runtime.Pop(-1)
	}
}

func FileOpenCommand(app *Application) func(runtime base.Runtime, cp *controls.CommandPalette) {
	return func(runtime base.Runtime, cp *controls.CommandPalette) {
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
					openFileChooser(app, runtime)
				},
				func(runtime base.Runtime) {
					// Noが選択された場合 - 保存せずにファイルを開く
					runtime.Pop(-1) // ダイアログを閉じる
					openFileChooser(app, runtime)
				},
			)
			confirmDialogUI := tui.WithCenter(tui.WithFrame(confirmDialog), 60, 15)

			runtime.Push(tui.Layer{
				Control: confirmDialogUI,
			})
			return
		}

		openFileChooser(app, runtime)
	}
}

func FileSaveCommand(app *Application) func(runtime base.Runtime, cp *controls.CommandPalette) {
	return func(runtime base.Runtime, cp *controls.CommandPalette) {
		if app.filePath != "" {
			app.saveFile()
			runtime.Pop(-1)
			return
		}
		// ファイルパスが設定されていない場合は「名前を付けて保存」と同じ動作
		showSaveAsDialog(app, runtime)
	}
}

func FileSaveAsCommand(app *Application) func(runtime base.Runtime, cp *controls.CommandPalette) {
	return func(runtime base.Runtime, cp *controls.CommandPalette) {
		showSaveAsDialog(app, runtime)
	}
}

func ChatMessage(app *Application) func(runtime base.Runtime, cp *controls.CommandPalette) {
	return func(runtime base.Runtime, cp *controls.CommandPalette) {

		inputDialog := controls.NewInputDialog(
			"Chat",
			"Enter prompt:",
			"",
			func(runtime base.Runtime, prompt string) {
				// OKが選択された場合
				if prompt != "" {
					go func() {
						message, err := app.chatManager.Post(context.TODO(), prompt)
						if err == nil {
							// 確認ダイアログを表示
							confirmDialog := controls.NewConfirmationDialog(
								"Response",
								message,
								func(runtime base.Runtime) {
									runtime.Pop(-1) // ダイアログを閉じる
								},
								func(runtime base.Runtime) {
									// Noが選択された場合 - 保存せずにファイルを開く
									runtime.Pop(-1) // ダイアログを閉じる
								},
							)
							confirmDialogUI := tui.WithCenter(tui.WithFrame(confirmDialog), 60, 15)

							runtime.Push(tui.Layer{
								Control: confirmDialogUI,
								OnPop: func(returnCode int) {
									runtime.Pop(0)
								},
							})
						} else {
							runtime.Pop(0) // ダイアログを閉じる
						}
					}()
				} else {
					runtime.Pop(0) // ダイアログを閉じる
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
			OnPop: func(returnCode int) {
				if returnCode == 0 {
					runtime.Pop(-1) // コマンドパレットも閉じる
				}
			},
		})
	}
}

func showSaveAsDialog(app *Application, runtime base.Runtime) {
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
				app.saveFileAs(filename)
				runtime.Pop(0) // ダイアログを閉じる
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
		OnPop: func(returnCode int) {
			if returnCode == 0 {
				runtime.Pop(-1) // コマンドパレットも閉じる
			}
		},
	})
}

func openFileChooser(app *Application, runtime base.Runtime) {
	// 現在のディレクトリを取得
	currentDir, err := os.Getwd()
	if err != nil {
		currentDir = "."
	}

	// ファイルチューザーを作成
	fileChooser := controls.NewFileChooser(
		currentDir,
		func(runtime base.Runtime, selectedFile string) {
			app.openFile(selectedFile)

			// ファイルチューザーを閉じる
			runtime.Pop(0)
		},
		func(runtime base.Runtime) {
			// キャンセル時の処理
			runtime.Pop(1)
		},
	)
	fileChooserUI := tui.WithCenter(tui.WithFrame(fileChooser), 80, 20)

	runtime.Push(tui.Layer{
		Control: fileChooserUI,
		OnPop: func(returnCode int) {
			if returnCode == 0 {
				runtime.Pop(-1)
			}
		},
	})
}
