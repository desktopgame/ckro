package model

import "strings"

type LineContainerElement struct {
	Elements      []Element
	StartPosition Position
	EndPosition   Position
}

func (l *LineContainerElement) GetStyle() *Style {
	return nil
}

func (l *LineContainerElement) GetText() string {
	sb := strings.Builder{}
	for _, e := range l.Elements {
		sb.WriteString(e.GetText())
	}
	return sb.String()
}

func (l *LineContainerElement) GetStartPosition() Position {
	return l.StartPosition
}

func (l *LineContainerElement) GetEndPosition() Position {
	return l.EndPosition
}

func (l *LineContainerElement) GetElement(index int) Element {
	return l.Elements[index]
}

func (l *LineContainerElement) GetElementCount() int {
	return len(l.Elements)
}
