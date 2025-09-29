package view

type BytePositionConverter interface {
	ConvertLocalPos(ctx Context, textLayout *TextLayout, localBytePos int) int
}
