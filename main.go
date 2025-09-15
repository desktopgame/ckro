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

	box2 := tui.Box{}
	box2.Init(tui.Vertical)

	t1 := tui.Tile{}
	t1.Init()
	t1.MinimumWidth = 10
	t1.TextBox.TextFrame()

	ta := tui.Tile{}
	ta.Init()
	ta.MinimumWidth = 3
	ta.FlexibleWidth = true
	ta.FlexibleHeigght = true
	ta.TextBox.TextFrame()

	modeline := tui.Tile{}
	modeline.Init()
	modeline.FlexibleWidth = true
	modeline.MinimumHeight = 6
	modeline.TextBox.TextFrame()

	box2.Controls = append(box2.Controls, &ta)
	box2.Controls = append(box2.Controls, &modeline)

	box.Controls = append(box.Controls, &t1)
	box.Controls = append(box.Controls, &box2)

	w, h := s.Size()
	box.Layout(w, h)

	s.Show()

	// inputBuffer := []rune{}

	for {
		s.Clear()

		t1.TextBox.Draw(s)
		ta.TextBox.Draw(s)
		modeline.TextBox.Draw(s)

		s.Show()

		ev := s.PollEvent()
		switch e := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
			w, h = s.Size()
			box.Layout(w, h)

			t1.TextBox.TextFrame()
			ta.TextBox.TextFrame()
			modeline.TextBox.TextFrame()
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
