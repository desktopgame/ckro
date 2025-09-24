package model

type InlineElement struct {
	Text          string
	Style         *Style
	StartPosition Position
	EndPosition   Position
}

func (il *InlineElement) GetStyle() *Style {
	return il.Style
}

func (il *InlineElement) GetText() string {
	return il.Text
}

func (il *InlineElement) GetStartPosition() Position {
	return il.StartPosition
}

func (il *InlineElement) GetEndPosition() Position {
	return il.EndPosition
}

func (il *InlineElement) GetElement(index int) Element {
	return nil
}

func (il *InlineElement) GetElementCount() int {
	return 0
}
