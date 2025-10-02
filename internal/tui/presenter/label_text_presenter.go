package presenter

import (
	"strings"

	"github.com/desktopgame/ckro/internal/text"
	"github.com/gdamore/tcell/v2"
)

type LabelTextPresenter struct {
	Text        string
	AlignCenter bool
}

func (label *LabelTextPresenter) Present(view View) {
	view.TextClear()

	sb := strings.Builder{}
	if label.AlignCenter {
		width := view.GetWidth()
		lines := view.GetHeight()
		if lines > 3 {
			for i := 0; i < lines/2; i++ {
				// view.InsertString("\n")
				sb.WriteString("\n")
			}

			length := text.DisplayWidth(label.Text)
			if length >= width {
				// view.InsertString(label.Text)
				sb.WriteString(label.Text)
			} else {
				for i := 0; i < (width-length)/2; i++ {
					// view.InsertString(" ")
					sb.WriteString(" ")
				}
				// view.InsertString(label.Text)
				sb.WriteString(label.Text)
			}

		} else {
			// view.InsertString(label.Text)
			sb.WriteString(label.Text)
		}
	} else {
		// view.InsertString(label.Text)
		sb.WriteString(label.Text)
	}
	view.GetDocument().ReplaceAll(strings.NewReader(sb.String()))
}

func (label *LabelTextPresenter) Handle(view View, ev tcell.Event) {
}

func (label *LabelTextPresenter) ShowCursor() bool {
	return false
}

func (label *LabelTextPresenter) IsFocusable() bool {
	return false
}
