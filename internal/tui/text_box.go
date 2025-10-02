package tui

import (
	"iter"
	"strings"

	"github.com/desktopgame/ckro/internal/text"
	"github.com/desktopgame/ckro/internal/tui/extensions/litemark"
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/desktopgame/ckro/internal/tui/presenter"
	"github.com/desktopgame/ckro/internal/tui/view"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

type textSegment struct {
	textLayout *view.TextLayout
	segment    presenter.Segment
}

// TextBox is editable text widget.
// TextBox has region of rect, rendering text within that range.
// long line is always wrap at right end of region.
type TextBox struct {
	Document     model.Document
	X            int
	Y            int
	Width        int
	Height       int
	TextEngine   TextEngine
	ShowCursor   bool
	scrollX      int
	scrollY      int
	viewPosition int

	renderCache TextRenderCache
	bytePos     view.CharacterReference
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
	tb.TextEngine = &LitemarkEngine{}
	tb.scrollX = 0
	tb.scrollY = 0
}

// CursorPosition returns position of cursor.
// TODO: refactor
func (tb *TextBox) CursorPosition() (X int, Y int, Rune rune, Combine []rune) {
	ctx := view.Context{
		Resolver: tb.TextEngine,
		Document: tb.Document,
	}

	tb.renderCache.Update(ctx, tb.Width)
	_, ei, _, eoff := tb.renderCache.Stats(tb.viewPosition)
	//view := tb.TextEngine.Resolve(tb.renderCache.GetElement(ei))
	//y := 0
	//for i := 0; i < ei; i++ {
	//	y += tb.renderCache.GetLayout(i).Height
	//}
	//_, _ := view.ConvertPos(ctx, tb.renderCache.GetLayout(ei), eoff)
	// ax := vlx
	//ay := y + vly

	//buf := tb.Document.GetBuffer()
	//cursorRow := ay
	// cursorCol := ax

	//if cursorRow >= buf.GetLineCount() {
	//	return 0, 0, ' ', nil
	//}

	// カーソルがある行までの画面行数を計算
	screenY := 0
	//for i := 0; i < cursorRow; i++ {
	//	line := buf.GetLineAt(i).GetContent()
	//	screenY += tb.calculateWrappedLines(line)
	//}
	for i := 0; i < ei; i++ {
		screenY += tb.renderCache.GetLayout(i).Height
	}

	// カーソルがある行での位置を正確に計算
	// cursorLine := buf.GetLineAt(cursorRow).GetContent()
	currentView := tb.TextEngine.Resolve(tb.renderCache.GetElement(ei))
	relx, rely := currentView.ConvertPos(ctx, tb.renderCache.GetLayout(ei), eoff)
	screenX := relx
	screenY += rely
	//screenX, additionalRows :=
	//screenY += additionalRows

	charRef := currentView.ConvertModel(ctx, tb.renderCache.GetLayout(ei), eoff)

	if charRef.Bytes == 0 {
		return screenX, screenY, ' ', nil
	}

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
	_, cursor, _, _ := tb.CursorPosition()
	// cursor := tb.Document.GetCursorRow()
	// lc := tb.WrappedLineCount()

	startY := tb.scrollY
	endY := startY + tb.Height

	if cursor >= endY {
		for cursor >= endY {
			tb.scrollY++

			startY = tb.scrollY
			endY = startY + tb.Height
		}
	} else if cursor <= startY {
		if cursor == 0 {
			tb.scrollY = 0
		} else {
			for cursor <= startY && cursor > 0 {
				tb.scrollY--

				startY = tb.scrollY
			}
		}
	} else if cursor > startY && cursor < endY {
		if cursor < tb.Height {
			tb.scrollY = 0
		}
	}
}

// CursorReset is reset cursor position and scroll.
func (tb *TextBox) CursorReset() {
	// TODO: impl
	// tb.Document.MoveReset()
	tb.scrollX = 0
	tb.scrollY = 0
}

// TextFrame is print a frame of TextBox region.
func (tb *TextBox) TextFrame() {
	tb.Document.Clear()
	tb.MoveReset()

	w := tb.Width
	h := tb.Height

	tb.InsertString("*")
	for i := 0; i < w-2; i++ {
		tb.InsertString("-")
	}
	tb.InsertString("*")
	tb.InsertString("\n")

	for i := 0; i < h-2; i++ {
		tb.InsertString("|")
		for j := 0; j < w-2; j++ {
			tb.InsertString(" ")
		}
		tb.InsertString("|")
		tb.InsertString("\n")
	}

	tb.InsertString("*")
	for i := 0; i < w-2; i++ {
		tb.InsertString("-")
	}
	tb.InsertString("*")

	tb.MoveReset()
}

// TextVertical is print vertical line.
func (tb *TextBox) TextVertical() {
	tb.Document.Clear()
	tb.MoveReset()
	h := tb.Height
	for i := 0; i < h; i++ {
		tb.InsertString("|\n")
	}
	tb.RemoveChar()
	tb.MoveReset()
}

// TextHorizontal is print horizontal line.
func (tb *TextBox) TextHorizontal() {
	tb.Document.Clear()
	tb.MoveReset()
	w := tb.Width
	for i := 0; i < w; i++ {
		tb.InsertString("-")
	}

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

	ctx := view.Context{
		Resolver: tb.TextEngine,
		Document: tb.Document,
	}
	tb.renderCache.Update(ctx, tb.Width)

	// バッファの内容を描画
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
			view := tb.TextEngine.Resolve(textSegment.TextLayout.Element)
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

	if tb.isStyled() {
		//_, ei, _, eoff := tb.currentViewState(ctx, tb.elements)
		//view := tb.TextEngine.Resolve(tb.elements[ei])
		//y := 0
		//for i := 0; i < ei; i++ {
		//	y += tb.layoutCache[i].Height
		//}
		//vlx, vly := view.ConvertPos(ctx, tb.elements[ei], eoff)
		//ax := vlx
		//ay := y + vly

		screenX, cursorRow, currentRune, combining := tb.CursorPosition()

		// カーソル位置の文字を反転表示
		cursorStyle := def.Reverse(true)
		clip.SetCursor(screenX, cursorRow-tb.scrollY, currentRune, combining, cursorStyle)

		// 全角文字の場合、隣接するセルもカーソル表示
		if currentRune != ' ' {
			width := runewidth.RuneWidth(currentRune)
			if width == 2 {
				// 隣接するセルにもカーソルを表示（空文字で反転）
				clip.SetCursor(screenX+1, cursorRow-tb.scrollY, 0, nil, cursorStyle)
			}
		}
		return
	}

	screenX, cursorRow, currentRune, combining := tb.CursorPosition()

	// カーソル位置の文字を反転表示
	cursorStyle := def.Reverse(true)
	clip.SetCursor(screenX, cursorRow-tb.scrollY, currentRune, combining, cursorStyle)

	// 全角文字の場合、隣接するセルもカーソル表示
	if currentRune != ' ' {
		width := runewidth.RuneWidth(currentRune)
		if width == 2 {
			// 隣接するセルにもカーソルを表示（空文字で反転）
			clip.SetCursor(screenX+1, cursorRow-tb.scrollY, 0, nil, cursorStyle)
		}
	}

}

func (tb *TextBox) isStyled() bool {
	_, ok := tb.Document.(*model.PlainDocument)
	return !ok
}

func (tb *TextBox) move(dir int) {
	ctx := view.Context{
		Resolver: tb.TextEngine,
		Document: tb.Document,
	}

	tb.renderCache.Update(ctx, tb.Width)

	_, elementIndex, elementStart, oldLocalViewPos := tb.renderCache.Stats(tb.viewPosition)

	tview := tb.TextEngine.Resolve(tb.renderCache.GetElement(elementIndex))
	var newLocalViewPos int
	switch dir {
	case 0:
		newLocalViewPos = tview.MoveLeft(ctx, tb.renderCache.GetElement(elementIndex), oldLocalViewPos)
	case 1:
		newLocalViewPos = tview.MoveRight(ctx, tb.renderCache.GetElement(elementIndex), oldLocalViewPos)
	case 2:
		newLocalViewPos = tview.MoveUp(ctx, tb.renderCache.GetElement(elementIndex), oldLocalViewPos)
	case 3:
		newLocalViewPos = tview.MoveDown(ctx, tb.renderCache.GetElement(elementIndex), oldLocalViewPos)
	}

	if newLocalViewPos == -1 {
		tvLen := tview.MoveLength(ctx, tb.renderCache.GetElement(elementIndex))
		if dir == 3 {

			if pLinebaseView, ok := tview.(view.LinebaseTextView); ok && elementIndex+1 < tb.renderCache.GetItemCount() {

				nextView := tb.TextEngine.Resolve(tb.renderCache.GetElement(elementIndex + 1))
				relx := pLinebaseView.ConvertRelativeX(ctx, tb.renderCache.GetLayout(elementIndex), oldLocalViewPos)

				if linebaseTV, ok := nextView.(view.LinebaseTextView); ok {
					offset := linebaseTV.MoveFirstLine(ctx, tb.renderCache.GetElement(elementIndex+1), relx)

					tb.viewPosition = elementStart + tvLen + offset

					bPos := linebaseTV.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex+1), offset)
					tb.bytePos = bPos
				} else {
					tb.viewPosition = elementStart + tvLen

					// bPos := tview.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex), tvLen)
					bPos := nextView.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex+1), 0)
					tb.bytePos = bPos
				}
			} else {
				if elementIndex+1 < tb.renderCache.GetItemCount() {

					nextView := tb.TextEngine.Resolve(tb.renderCache.GetElement(elementIndex + 1))
					tb.viewPosition = elementStart + tvLen

					// bPos := tview.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex), tvLen)
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

				nView := tb.TextEngine.Resolve(tb.renderCache.GetElement(elementIndex + 1))
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
				pView := tb.TextEngine.Resolve(tb.renderCache.GetElement(elementIndex - 1))
				pViewLen := pView.MoveLength(ctx, tb.renderCache.GetElement(elementIndex-1))
				bPos := pView.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex-1), pViewLen-1)
				tb.bytePos = bPos
			} else {
				bPos := tview.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex), 0)
				tb.bytePos = bPos
			}
		} else if dir == 2 {

			if pLinebaseView, ok := tview.(view.LinebaseTextView); ok && elementIndex > 0 {

				prevView := tb.TextEngine.Resolve(tb.renderCache.GetElement(elementIndex - 1))
				relx := pLinebaseView.ConvertRelativeX(ctx, tb.renderCache.GetLayout(elementIndex), oldLocalViewPos)

				if linebaseTV, ok := prevView.(view.LinebaseTextView); ok {
					offset := linebaseTV.MoveLastLine(ctx, tb.renderCache.GetElement(elementIndex-1), relx)
					l := linebaseTV.MoveLength(ctx, tb.renderCache.GetElement(elementIndex-1))

					tb.viewPosition = elementStart - l + offset

					bPos := linebaseTV.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex-1), offset)
					tb.bytePos = bPos
				} else {
					tb.viewPosition = max(elementStart-1, 0)

					if elementIndex > 0 {
						pView := tb.TextEngine.Resolve(tb.renderCache.GetElement(elementIndex - 1))
						pViewLen := pView.MoveLength(ctx, tb.renderCache.GetElement(elementIndex-1))
						// bPos := tview.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex), tvLen)
						bPos := pView.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex-1), pViewLen-1)
						tb.bytePos = bPos
					}
				}
			} else {
				tb.viewPosition = max(elementStart-1, 0)

				if elementIndex > 0 {
					pView := tb.TextEngine.Resolve(tb.renderCache.GetElement(elementIndex - 1))
					pViewLen := pView.MoveLength(ctx, tb.renderCache.GetElement(elementIndex-1))
					// bPos := tview.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex), tvLen)
					bPos := pView.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex-1), pViewLen-1)
					tb.bytePos = bPos
				}
			}
		}
	} else {
		// moves := newLocalViewPos - oldLocalViewPos
		tb.viewPosition = elementStart + newLocalViewPos
		bPos := tview.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex), newLocalViewPos)
		tb.bytePos = bPos
	}

	//tb.viewPosition = min(max(tb.viewPosition, 0), ttl-1)
}

func (tb *TextBox) modelToView() (ElementIndex int, ViewStart int, ViewLocalPos int) {
	ctx := view.Context{
		Resolver: tb.TextEngine,
		Document: tb.Document,
	}

	if tb.renderCache.GetItemCount() == 0 {
		return 0, 0, 0
	}

	elementIndex := -1
	for i := 0; i < tb.renderCache.GetItemCount(); i++ {
		element := tb.renderCache.GetElement(i)
		r := element.GetRange(0)
		st := r.StartPosition
		ed := r.EndPosition

		// 空行のフォロー
		if st.Row == ed.Row && st.Column == ed.Column {
			if tb.bytePos.StartPosition.Row == st.Row && tb.bytePos.StartPosition.Column == st.Column {
				elementIndex = i
				break
			}
		}

		if tb.bytePos.StartPosition.Row >= st.Row && tb.bytePos.StartPosition.Row <= ed.Row {
			if tb.bytePos.Bytes == 0 {
				col := max(tb.bytePos.StartPosition.Column-1, 0)
				if col >= st.Column && (col < ed.Column || ed.Row > st.Row) {
					elementIndex = i
					break
				}
			} else {
				if tb.bytePos.StartPosition.Column >= st.Column && (tb.bytePos.StartPosition.Column < ed.Column || ed.Row > st.Row) {
					elementIndex = i
					break
				}
			}
		}
	}

	viewStart := 0
	for i := 0; i < elementIndex; i++ {
		element := tb.renderCache.GetElement(i)
		textView := tb.TextEngine.Resolve(element)

		viewStart += textView.MoveLength(ctx, element)
	}

	element := tb.renderCache.GetElement(elementIndex)
	textView := tb.TextEngine.Resolve(element)
	viewLocalPos := textView.ConvertViewLocalPos(ctx, tb.renderCache.GetLayout(elementIndex), tb.bytePos.StartPosition)

	return elementIndex, viewStart, viewLocalPos
}

func (tb *TextBox) viewToModel() view.CharacterReference {
	ctx := view.Context{
		Resolver: tb.TextEngine,
		Document: tb.Document,
	}

	_, ei, _, _ := tb.renderCache.Stats(tb.viewPosition)
	element := tb.renderCache.GetElement(ei)
	textView := tb.TextEngine.Resolve(element)
	return textView.ConvertModel(ctx, tb.renderCache.GetLayout(ei), 0)
}

func (tb *TextBox) InsertString(s string) {
	ctx := view.Context{
		Resolver: tb.TextEngine,
		Document: tb.Document,
	}

	tb.renderCache.Update(ctx, tb.Width)

	_, ei, _, _ := tb.renderCache.Stats(tb.viewPosition)
	if ge, ok := tb.renderCache.GetElement(ei).(*model.GhostElement); ok {
		lines := strings.Repeat("\n", ge.Index+1)
		tb.viewPosition -= ge.Index + 1
		lastElement := tb.renderCache.GetElement(tb.renderCache.GetItemCount() - 10 - 1)
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

	elementIndex, viewStart, viewLocalPos := tb.modelToView()
	element := tb.renderCache.GetElement(elementIndex)
	textView := tb.TextEngine.Resolve(element)

	if viewLocalPos == 0 {
		if strings.HasSuffix(s, "\n") {
			if _, ok := element.(*litemark.HeadingElement); ok {
				if elementIndex == 0 {
					tb.bytePos.StartPosition.Row = 0
					tb.bytePos.StartPosition.Column = 0

				} else {
					//element = tb.renderCache.GetElement(elementIndex - 1)
					//textView = tb.TextEngine.Resolve(element)
					//viewLocalPos = textView.MoveLength(ctx, element) - 1
					//bPos := textView.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex-1), viewLocalPos)
					//tb.bytePos = bPos
					tb.bytePos.StartPosition.Column = 0
					tb.bytePos.Bytes = 0
				}
			}
		}
	}

	tb.Document.InsertString(tb.bytePos.StartPosition.Row, tb.bytePos.StartPosition.Column, s)

	breakLine := strings.Contains(s, "\n")
	if strings.HasPrefix(s, "\n") {
		tb.bytePos.Bytes = 0
	}

	tb.renderCache.Update(ctx, tb.Width)

	elementIndex, viewStart, viewLocalPos = tb.modelToView()
	element = tb.renderCache.GetElement(elementIndex)
	textView = tb.TextEngine.Resolve(element)

	moves := text.GraphemeLength(s)
	for i := 0; i < moves; i++ {
		viewLocalPos = textView.MoveRight(ctx, element, viewLocalPos)

		if viewLocalPos == -1 {
			// 文字挿入によってビューが分割された場合
			// 移動可能な回数が1回かつテキストが存在しない場合は空行とみなす
			// その場合には次の行へ降りる
			// 移動可能回数が1回かつテキストが存在する場合は行を継続する
			// 移動可能回数が2回以上の場合、そのビューの開始位置までジャンプする
			if textView.MoveLength(ctx, element) == 1 {
				if len(ctx.GetText(element)) == 0 {
					tb.viewPosition = viewStart + 1
				} else {
					tb.viewPosition = viewStart
				}
			} else {
				if _, ok := textView.(*litemark.TextView); ok && !breakLine {
					// *a* このときは-1
					// *a*NL このときは0
					tb.viewPosition = viewStart + textView.MoveLength(ctx, element) - 1
				} else {
					tb.viewPosition = viewStart + textView.MoveLength(ctx, element)
				}

				if tb.viewPosition >= tb.renderCache.Total() {
					tb.viewPosition = tb.renderCache.Total() - 1
				}
			}
			_, elementIndex, viewStart, viewLocalPos = tb.renderCache.Stats(tb.viewPosition)
			element = tb.renderCache.GetElement(elementIndex)
			textView = tb.TextEngine.Resolve(element)

			bPos := textView.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex), viewLocalPos)
			tb.bytePos = bPos

			// tb.viewPosition = viewStart + textView.MoveLength(ctx, element)
			//break
		} else {
			tb.viewPosition = viewStart + viewLocalPos

			if viewLocalPos == textView.MoveLength(ctx, element) {

				if elementIndex+1 < tb.renderCache.GetItemCount() {
					element = tb.renderCache.GetElement(elementIndex + 1)
					textView = tb.TextEngine.Resolve(element)
					tb.bytePos = textView.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex+1), 0)
				} else {
					bPos := textView.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex), viewLocalPos)
					tb.bytePos = bPos
				}
			} else {
				//tb.bytePos = textView.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex), viewLocalPos).StartPosition

				bPos := textView.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex), viewLocalPos)
				tb.bytePos = bPos
			}
		}
	}
}

func (tb *TextBox) RemoveChar() {
	if tb.viewPosition == 0 {
		return
	}

	_, ei, _, _ := tb.renderCache.Stats(tb.viewPosition)
	if ge, ok := tb.renderCache.GetElement(ei).(*model.GhostElement); ok {
		tb.viewPosition -= ge.Index + 1
		return
	}

	ctx := view.Context{
		Resolver: tb.TextEngine,
		Document: tb.Document,
	}

	tb.renderCache.Update(ctx, tb.Width)

	elementIndex, viewStart, viewLocalPos := tb.modelToView()
	element := tb.renderCache.GetElement(elementIndex)
	textView := tb.TextEngine.Resolve(element)

	if textView.MoveLength(ctx, element) == 2 {
		if _, ok := element.(*litemark.HeadingElement); ok {
			r := element.GetRange(0)
			bPos := view.CharacterReference{
				StartPosition: r.StartPosition,
				Bytes:         (r.EndPosition.Column - r.StartPosition.Column),
			}
			tb.bytePos = bPos
			tb.Document.Remove(bPos.StartPosition.Row, bPos.StartPosition.Column, max(bPos.Bytes, 1))
			tb.viewPosition = viewStart
			return
		} else if txView, ok := textView.(*litemark.TextView); ok {
			combine, bPos := txView.RemoveCombine(ctx, tb.renderCache.GetLayout(elementIndex), viewLocalPos)
			if combine {
				elementIndex, viewStart, viewLocalPos = tb.modelToView()
				element = tb.renderCache.GetElement(elementIndex)
				textView = tb.TextEngine.Resolve(element)
				tb.bytePos = bPos
				tb.Document.Remove(bPos.StartPosition.Row, bPos.StartPosition.Column, bPos.Bytes)
				tb.bytePos.Bytes = 0
				tb.renderCache.Update(ctx, tb.Width)

				elementIndex, viewStart, viewLocalPos = tb.modelToView()
				tb.viewPosition = viewStart + viewLocalPos
				return
			}
		}
	} else if _, ok := element.(*litemark.CodeBlockElement); ok && viewLocalPos == 0 {
		r := element.GetRange(0)
		sg := tb.Document.Read(r)
		bPos := view.CharacterReference{
			StartPosition: model.Position{
				Row:    r.StartPosition.Row,
				Column: sg.GetSpan(0).EndColumn - 1,
			},
			Bytes: 1,
		}
		tb.bytePos = bPos
		tb.Document.Remove(bPos.StartPosition.Row, bPos.StartPosition.Column, max(bPos.Bytes, 1))
		//		tb.bytePos.StartPosition.Column++
		tb.bytePos.Bytes = 0

		tb.renderCache.Update(ctx, tb.Width)

		_, vs, vl := tb.modelToView()
		tb.viewPosition = vs + vl
		return
	} else if txView, ok := textView.(*litemark.TextView); ok {
		combine, bPos := txView.RemoveCombine(ctx, tb.renderCache.GetLayout(elementIndex), viewLocalPos)
		if combine {
			elementIndex, viewStart, viewLocalPos = tb.modelToView()
			element = tb.renderCache.GetElement(elementIndex)
			textView = tb.TextEngine.Resolve(element)
			tb.bytePos = bPos
			tb.Document.Remove(bPos.StartPosition.Row, bPos.StartPosition.Column, bPos.Bytes)
			tb.bytePos.Bytes = 0
			tb.renderCache.Update(ctx, tb.Width)

			elementIndex, viewStart, viewLocalPos = tb.modelToView()
			tb.viewPosition = viewStart + viewLocalPos
			return
		}
	}

	newViewLocalPos := textView.MoveLeft(ctx, element, viewLocalPos)
	if newViewLocalPos == -1 {
		element = tb.renderCache.GetElement(elementIndex - 1)

		if _, ok := element.(*litemark.CodeBlockElement); ok {
			r := element.GetRange(0)
			bPos := view.CharacterReference{
				StartPosition: model.Position{
					Row:    r.EndPosition.Row,
					Column: r.EndPosition.Column - 1,
				},
				Bytes: 1,
			}
			tb.bytePos = bPos
			tb.Document.Remove(bPos.StartPosition.Row, bPos.StartPosition.Column, max(bPos.Bytes, 1))

			//tb.bytePos.StartPosition.Column--
			tb.bytePos.Bytes = 0

			tb.renderCache.Update(ctx, tb.Width)

			_, vs, vl := tb.modelToView()
			tb.viewPosition = vs + vl
			return
		}
		if _, ok := element.(*litemark.HorizontalLineElement); ok {
			r := element.GetRange(0)
			bPos := view.CharacterReference{
				StartPosition: model.Position{
					Row:    r.EndPosition.Row,
					Column: r.EndPosition.Column - 1,
				},
				Bytes: 1,
			}
			tb.bytePos = bPos
			tb.Document.Remove(bPos.StartPosition.Row, bPos.StartPosition.Column, max(bPos.Bytes, 1))

			//tb.bytePos.StartPosition.Column--
			tb.bytePos.Bytes = 0

			tb.renderCache.Update(ctx, tb.Width)

			_, vs, vl := tb.modelToView()
			tb.viewPosition = vs + vl
			return
		}
		prevView := tb.TextEngine.Resolve(element)
		bPos := prevView.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex-1), prevView.MoveLength(ctx, element)-1)

		charRef := prevView.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex-1), prevView.MoveLength(ctx, element)-1)
		tb.bytePos = charRef

		tb.Document.Remove(bPos.StartPosition.Row, bPos.StartPosition.Column, max(bPos.Bytes, 1))

	} else {
		bPos := textView.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex), newViewLocalPos)

		charRef := textView.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex), newViewLocalPos)
		tb.bytePos = charRef

		tb.Document.Remove(bPos.StartPosition.Row, bPos.StartPosition.Column, max(bPos.Bytes, 1))
	}

	tb.bytePos.Bytes = 0

	tb.renderCache.Update(ctx, tb.Width)

	elementIndex = -1
	for i := 0; i < tb.renderCache.GetItemCount(); i++ {
		element := tb.renderCache.GetElement(i)
		r := element.GetRange(0)
		st := r.StartPosition
		ed := r.EndPosition

		// 空行のフォロー
		if st.Row == ed.Row && st.Column == ed.Column {
			if tb.bytePos.StartPosition.Row == st.Row && tb.bytePos.StartPosition.Column == st.Column {
				elementIndex = i
				break
			}
		}

		if tb.bytePos.StartPosition.Row >= st.Row && tb.bytePos.StartPosition.Row <= ed.Row {
			if tb.bytePos.Bytes == 0 {
				col := max(tb.bytePos.StartPosition.Column-1, 0)
				if col >= st.Column && (col < ed.Column || ed.Row > st.Row) {
					elementIndex = i
					break
				}
			} else {
				if tb.bytePos.StartPosition.Column >= st.Column && (tb.bytePos.StartPosition.Column < ed.Column || ed.Row > st.Row) {
					elementIndex = i
					break
				}
			}
		}
	}

	viewStart = 0
	for i := 0; i < elementIndex; i++ {
		element := tb.renderCache.GetElement(i)
		textView := tb.TextEngine.Resolve(element)

		viewStart += textView.MoveLength(ctx, element)
	}

	element = tb.renderCache.GetElement(elementIndex)
	textView = tb.TextEngine.Resolve(element)
	viewLocalPos = textView.ConvertViewLocalPos(ctx, tb.renderCache.GetLayout(elementIndex), tb.bytePos.StartPosition)

	tb.viewPosition = viewStart + viewLocalPos
}

func (tb *TextBox) MoveLeft() {
	tb.move(0)
}

func (tb *TextBox) MoveRight() {
	tb.move(1)
}

func (tb *TextBox) MoveUp() {
	tb.move(2)
}

func (tb *TextBox) MoveDown() {
	tb.move(3)
}

func (tb *TextBox) MoveReset() {
	ctx := view.Context{
		Resolver: tb.TextEngine,
		Document: tb.Document,
	}

	tb.renderCache.Update(ctx, tb.Width)
	tb.viewPosition = 0
	tb.bytePos = tb.viewToModel()
}

// BreakIter returns segment array by line, in consideration a wrap.
func (tb *TextBox) BreakIter() iter.Seq[presenter.Segment] {
	//buf := tb.Document.GetBuffer()
	//sb := strings.Builder{}

	return func(yield func(presenter.Segment) bool) {
		ctx := view.Context{
			Resolver: tb.TextEngine,
			Document: tb.Document,
		}
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
