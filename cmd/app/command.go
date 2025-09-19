package main

import (
	"log"
	"os"

	"github.com/desktopgame/ckro/internal/tui"
	"github.com/desktopgame/ckro/internal/tui/base"
	"github.com/desktopgame/ckro/internal/tui/controls"
)

func FileNewCommand(app *Application) func(cp *controls.CommandPalette, stackable base.Stackable) {
	return func(cp *controls.CommandPalette, stackable tui.Stackable) {
		if app.modified {
			// 確認ダイアログを表示
			confirmDialog := controls.NewConfirmationDialog(
				"Unsaved Changes",
				"You have unsaved changes.\nDo you want to save before creating a new file?",
				func(stackable base.Stackable) {
					// Yesが選択された場合 - 保存してからファイルを開く
					// TODO: 保存処理を実装
					log.Println("Save and open new file")
					stackable.Pop(-1) // ダイアログを閉じる
					app.newFile()
				},
				func(stackable base.Stackable) {
					// Noが選択された場合 - 保存せずにファイルを開く
					stackable.Pop(-1) // ダイアログを閉じる
					app.newFile()
				},
			)
			confirmDialogUI := tui.WithCenter(tui.WithFrame(confirmDialog), 60, 15)

			stackable.Push(tui.Layer{
				Control: confirmDialogUI,
				OnPop: func(returnCode int) {
					stackable.Pop(-1)
				},
			})
			return
		}
		app.newFile()
		stackable.Pop(-1)
	}
}

func FileOpenCommand(app *Application) func(cp *controls.CommandPalette, stackable base.Stackable) {
	return func(cp *controls.CommandPalette, stackable tui.Stackable) {
		if app.modified {
			// 確認ダイアログを表示
			confirmDialog := controls.NewConfirmationDialog(
				"Unsaved Changes",
				"You have unsaved changes.\nDo you want to save before opening a new file?",
				func(stackable base.Stackable) {
					// Yesが選択された場合 - 保存してからファイルを開く
					// TODO: 保存処理を実装
					log.Println("Save and open new file")
					stackable.Pop(-1) // ダイアログを閉じる
					openFileChooser(app, stackable)
				},
				func(stackable base.Stackable) {
					// Noが選択された場合 - 保存せずにファイルを開く
					stackable.Pop(-1) // ダイアログを閉じる
					openFileChooser(app, stackable)
				},
			)
			confirmDialogUI := tui.WithCenter(tui.WithFrame(confirmDialog), 60, 15)

			stackable.Push(tui.Layer{
				Control: confirmDialogUI,
			})
			return
		}

		openFileChooser(app, stackable)
	}
}

func openFileChooser(app *Application, stackable tui.Stackable) {
	// 現在のディレクトリを取得
	currentDir, err := os.Getwd()
	if err != nil {
		currentDir = "."
	}

	// ファイルチューザーを作成
	fileChooser := controls.NewFileChooser(
		currentDir,
		func(stackable tui.Stackable, selectedFile string) {
			app.openFile(selectedFile)

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
}
