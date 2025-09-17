package presenter

import (
	"github.com/gdamore/tcell/v2"
)

type ListTextPresenter struct {
	Items         []string // リスト項目
	SelectedIndex int      // 選択されている項目のインデックス
	CursorChar    rune     // カーソル文字（デフォルト: '>'）
	Prefix        string   // 各項目の前に付けるプレフィックス（デフォルト: " "）
}

func (lp *ListTextPresenter) Present(view View) {
	view.TextClear()
	doc := view.GetDocument()

	// デフォルト値の設定
	cursorChar := lp.CursorChar
	if cursorChar == 0 {
		cursorChar = '>'
	}

	prefix := lp.Prefix
	if prefix == "" {
		prefix = " "
	}

	// リスト項目を描画
	for i, item := range lp.Items {
		// カーソル文字または空白を挿入
		if i == lp.SelectedIndex {
			doc.InsertString(string(cursorChar))
		} else {
			doc.InsertString(" ")
		}

		// プレフィックスと項目テキストを挿入
		doc.InsertString(prefix + item)

		// 最後の項目でなければ改行
		if i < len(lp.Items)-1 {
			doc.InsertLine()
		}
	}

	// カーソル位置を選択された項目に設定
	lp.setCursorToSelectedItem(view)
}

// setCursorToSelectedItem sets cursor to the selected item
func (lp *ListTextPresenter) setCursorToSelectedItem(view View) {
	doc := view.GetDocument()
	doc.MoveReset()

	// 選択された行まで移動
	for i := 0; i < lp.SelectedIndex && i < len(lp.Items)-1; i++ {
		doc.MoveDown()
	}

	// 行の先頭に移動（カーソル文字の位置）
	for doc.GetCursorColumn() > 0 {
		doc.MoveLeft()
	}
}

func (lp *ListTextPresenter) Handle(view View, ev tcell.Event) {
	if keyEvent, ok := ev.(*tcell.EventKey); ok {
		switch keyEvent.Key() {
		case tcell.KeyUp:
			// 上の項目を選択
			if lp.SelectedIndex > 0 {
				lp.SelectedIndex--
				lp.Present(view)
			}
		case tcell.KeyDown:
			// 下の項目を選択
			if lp.SelectedIndex < len(lp.Items)-1 {
				lp.SelectedIndex++
				lp.Present(view)
			}
		case tcell.KeyHome:
			// 最初の項目を選択
			if len(lp.Items) > 0 {
				lp.SelectedIndex = 0
				lp.Present(view)
			}
		case tcell.KeyEnd:
			// 最後の項目を選択
			if len(lp.Items) > 0 {
				lp.SelectedIndex = len(lp.Items) - 1
				lp.Present(view)
			}
		}
	}

	// カーソル更新
	view.CursorUpdate()
}

func (lp *ListTextPresenter) ShowCursor() bool {
	return false // リストにはカーソルを表示
}

func (lp *ListTextPresenter) IsFocusable() bool {
	return true // リストはフォーカス可能
}

// GetSelectedItem returns the currently selected item
func (lp *ListTextPresenter) GetSelectedItem() string {
	if lp.SelectedIndex >= 0 && lp.SelectedIndex < len(lp.Items) {
		return lp.Items[lp.SelectedIndex]
	}
	return ""
}

// GetSelectedIndex returns the currently selected index
func (lp *ListTextPresenter) GetSelectedIndex() int {
	return lp.SelectedIndex
}

// SetSelectedIndex sets the selected index
func (lp *ListTextPresenter) SetSelectedIndex(index int) {
	if index >= 0 && index < len(lp.Items) {
		lp.SelectedIndex = index
	}
}

// AddItem adds a new item to the list
func (lp *ListTextPresenter) AddItem(item string) {
	lp.Items = append(lp.Items, item)
}

// RemoveItem removes an item at the specified index
func (lp *ListTextPresenter) RemoveItem(index int) {
	if index >= 0 && index < len(lp.Items) {
		lp.Items = append(lp.Items[:index], lp.Items[index+1:]...)

		// 選択インデックスを調整
		if lp.SelectedIndex >= len(lp.Items) && len(lp.Items) > 0 {
			lp.SelectedIndex = len(lp.Items) - 1
		} else if len(lp.Items) == 0 {
			lp.SelectedIndex = 0
		}
	}
}

// Clear removes all items from the list
func (lp *ListTextPresenter) Clear() {
	lp.Items = nil
	lp.SelectedIndex = 0
}
