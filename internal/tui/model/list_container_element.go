package model

import "strings"

type ListContainerElement struct {
	Elements      []Element
	StartPosition Position
	EndPosition   Position
}

func (l *ListContainerElement) GetStyle() *Style {
	return nil
}

func (l *ListContainerElement) GetText() string {
	sb := strings.Builder{}
	for _, e := range l.Elements {
		sb.WriteString(e.GetText())
	}
	return sb.String()
}

func (l *ListContainerElement) GetStartPosition() Position {
	return l.StartPosition
}

func (l *ListContainerElement) GetEndPosition() Position {
	return l.EndPosition
}

func (l *ListContainerElement) GetElement(index int) Element {
	return l.Elements[index]
}

func (l *ListContainerElement) GetElementCount() int {
	return len(l.Elements)
}
