package tui

import (
	"github.com/desktopgame/ckro/internal/text"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

type TextBox struct {
	Document   *Document
	X          int
	Y          int
	Width      int
	Height     int
	ShowCursor bool
	ScrollX    int
	ScrollY    int
}

func (tb *TextBox) Init() {
	tb.Document = &Document{}
	tb.Document.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 20
	tb.Height = 6
	tb.ScrollX = 0
	tb.ScrollY = 0
}

func (tb *TextBox) WrappedLineCount() int {
	lc := 0
	buf := tb.Document.GetBuffer()

	for i := 0; i < buf.GetLineCount(); i++ {
		line := buf.GetLineAt(i)
		lineContent := line.GetContent()
		lineWidth := text.DisplayWidth(lineContent)

		if lineWidth <= tb.Width {
			lc++
		} else {
			lc += (lineWidth / tb.Width)

			if lineWidth%tb.Width > 0 {
				lc++
			}
		}
	}
	return lc
}

func (tb *TextBox) CursorPosition() (X int, Y int, Rune rune, Combine []rune) {
	// カーソルを表示
	buf := tb.Document.GetBuffer()
	cursorRow := tb.Document.GetCursorRow()
	cursorCol := tb.Document.GetCursorColumn()

	// Documentのカーソル位置を画面座標に変換
	var screenX int
	var currentRune rune = ' '
	var combining []rune

	if cursorRow < buf.GetLineCount() {
		// カーソル位置の文字を取得（空行や行末の場合はスペース）
		cursorLine := ""
		currentRow := tb.ScrollY
		substringFrom := 0
		substringTo := 0

		for i := tb.ScrollY; i < cursorRow; i++ {
			line := buf.GetLineAt(i).GetContent()
			cursorLine = line
			screenX = text.DisplayWidth(line)
			substringFrom = 0
			substringTo = len(cursorLine)

			currentRow++

			if screenX >= tb.Width {
				remain := screenX - tb.Width
				screenX = 0
				// cursorRow++
				currentRow++

				substringFrom = tb.Width

				for remain > 0 {
					consume := min(remain, tb.Width)
					screenX += consume
					remain -= consume
					substringFrom += consume

					if screenX >= tb.Width {
						screenX = 0
						// cursorRow++
						currentRow++
					}
				}
			}
		}
		cursorRow = currentRow
		lineSub := tb.Document.GetBuffer().GetLineAt(tb.Document.GetCursorRow()).GetContent()
		screenX = text.DisplayPos(lineSub, cursorCol)
		if screenX >= tb.Width {
			charOffset := 0
			charsTotal := 0
			for screenX >= tb.Width {
				chars := 0

				for chars < tb.Width {
					ch := text.GraphemeSubString(lineSub, charOffset, charOffset+1)
					chw := text.DisplayWidth(ch)
					chars += chw
					charOffset++

					if chars > tb.Width {
						chars -= chw
						charOffset--
						break
					}
				}

				cursorRow++
				charsTotal += chars
				screenX -= chars
			}
			screenX = text.DisplayPos(lineSub, cursorCol) - charsTotal
		}
		// screenX = cursorCol

		cursorLineRange := cursorLine[substringFrom:substringTo]
		if cursorCol >= text.GraphemeLength(cursorLineRange) {
			// 行末またはそれを超えた位置
			currentRune = ' '
			combining = nil
		} else {
			currentRune, combining = text.DisplayRunesAt(cursorLineRange, screenX)
		}
		// cursorCol = screenX

	} else {
		// 無効な行の場合
		screenX = 0
		currentRune = ' '
		combining = nil
	}

	return screenX, cursorRow, currentRune, combining
}

func (tb *TextBox) CursorUpdate() {
	_, cursor, _, _ := tb.CursorPosition()
	// cursor := tb.Document.GetCursorRow()
	// lc := tb.WrappedLineCount()

	startY := tb.ScrollY
	endY := startY + tb.Height

	if cursor >= endY {
		for cursor >= endY {
			tb.ScrollY++

			startY = tb.ScrollY
			endY = startY + tb.Height
		}
	} else if cursor <= startY {
		for cursor <= startY && cursor > 0 {
			tb.ScrollY--

			startY = tb.ScrollY
		}
	}
}

func (tb *TextBox) TextFrame() {
	tb.Document.Init()

	w := tb.Width
	h := tb.Height

	tb.Document.InsertString("*")
	for i := 0; i < w-2; i++ {
		tb.Document.InsertString("-")
	}
	tb.Document.InsertString("*")
	tb.Document.InsertLine()

	for i := 0; i < h-2; i++ {
		tb.Document.InsertString("|")
		for j := 0; j < w-2; j++ {
			tb.Document.InsertString(" ")
		}
		tb.Document.InsertString("|")
		tb.Document.InsertLine()
	}

	tb.Document.InsertString("*")
	for i := 0; i < w-2; i++ {
		tb.Document.InsertString("-")
	}
	tb.Document.InsertString("*")

	tb.Document.MoveReset()
}

func (tb *TextBox) TextClear() {
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

	startY := tb.ScrollY
	endY := min(startY+tb.Height, buf.GetLineCount())
	drawY := 0
	for i := startY; i < endY; i++ {
		line := buf.GetLineAt(i).GetContent()
		x := 0

		// unisegを使ってgrapheme clusterごとに処理
		clusters := text.GraphemeClusters(line)
		for _, cluster := range clusters {
			runes := []rune(cluster)

			if len(runes) > 0 {
				// 最初のruneをメインとして設定
				mainRune := runes[0]
				var combining []rune

				// 残りのruneをcombining charactersとして設定
				if len(runes) > 1 {
					combining = runes[1:]
				}
				width := runewidth.RuneWidth(mainRune)

				if x+width > tb.Width {
					drawY++
					x = 0
				}
				if drawY >= tb.Height {
					break
				}

				clip.SetContent(x, drawY, mainRune, combining, def)

				// 全角文字の場合、次のセルを空にする
				if width == 2 {
					x++
					clip.SetContent(x, drawY, 0, nil, def)
				}
			}
			x++
			if x > tb.Width {
				drawY++
				x = 0
			}
			if drawY >= tb.Height {
				break
			}
		}
		drawY++
	}

	if !tb.ShowCursor {
		return
	}

	screenX, cursorRow, currentRune, combining := tb.CursorPosition()

	// カーソル位置の文字を反転表示
	cursorStyle := def.Reverse(true)
	clip.SetContent(screenX, cursorRow-tb.ScrollY, currentRune, combining, cursorStyle)

	// 全角文字の場合、隣接するセルもカーソル表示
	if currentRune != ' ' {
		width := runewidth.RuneWidth(currentRune)
		if width == 2 {
			// 隣接するセルにもカーソルを表示（空文字で反転）
			clip.SetContent(screenX+1, cursorRow-tb.ScrollY, 0, nil, cursorStyle)
		}
	}

}
