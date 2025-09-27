package litemark

type Scanner struct {
	Reader    Reader
	lineIndex int
}

func (s *Scanner) Next() string {
	line := s.Reader.GetLine(s.lineIndex)
	s.lineIndex++
	return line
}

func (s *Scanner) Ready() bool {
	return s.lineIndex < s.Reader.LineCount()
}
