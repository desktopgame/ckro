package controls

import (
	"strings"

	"github.com/desktopgame/ckro/internal/tui"
	"github.com/desktopgame/ckro/internal/tui/base"
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/desktopgame/ckro/internal/tui/presenter"
	"github.com/gdamore/tcell/v2"
)

type CommandPalette struct {
	x, y          int
	Width, Height int

	searchInput *tui.Tile
	commandList *tui.Tile
	paletteBox  *tui.Box

	allCommands      []Command
	filteredCommands []Command
	inputFocused     bool
}

func NewCommandPalette(commands []Command) *CommandPalette {
	palette := &CommandPalette{}
	palette.Init(commands)
	return palette
}

func (cp *CommandPalette) Init(commands []Command) {
	cp.allCommands = commands
	cp.filteredCommands = make([]Command, len(commands))
	copy(cp.filteredCommands, commands)
	cp.inputFocused = true

	// 検索入力フィールド
	cp.searchInput = tui.NewEditTile()
	cp.searchInput.FlexibleWidth = true
	cp.searchInput.MinimumHeight = 1

	// コマンドリスト
	cp.commandList = tui.NewListTile(cp.GetLabels())
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

func (cp *CommandPalette) Draw(g *base.Graphics) {
	cp.paletteBox.Draw(g)
}

func (cp *CommandPalette) Handle(ev base.Event) {
	if keyEvent, ok := ev.GetSource().(*tcell.EventKey); ok {
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
					selectedIndex := listPresenter.GetSelectedIndex()
					if selectedIndex >= 0 {
						cp.filteredCommands[selectedIndex].Execute(ev.GetRuntime(), cp)
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
		case tcell.KeyEscape:
			ev.GetRuntime().Pop(-1)
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
	cp.searchInput.TextBox.InsertString(string(r))
}

func (cp *CommandPalette) removeLastChar() {
	// 検索入力フィールドから最後の文字を削除
	cp.searchInput.TextBox.RemoveChar()
}

func (cp *CommandPalette) getSearchQuery() string {
	// 検索入力フィールドの内容を取得
	tb := cp.searchInput.TextBox
	sg := tb.Document.Read(model.Range{
		StartPosition: model.Position{
			Row:    0,
			Column: 0,
		},
		EndPosition: model.Position{
			Row:    0,
			Column: tb.Document.GetLineBytes(0),
		},
	})
	return sg.GetLine(0)
}

func (cp *CommandPalette) filterCommands() {
	query := strings.ToLower(cp.getSearchQuery())
	cp.filteredCommands = nil

	for _, command := range cp.allCommands {
		if strings.Contains(strings.ToLower(command.GetLabel()), query) {
			cp.filteredCommands = append(cp.filteredCommands, command)
		}
	}

	// リストを更新
	if listPresenter, ok := cp.commandList.TextPresenter.(*presenter.ListTextPresenter); ok {
		listPresenter.Items = cp.GetLabels()
		listPresenter.SelectedIndex = 0
	}
}

func (cp *CommandPalette) GetLabels() []string {
	var labels []string
	for i := 0; i < len(cp.filteredCommands); i++ {
		labels = append(labels, cp.filteredCommands[i].GetLabel())
	}
	return labels
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
