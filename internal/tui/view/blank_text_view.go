package view

type BlankTextView interface {
	IsBlank(ctx Context, textLayout *TextLayout) bool
}
