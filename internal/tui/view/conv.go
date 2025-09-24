package view

import (
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/gdamore/tcell/v2"
)

func ConvertStyle(s *model.Style) tcell.Style {
	style := tcell.StyleDefault

	if s.IsBold {
		style = style.Bold(true)
	}
	if s.IsItalic {
		style = style.Italic(true)
	}
	if s.IsUnderline {
		style = style.Underline(true)
	}

	if s.Foreground != model.Default {
		switch s.Foreground {
		case model.White:
			style = style.Foreground(tcell.ColorWhite)
		case model.Black:
			style = style.Foreground(tcell.ColorBlack)
		case model.Red:
			style = style.Foreground(tcell.ColorRed)
		case model.Green:
			style = style.Foreground(tcell.ColorGreen)
		case model.Blue:
			style = style.Foreground(tcell.ColorBlue)
		}
	}

	if s.Foreground != model.Default {
		switch s.Foreground {
		case model.White:
			style = style.Background(tcell.ColorWhite)
		case model.Black:
			style = style.Background(tcell.ColorBlack)
		case model.Red:
			style = style.Background(tcell.ColorRed)
		case model.Green:
			style = style.Background(tcell.ColorGreen)
		case model.Blue:
			style = style.Background(tcell.ColorBlue)
		}
	}

	return style
}
