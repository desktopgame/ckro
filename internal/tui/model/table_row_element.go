package model

type TableRowElement struct {
	Elements      []Element
	StartPosition Position
	EndPosition   Position
}

func (t *TableRowElement) GetStyle() *Style {
	return nil
}

func (t *TableRowElement) GetText() string {
	return "AAA"
}

func (t *TableRowElement) GetStartPosition() Position {
	return t.StartPosition
}

func (t *TableRowElement) GetEndPosition() Position {
	return t.EndPosition
}

func (t *TableRowElement) GetElement(index int) Element {
	return t.Elements[index]
}

func (t *TableRowElement) GetElementCount() int {
	return len(t.Elements)
}
