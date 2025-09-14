package main

import (
	"log"

	"github.com/desktopgame/ckro/internal/text"
	"github.com/desktopgame/ckro/internal/tui"
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

	// 画面に一言
	msg :=
		`
Hello, world1
👨‍👩‍👧‍👦
Hello, world2
あいうえお
`

	tb := tui.TextBox{}
	tb.Init()
	tb.X = 0
	tb.Y = 0
	tb.Width = 20
	tb.Height = 6
	tb.Document.InsertString(msg)
	tb.HasFocus = true

	s.Show()

	inputBuffer := []rune{}

	for {
		s.Clear()

		tb.Draw(s)

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
				tb.Document.MoveUp()
			case tcell.KeyDown:
				tb.Document.MoveDown()
			case tcell.KeyLeft:
				tb.Document.MoveLeft()
			case tcell.KeyRight:
				tb.Document.MoveRight()
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				tb.Document.RemoveChar()
			case tcell.KeyEnter:
				tb.Document.InsertLine()
			case tcell.KeyRune:
				// 通常の文字入力
				inputBuffer = append(inputBuffer, e.Rune())
				inputString := string(inputBuffer)
				if text.GraphemeLength(inputString) == 1 {
					tb.Document.InsertString(inputString)
					inputBuffer = []rune{}
				}
				//doc.InsertString("あ")
			}
		}

		tb.UpdateCursor()
	}
}
