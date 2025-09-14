package tui

import (
	"github.com/desktopgame/ckro/internal/text"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/uniseg"
)

type TextBox struct {
	Document *Document
	X        int
	Y        int
	Width    int
	Height   int
	HasFocus bool
}

func (tb *TextBox) Init() {
	tb.Document = &Document{}
	tb.Document.Init()
}

func (tb *TextBox) Draw(s tcell.Screen) {
	// バッファの内容を描画
	clip := Clip{
		Screen: s,
		X:      tb.X,
		Y:      tb.Y,
		Width:  tb.Width,
		Height: tb.Height,
	}
	def := tcell.StyleDefault
	buf := tb.Document.GetBuffer()
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

				clip.SetContent(x, y, mainRune, combining, def)

				// 全角文字の場合、次のセルを空にする
				width := runewidth.RuneWidth(mainRune)
				if width == 2 {
					x++
					if x < 80 { // 画面幅の制限内で
						clip.SetContent(x, y, 0, nil, def)
					}
				}
			}
			x++
		}
		y++
	}

	if !tb.HasFocus {
		return
	}

	// カーソルを表示
	cursorRow := tb.Document.GetCursorRow()
	cursorCol := tb.Document.GetCursorColumn()

	// Documentのカーソル位置を画面座標に変換
	var screenX int
	var currentRune rune = ' '
	var combining []rune

	if cursorRow < buf.GetLineCount() {
		line := buf.GetLineAt(cursorRow).GetContent()
		screenX = text.DisplayPos(line, cursorCol)

		// カーソル位置の文字を取得（空行や行末の場合はスペース）
		if cursorCol >= text.GraphemeLength(line) {
			// 行末またはそれを超えた位置
			currentRune = ' '
			combining = nil
		} else {
			currentRune, combining = text.DisplayRunesAt(line, screenX)
		}
	} else {
		// 無効な行の場合
		screenX = 0
		currentRune = ' '
		combining = nil
	}

	// カーソル位置の文字を反転表示
	cursorStyle := def.Reverse(true)
	clip.SetContent(screenX, cursorRow, currentRune, combining, cursorStyle)

	// 全角文字の場合、隣接するセルもカーソル表示
	if currentRune != ' ' {
		width := runewidth.RuneWidth(currentRune)
		if width == 2 {
			// 隣接するセルにもカーソルを表示（空文字で反転）
			clip.SetContent(screenX+1, cursorRow, 0, nil, cursorStyle)
		}
	}

}
