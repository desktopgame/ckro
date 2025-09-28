package view

import (
	"github.com/desktopgame/ckro/internal/tui/model"
)

type BytePositionConverter interface {
	ConvertLocalPos(ctx Context, e model.Element, localBytePos int) int
}
