package model

type TableContainerElement struct {
	Elements      []Element
	StartPosition Position
	EndPosition   Position
}

func (t *TableContainerElement) GetStyle() *Style {
	return nil
}

func (t *TableContainerElement) GetText() string {
	return "CYZ"
}

func (t *TableContainerElement) GetStartPosition() Position {
	return t.StartPosition
}

func (t *TableContainerElement) GetEndPosition() Position {
	return t.EndPosition
}

func (t *TableContainerElement) GetElement(index int) Element {
	return t.Elements[index]
}

func (t *TableContainerElement) GetElementCount() int {
	return len(t.Elements)
}
