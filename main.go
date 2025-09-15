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

	hbox := tui.Box{}
	hbox.Init(tui.Horizontal)

	vbox := tui.Box{}
	vbox.Init(tui.Vertical)

	tree := tui.Tile{}
	tree.Init()
	tree.MinimumWidth = 20
	tree.FlexibleHeight = true
	// tree.TextBox.TextFrame()

	sep_tree := tui.Tile{}
	sep_tree.Init()
	sep_tree.MinimumWidth = 1
	sep_tree.FlexibleHeight = true
	sep_tree.TextBox.TextVertical()

	textArea := tui.Tile{}
	textArea.Init()
	textArea.MinimumWidth = 3
	textArea.FlexibleWidth = true
	textArea.FlexibleHeight = true
	// textArea.TextBox.TextFrame()

	sep_textArea := tui.Tile{}
	sep_textArea.Init()
	sep_textArea.MinimumHeight = 1
	sep_textArea.FlexibleWidth = true
	sep_textArea.TextBox.TextHorizontal()

	modeline := tui.Tile{}
	modeline.Init()
	modeline.FlexibleWidth = true
	modeline.MinimumHeight = 3
	// modeline.TextBox.TextFrame()

	sep_modeline := tui.Tile{}
	sep_modeline.Init()
	sep_modeline.MinimumHeight = 1
	sep_modeline.FlexibleWidth = true
	sep_modeline.TextBox.TextHorizontal()

	minibuffer := tui.Tile{}
	minibuffer.Init()
	minibuffer.FlexibleWidth = true
	minibuffer.MinimumHeight = 3
	// minibuffer.TextBox.TextFrame()

	vbox.Controls = append(vbox.Controls, &textArea)
	vbox.Controls = append(vbox.Controls, &sep_textArea)
	vbox.Controls = append(vbox.Controls, &modeline)
	vbox.Controls = append(vbox.Controls, &sep_modeline)
	vbox.Controls = append(vbox.Controls, &minibuffer)

	hbox.Controls = append(hbox.Controls, &tree)
	hbox.Controls = append(hbox.Controls, &sep_tree)
	hbox.Controls = append(hbox.Controls, &vbox)

	w, h := s.Size()
	hbox.Layout(w, h)

	s.Show()

	// inputBuffer := []rune{}

	for {
		s.Clear()

		tree.TextBox.Draw(s)
		sep_tree.TextBox.Draw(s)
		sep_textArea.TextBox.Draw(s)
		textArea.TextBox.Draw(s)
		modeline.TextBox.Draw(s)
		sep_modeline.TextBox.Draw(s)
		minibuffer.TextBox.Draw(s)

		s.Show()

		ev := s.PollEvent()
		switch e := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
			w, h = s.Size()
			hbox.Layout(w, h)

			// tree.TextBox.TextFrame()
			// textArea.TextBox.TextFrame()
			// modeline.TextBox.TextFrame()
			// minibuffer.TextBox.TextFrame()
			sep_tree.TextBox.TextVertical()
			sep_textArea.TextBox.TextHorizontal()
			sep_modeline.TextBox.TextHorizontal()
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
