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

	tb1 := tui.TextBox{}
	tb1.Init()
	tb1.X = 0
	tb1.Y = 0
	tb1.Width = 20
	tb1.Height = 5
	tb1.Document.InsertString(msg)
	tb1.HasFocus = true

	tb2 := tui.TextBox{}
	tb2.Init()
	tb2.X = 0
	tb2.Y = 10
	tb2.Width = 20
	tb2.Height = 5
	tb2.Document.InsertString(msg)
	tb2.HasFocus = false

	tb := &tb1
	tbi := 1

	s.Show()

	inputBuffer := []rune{}

	for {
		s.Clear()

		tb1.Draw(s)
		tb2.Draw(s)

		s.Show()

		ev := s.PollEvent()
		switch e := ev.(type) {
		case *tcell.EventResize:
			s.Sync() // リサイズ時に再同期
		case *tcell.EventKey:
			switch e.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				return
			case tcell.KeyTAB:
				switch tbi {
				case 1:
					tb1.HasFocus = false
					tb2.HasFocus = true
					tb = &tb2
					tbi = 2
				case 2:
					tb2.HasFocus = false
					tb1.HasFocus = true
					tb = &tb1
					tbi = 1
				}
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
	}
}
