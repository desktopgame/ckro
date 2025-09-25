package model

type PlainElement struct {
	Text          string
	Style         *Style
	StartPosition Position
	EndPosition   Position
}

func (p *PlainElement) GetStyle() *Style {
	return nil
}

func (p *PlainElement) GetText() string {
	return p.Text
}

func (p *PlainElement) GetStartPosition() Position {
	return p.StartPosition
}

func (p *PlainElement) GetEndPosition() Position {
	return p.EndPosition
}

func (p *PlainElement) GetElement(index int) Element {
	return nil
}

func (p *PlainElement) GetElementCount() int {
	return 0
}
