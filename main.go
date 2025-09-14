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
		buf := doc.GetBuffer()
		y := 0
		for i := 0; i < buf.GetLineCount(); i++ {
			for x, r := range buf.GetLineAt(i).GetContent() {
				s.SetContent(x, y, r, nil, def)
			}
			y++
		}

		ev := s.PollEvent()
		switch e := ev.(type) {
		case *tcell.EventResize:
			s.Sync() // リサイズ時に再同期
		case *tcell.EventKey:
			if e.Key() == tcell.KeyEscape || e.Key() == tcell.KeyCtrlC {
				return
			}
		}
	}
}
