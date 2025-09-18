package controls

import (
	"strings"

	"github.com/desktopgame/ckro/internal/tui"
	"github.com/desktopgame/ckro/internal/tui/presenter"
	"github.com/gdamore/tcell/v2"
)

type CommandPalette struct {
	x, y          int
	Width, Height int

	searchInput *tui.Tile
	commandList *tui.Tile
	paletteBox  *tui.Box

	allCommands      []string
	filteredCommands []string
	inputFocused     bool
	onCommandExecute func(string)
}

func NewCommandPalette(commands []string, onExecute func(string)) *CommandPalette {
	palette := &CommandPalette{}
	palette.Init(commands, onExecute)
	return palette
}

func (cp *CommandPalette) Init(commands []string, onExecute func(string)) {
	cp.allCommands = commands
	cp.filteredCommands = make([]string, len(commands))
	copy(cp.filteredCommands, commands)
	cp.inputFocused = true
	cp.onCommandExecute = onExecute

	// 検索入力フィールド
	cp.searchInput = tui.NewEditTile()
	cp.searchInput.FlexibleWidth = true
	cp.searchInput.MinimumHeight = 1

	// コマンドリスト
	cp.commandList = tui.NewListTile(cp.filteredCommands)
	cp.commandList.FlexibleWidth = true
	cp.commandList.FlexibleHeight = true

	// 垂直レイアウトで組み合わせ
	cp.paletteBox = tui.NewVBox(
		cp.searchInput,
		tui.NewHorizontalSeparator(),
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

func (cp *CommandPalette) Traverse(fm *tui.FocusManager) {
	// CommandPalette自体をFocusableとして登録
	if cp.IsFocusable() {
		fm.Register(cp)
	}
	// 子コントロールは登録しない（CommandPaletteが全てのイベントを処理）
}

func (cp *CommandPalette) Focus(on bool) {
	cp.searchInput.TextBox.ShowCursor = true
}

func (cp *CommandPalette) SubFocusFirst() {
	cp.inputFocused = true
	cp.searchInput.TextBox.ShowCursor = true
}

func (cp *CommandPalette) SubFocusPrev() bool {
	if !cp.inputFocused {
		cp.inputFocused = true
		cp.searchInput.TextBox.ShowCursor = true
	}
	return false
}

func (cp *CommandPalette) SubFocusNext() bool {
	if cp.inputFocused {
		cp.inputFocused = false
		cp.searchInput.TextBox.ShowCursor = false
	}
	return false
}

func (cp *CommandPalette) SubFocusLast() {
	cp.inputFocused = false
	cp.searchInput.TextBox.ShowCursor = false
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
