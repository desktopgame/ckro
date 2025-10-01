package litemark

type Reader interface {
	GetLineString(lineIndex int) string
	GetLineCount() int
}
