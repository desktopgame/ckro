package view

import "github.com/desktopgame/ckro/internal/tui/model"

// TextView is render a part of text.
// TextView is should be stateless, and every state is given from parameter.
type TextView interface {
	// Layout is layout this view and contained sub views.
	// x and y is origin of this view, aim left top.
	// w and h is given layout size.
	Layout(ctx Context, textLayout *TextLayout, x, y, w, h int)

	// Draw is rendernig content.
	Draw(ctx Context, textLayout *TextLayout, renderer Renderer)

	// Measure returns TextLayout for this view and contained sub views.
	// calculate only minimum size, at this point.
	Measure(ctx Context, e model.Element, width int, height int) *TextLayout

	// MoveLength returns movable count for cursor in this view.
	// inclusive a subview move count if this view is inclusive a subview.
	MoveLength(ctx Context, textLayout *TextLayout) int

	// MoveUp returns position after cursor up.
	MoveUp(ctx Context, textLayout *TextLayout, viewLocalPos int) int

	// MoveDown returns position after cursor down.
	MoveDown(ctx Context, textLayout *TextLayout, viewLocalPos int) int

	// MoveLeft returns position after cursor left.
	MoveLeft(ctx Context, textLayout *TextLayout, viewLocalPos int) int

	// MoveRight returns position after cursor right.
	MoveRight(ctx Context, textLayout *TextLayout, viewLocalPos int) int

	// ConvertPos returns 2D position from view local position.
	ConvertPos(ctx Context, textLayout *TextLayout, viewLocalPos int) (ViewLocalX int, ViewLocalY int)

	// ConvertModel returns character at specified view local position.
	ConvertModel(ctx Context, textLayout *TextLayout, viewLocalPos int) CharacterReference

	// ConvertViewLocalPos returns view local position at specified byte position.
	// move to inward if byte position is aim to decoration.
	ConvertViewLocalPos(ctx Context, textLayout *TextLayout, bytePos model.Position) int

	ShouldBeforeInsertionNewLineOnLineBegin(ctx Context, textLayout *TextLayout, viewLocalPos int) bool
	ShouldRemoveWithLine(ctx Context, textLayout *TextLayout, viewLocalPos int) (int, bool)
	ShouldRemoveWithSpecifiedColumnAfter(ctx Context, textLayout *TextLayout, viewLocalPos int) (model.Position, bool)
	ShouldRemoveWithSpecifiedRangeLines(ctx Context, textLayout *TextLayout, viewLocalPos int) (model.Range, bool)
	ShouldRemoveWithSpecifiedRangeColumns(ctx Context, textLayout *TextLayout, viewLocalPos int) (model.Range, bool)
	ShouldRemoveLastCharacter(ctx Context, textLayout *TextLayout, viewLocalPos int) (model.Element, bool)
}
