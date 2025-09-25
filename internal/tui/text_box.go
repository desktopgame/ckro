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

type textSegment struct {
	textLayout *view.TextLayout
	segment    presenter.Segment
}

// TextBox is editable text widget.
// TextBox has region of rect, rendering text within that range.
// long line is always wrap at right end of region.
type TextBox struct {
	Document   model.Document
	X          int
	Y          int
	Width      int
	Height     int
	TextEngine TextEngine
	ShowCursor bool
	scrollX    int
	scrollY    int
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
	tb.TextEngine = &MarkdownTextEngine{}
	tb.scrollX = 0
	tb.scrollY = 0
}

// CursorPosition returns position of cursor.
// TODO: refactor
func (tb *TextBox) CursorPosition() (X int, Y int, Rune rune, Combine []rune) {
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
			view.Draw(tb.TextEngine, textSegment.TextLayout, &clip)
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

// BreakIter returns segment array by line, in consideration a wrap.
func (tb *TextBox) BreakIter() iter.Seq[presenter.Segment] {
	//buf := tb.Document.GetBuffer()
	//sb := strings.Builder{}

	return func(yield func(presenter.Segment) bool) {
		elements := tb.Document.Render()

		entries := []*view.TextLayout{}
		for i := 0; i < len(elements); i++ {
			element := elements[i]
			view := tb.TextEngine.Resolve(element)

			newLayout := view.Layout(tb.TextEngine, element, tb.Width)
			entries = append(entries, newLayout)
		}

		viewLine := 0
		for i := 0; i < len(entries); i++ {
			entry := entries[i]
			textView := tb.TextEngine.Resolve(entry.Element)
			height := textView.Height(tb.TextEngine, entry)

			lineWrap := false
			for j := 0; j < height; j++ {
				if textView.Width(tb.TextEngine, entry, j) > tb.Width {
					lineWrap = true
					break
				}
			}

			if lineWrap {
				lineCount := entry.Element.GetEndPosition().Row - entry.Element.GetStartPosition().Row + 1
				for j := 0; j < lineCount; j++ {
					lineNo := entry.Element.GetStartPosition().Row + j
					line := tb.Document.GetBuffer().GetLineAt(lineNo)
					lineWidth := text.DisplayWidth(line.GetContent())

					if lineWidth <= tb.Width {
						//*
						segment := presenter.Segment{
							TextLayout: &view.TextLayout{
								Element: &model.ParagraphElement{
									Line: line.GetContent(),
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
						startX := 0
						sb := strings.Builder{}

						clusters := text.GraphemeClusters(line.GetContent())
						for _, cluster := range clusters {
							runes := []rune(cluster)

							if cluster == "\t" {
								w := text.TabWidth - (x % text.TabWidth)
								// w := text.TabWidth
								if x+w > tb.Width {
									segment := presenter.Segment{
										TextLayout: &view.TextLayout{
											Element: &model.ParagraphElement{
												Line: sb.String(),
												StartPosition: model.Position{
													Row:    lineNo,
													Column: startX,
												},
												EndPosition: model.Position{
													Row:    lineNo,
													Column: x,
												},
											},
										},
										ModelLine: lineNo,
										ViewLine:  viewLine,
									}
									if !yield(segment) {
										return
									}
									sb.Reset()

									viewLine++
									startX = x
									x = 0
								}
								sb.WriteString(cluster)
								x += w
							} else if len(runes) > 0 {
								mainRune := runes[0]
								width := runewidth.RuneWidth(mainRune)

								if x+width > tb.Width {
									segment := presenter.Segment{
										TextLayout: &view.TextLayout{
											Element: &model.ParagraphElement{
												Line: sb.String(),
												StartPosition: model.Position{
													Row:    lineNo,
													Column: startX,
												},
												EndPosition: model.Position{
													Row:    lineNo,
													Column: x,
												},
											},
										},
										ModelLine: lineNo,
										ViewLine:  viewLine,
									}
									if !yield(segment) {
										return
									}
									sb.Reset()

									viewLine++
									startX = x
									x = 0
								}
								sb.WriteString(cluster)
								x += width
							}
							if x > tb.Width {
								segment := presenter.Segment{
									TextLayout: &view.TextLayout{
										Element: &model.ParagraphElement{
											Line: sb.String(),
											StartPosition: model.Position{
												Row:    lineNo,
												Column: startX,
											},
											EndPosition: model.Position{
												Row:    lineNo,
												Column: x,
											},
										},
									},
									ModelLine: lineNo,
									ViewLine:  viewLine,
								}
								if !yield(segment) {
									return
								}
								sb.Reset()

								viewLine++
								startX = x
								x = 0
							}
							//sb.WriteString(cluster)
						}
						segment := presenter.Segment{
							TextLayout: &view.TextLayout{
								Element: &model.ParagraphElement{
									Line:          sb.String(),
									StartPosition: model.Position{},
									EndPosition:   model.Position{},
								},
							},
							ModelLine: lineNo,
							ViewLine:  viewLine,
						}
						if !yield(segment) {
							return
						}
						sb.Reset()

						viewLine++
						//if drawY-tb.scrollY >= tb.Height {
						//	break
						//}
					}
				}
			} else {
				for j := 0; j < height; j++ {
					segment := presenter.Segment{
						TextLayout:    entry,
						ModelLine:     entry.Element.GetStartPosition().Row,
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
