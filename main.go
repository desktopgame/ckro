package main

import (
	"log"

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

	box := tui.Box{}
	box.Init(tui.Horizontal)

	t1 := tui.Tile{}
	t1.Init()
	t1.MinimumWidth = 10
	t1.TextBox.TextFrame()

	t2 := tui.Tile{}
	t2.Init()
	t2.MinimumWidth = 3
	t2.FlexibleWidth = true
	t2.TextBox.TextFrame()

	box.Controls = append(box.Controls, &t1)
	box.Controls = append(box.Controls, &t2)

	w, h := s.Size()
	box.Layout(w, h)

	s.Show()

	// inputBuffer := []rune{}

	for {
		s.Clear()

		t1.TextBox.Draw(s)
		t2.TextBox.Draw(s)

		s.Show()

		ev := s.PollEvent()
		switch e := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
			w, h = s.Size()
			box.Layout(w, h)

			t1.TextBox.TextFrame()
			t2.TextBox.TextFrame()
		case *tcell.EventKey:
			switch e.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				return
			case tcell.KeyUp:
				// tb.Document.MoveUp()
			case tcell.KeyDown:
				// tb.Document.MoveDown()
			case tcell.KeyLeft:
				// tb.Document.MoveLeft()
			case tcell.KeyRight:
				// tb.Document.MoveRight()
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				// tb.Document.RemoveChar()
			case tcell.KeyEnter:
				// tb.Document.InsertLine()
			case tcell.KeyRune:
				// inputBuffer = append(inputBuffer, e.Rune())
				// inputString := string(inputBuffer)
				// if text.GraphemeLength(inputString) == 1 {
				// 	tb.Document.InsertString(inputString)
				// 	inputBuffer = []rune{}
				// }
			}
		}

		// tb.CursorUpdate()
	}
}
