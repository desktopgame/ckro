package main

import (
	"log"

	"github.com/desktopgame/ckro/internal/tui"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/uniseg"
)

// Documentのカーソル位置（rune単位）を画面座標（grapheme cluster単位）に変換
func documentToScreenPos(line string, docColumn int) int {
	if docColumn <= 0 {
		return 0
	}

	runes := []rune(line)
	if docColumn >= len(runes) {
		return runewidth.StringWidth(line)
	}

	targetStr := string(runes[:docColumn])
	return runewidth.StringWidth(targetStr)
}

// 画面座標（grapheme cluster単位）でのカーソル位置の文字を取得
func getCharAtScreenPos(line string, screenX int) (rune, []rune) {
	currentX := 0
	gr := uniseg.NewGraphemes(line)
	for gr.Next() {
		if currentX == screenX {
			cluster := gr.Str()
			runes := []rune(cluster)
			if len(runes) > 0 {
				mainRune := runes[0]
				var combining []rune
				if len(runes) > 1 {
					combining = runes[1:]
				}
				return mainRune, combining
			}
		}
		currentX++
	}
	return ' ', nil
}

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

	// 画面に一言
	msg :=
		`
Hello, world1
👨‍👩‍👧‍👦
Hello, world2
`

	doc := tui.Document{}
	doc.Init()
	doc.InsertString(msg)

	s.Show()

	// イベントループ
	for {
		// 画面をクリア
		s.Clear()

		// バッファの内容を描画
		buf := doc.GetBuffer()
		y := 0
		for i := 0; i < buf.GetLineCount(); i++ {
			line := buf.GetLineAt(i).GetContent()
			x := 0

			// unisegを使ってgrapheme clusterごとに処理
			gr := uniseg.NewGraphemes(line)
			for gr.Next() {
				cluster := gr.Str()
				runes := []rune(cluster)

				if len(runes) > 0 {
					// 最初のruneをメインとして設定
					mainRune := runes[0]
					var combining []rune

					// 残りのruneをcombining charactersとして設定
					if len(runes) > 1 {
						combining = runes[1:]
					}

					s.SetContent(x, y, mainRune, combining, def)
				}
				x++
			}
			y++
		}

		// カーソルを表示
		cursorRow := doc.GetCursorRow()
		cursorCol := doc.GetCursorColumn()

		// Documentのカーソル位置を画面座標に変換
		var screenX int
		var currentRune rune = ' '
		var combining []rune

		if cursorRow < buf.GetLineCount() {
			line := buf.GetLineAt(cursorRow).GetContent()
			screenX = documentToScreenPos(line, cursorCol)
			currentRune, combining = getCharAtScreenPos(line, screenX)
		}

		// カーソル位置の文字を反転表示
		cursorStyle := def.Reverse(true)
		s.SetContent(screenX, cursorRow, currentRune, combining, cursorStyle)

		s.Show()

		ev := s.PollEvent()
		switch e := ev.(type) {
		case *tcell.EventResize:
			s.Sync() // リサイズ時に再同期
		case *tcell.EventKey:
			switch e.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				return
			case tcell.KeyUp:
				doc.MoveUp()
			case tcell.KeyDown:
				doc.MoveDown()
			case tcell.KeyLeft:
				doc.MoveLeft()
			case tcell.KeyRight:
				doc.MoveRight()
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				doc.RemoveChar()
			case tcell.KeyEnter:
				doc.InsertLine()
			case tcell.KeyRune:
				// 通常の文字入力
				doc.InsertString(string(e.Rune()))
			}
		}
	}
}
