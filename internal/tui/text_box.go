package tui

import (
	"iter"

	"github.com/desktopgame/ckro/internal/text"
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/desktopgame/ckro/internal/tui/presenter"
	"github.com/desktopgame/ckro/internal/tui/view"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

type Affinity int // 既存：インライン境界の内/外
const (
	Upstream Affinity = iota
	Downstream
)

type LineAffinity int // 新規：行境界（ハード改行）の手前/後ろ
const (
	BeforeBreak LineAffinity = iota // 現行行の末（改行の手前）に“寄りつく”
	AfterBreak                      // 次行先頭（改行の後）に“寄りつく”
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

	renderCache  TextRenderCache
	bytePos      view.CharacterReference
	affinity     Affinity
	lineAffinity LineAffinity
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
	view := tb.TextEngine.Resolve(tb.renderCache.GetElement(ei))
	y := 0
	for i := 0; i < ei; i++ {
		y += tb.renderCache.GetLayout(i).Height
	}
	_, vly := view.ConvertPos(ctx, tb.renderCache.GetLayout(ei), eoff)
	// ax := vlx
	ay := y + vly

	buf := tb.Document.GetBuffer()
	cursorRow := ay
	// cursorCol := ax

	if cursorRow >= buf.GetLineCount() {
		return 0, 0, ' ', nil
	}

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

// calculateCursorPosition is calculates the exact screen position considering line wrapping
func (tb *TextBox) calculateCursorPosition(line string, cursorCol int) (screenX int, additionalRows int) {
	if tb.Width <= 0 {
		return 0, 0
	}

	currentX := 0
	currentRow := 0
	clusters := text.GraphemeClusters(line)

	// カーソルが行末を超えている場合の処理
	if cursorCol >= len(clusters) {
		// 全ての文字を処理してから、カーソル位置を決定
		for _, cluster := range clusters {
			clusterWidth := tb.calculateClusterWidth(cluster, currentX)

			// 現在の行に収まるかチェック
			if currentX+clusterWidth > tb.Width {
				// 次の行に移動
				currentRow++
				currentX = 0
				// 行が変わったので幅を再計算
				clusterWidth = tb.calculateClusterWidth(cluster, currentX)
			}

			currentX += clusterWidth
		}

		// 行末の場合、最後の文字の後の位置
		if currentX >= tb.Width {
			currentRow++
			currentX = 0
		}

		return currentX, currentRow
	}

	// 通常の処理：指定された位置まで
	for i := 0; i < cursorCol; i++ {
		cluster := clusters[i]
		clusterWidth := tb.calculateClusterWidth(cluster, currentX)

		// 現在の行に収まるかチェック
		if currentX+clusterWidth > tb.Width {
			// 次の行に移動
			currentRow++
			currentX = 0
			// 行が変わったので幅を再計算
			clusterWidth = tb.calculateClusterWidth(cluster, currentX)
		}

		currentX += clusterWidth
	}

	// カーソルが特定の文字（タブなど）の上にある場合の特別処理
	if cursorCol < len(clusters) {
		cluster := clusters[cursorCol]
		if cluster == "\t" {
			// タブの場合、タブが次の行に移動するかチェック
			clusterWidth := tb.calculateClusterWidth(cluster, currentX)
			if currentX+clusterWidth > tb.Width {
				// タブが次の行に移動する場合、次の行の先頭を返す
				return 0, currentRow + 1
			}
		}
	}

	return currentX, currentRow
}

// calculateClusterWidth is calculates the display width of a cluster considering tabs
func (tb *TextBox) calculateClusterWidth(cluster string, currentX int) int {
	if cluster == "\t" {
		return text.TabWidth - (currentX % text.TabWidth)
	}
	return runewidth.StringWidth(cluster)
}

// calculateWrappedLines is calculates how many screen lines a text line takes
func (tb *TextBox) calculateWrappedLines(line string) int {
	if tb.Width <= 0 {
		return 1
	}

	currentX := 0
	currentRow := 1
	clusters := text.GraphemeClusters(line)

	for _, cluster := range clusters {
		clusterWidth := tb.calculateClusterWidth(cluster, currentX)

		// 現在の行に収まるかチェック
		if currentX+clusterWidth > tb.Width {
			// 次の行に移動
			currentRow++
			currentX = 0
			// 行が変わったので幅を再計算
			clusterWidth = tb.calculateClusterWidth(cluster, currentX)
		}

		currentX += clusterWidth
	}

	return currentRow
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
	tb.Document.MoveReset()
	tb.scrollX = 0
	tb.scrollY = 0
}

// TextFrame is print a frame of TextBox region.
func (tb *TextBox) TextFrame() {
	tb.Document.Clear()

	w := tb.Width
	h := tb.Height

	tb.Document.InsertString("*")
	for i := 0; i < w-2; i++ {
		tb.Document.InsertString("-")
	}
	tb.Document.InsertString("*")
	tb.Document.InsertLine()

	for i := 0; i < h-2; i++ {
		tb.Document.InsertString("|")
		for j := 0; j < w-2; j++ {
			tb.Document.InsertString(" ")
		}
		tb.Document.InsertString("|")
		tb.Document.InsertLine()
	}

	tb.Document.InsertString("*")
	for i := 0; i < w-2; i++ {
		tb.Document.InsertString("-")
	}
	tb.Document.InsertString("*")

	tb.Document.MoveReset()
}

// TextVertical is print vertical line.
func (tb *TextBox) TextVertical() {
	tb.Document.Clear()

	h := tb.Height

	for i := 0; i < h; i++ {
		tb.Document.InsertString("|\n")
	}
	tb.Document.RemoveChar()
	tb.Document.MoveReset()
}

// TextHorizontal is print horizontal line.
func (tb *TextBox) TextHorizontal() {
	tb.Document.Clear()

	w := tb.Width

	for i := 0; i < w; i++ {
		tb.Document.InsertString("-")
	}
	tb.Document.MoveReset()

}

// TextClear is do reset to content.
func (tb *TextBox) TextClear() {
	tb.Document.Clear()
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

			if pLinebaseView, ok := tview.(view.LinebaseTextView); ok {

				nextView := tb.TextEngine.Resolve(tb.renderCache.GetElement(elementIndex + 1))
				relx := pLinebaseView.ConvertRelativeX(ctx, tb.renderCache.GetElement(elementIndex), oldLocalViewPos)

				if linebaseTV, ok := nextView.(view.LinebaseTextView); ok {
					offset := linebaseTV.MoveFirstLine(ctx, tb.renderCache.GetElement(elementIndex+1), relx)

					tb.viewPosition = elementStart + tvLen + offset

					bPos := tview.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex), tvLen+offset)
					tb.bytePos = bPos

					if bPos.Bytes == 0 {
						tb.lineAffinity = AfterBreak
					} else {
						tb.lineAffinity = BeforeBreak
					}
				} else {
					tb.viewPosition = elementStart + tvLen

					// bPos := tview.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex), tvLen)
					bPos := nextView.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex+1), 0)
					tb.bytePos = bPos

					if bPos.Bytes == 0 {
						tb.lineAffinity = AfterBreak
					} else {
						tb.lineAffinity = BeforeBreak
					}
				}
			} else {
				tb.viewPosition = elementStart + tvLen

				bPos := tview.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex), tvLen)
				tb.bytePos = bPos

				if bPos.Bytes == 0 {
					tb.lineAffinity = AfterBreak
				} else {
					tb.lineAffinity = BeforeBreak
				}
			}
		} else if dir == 1 {
			tb.viewPosition = elementStart + tvLen

			bPos := tview.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex), tvLen)
			tb.bytePos = bPos

			if bPos.Bytes == 0 {
				tb.lineAffinity = AfterBreak
			}
		}

		if dir == 0 {
			tb.viewPosition = max(elementStart-1, 0)
		} else if dir == 2 {

			if pLinebaseView, ok := tview.(view.LinebaseTextView); ok {

				prevView := tb.TextEngine.Resolve(tb.renderCache.GetElement(elementIndex - 1))
				relx := pLinebaseView.ConvertRelativeX(ctx, tb.renderCache.GetElement(elementIndex), oldLocalViewPos)

				if linebaseTV, ok := prevView.(view.LinebaseTextView); ok {
					offset := linebaseTV.MoveLastLine(ctx, tb.renderCache.GetElement(elementIndex-1), relx)
					l := linebaseTV.MoveLength(ctx, tb.renderCache.GetElement(elementIndex-1))

					tb.viewPosition = elementStart - l + offset

					bPos := linebaseTV.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex), offset)
					tb.bytePos = bPos

					if bPos.Bytes == 0 {
						tb.lineAffinity = AfterBreak
					} else {
						tb.lineAffinity = BeforeBreak
					}
				} else {
					tb.viewPosition = max(elementStart-1, 0)

					if elementIndex > 0 {
						pView := tb.TextEngine.Resolve(tb.renderCache.GetElement(elementIndex - 1))
						pViewLen := pView.MoveLength(ctx, tb.renderCache.GetElement(elementIndex-1))
						// bPos := tview.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex), tvLen)
						bPos := pView.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex-1), pViewLen-1)
						tb.bytePos = bPos

						if bPos.Bytes == 0 {
							tb.lineAffinity = AfterBreak
						}
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

					if bPos.Bytes == 0 {
						tb.lineAffinity = AfterBreak
					}
				}
			}
		}
	} else {
		// moves := newLocalViewPos - oldLocalViewPos
		tb.viewPosition = elementStart + newLocalViewPos
		bPos := tview.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex), newLocalViewPos)
		tb.bytePos = bPos

		if bPos.Bytes == 0 {
			tb.lineAffinity = AfterBreak
		} else {
			tb.lineAffinity = BeforeBreak
		}
	}

	//tb.viewPosition = min(max(tb.viewPosition, 0), ttl-1)
}

func (tb *TextBox) InsertString(s string) {
	ctx := view.Context{
		Resolver: tb.TextEngine,
		Document: tb.Document,
	}

	tb.Document.WriteString(tb.bytePos.StartPosition.Row, tb.bytePos.StartPosition.Column, s)
	tb.renderCache.Update(ctx, tb.Width)

	elementIndex := -1
	beforeView := false
	for i := 0; i < tb.renderCache.GetItemCount(); i++ {
		element := tb.renderCache.GetElement(i)
		r := element.GetRange(0)
		st := r.StartPosition
		ed := r.EndPosition

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

	if tb.lineAffinity == BeforeBreak {
		beforeView = true
	}

	viewStart := 0
	for i := 0; i < elementIndex; i++ {
		element := tb.renderCache.GetElement(i)
		textView := tb.TextEngine.Resolve(element)

		viewStart += textView.MoveLength(ctx, element)
	}

	element := tb.renderCache.GetElement(elementIndex)
	textView := tb.TextEngine.Resolve(element)

	var viewLocalPos int
	if beforeView {
		// viewLocalPos = textView.MoveLength(ctx, tb.renderCache.GetElement(elementIndex)) - 1
		viewLocalPos = textView.ConvertViewLocalPos(ctx, tb.renderCache.GetLayout(elementIndex), tb.bytePos.StartPosition)
	} else {
		viewLocalPos = textView.ConvertViewLocalPos(ctx, tb.renderCache.GetLayout(elementIndex), tb.bytePos.StartPosition)
	}

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
				if blankable, ok := textView.(view.BlankTextView); ok && blankable.IsBlank(ctx, tb.renderCache.GetLayout(elementIndex)) {
					tb.viewPosition = viewStart + 1
				} else {
					tb.viewPosition = viewStart
					tb.lineAffinity = BeforeBreak
				}
			} else {
				tb.viewPosition = viewStart + textView.MoveLength(ctx, element)
			}
			_, elementIndex, viewStart, viewLocalPos = tb.renderCache.Stats(tb.viewPosition)
			element = tb.renderCache.GetElement(elementIndex)
			textView = tb.TextEngine.Resolve(element)

			bPos := textView.ConvertModel(ctx, tb.renderCache.GetLayout(elementIndex), viewLocalPos)
			if bPos.Bytes == 0 {
				//bPos.StartPosition.Row++
				//bPos.StartPosition.Column = 0
				tb.lineAffinity = BeforeBreak
			}
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
				if bPos.Bytes == 0 {
					//bPos.StartPosition.Row++
					//bPos.StartPosition.Column = 0
					tb.lineAffinity = BeforeBreak
				}
				tb.bytePos = bPos
			}
		}
	}
}

func (tb *TextBox) RemoveChar() {
	if tb.viewPosition == 0 {
		return
	}

	ctx := view.Context{
		Resolver: tb.TextEngine,
		Document: tb.Document,
	}

	tb.renderCache.Update(ctx, tb.Width)

	_, elementIndex, viewStart, viewLocalPos := tb.renderCache.Stats(tb.viewPosition)
	element := tb.renderCache.GetElement(elementIndex)
	tl := tb.renderCache.GetLayout(elementIndex)
	textView := tb.TextEngine.Resolve(element)

	newViewLocalPos := textView.MoveLeft(ctx, element, viewLocalPos)

	if newViewLocalPos >= 0 {
		viewLocalPos = newViewLocalPos
		position := textView.ConvertModel(ctx, tl, viewLocalPos)
		position.Bytes = max(position.Bytes, 1)

		tb.Document.Remove(position.StartPosition.Row, position.StartPosition.Column, position.Bytes)
		tb.renderCache.Update(ctx, tb.Width)
		tb.viewPosition = viewStart + viewLocalPos
	} else {
		element = tb.renderCache.GetElement(elementIndex - 1)
		tl = tb.renderCache.GetLayout(elementIndex - 1)
		textView = tb.TextEngine.Resolve(element)
		viewLocalPos = textView.MoveLength(ctx, element) - 1

		position := textView.ConvertModel(ctx, tl, viewLocalPos)
		position.Bytes = max(position.Bytes, 1)

		tb.Document.Remove(position.StartPosition.Row, position.StartPosition.Column, position.Bytes)
		tb.renderCache.Update(ctx, tb.Width)
		tb.viewPosition = viewStart - 1
	}
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
	tb.viewPosition = 0
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
				segment := presenter.Segment{
					TextLayout:    entry,
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

// WrappedLineCount returns count of lines, in consideration a wrap.
func (tb *TextBox) WrappedLineCount() int {
	lc := 0
	buf := tb.Document.GetBuffer()

	for i := 0; i < buf.GetLineCount(); i++ {
		line := buf.GetLineAt(i)
		lineContent := line.GetContent()
		lc += tb.calculateWrappedLines(lineContent)
	}
	return lc
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
