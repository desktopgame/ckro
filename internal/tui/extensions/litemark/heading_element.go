package litemark

import "github.com/desktopgame/ckro/internal/tui/model"

type HeadingElement struct {
	Text          string
	Style         *model.Style
	Level         int
	StartPosition model.Position
	EndPosition   model.Position
}

func (h *HeadingElement) GetStyle() *model.Style {
	return h.Style
}

func (h *HeadingElement) GetText() string {
	return h.Text
}

func (h *HeadingElement) GetStartPosition() model.Position {
	return h.StartPosition
}

func (h *HeadingElement) GetEndPosition() model.Position {
	return h.EndPosition
}

func (h *HeadingElement) GetElement(index int) model.Element {
	return nil
}

func (h *HeadingElement) GetElementCount() int {
	return 0
}
