package model

type FoldBlockElement struct {
	Ranges   []Range
	Children []Element
	Lang     string
}

func (f *FoldBlockElement) GetRange(index int) Range {
	return f.Ranges[index]
}

func (f *FoldBlockElement) GetRangeCount() int {
	return len(f.Ranges)
}

func (f *FoldBlockElement) GetElement(index int) Element {
	return f.Children[index]
}

func (f *FoldBlockElement) GetElementCount() int {
	return len(f.Children)
}
