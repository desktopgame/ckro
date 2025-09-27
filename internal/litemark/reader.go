package litemark

type Reader interface {
	GetLine(lineIndex int) string
	LineCount() int
}
