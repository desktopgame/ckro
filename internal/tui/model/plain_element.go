package model

type PlainElement struct {
	Range Range
}

func (p *PlainElement) GetRange(index int) Range {
	return p.Range
}

func (p *PlainElement) GetRangeCount() int {
	return 1
}

func (p *PlainElement) GetElement(index int) Element {
	return nil
}

func (p *PlainElement) GetElementCount() int {
	return 0
}
