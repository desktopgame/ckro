package model

type StyledParagraphElement struct {
	Line          string
	StartPosition Position
	EndPosition   Position
	Style         *Style
}

func (s *StyledParagraphElement) GetStyle() *Style {
	return s.Style
}

func (s *StyledParagraphElement) GetText() string {
	return s.Line
}

func (s *StyledParagraphElement) GetStartPosition() Position {
	return s.StartPosition
}

func (s *StyledParagraphElement) GetEndPosition() Position {
	return s.EndPosition
}

func (s *StyledParagraphElement) GetElement(index int) Element {
	return nil
}

func (s *StyledParagraphElement) GetElementCount() int {
	return 0
}
