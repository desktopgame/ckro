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

	hbox := tui.Box{}
	hbox.Init(tui.Horizontal)

	vbox := tui.Box{}
	vbox.Init(tui.Vertical)

	tree := tui.Tile{}
	tree.Init()
	tree.MinimumWidth = 20
	tree.FlexibleHeight = true
	// tree.TextBox.TextFrame()

	treeSeparator := tui.Tile{}
	treeSeparator.Init()
	treeSeparator.MinimumWidth = 1
	treeSeparator.FlexibleHeight = true
	treeSeparator.TextBox.TextVertical()

	textArea := tui.Tile{}
	textArea.Init()
	textArea.MinimumWidth = 3
	textArea.FlexibleWidth = true
	textArea.FlexibleHeight = true
	textArea.TextBox.ShowCursor = true
	// textArea.TextBox.TextFrame()

	textAreaSeparator := tui.Tile{}
	textAreaSeparator.Init()
	textAreaSeparator.MinimumHeight = 1
	textAreaSeparator.FlexibleWidth = true
	textAreaSeparator.TextBox.TextHorizontal()

	modeline := tui.Tile{}
	modeline.Init()
	modeline.FlexibleWidth = true
	modeline.MinimumHeight = 3
	// modeline.TextBox.TextFrame()

	modelineSeparator := tui.Tile{}
	modelineSeparator.Init()
	modelineSeparator.MinimumHeight = 1
	modelineSeparator.FlexibleWidth = true
	modelineSeparator.TextBox.TextHorizontal()

	minibuffer := tui.Tile{}
	minibuffer.Init()
	minibuffer.FlexibleWidth = true
	minibuffer.MinimumHeight = 3
	// minibuffer.TextBox.TextFrame()

	vbox.Controls = append(vbox.Controls, &textArea)
	vbox.Controls = append(vbox.Controls, &textAreaSeparator)
	vbox.Controls = append(vbox.Controls, &modeline)
	vbox.Controls = append(vbox.Controls, &modelineSeparator)
	vbox.Controls = append(vbox.Controls, &minibuffer)

	hbox.Controls = append(hbox.Controls, &tree)
	hbox.Controls = append(hbox.Controls, &treeSeparator)
	hbox.Controls = append(hbox.Controls, &vbox)

	w, h := s.Size()
	hbox.Layout(w, h)

	s.Show()

	inputBuffer := []rune{}

	for {
		s.Clear()

		tree.TextBox.Draw(s)
		treeSeparator.TextBox.Draw(s)
		textAreaSeparator.TextBox.Draw(s)
		textArea.TextBox.Draw(s)
		modeline.TextBox.Draw(s)
		modelineSeparator.TextBox.Draw(s)
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
			treeSeparator.TextBox.TextVertical()
			textAreaSeparator.TextBox.TextHorizontal()
			modelineSeparator.TextBox.TextHorizontal()
		case *tcell.EventKey:
			switch e.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				return
			case tcell.KeyUp:
				textArea.TextBox.Document.MoveUp()
			case tcell.KeyDown:
				textArea.TextBox.Document.MoveDown()
			case tcell.KeyLeft:
				textArea.TextBox.Document.MoveLeft()
			case tcell.KeyRight:
				textArea.TextBox.Document.MoveRight()
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				textArea.TextBox.Document.RemoveChar()
			case tcell.KeyEnter:
				textArea.TextBox.Document.InsertLine()
			case tcell.KeyRune:
				inputBuffer = append(inputBuffer, e.Rune())
				inputString := string(inputBuffer)
				if text.GraphemeLength(inputString) == 1 {
					textArea.TextBox.Document.InsertString(inputString)
					inputBuffer = []rune{}
				}
			}
		}

		textArea.TextBox.CursorUpdate()
	}
}
