package model

type GhostElement struct {
	Range Range
	Index int
}

func (g *GhostElement) GetRange(index int) Range {
	return g.Range
}

func (g *GhostElement) GetRangeCount() int {
	return 1
}

func (g *GhostElement) GetElement(index int) Element {
	return nil
}

func (g *GhostElement) GetElementCount() int {
	return 0
}
