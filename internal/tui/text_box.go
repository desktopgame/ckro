package tui

import (
	"iter"
	"strings"

	"github.com/desktopgame/ckro/internal/text"
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/desktopgame/ckro/internal/tui/presenter"
	"github.com/desktopgame/ckro/internal/tui/view"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

// TextBox is editable text widget.
// TextBox has region of rect, rendering text within that range.
// long line is always wrap at right end of region.
type TextBox struct {
	Document     model.Document
	X            int
	Y            int
	Width        int
	Height       int
	ViewResolver view.TextViewResolver
	ShowCursor   bool
	scrollX      int
	scrollY      int
	viewPosition int

	renderCache TextRenderCache
	bytePos     view.CharacterReference

	foldManager FoldManager

	textSelection        view.TextSelection
	textSelectionEnabled bool
}

// Init is initialize TextBox.
func (tb *TextBox) Init() {
	pd := &model.PlainDocument{}
	pd.Init()

	tb.Document = pd
	tb.X = 0
	tb.Y = 0
	tb.Width = 20
	tb.Height = 6
	tb.ViewResolver = &PlainTextViewResolver{}
	tb.scrollX = 0
	tb.scrollY = 0
}

func (tb *TextBox) context() view.Context {
	return view.Context{
		Resolver:      tb.ViewResolver,
		Document:      tb.Document,
		FoldManager:   &tb.foldManager,
		TextSelection: tb.textSelection,
	}
}

// CursorPosition returns position of cursor.
func (tb *TextBox) CursorPosition() (X int, Y int, Rune rune, Combine []rune) {
	// update elements
	ctx := tb.context()
	tb.renderCache.Update(ctx, tb.Width)

	// get view info from viewPosition
	_, ei, _, vl := tb.renderCache.Stats(tb.viewPosition)

	// calculate a total line count until just before on current view
	screenY := 0
	for i := 0; i < ei; i++ {
		screenY += tb.renderCache.GetLayout(i).Height
	}

	// calculate local position on current view
	currentView := tb.ViewResolver.Resolve(tb.renderCache.GetLayout(ei).Element)
	vlx, vly := currentView.ConvertPos(ctx, tb.renderCache.GetLayout(ei), vl)

	// determine cursor position
	screenX := vlx
	screenY += vly

	// get character on current cursor
	charRef := currentView.ConvertModel(ctx, tb.renderCache.GetLayout(ei), vl)

	// return space if line end or linewrap.
	if charRef.Bytes == 0 || charRef.LineWrap {
		return screenX, screenY, ' ', nil
	}

	// otherwise, return character on current cursor
	r := model.Range{
		StartPosition: charRef.StartPosition,
		EndPosition: model.Position{
			Row:    charRef.StartPosition.Row,
			Column: charRef.StartPosition.Column + charRef.Bytes,
		},
	}
	sg := tb.Document.Read(r)

	char := sg.GetLine(0)

	runes := []rune(char)
	currentRune := runes[0]
	combining := runes[1:]

	return screenX, screenY, currentRune, combining
}

// CursorUpdate is scroll to until cursor visible
func (tb *TextBox) CursorUpdate() {
	_, cursorY, _, _ := tb.CursorPosition()

	startY := tb.scrollY
	endY := startY + tb.Height

	if cursorY >= endY {
		for cursorY >= endY {
			tb.scrollY++

			startY = tb.scrollY
			endY = startY + tb.Height
		}
	} else if cursorY <= startY {
		if cursorY == 0 {
			tb.scrollY = 0
		} else {
			for cursorY <= startY && cursorY > 0 {
				tb.scrollY--

				startY = tb.scrollY
			}
		}
	} else if cursorY > startY && cursorY < endY {
		if cursorY < tb.Height {
			tb.scrollY = 0
		}
	}
}

// CursorReset is reset cursor position and scroll.
func (tb *TextBox) CursorReset() {
	tb.MoveReset()
}

// TextFrame is print a frame of TextBox region.
func (tb *TextBox) TextFrame() {
	tb.Document.Clear()
	tb.MoveReset()

	w := tb.Width
	h := tb.Height
	sb := strings.Builder{}

	sb.WriteString("*")
	for i := 0; i < w-2; i++ {
		sb.WriteString("-")
	}
	sb.WriteString("*")
	sb.WriteString("\n")

	for i := 0; i < h-2; i++ {
		sb.WriteString("|")
		for j := 0; j < w-2; j++ {
			sb.WriteString(" ")
		}
		sb.WriteString("|")
		sb.WriteString("\n")
	}

	sb.WriteString("*")
	for i := 0; i < w-2; i++ {
		sb.WriteString("-")
	}
	sb.WriteString("*")

	tb.Document.ReplaceAll(strings.NewReader(sb.String()))
	tb.MoveReset()
}

// TextVertical is print vertical line.
func (tb *TextBox) TextVertical() {
	tb.Document.Clear()
	tb.MoveReset()
	sb := strings.Builder{}
	h := tb.Height
	for i := 0; i < h; i++ {
		sb.WriteString("|")

		if i < h-1 {
			sb.WriteString("\n")
		}
	}

	tb.Document.ReplaceAll(strings.NewReader(sb.String()))
	tb.MoveReset()
}

// TextHorizontal is print horizontal line.
func (tb *TextBox) TextHorizontal() {
	tb.Document.Clear()
	tb.MoveReset()
	sb := strings.Builder{}
	w := tb.Width
	for i := 0; i < w; i++ {
		sb.WriteString("-")
	}

	tb.Document.ReplaceAll(strings.NewReader(sb.String()))
	tb.MoveReset()
}

// TextClear is do reset to content.
func (tb *TextBox) TextClear() {
	tb.Document.Clear()
	tb.MoveReset()
}

// Draw is render content.
func (tb *TextBox) Draw(g *Graphics) {
	if tb.Width == 0 || tb.Height == 0 {
		return
	}

	// update elements
	ctx := tb.context()
	tb.renderCache.Update(ctx, tb.Width)
	tb.foldManager.Refresh(tb.Document)

	// draw all view when visible in terminal window.
	clip := Clip{
		Graphics:   g,
		X:          tb.X,
		Y:          tb.Y,
		Width:      tb.Width,
		Height:     tb.Height,
		FirstLineY: tb.scrollY,
	}
	cursor := clip
	def := tcell.StyleDefault

	for textSegment := range tb.BreakIter() {
		if textSegment.LocalViewLine == 0 {
			view := tb.ViewResolver.Resolve(textSegment.TextLayout.Element)
			view.Draw(ctx, textSegment.TextLayout, &clip)
			clip.offsetY++
		} else {
			clip.offsetY++
		}
	}

	if !tb.ShowCursor {
		return
	}

	clip = cursor

	screenX, cursorRow, currentRune, combining := tb.CursorPosition()

	// draw character on current cursor by reverse color
	cursorStyle := def.Reverse(true)
	clip.SetCursor(screenX, cursorRow-tb.scrollY, currentRune, combining, cursorStyle)

	// fill neighbor cells, if two width character
	if currentRune != ' ' {
		width := runewidth.RuneWidth(currentRune)
		if width == 2 {
			clip.SetCursor(screenX+1, cursorRow-tb.scrollY, 0, nil, cursorStyle)
		}
	}
}

func (tb *TextBox) modelToView() (ElementIndex int, ViewStart int, ViewLocalPos int) {
	ctx := tb.context()
	if tb.renderCache.GetItemCount() == 0 {
		return 0, 0, 0
	}

	elementIndex := -1
	for i := 0; i < tb.renderCache.GetItemCount(); i++ {
		element := tb.renderCache.GetElement(i)
		r := element.GetRange(0)
		st := r.StartPosition
		ed := r.EndPosition

		// special support for blank line
		// in normally, range is Half-open section
		// but, length is zero when blank line
		// inclusive the end column in this case
		if st.Row == ed.Row && st.Column == ed.Column {
			if tb.bytePos.StartPosition.Row == st.Row && tb.bytePos.StartPosition.Column == st.Column {
				// should be zero length when cursor at blank line
				if tb.bytePos.Bytes > 0 {
					panic("should be zero length when cursor at blank line")
				}
				elementIndex = i
				break
			}
		}

		if tb.bytePos.StartPosition.Row >= st.Row && tb.bytePos.StartPosition.Row <= ed.Row {
			if tb.bytePos.Bytes == 0 {
				// when cursor at line end
				if tb.bytePos.StartPosition.Column >= st.Column && (tb.bytePos.StartPosition.Column <= ed.Column || ed.Row > st.Row) {
					elementIndex = i
					break
				}
			} else {
				// otherwise, judge by Half-open section
				if tb.bytePos.StartPosition.Column >= st.Column && (tb.bytePos.StartPosition.Column < ed.Column || ed.Row > st.Row) {
					elementIndex = i
					break
				}
			}
		}
	}

	viewStart := 0
	for i := 0; i < elementIndex; i++ {
		layout := tb.renderCache.GetLayout(i)
		textView := tb.ViewResolver.Resolve(layout.Element)

		viewStart += textView.MoveLength(ctx, layout)
	}

	element := tb.renderCache.GetElement(elementIndex)
	textView := tb.ViewResolver.Resolve(element)
	viewLocalPos := textView.ConvertViewLocalPos(ctx, tb.renderCache.GetLayout(elementIndex), tb.bytePos.StartPosition)

	return elementIndex, viewStart, viewLocalPos
}

func (tb *TextBox) viewToModel() view.CharacterReference {
	ctx := tb.context()
	_, ei, _, vl := tb.renderCache.Stats(tb.viewPosition)
	element := tb.renderCache.GetElement(ei)
	textView := tb.ViewResolver.Resolve(element)
	return textView.ConvertModel(ctx, tb.renderCache.GetLayout(ei), vl)
}

// InsertString is insert string into current byte position.
func (tb *TextBox) InsertString(s string) {
	if !tb.CanEdit() {
		return
	}

	// update elements
	ctx := tb.context()
	tb.renderCache.Update(ctx, tb.Width)

	// get view info from viewPosition
	_, ei, _, _ := tb.renderCache.Stats(tb.viewPosition)

	// special support for "ghost element"
	// "ghost element" is locatable a cursor, but does not exist string
	// blank line inserted when type text on this element
	if ge, ok := tb.renderCache.GetElement(ei).(*model.GhostElement); ok {
		lines := strings.Repeat("\n", ge.Index+1)
		tb.viewPosition -= ge.Index + 1
		lastElement := tb.renderCache.GetElement(tb.renderCache.GetItemCount() - tb.renderCache.Ghosts() - 1)
		r := lastElement.GetRange(0)
		tb.bytePos = view.CharacterReference{
			StartPosition: model.Position{
				Row:    r.EndPosition.Row,
				Column: r.EndPosition.Column,
			},
			Bytes: 0,
		}
		tb.InsertString(lines)
		tb.InsertString(s)
		return
	}

	// get view before edit
	elementIndex, viewStart, viewLocalPos := tb.modelToView()
	layout := tb.renderCache.GetLayout(elementIndex)
	textView := tb.ViewResolver.Resolve(layout.Element)

	// control a insert position, if exist hidden string in line starts
	if strings.HasSuffix(s, "\n") {
		if textView.ShouldBeforeInsertionNewLineOnLineBegin(ctx, layout, viewLocalPos) {
			tb.bytePos.StartPosition.Column = 0
		}
	}

	// update document
	tb.Document.InsertString(tb.bytePos.StartPosition.Row, tb.bytePos.StartPosition.Column, s)
	tb.renderCache.Update(ctx, tb.Width)

	// calculate new model position
	insertedPos := tb.bytePos.StartPosition
	for i := 0; i < len(s); i++ {
		b := s[i]

		if b == '\n' {
			insertedPos.Row++
			insertedPos.Column = 0
		} else {
			insertedPos.Column++
		}
	}

	// update view position
	tb.bytePos = view.CharacterReference{
		StartPosition: insertedPos,
		Bytes:         tb.bytePos.Bytes,
	}
	_, viewStart, viewLocalPos = tb.modelToView()
	tb.viewPosition = viewStart + viewLocalPos
	tb.bytePos = tb.viewToModel()
}

// RemoveChar is remove character from current byte position to backwards.
func (tb *TextBox) RemoveChar() {
	if !tb.CanEdit() {
		return
	}

	// update elements
	ctx := tb.context()
	tb.renderCache.Update(ctx, tb.Width)

	// get view info from viewPosition
	_, ei, _, _ := tb.renderCache.Stats(tb.viewPosition)

	// special support for "ghost element"
	// "ghost element" is locatable a cursor, but does not exist string
	// in this case, cant remove string
	if ge, ok := tb.renderCache.GetElement(ei).(*model.GhostElement); ok {
		tb.viewPosition -= ge.Index + 1
		return
	}

	// get view before edit
	elementIndex, viewStart, viewLocalPos := tb.modelToView()
	layout := tb.renderCache.GetLayout(elementIndex)
	textView := tb.ViewResolver.Resolve(layout.Element)

	// special supports...
	// can't define perfect completely "general remove operation" when text editor is handle a rich content
	// so, process some edge cases in here
	if lineIndex, ok := textView.ShouldRemoveWithLine(ctx, layout, viewLocalPos); ok {
		bPos := view.CharacterReference{
			StartPosition: model.Position{
				Row:    lineIndex,
				Column: 0,
			},
			Bytes: tb.Document.GetLineBytes(lineIndex),
		}
		tb.bytePos = bPos
		tb.Document.Remove(bPos.StartPosition.Row, bPos.StartPosition.Column, max(bPos.Bytes, 1))

		tb.bytePos.Bytes = 0
		_, viewStart, viewLocalPos = tb.modelToView()
		tb.viewPosition = viewStart + viewLocalPos
		return
	} else if rng, ok := textView.ShouldRemoveWithSpecifiedRangeColumns(ctx, layout, viewLocalPos); ok {
		bPos := view.CharacterReference{
			StartPosition: model.Position{
				Row:    rng.StartPosition.Row,
				Column: rng.StartPosition.Column,
			},
			Bytes: rng.EndPosition.Column - rng.StartPosition.Column,
		}
		tb.bytePos = bPos
		tb.Document.Remove(bPos.StartPosition.Row, bPos.StartPosition.Column, max(bPos.Bytes, 1))
		tb.renderCache.Update(ctx, tb.Width)

		tb.bytePos.Bytes = 0
		_, vs, vl := tb.modelToView()
		tb.viewPosition = vs + vl
		return
	} else if pos, ok := textView.ShouldRemoveWithSpecifiedColumnAfter(ctx, layout, viewLocalPos); ok {
		bPos := view.CharacterReference{
			StartPosition: model.Position{
				Row:    pos.Row,
				Column: pos.Column,
			},
			Bytes: tb.Document.GetLineBytes(pos.Row) - pos.Column,
		}
		tb.bytePos = bPos
		tb.Document.Remove(bPos.StartPosition.Row, bPos.StartPosition.Column, max(bPos.Bytes, 1))
		tb.renderCache.Update(ctx, tb.Width)

		tb.bytePos.Bytes = 0
		_, vs, vl := tb.modelToView()
		tb.viewPosition = vs + vl
		return
	} else if rng, ok := textView.ShouldRemoveWithSpecifiedRangeLines(ctx, layout, viewLocalPos); ok {
		sg := tb.Document.Read(rng)
		for i := sg.GetLineCount() - 1; i >= 0; i-- {
			sp := sg.GetSpan(i)
			tb.Document.Remove(rng.StartPosition.Row+i, sp.StartColumn, sp.EndColumn-sp.StartColumn)
		}

		bPos := view.CharacterReference{
			StartPosition: rng.StartPosition,
			Bytes:         0,
		}
		tb.bytePos = bPos

		tb.renderCache.Update(ctx, tb.Width)

		_, vs, vl := tb.modelToView()
		tb.viewPosition = vs + vl
		return
	}

	if tb.viewPosition == 0 {
		return
	}

	// fallback to "general remove operation"
	newViewLocalPos := textView.MoveLeft(ctx, layout, viewLocalPos)

	// can't move left by current view
	// in other words, curosr into previous view
	if newViewLocalPos == -1 {
		layout = tb.renderCache.GetLayout(elementIndex - 1)
		prevView := tb.ViewResolver.Resolve(layout.Element)

		// special suports...
		// remove last character of previous view, inclusive invisible content
		if elem, ok := prevView.ShouldRemoveLastCharacter(ctx, layout, prevView.MoveLength(ctx, layout)-1); ok {
			r := elem.GetRange(0)
			bPos := view.CharacterReference{
				StartPosition: model.Position{
					Row:    r.EndPosition.Row,
					Column: r.EndPosition.Column - 1,
				},
				Bytes: 1,
			}
			tb.bytePos = bPos
			tb.Document.Remove(bPos.StartPosition.Row, bPos.StartPosition.Column, max(bPos.Bytes, 1))
			tb.renderCache.Update(ctx, tb.Width)

			tb.bytePos.Bytes = 0
			_, vs, vl := tb.modelToView()
			tb.viewPosition = vs + vl
			return
		}

		bPos := prevView.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex-1), prevView.MoveLength(ctx, layout)-1)
		tb.bytePos = bPos

		tb.Document.Remove(bPos.StartPosition.Row, bPos.StartPosition.Column, max(bPos.Bytes, 1))
		tb.bytePos.Bytes = 0
	} else {
		// special suports...
		// remove last character of previous view, inclusive invisible content
		if elem, ok := textView.ShouldRemoveLastCharacter(ctx, layout, newViewLocalPos); ok {
			r := elem.GetRange(0)
			bPos := view.CharacterReference{
				StartPosition: model.Position{
					Row:    r.EndPosition.Row,
					Column: r.EndPosition.Column - 1,
				},
				Bytes: 1,
			}
			tb.bytePos = bPos
			tb.Document.Remove(bPos.StartPosition.Row, bPos.StartPosition.Column, max(bPos.Bytes, 1))
			tb.renderCache.Update(ctx, tb.Width)

			tb.bytePos.Bytes = 0
			_, vs, vl := tb.modelToView()
			tb.viewPosition = vs + vl
			return
		}

		bytes := tb.bytePos.Bytes
		bPos := textView.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex), newViewLocalPos)
		tb.bytePos = bPos

		tb.Document.Remove(bPos.StartPosition.Row, bPos.StartPosition.Column, max(bPos.Bytes, 1))

		tb.bytePos.Bytes = bytes
	}

	tb.renderCache.Update(ctx, tb.Width)

	_, vs, vl := tb.modelToView()
	tb.viewPosition = vs + vl
}

// RemoveSelection is remove selected text.
func (tb *TextBox) RemoveSelection() {
	if tb.textSelection.IsZero() {
		return
	}
	first, last := tb.textSelection.Ordered()

	if first.StartPosition.Row == last.StartPosition.Row {
		// if selected range is single line
		removeStart := first.StartPosition.Column
		removeBytes := last.StartPosition.Column - first.StartPosition.Column
		tb.Document.Remove(first.StartPosition.Row, removeStart, removeBytes)
	} else {
		// if selected range is multiple line
		if last.StartPosition.Row-first.StartPosition.Row >= 2 {
			removeBytes := last.StartPosition.Column
			tb.Document.Remove(last.StartPosition.Row, 0, removeBytes)

			for i := last.StartPosition.Row - 1; i > first.StartPosition.Row; i-- {
				lineLen := tb.Document.GetLineBytes(i)
				tb.Document.Remove(i, 0, lineLen)
				tb.Document.Remove(i, 0, 1)
			}

			removeBytes = first.StartPosition.Column
			lineLen := tb.Document.GetLineBytes(first.StartPosition.Row)
			tb.Document.Remove(first.StartPosition.Row, removeBytes, lineLen-removeBytes)

			tb.Document.Remove(first.StartPosition.Row, removeBytes, 1)
		} else {
			removeBytes := last.StartPosition.Column
			r := model.Range{
				StartPosition: model.Position{
					Row:    last.StartPosition.Row,
					Column: removeBytes,
				},
				EndPosition: model.Position{
					Row:    last.StartPosition.Row,
					Column: tb.Document.GetLineBytes(last.StartPosition.Row),
				},
			}
			sg := tb.Document.Read(r)
			copy := sg.GetLine(0)

			tb.Document.Remove(last.StartPosition.Row, 0, tb.Document.GetLineBytes(last.StartPosition.Row))
			tb.Document.Remove(last.StartPosition.Row, 0, 1)

			removeBytes = first.StartPosition.Column
			lineLen := tb.Document.GetLineBytes(first.StartPosition.Row)
			tb.Document.Remove(first.StartPosition.Row, removeBytes, lineLen-removeBytes)
			tb.Document.InsertString(first.StartPosition.Row, removeBytes, copy)
		}
	}
	tb.renderCache.Update(tb.context(), tb.Width)

	tb.bytePos = view.CharacterReference{
		StartPosition: first.StartPosition,
		Bytes:         last.Bytes,
	}
	_, vs, vl := tb.modelToView()
	tb.viewPosition = vs + vl
}

// CanEdit returns can editable.
func (tb *TextBox) CanEdit() bool {
	ctx := tb.context()
	tb.renderCache.Update(ctx, tb.Width)

	_, ei, _, _ := tb.renderCache.Stats(tb.viewPosition)
	element := tb.renderCache.GetElement(ei)

	if _, ok := element.(*model.FoldBlockElement); ok {
		if tb.foldManager.IsFolded(tb.Document, element) {
			return false
		}
		return true
	}
	return true
}

// Submit is try submit process on current TextView.
// returns true if execute in actual.
func (tb *TextBox) Submit() bool {
	ctx := tb.context()
	tb.renderCache.Update(ctx, tb.Width)

	_, ei, _, vl := tb.renderCache.Stats(tb.viewPosition)
	element := tb.renderCache.GetElement(ei)

	if fold, ok := element.(*model.FoldBlockElement); ok {
		if vl > 0 {
			if tb.foldManager.IsFolded(tb.Document, element) {
				return true
			}
			return false
		}
		tb.foldManager.ToggleFold(tb.Document, fold)
		tb.renderCache.ForceUpdate(ctx, tb.Width)

		_, vs, vl := tb.modelToView()
		tb.viewPosition = vs + vl
		return true
	}
	return false
}

// SelectionStart is start text selection.
// select text on every times to call MoveXxx, until call to SelectionEnd
func (tb *TextBox) SelectionStart() {
	tb.textSelection.FromPos = tb.bytePos
	tb.textSelectionEnabled = true
}

// SelectionEnd is stop text selection.
func (tb *TextBox) SelectionEnd() {
	tb.textSelectionEnabled = false
	tb.textSelection = view.TextSelection{}
}

func (tb *TextBox) move(dir int) {
	ctx := tb.context()
	tb.renderCache.Update(ctx, tb.Width)

	_, elementIndex, elementStart, oldLocalViewPos := tb.renderCache.Stats(tb.viewPosition)

	tview := tb.ViewResolver.Resolve(tb.renderCache.GetElement(elementIndex))
	var newLocalViewPos int
	switch dir {
	case 0:
		newLocalViewPos = tview.MoveLeft(ctx, tb.renderCache.GetLayout(elementIndex), oldLocalViewPos)
	case 1:
		newLocalViewPos = tview.MoveRight(ctx, tb.renderCache.GetLayout(elementIndex), oldLocalViewPos)
	case 2:
		newLocalViewPos = tview.MoveUp(ctx, tb.renderCache.GetLayout(elementIndex), oldLocalViewPos)
	case 3:
		newLocalViewPos = tview.MoveDown(ctx, tb.renderCache.GetLayout(elementIndex), oldLocalViewPos)
	}

	if newLocalViewPos == -1 {
		tvLen := tview.MoveLength(ctx, tb.renderCache.GetLayout(elementIndex))
		if dir == 3 {
			if pLinebaseView, ok := tview.(view.LinebaseTextView); ok && elementIndex+1 < tb.renderCache.GetItemCount() {
				nextView := tb.ViewResolver.Resolve(tb.renderCache.GetElement(elementIndex + 1))
				relx := pLinebaseView.ConvertRelativeX(ctx, tb.renderCache.GetLayout(elementIndex), oldLocalViewPos)

				if linebaseTV, ok := nextView.(view.LinebaseTextView); ok {
					offset := linebaseTV.MoveFirstLine(ctx, tb.renderCache.GetLayout(elementIndex+1), relx)

					tb.viewPosition = elementStart + tvLen + offset

					bPos := linebaseTV.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex+1), offset)
					tb.bytePos = bPos
				} else {
					tb.viewPosition = elementStart + tvLen

					bPos := nextView.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex+1), 0)
					tb.bytePos = bPos
				}
			} else {
				if elementIndex+1 < tb.renderCache.GetItemCount() {
					nextView := tb.ViewResolver.Resolve(tb.renderCache.GetElement(elementIndex + 1))
					tb.viewPosition = elementStart + tvLen

					bPos := nextView.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex+1), 0)
					tb.bytePos = bPos
				} else {

					tb.viewPosition = elementStart + tvLen - 1

					bPos := tview.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex), tvLen-1)
					tb.bytePos = bPos
				}
			}
		} else if dir == 1 {
			if elementIndex+1 < tb.renderCache.GetItemCount() {
				tb.viewPosition = elementStart + tvLen

				nView := tb.ViewResolver.Resolve(tb.renderCache.GetElement(elementIndex + 1))
				bPos := nView.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex+1), 0)
				tb.bytePos = bPos
			} else {
				tb.viewPosition = elementStart + tvLen - 1
				bPos := tview.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex), tvLen-1)
				tb.bytePos = bPos
			}
		}

		if dir == 0 {
			tb.viewPosition = max(elementStart-1, 0)

			if elementIndex > 0 {
				pView := tb.ViewResolver.Resolve(tb.renderCache.GetElement(elementIndex - 1))
				pViewLen := pView.MoveLength(ctx, tb.renderCache.GetLayout(elementIndex-1))
				bPos := pView.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex-1), pViewLen-1)
				tb.bytePos = bPos
			} else {
				bPos := tview.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex), 0)
				tb.bytePos = bPos
			}
		} else if dir == 2 {
			if pLinebaseView, ok := tview.(view.LinebaseTextView); ok && elementIndex > 0 {
				prevView := tb.ViewResolver.Resolve(tb.renderCache.GetElement(elementIndex - 1))
				relx := pLinebaseView.ConvertRelativeX(ctx, tb.renderCache.GetLayout(elementIndex), oldLocalViewPos)

				if linebaseTV, ok := prevView.(view.LinebaseTextView); ok {
					offset := linebaseTV.MoveLastLine(ctx, tb.renderCache.GetLayout(elementIndex-1), relx)
					l := linebaseTV.MoveLength(ctx, tb.renderCache.GetLayout(elementIndex-1))

					tb.viewPosition = elementStart - l + offset

					bPos := linebaseTV.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex-1), offset)
					tb.bytePos = bPos
				} else {
					tb.viewPosition = max(elementStart-1, 0)

					if elementIndex > 0 {
						pView := tb.ViewResolver.Resolve(tb.renderCache.GetElement(elementIndex - 1))
						pViewLen := pView.MoveLength(ctx, tb.renderCache.GetLayout(elementIndex-1))
						bPos := pView.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex-1), pViewLen-1)
						tb.bytePos = bPos
					}
				}
			} else {
				tb.viewPosition = max(elementStart-1, 0)

				if elementIndex > 0 {
					pView := tb.ViewResolver.Resolve(tb.renderCache.GetElement(elementIndex - 1))
					pViewLen := pView.MoveLength(ctx, tb.renderCache.GetLayout(elementIndex-1))
					bPos := pView.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex-1), pViewLen-1)
					tb.bytePos = bPos
				}
			}
		}
	} else {
		tb.viewPosition = elementStart + newLocalViewPos
		bPos := tview.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex), newLocalViewPos)
		tb.bytePos = bPos
	}

	// selection update
	if tb.textSelectionEnabled {
		tb.textSelection.ToPos = tb.bytePos
	}
}

// MoveLeft is move cursor to left.
func (tb *TextBox) MoveLeft() {
	tb.move(0)
}

// MoveRight is move cursor to right.
func (tb *TextBox) MoveRight() {
	tb.move(1)
}

// MoveUp is move cursor to up.
func (tb *TextBox) MoveUp() {
	tb.move(2)
}

// MoveDown is move cursor to down.
func (tb *TextBox) MoveDown() {
	tb.move(3)
}

// MoveLineStart is move cursor to starts of current line.
func (tb *TextBox) MoveLineStart() {
	ctx := tb.context()
	tb.renderCache.Update(ctx, tb.Width)

	_, ei, estart, eoff := tb.renderCache.Stats(tb.viewPosition)
	layout := tb.renderCache.GetLayout(ei)
	tl := tb.renderCache.GetLayout(ei)
	textView := tb.ViewResolver.Resolve(layout.Element)

	lx, _ := textView.ConvertPos(ctx, tl, eoff)
	viewLocalPos := eoff
	for lx > 0 {
		nextLocalPos := textView.MoveLeft(ctx, layout, viewLocalPos)
		if nextLocalPos == -1 {
			break
		}
		viewLocalPos = nextLocalPos
	}
	tb.viewPosition = estart + viewLocalPos
	tb.bytePos = tb.viewToModel()
}

// MoveLineEnd is move cursor to ends of current line.
func (tb *TextBox) MoveLineEnd() {
	ctx := tb.context()
	tb.renderCache.Update(ctx, tb.Width)

	_, ei, estart, eoff := tb.renderCache.Stats(tb.viewPosition)
	layout := tb.renderCache.GetLayout(ei)
	textView := tb.ViewResolver.Resolve(layout.Element)

	viewLocalPos := eoff
	for {
		nextLocalPos := textView.MoveRight(ctx, layout, viewLocalPos)
		if nextLocalPos == -1 {
			break
		}
		viewLocalPos = nextLocalPos
	}
	tb.viewPosition = estart + viewLocalPos
	tb.bytePos = tb.viewToModel()
}

// MoveTextStart is move cursor to starts of text.
func (tb *TextBox) MoveTextStart() {
	tb.renderCache.Update(tb.context(), tb.Width)

	tb.viewPosition = 0
	tb.bytePos = tb.viewToModel()
}

// MoveTextEnd is move cursor to ends of text.
func (tb *TextBox) MoveTextEnd() {
	ctx := tb.context()
	tb.renderCache.Update(ctx, tb.Width)

	vp := 0
	for i := 0; i < tb.renderCache.GetItemCount()-tb.renderCache.Ghosts(); i++ {
		layout := tb.renderCache.GetLayout(i)
		element := layout.Element
		textView := tb.ViewResolver.Resolve(element)
		vp += textView.MoveLength(ctx, layout)
	}
	tb.viewPosition = vp - 1
	tb.bytePos = tb.viewToModel()
}

// MoveReset is reset cursor
func (tb *TextBox) MoveReset() {
	ctx := tb.context()
	tb.renderCache.Update(ctx, tb.Width)
	tb.viewPosition = 0
	tb.bytePos = tb.viewToModel()
	tb.scrollX = 0
	tb.scrollY = 0
}

// FindPrev is finding text to backwards direction, and move to cursor it.
func (tb *TextBox) FindPrev(s string) bool {
	r := model.Range{
		StartPosition: model.Position{
			Row:    0,
			Column: 0,
		},
		EndPosition: model.Position{
			Row:    tb.Document.GetLineCount() - 1,
			Column: tb.Document.GetLineBytes(tb.Document.GetLineCount() - 1),
		},
	}
	sg := tb.Document.Read(r)
	lines := strings.Split(s, "\n")
	linePos := len(lines) - 1
	success := false

	for i := tb.bytePos.StartPosition.Row; i >= 0; i-- {
		line := sg.GetLine(i)

		if len(lines) > 1 {
			if linePos == 0 {
				p := strings.LastIndex(line, lines[0])
				if p >= 0 {
					tb.bytePos = view.CharacterReference{
						StartPosition: model.Position{
							Row:    i,
							Column: len(line[0:p]),
						},
						Bytes: len(text.GraphemeClusters(lines[0])[0]),
					}
					success = true
					break
				} else {
					linePos = len(lines) - 1
				}
			} else if linePos == len(lines)-1 {
				findPos := len(line)
				if i == tb.bytePos.StartPosition.Row {
					if tb.bytePos.Bytes > 0 {
						findPos = tb.bytePos.StartPosition.Column
					}
				}

				if strings.HasPrefix(line[0:findPos], lines[linePos]) {
					linePos--
				} else {
					linePos = len(lines) - 1
				}
			} else {
				if line == lines[linePos] {
					linePos--
				} else {
					linePos = len(lines) - 1
				}
			}
		} else {
			findPos := len(line)
			if i == tb.bytePos.StartPosition.Row {
				if tb.bytePos.Bytes > 0 {
					findPos = tb.bytePos.StartPosition.Column
				}
			}
			p := strings.LastIndex(line[0:findPos], lines[linePos])
			if p >= 0 {
				tb.bytePos = view.CharacterReference{
					StartPosition: model.Position{
						Row:    i,
						Column: len(line[0:p]),
					},
					Bytes: len(text.GraphemeClusters(lines[0])[0]),
				}
				success = true
				break
			}
		}
	}
	if success {
		_, vs, vl := tb.modelToView()
		tb.viewPosition = vs + vl
	}
	return success
}

// FindNext is finding text to forwards direction, and move to cursor it.
func (tb *TextBox) FindNext(s string) bool {
	r := model.Range{
		StartPosition: model.Position{
			Row:    0,
			Column: 0,
		},
		EndPosition: model.Position{
			Row:    tb.Document.GetLineCount() - 1,
			Column: tb.Document.GetLineBytes(tb.Document.GetLineCount() - 1),
		},
	}
	sg := tb.Document.Read(r)
	lines := strings.Split(s, "\n")
	linePos := 0
	findRow := 0
	findCol := 0
	success := false

	for i := tb.bytePos.StartPosition.Row; i < tb.Document.GetLineCount(); i++ {
		line := sg.GetLine(i)

		if len(lines) > 1 {
			if linePos == 0 {
				findPos := 0
				if i == tb.bytePos.StartPosition.Row {
					if tb.bytePos.Bytes > 0 {
						findPos = tb.bytePos.StartPosition.Column + 1
					}
				}

				if strings.HasSuffix(line[findPos:], lines[0]) {
					p := strings.LastIndex(line[findPos:], lines[0])
					if p >= 0 {
						findRow = i
						findCol = len(line[0 : findPos+p])
						linePos++
					} else {
						linePos = 0
						findRow = 0
						findCol = 0
					}
				} else {
					linePos = 0
					findRow = 0
					findCol = 0
				}
			} else if linePos == len(lines)-1 {
				bytes := len(text.GraphemeClusters(lines[0])[0])
				if strings.HasPrefix(line, lines[linePos]) {
					tb.bytePos = view.CharacterReference{
						StartPosition: model.Position{
							Row:    findRow,
							Column: findCol,
						},
						Bytes: bytes,
					}
					success = true
					break
				} else {
					linePos = 0
					findRow = 0
					findCol = 0
				}
			} else {
				if line == lines[linePos] {
					linePos++
				} else {
					linePos = 0
					findRow = 0
					findCol = 0
				}
			}
		} else {
			findPos := 0
			if i == tb.bytePos.StartPosition.Row {
				if tb.bytePos.Bytes > 0 {
					findPos = tb.bytePos.StartPosition.Column + 1
				}
			}
			p := strings.Index(line[findPos:], lines[linePos])
			if p >= 0 {
				tb.bytePos = view.CharacterReference{
					StartPosition: model.Position{
						Row:    i,
						Column: len(line[0:p]),
					},
					Bytes: len(text.GraphemeClusters(lines[0])[0]),
				}
				success = true
				break
			}
		}
	}
	if success {
		_, vs, vl := tb.modelToView()
		tb.viewPosition = vs + vl
	}
	return success
}

// Replace is replace string.
func (tb *TextBox) Replace(length int, s string) {
	tb.Document.Remove(tb.bytePos.StartPosition.Row, tb.bytePos.StartPosition.Column, length)
	tb.Document.InsertString(tb.bytePos.StartPosition.Row, tb.bytePos.StartPosition.Column, s)

	startPosition := model.Position{
		Row:    tb.bytePos.StartPosition.Row,
		Column: tb.bytePos.StartPosition.Column,
	}

	for i := 0; i < len(s); i++ {
		b := s[i]

		if b == '\n' {
			startPosition.Row++
			startPosition.Column = 0
		} else {
			startPosition.Column++
		}
	}

	r := model.Range{
		StartPosition: startPosition,
		EndPosition: model.Position{
			Row:    startPosition.Row,
			Column: tb.Document.GetLineBytes(startPosition.Row),
		},
	}

	if r.IsZero() {
		tb.bytePos = view.CharacterReference{
			StartPosition: r.StartPosition,
			Bytes:         0,
		}
	} else {
		sg := tb.Document.Read(r)
		str := sg.GetLine(0)
		tb.bytePos = view.CharacterReference{
			StartPosition: r.StartPosition,
			Bytes:         len(text.GraphemeClusters(str)[0]),
		}
	}
	tb.renderCache.Update(tb.context(), tb.Width)
	_, vs, vl := tb.modelToView()
	tb.viewPosition = vs + vl
}

// BreakIter returns segment array by line, in consideration a wrap.
func (tb *TextBox) BreakIter() iter.Seq[presenter.Segment] {
	//buf := tb.Document.GetBuffer()
	//sb := strings.Builder{}

	return func(yield func(presenter.Segment) bool) {
		ctx := tb.context()
		tb.renderCache.Update(ctx, tb.Width)

		viewLine := 0
		for i := 0; i < tb.renderCache.GetItemCount(); i++ {
			entry := tb.renderCache.GetLayout(i)
			element := entry.Element

			for j := 0; j < entry.Height; j++ {
				_, isGhost := element.(*model.GhostElement)
				segment := presenter.Segment{
					TextLayout:    entry,
					IsGhostLine:   isGhost,
					ModelLine:     element.GetRange(0).StartPosition.Row,
					ViewLine:      viewLine,
					LocalViewLine: j,
				}
				if !yield(segment) {
					return
				}
				viewLine++
			}
		}
	}
}

// GetDocument returns Document.
func (tb *TextBox) GetDocument() model.Document {
	return tb.Document
}

// GetWidth returns width of TextBox region.
func (tb *TextBox) GetWidth() int {
	return tb.Width
}

// GetHeight returns height of TextBox region.
func (tb *TextBox) GetHeight() int {
	return tb.Height
}

// GetScrollX returns scroll amount by horizontal.
func (tb *TextBox) GetScrollX() int {
	return tb.scrollX
}

// GetScrollY returns scroll amount by vertical.
func (tb *TextBox) GetScrollY() int {
	return tb.scrollY
}

// GetViewHeight returns height of total TextView.
func (tb *TextBox) GetViewHeight() int {
	ctx := tb.context()
	tb.renderCache.Update(ctx, tb.Width)

	h := 0
	for i := 0; i < tb.renderCache.GetItemCount(); i++ {
		h += tb.renderCache.GetLayout(i).Height
	}
	return h
}

// GetBytePosition returns byte position of cursor.
func (tb *TextBox) GetBytePosition() view.CharacterReference {
	return tb.bytePos
}

// GetViewPosition returns view position of cursor.
func (tb *TextBox) GetViewPosition() int {
	return tb.viewPosition
}

// GetSelection returns selected range.
func (tb *TextBox) GetSelection() view.TextSelection {
	return tb.textSelection
}
