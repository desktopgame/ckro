package tui

import (
	"strings"

	"github.com/desktopgame/ckro/internal/tui/presenter"
	"github.com/gdamore/tcell/v2"
)

type CommandPalette struct {
	x, y          int
	Width, Height int

	searchInput *Tile
	commandList *Tile
	paletteBox  *Box

	allCommands      []string
	filteredCommands []string
	inputFocused     bool
	onCommandExecute func(string)
}

func (cp *CommandPalette) Init(commands []string, onExecute func(string)) {
	cp.allCommands = commands
	cp.filteredCommands = make([]string, len(commands))
	copy(cp.filteredCommands, commands)
	cp.inputFocused = true
	cp.onCommandExecute = onExecute

	// 検索入力フィールド
	cp.searchInput = NewEditTile()
	cp.searchInput.FlexibleWidth = true
	cp.searchInput.MinimumHeight = 1
	cp.searchInput.TextBox.ShowCursor = true

	// コマンドリスト
	cp.commandList = NewListTile(cp.filteredCommands)
	cp.commandList.FlexibleWidth = true
	cp.commandList.FlexibleHeight = true

	// 垂直レイアウトで組み合わせ
	cp.paletteBox = NewVBox(
		cp.searchInput,
		NewHorizontalSeparator(),
		cp.commandList,
	)
}

func (cp *CommandPalette) Move(x, y int) {
	cp.x = x
	cp.y = y
	cp.paletteBox.Move(x, y)
}

func (cp *CommandPalette) Layout(width, height int) {
	cp.Width = width
	cp.Height = height
	cp.paletteBox.Layout(width, height)
}

func (cp *CommandPalette) MinimumSize(width, height int) (int, int) {
	return cp.paletteBox.MinimumSize(width, height)
}

func (cp *CommandPalette) Update() {
	cp.paletteBox.Update()
}

func (cp *CommandPalette) Draw(screen tcell.Screen) {
	cp.paletteBox.Draw(screen)
}

func (cp *CommandPalette) Handle(ev tcell.Event) {
	if keyEvent, ok := ev.(*tcell.EventKey); ok {
		switch keyEvent.Key() {
		case tcell.KeyUp:
			if !cp.inputFocused {
				// リストの上移動
				if listPresenter, ok := cp.commandList.TextPresenter.(*presenter.ListTextPresenter); ok {
					if listPresenter.SelectedIndex > 0 {
						listPresenter.SelectedIndex--
					}
				}
				return // イベントを消費
			}
		case tcell.KeyDown:
			if !cp.inputFocused {
				// リストの下移動
				if listPresenter, ok := cp.commandList.TextPresenter.(*presenter.ListTextPresenter); ok {
					if listPresenter.SelectedIndex < len(listPresenter.Items)-1 {
						listPresenter.SelectedIndex++
					}
				}
				return // イベントを消費
			}
		case tcell.KeyTab:
			// フォーカス切り替え
			cp.inputFocused = !cp.inputFocused
			if cp.inputFocused {
				cp.searchInput.TextBox.ShowCursor = true
			} else {
				cp.searchInput.TextBox.ShowCursor = false
				// リストの最初の項目を選択
				if listPresenter, ok := cp.commandList.TextPresenter.(*presenter.ListTextPresenter); ok {
					listPresenter.SelectedIndex = 0
				}
			}
			return // イベントを消費
		case tcell.KeyEnter:
			if !cp.inputFocused {
				// コマンド実行
				if listPresenter, ok := cp.commandList.TextPresenter.(*presenter.ListTextPresenter); ok {
					selectedCommand := listPresenter.GetSelectedItem()
					if cp.onCommandExecute != nil && selectedCommand != "" {
						cp.onCommandExecute(selectedCommand)
					}
				}
				return // イベントを消費
			}
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			if cp.inputFocused {
				// 検索クエリから文字を削除
				cp.removeLastChar()
				cp.filterCommands()
				return // イベントを消費
			}
		default:
			// 文字入力
			if cp.inputFocused && keyEvent.Rune() != 0 {
				cp.addChar(keyEvent.Rune())
				cp.filterCommands()
				return // イベントを消費
			}
		}
	}

	// フォーカスされているコントロールにイベントを転送
	if cp.inputFocused {
		cp.searchInput.Handle(ev)
	} else {
		cp.commandList.Handle(ev)
	}
}

func (cp *CommandPalette) addChar(r rune) {
	// 検索入力フィールドに文字を追加
	doc := cp.searchInput.TextBox.GetDocument()
	doc.InsertString(string(r))
}

func (cp *CommandPalette) removeLastChar() {
	// 検索入力フィールドから最後の文字を削除
	doc := cp.searchInput.TextBox.GetDocument()
	if doc.GetCursorColumn() > 0 {
		doc.RemoveChar()
	}
}

func (cp *CommandPalette) getSearchQuery() string {
	// 検索入力フィールドの内容を取得
	doc := cp.searchInput.TextBox.GetDocument()
	buffer := doc.GetBuffer()
	if buffer.GetLineCount() > 0 {
		return buffer.GetLineAt(0).GetContent()
	}
	return ""
}

func (cp *CommandPalette) filterCommands() {
	query := strings.ToLower(cp.getSearchQuery())
	cp.filteredCommands = nil

	for _, command := range cp.allCommands {
		if strings.Contains(strings.ToLower(command), query) {
			cp.filteredCommands = append(cp.filteredCommands, command)
		}
	}

	// リストを更新
	if listPresenter, ok := cp.commandList.TextPresenter.(*presenter.ListTextPresenter); ok {
		listPresenter.Items = cp.filteredCommands
		listPresenter.SelectedIndex = 0
	}
}

func (cp *CommandPalette) Traverse(fm *FocusManager) {
	// CommandPalette自体をFocusableとして登録
	if cp.IsFocusable() {
		fm.Register(cp)
	}
	// 子コントロールは登録しない（CommandPaletteが全てのイベントを処理）
}

// Focusableインターフェースの実装
func (cp *CommandPalette) GetTextBox() *TextBox {
	// 現在フォーカスされているコントロールのTextBoxを返す
	if cp.inputFocused {
		return cp.searchInput.TextBox
	} else {
		return cp.commandList.TextBox
	}
}

func (cp *CommandPalette) GetTextPresenter() TextPresenter {
	// 現在フォーカスされているコントロールのTextPresenterを返す
	if cp.inputFocused {
		return cp.searchInput.TextPresenter
	} else {
		return cp.commandList.TextPresenter
	}
}

func (cp *CommandPalette) IsFocusable() bool {
	return true
}

func (cp *CommandPalette) IsFlexibleWidth() bool {
	return cp.paletteBox.IsFlexibleWidth()
}

func (cp *CommandPalette) IsFlexibleHeight() bool {
	return cp.paletteBox.IsFlexibleHeight()
}
