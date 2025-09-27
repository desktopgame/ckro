package litemark

type StringReader struct {
	Source []string
}

func (s *StringReader) GetLine(lineIndex int) string {
	return s.Source[lineIndex]
}

func (s *StringReader) LineCount() int {
	return len(s.Source)
}
