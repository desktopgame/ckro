package litemark

type StringReader struct {
	Source []string
}

func (s *StringReader) GetLineString(lineIndex int) string {
	return s.Source[lineIndex]
}

func (s *StringReader) GetLineCount() int {
	return len(s.Source)
}
