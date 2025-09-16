package presenter

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

type LineNumberTextPresenter struct {
	TargetView View // 行番号を表示する対象のView
}

func (ln *LineNumberTextPresenter) Present(view View) {
	if ln.TargetView == nil {
		return
	}

	view.TextClear()
	doc := view.GetDocument()
	targetDoc := ln.TargetView.GetDocument()

	// 対象ドキュメントの行数を取得
	lineCount := targetDoc.GetBuffer().GetLineCount()
	if lineCount == 0 {
		lineCount = 1 // 最低1行は表示
	}

	// 行番号の桁数を計算（最大行数に基づく）
	maxDigits := len(fmt.Sprintf("%d", lineCount))

	// BreakIterを使って折り返しを考慮した行番号を生成
	scrollY := ln.TargetView.GetScrollY()
	segments := make([]Segment, 0)

	// 全セグメントを収集
	for segment := range ln.TargetView.BreakIter() {
		segments = append(segments, segment)
	}

	// スクロール範囲内のセグメントのみ処理
	for i, segment := range segments {
		if segment.ViewLine >= scrollY {
			// 論理行番号を表示（1ベース）
			lineNumber := fmt.Sprintf("%*d", maxDigits, segment.ModelLine+1)
			doc.InsertString(lineNumber)

			// 最後のセグメントでない場合は改行
			if i < len(segments)-1 {
				doc.InsertLine()
			}
		}
	}

	// カーソル位置を対象ビューと同期
	ln.syncCursorPosition(view, targetDoc)
}

// syncCursorPosition synchronizes cursor position with target view
func (ln *LineNumberTextPresenter) syncCursorPosition(view View, targetDoc interface{}) {
	// 対象ドキュメントのカーソル行を取得
	if doc, ok := targetDoc.(interface{ GetCursorRow() int }); ok {
		cursorRow := doc.GetCursorRow()

		// 行番号ビューのカーソルを同じ行に移動
		lineNumDoc := view.GetDocument()
		lineNumDoc.MoveReset()

		for i := 0; i < cursorRow && i < lineNumDoc.GetBuffer().GetLineCount()-1; i++ {
			lineNumDoc.MoveDown()
		}

		// 行の先頭に移動
		for lineNumDoc.GetCursorColumn() > 0 {
			lineNumDoc.MoveLeft()
		}
	}
}

func (ln *LineNumberTextPresenter) Handle(view View, ev tcell.Event) {
	// 行番号は編集不可なので、イベントは処理しない
	// ただし、カーソル更新は行う
	view.CursorUpdate()
}

func (ln *LineNumberTextPresenter) ShowCursor() bool {
	return false // 行番号にはカーソルを表示しない
}

func (ln *LineNumberTextPresenter) IsFocusable() bool {
	return false // 行番号はフォーカス不可
}
