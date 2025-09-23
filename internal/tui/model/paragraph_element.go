package model

type ParagraphElement struct {
	Line string
}

func (p *ParagraphElement) GetText() string {
	return p.Line
}

func (p *ParagraphElement) GetElement(index int) Element {
	return nil
}

func (p *ParagraphElement) GetElementCount() int {
	return 0
}
