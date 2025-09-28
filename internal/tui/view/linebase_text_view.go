package view

import "github.com/desktopgame/ckro/internal/tui/model"

type LinebaseTextView interface {
	TextView

	ConvertRelativeX(ctx Context, e model.Element, viewLocalPos int) int

	MoveFirstLine(ctx Context, e model.Element, relX int) int

	MoveLastLine(ctx Context, e model.Element, relX int) int
}
