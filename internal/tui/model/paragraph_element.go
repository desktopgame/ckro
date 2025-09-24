package model

type ParagraphElement struct {
	Line          string
	StartPosition Position
	EndPosition   Position
}

func (p *ParagraphElement) GetStyle() *Style {
	return nil
}

func (p *ParagraphElement) GetText() string {
	return p.Line
}

func (p *ParagraphElement) GetStartPosition() Position {
	return p.StartPosition
}

func (p *ParagraphElement) GetEndPosition() Position {
	return p.EndPosition
}

func (p *ParagraphElement) GetElement(index int) Element {
	return nil
}

func (p *ParagraphElement) GetElementCount() int {
	return 0
}
