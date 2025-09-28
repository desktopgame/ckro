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

	documentVersion uint
	elements        []model.Element
	layoutCache     []*view.TextLayout
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
func (tb *TextBox) CursorPosition2() (X int, Y int, Rune rune, Combine []rune) {
	buf := tb.Document.GetBuffer()
	cursorRow := tb.Document.GetCursorRow()
	cursorCol := tb.Document.GetCursorColumn()

	if cursorRow >= buf.GetLineCount() {
		return 0, 0, ' ', nil
	}

	// カーソルがある行までの画面行数を計算
	screenY := 0
	for i := 0; i < cursorRow; i++ {
		line := buf.GetLineAt(i).GetContent()
		screenY += tb.calculateWrappedLines(line)
	}

	// カーソルがある行での位置を正確に計算
	cursorLine := buf.GetLineAt(cursorRow).GetContent()
	screenX, additionalRows := tb.calculateCursorPosition(cursorLine, cursorCol)
	screenY += additionalRows

	// カーソル位置の文字を取得
	var currentRune rune = ' '
	var combining []rune

	if cursorCol < text.GraphemeLength(cursorLine) {
		// カーソル位置に文字がある場合
		clusters := text.GraphemeClusters(cursorLine)
		if cursorCol < len(clusters) {
			cluster := clusters[cursorCol]
			runes := []rune(cluster)
			if len(runes) > 0 {
				currentRune = runes[0]
				if len(runes) > 1 {
					combining = runes[1:]
				}
			}
		}
	}

	return screenX, screenY, currentRune, combining
}

// CursorPosition returns position of cursor.
// TODO: refactor
func (tb *TextBox) CursorPosition() (X int, Y int, Rune rune, Combine []rune) {
	ctx := view.Context{
		Resolver: tb.TextEngine,
		Document: tb.Document,
	}

	_, ei, _, eoff := tb.currentViewState(ctx, tb.elements)
	view := tb.TextEngine.Resolve(tb.elements[ei])
	y := 0
	for i := 0; i < ei; i++ {
		y += tb.layoutCache[i].Height
	}
	vlx, vly := view.ConvertPos(ctx, tb.elements[ei], eoff)
	ax := vlx
	ay := y + vly

	buf := tb.Document.GetBuffer()
	cursorRow := ay
	cursorCol := ax

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
		screenY += tb.layoutCache[i].Height
	}

	// カーソルがある行での位置を正確に計算
	cursorLine := buf.GetLineAt(cursorRow).GetContent()
	relx, rely := tb.TextEngine.Resolve(tb.elements[ei]).ConvertPos(ctx, tb.elements[ei], eoff)
	screenX := relx
	screenY += rely
	//screenX, additionalRows :=
	//screenY += additionalRows

	// カーソル位置の文字を取得
	var currentRune rune = ' '
	var combining []rune

	if cursorCol < text.GraphemeLength(cursorLine) {
		// カーソル位置に文字がある場合
		clusters := text.GraphemeClusters(cursorLine)
		if cursorCol < len(clusters) {
			cluster := clusters[cursorCol]
			runes := []rune(cluster)
			if len(runes) > 0 {
				currentRune = runes[0]
				if len(runes) > 1 {
					combining = runes[1:]
				}
			}
		}
	}

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

	tb.layout()

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

	ctx := view.Context{
		Resolver: tb.TextEngine,
		Document: tb.Document,
	}

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

func (tb *TextBox) currentViewState(ctx view.Context, elements []model.Element) (TotalViewLen int, CurrentElementIndex int, CurrentElementStart int, CurrentElementOffset int) {
	totalViewLen := 0
	elementIndex := 0
	elementStart := -1
	oldLocalViewPos := 0
	for i, elem := range elements {
		viewStart := totalViewLen
		view := tb.TextEngine.Resolve(elem)
		viewLen := view.MoveLength(ctx, elem)
		viewEnd := viewStart + viewLen

		if tb.viewPosition >= viewStart && tb.viewPosition < viewEnd {
			elementIndex = i
			oldLocalViewPos = tb.viewPosition - viewStart
			elementStart = viewStart
		}
		totalViewLen += viewLen
	}
	return totalViewLen, elementIndex, elementStart, oldLocalViewPos
}

func (tb *TextBox) move(dir int) {
	tb.layout()

	ctx := view.Context{
		Resolver: tb.TextEngine,
		Document: tb.Document,
	}

	ttl, elementIndex, elementStart, oldLocalViewPos := tb.currentViewState(ctx, tb.elements)

	tview := tb.TextEngine.Resolve(tb.elements[elementIndex])
	var newLocalViewPos int
	switch dir {
	case 0:
		newLocalViewPos = tview.MoveLeft(ctx, tb.elements[elementIndex], oldLocalViewPos)
	case 1:
		newLocalViewPos = tview.MoveRight(ctx, tb.elements[elementIndex], oldLocalViewPos)
	case 2:
		newLocalViewPos = tview.MoveUp(ctx, tb.elements[elementIndex], oldLocalViewPos)
	case 3:
		newLocalViewPos = tview.MoveDown(ctx, tb.elements[elementIndex], oldLocalViewPos)
	}

	if newLocalViewPos == -1 {
		if dir == 3 {

			if pLinebaseView, ok := tview.(view.LinebaseTextView); ok {

				nextView := tb.TextEngine.Resolve(tb.elements[elementIndex+1])
				relx := pLinebaseView.ConvertRelativeX(ctx, tb.elements[elementIndex], oldLocalViewPos)

				if linebaseTV, ok := nextView.(view.LinebaseTextView); ok {
					offset := linebaseTV.MoveFirstLine(ctx, tb.elements[elementIndex+1], relx)

					tb.viewPosition = elementStart + tview.MoveLength(ctx, tb.elements[elementIndex]) + offset

				} else {
					tb.viewPosition = elementStart + tview.MoveLength(ctx, tb.elements[elementIndex])
				}
			} else {
				tb.viewPosition = elementStart + tview.MoveLength(ctx, tb.elements[elementIndex])
			}
		} else if dir == 1 {
			tb.viewPosition = elementStart + tview.MoveLength(ctx, tb.elements[elementIndex])
		}

		if dir == 0 {
			tb.viewPosition = max(elementStart-1, 0)
		} else if dir == 2 {

			if pLinebaseView, ok := tview.(view.LinebaseTextView); ok {

				prevView := tb.TextEngine.Resolve(tb.elements[elementIndex-1])
				relx := pLinebaseView.ConvertRelativeX(ctx, tb.elements[elementIndex], oldLocalViewPos)

				if linebaseTV, ok := prevView.(view.LinebaseTextView); ok {
					offset := linebaseTV.MoveLastLine(ctx, tb.elements[elementIndex-1], relx)
					l := linebaseTV.MoveLength(ctx, tb.elements[elementIndex-1])

					tb.viewPosition = elementStart - l + offset

				} else {
					tb.viewPosition = max(elementStart-1, 0)
				}
			} else {
				tb.viewPosition = max(elementStart-1, 0)
			}
		}
	} else {
		// moves := newLocalViewPos - oldLocalViewPos
		tb.viewPosition = elementStart + newLocalViewPos
	}

	tb.viewPosition = min(max(tb.viewPosition, 0), ttl-1)
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

func (tb *TextBox) layout() {
	if tb.documentVersion > 0 && tb.documentVersion == tb.Document.GetVersion() {
		return
	}

	ctx := view.Context{
		Resolver: tb.TextEngine,
		Document: tb.Document,
	}

	elements := tb.Document.Render()

	entries := []*view.TextLayout{}
	for i := 0; i < len(elements); i++ {
		element := elements[i]
		textView := tb.TextEngine.Resolve(element)
		tl := textView.MinimumSize(ctx, element, tb.Width, 9999)
		entries = append(entries, tl)
	}

	tb.elements = elements
	tb.layoutCache = entries
	tb.documentVersion = tb.Document.GetVersion()
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

		entries := tb.layoutCache
		viewLine := 0
		for _, entry := range entries {
			element := entry.Element

			lineWrap := entry.MinimumWidth > tb.Width

			if lineWrap {
				elementRange := element.GetRange(0)
				lineCount := elementRange.EndPosition.Row - elementRange.StartPosition.Row + 1
				for j := 0; j < lineCount; j++ {
					lineNo := elementRange.StartPosition.Row + j
					line := tb.Document.GetBuffer().GetLineAt(lineNo)
					lineWidth := text.DisplayWidth(line.GetContent())

					if lineWidth <= tb.Width {
						//*
						segment := presenter.Segment{
							TextLayout: &view.TextLayout{
								Element: &model.PlainElement{
									Range: model.Range{
										StartPosition: model.Position{
											Row:    lineNo,
											Column: 0,
										},
										EndPosition: model.Position{
											Row:    lineNo,
											Column: len(line.GetContent()),
										},
									},
								},
							},
							ModelLine: lineNo,
							ViewLine:  viewLine,
						}
						if !yield(segment) {
							return
						}
						viewLine++
						//*/
					} else {
						x := 0
						bytes := 0
						bytes2 := 0

						clusters := text.GraphemeClusters(line.GetContent())
						for _, cluster := range clusters {
							runes := []rune(cluster)

							if cluster == "\t" {
								w := text.TabWidth - (x % text.TabWidth)
								// w := text.TabWidth
								if x+w > tb.Width {
									segment := presenter.Segment{
										TextLayout: &view.TextLayout{
											Element: &model.PlainElement{
												Range: model.Range{
													StartPosition: model.Position{
														Row:    lineNo,
														Column: bytes,
													},
													EndPosition: model.Position{
														Row:    lineNo,
														Column: bytes2,
													},
												},
											},
										},
										ModelLine: lineNo,
										ViewLine:  viewLine,
									}
									if !yield(segment) {
										return
									}

									viewLine++
									bytes = bytes2
									x = 0
								}
								x += w
							} else if len(runes) > 0 {
								mainRune := runes[0]
								width := runewidth.RuneWidth(mainRune)

								if x+width > tb.Width {
									segment := presenter.Segment{
										TextLayout: &view.TextLayout{
											Element: &model.PlainElement{
												Range: model.Range{
													StartPosition: model.Position{
														Row:    lineNo,
														Column: bytes,
													},
													EndPosition: model.Position{
														Row:    lineNo,
														Column: bytes2,
													},
												},
											},
										},
										ModelLine: lineNo,
										ViewLine:  viewLine,
									}
									if !yield(segment) {
										return
									}

									viewLine++
									bytes = bytes2
									x = 0
								}
								x += width
							}

							bytes2 += len(cluster)
						}

						if bytes2 > bytes {
							segment := presenter.Segment{
								TextLayout: &view.TextLayout{
									Element: &model.PlainElement{
										Range: model.Range{
											StartPosition: model.Position{
												Row:    lineNo,
												Column: bytes,
											},
											EndPosition: model.Position{
												Row:    lineNo,
												Column: bytes2,
											},
										},
									},
								},
								ModelLine: lineNo,
								ViewLine:  viewLine,
							}
							if !yield(segment) {
								return
							}
							viewLine++
						}
					}
				}
			} else {
				textView := tb.TextEngine.Resolve(entry.Element)
				textView.Layout(ctx, entry, 0, viewLine, entry.MinimumWidth, entry.MinimumHeight)
				for j := 0; j < entry.MinimumHeight; j++ {
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
