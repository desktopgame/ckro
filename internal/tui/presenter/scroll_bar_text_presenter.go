package presenter

import (
	"github.com/gdamore/tcell/v2"
)

type ScrollBarTextPresenter struct {
	TargetView View // スクロールバーを表示する対象のView
}

func (sb *ScrollBarTextPresenter) Present(view View) {
	if sb.TargetView == nil {
		return
	}

	view.TextClear()
	doc := view.GetDocument()
	targetDoc := sb.TargetView.GetDocument()

	// 対象ビューのスクロール情報を取得
	scrollY := sb.TargetView.GetScrollY()
	viewHeight := sb.TargetView.GetHeight()

	// 対象ドキュメントの総行数を取得
	totalLines := targetDoc.GetBuffer().GetLineCount()
	if totalLines == 0 {
		totalLines = 1
	}

	// スクロールバーの高さ（ビューの高さと同じ）
	scrollBarHeight := viewHeight

	// スクロール可能な範囲を計算
	maxScrollY := max(0, totalLines-viewHeight)

	// スクロールバーを描画
	for i := 0; i < scrollBarHeight; i++ {
		var char rune

		if totalLines <= viewHeight {
			// スクロールが不要な場合（全体が表示されている）
			char = '│'
		} else {
			// スクロール位置を正規化（範囲外の値を修正）
			normalizedScrollY := max(0, min(scrollY, maxScrollY))

			// つまみのサイズを計算（表示範囲の割合に基づく）
			thumbSize := max(1, (viewHeight*scrollBarHeight)/totalLines)

			// つまみの開始位置を計算
			thumbStart := 0
			if maxScrollY > 0 {
				thumbStart = (normalizedScrollY * (scrollBarHeight - thumbSize)) / maxScrollY
			}
			thumbEnd := thumbStart + thumbSize

			if i >= thumbStart && i < thumbEnd {
				// スクロールバーのつまみ部分
				char = '█'
			} else {
				// スクロールバーの背景部分
				char = '░'
			}
		}

		doc.InsertString(string(char))

		if i < scrollBarHeight-1 {
			doc.InsertLine()
		}
	}

	// カーソル位置を対象ビューと同期（スクロールバーにはカーソルを表示しない）
	sb.syncCursorPosition(view, targetDoc)
}

// syncCursorPosition synchronizes cursor position with target view
func (sb *ScrollBarTextPresenter) syncCursorPosition(view View, targetDoc interface{}) {
	// 対象ドキュメントのカーソル行を取得
	if doc, ok := targetDoc.(interface{ GetCursorRow() int }); ok {
		cursorRow := doc.GetCursorRow()

		// スクロールバービューのカーソルを同じ行に移動
		scrollBarDoc := view.GetDocument()
		scrollBarDoc.MoveReset()

		for i := 0; i < cursorRow && i < scrollBarDoc.GetBuffer().GetLineCount()-1; i++ {
			scrollBarDoc.MoveDown()
		}

		// 行の先頭に移動
		for scrollBarDoc.GetCursorColumn() > 0 {
			scrollBarDoc.MoveLeft()
		}
	}
}

func (sb *ScrollBarTextPresenter) Handle(view View, ev tcell.Event) {
	// スクロールバーは編集不可なので、イベントは処理しない
	// ただし、カーソル更新は行う
	view.CursorUpdate()
}

func (sb *ScrollBarTextPresenter) ShowCursor() bool {
	return false // スクロールバーにはカーソルを表示しない
}

func (sb *ScrollBarTextPresenter) IsFocusable() bool {
	return false // スクロールバーはフォーカス不可
}
