package main

import (
	"log"
	"os"

	"github.com/desktopgame/ckro/internal/tui"
	"github.com/desktopgame/ckro/internal/tui/base"
	"github.com/desktopgame/ckro/internal/tui/controls"
)

func FileOpenCommand(app *Application) func(cp *controls.CommandPalette, stackable base.Stackable) {
	return func(cp *controls.CommandPalette, stackable tui.Stackable) {
		if app.modified {
			stackable.Pop(1)
			return
		}
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
				doc := app.textArea.TextBox.GetDocument()
				doc.Init()
				doc.InsertString(string(content))
				app.filePath = selectedFile
				app.modified = false
				app.textArea.TextBox.CursorReset()

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
}
