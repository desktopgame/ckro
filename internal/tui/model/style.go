package model

type Color int

const (
	Default = iota
	White
	Black
	Red
	Blue
	Green
)

type Style struct {
	IsBold      bool
	IsItalic    bool
	IsUnderline bool
	Foreground  Color
	Background  Color
}
