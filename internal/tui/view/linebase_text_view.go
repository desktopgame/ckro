package view

type LinebaseTextView interface {
	TextView

	ConvertRelativeX(ctx Context, textLayout *TextLayout, viewLocalPos int) int

	MoveFirstLine(ctx Context, textLayout *TextLayout, relX int) int

	MoveLastLine(ctx Context, textLayout *TextLayout, relX int) int
}
