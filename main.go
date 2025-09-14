package main

import (
	"log"

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
	msg := "Hello, tcell!  (ESC or Ctrl-C to quit)"
	for x, r := range msg {
		s.SetContent(x, 0, r, nil, def)
	}
	s.Show()

	// イベントループ
	for {
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
