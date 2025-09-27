package litemark

type Reader interface {
	GetLineAt(lineIndex int) string
	LineCount() int
}
