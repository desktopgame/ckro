package controls

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/desktopgame/ckro/internal/tui"
	"github.com/desktopgame/ckro/internal/tui/base"
	"github.com/desktopgame/ckro/internal/tui/presenter"
	"github.com/gdamore/tcell/v2"
)

type FolderChooser struct {
	x, y          int
	Width, Height int

	pathLabel  *tui.Tile
	folderList *tui.Tile
	chooserBox *tui.Box
	helpLabel  *tui.Tile

	currentPath    string
	folders        []FolderEntry
	selectedIndex  int
	onFolderSelect func(base.Runtime, string)
	onCancel       func(base.Runtime)
}

type FolderEntry struct {
	Name string
	Path string
}

func NewFolderChooser(initialPath string, onFolderSelect func(base.Runtime, string), onCancel func(base.Runtime)) *FolderChooser {
	fc := &FolderChooser{
		currentPath:    initialPath,
		onFolderSelect: onFolderSelect,
		onCancel:       onCancel,
	}
	fc.Init()
	return fc
}

func (fc *FolderChooser) Init() {
	// パス表示ラベル
	fc.pathLabel = tui.NewLabelTile(fc.currentPath)
	fc.pathLabel.FlexibleWidth = true
	fc.pathLabel.MinimumHeight = 1

	// フォルダリスト
	fc.folderList = tui.NewListTile([]string{})
	fc.folderList.FlexibleWidth = true
	fc.folderList.FlexibleHeight = true

	// ヘルプラベル
	fc.helpLabel = tui.NewLabelTile("Enter: Select Folder | Right/Space: Expand | Left/Backspace: Up | Esc: Cancel")
	fc.helpLabel.FlexibleWidth = true
	fc.helpLabel.MinimumHeight = 1

	// 垂直レイアウトで組み合わせ
	fc.chooserBox = tui.NewVBox(
		fc.pathLabel,
		tui.NewHorizontalSeparator(),
		fc.folderList,
		tui.NewHorizontalSeparator(),
		fc.helpLabel,
	)

	// 初期ディレクトリを読み込み
	fc.loadDirectory()
}

func (fc *FolderChooser) loadDirectory() {
	fc.folders = nil
	fc.selectedIndex = 0

	// 現在のディレクトリを読み込み
	entries, err := os.ReadDir(fc.currentPath)
	if err != nil {
		// エラーの場合は空のリストを表示
		fc.updateFolderList()
		return
	}

	// 親ディレクトリへのエントリを追加（ルートディレクトリでない場合）
	if fc.currentPath != "/" && fc.currentPath != "" {
		fc.folders = append(fc.folders, FolderEntry{
			Name: "..",
			Path: filepath.Dir(fc.currentPath),
		})
	}

	// ディレクトリのみを追加
	var dirs []FolderEntry

	for _, entry := range entries {
		if entry.IsDir() {
			// 隠しフォルダをスキップ
			if strings.HasPrefix(entry.Name(), ".") && entry.Name() != ".." {
				continue
			}

			folderEntry := FolderEntry{
				Name: entry.Name(),
				Path: filepath.Join(fc.currentPath, entry.Name()),
			}
			dirs = append(dirs, folderEntry)
		}
	}

	// ディレクトリをソート
	sort.Slice(dirs, func(i, j int) bool {
		return strings.ToLower(dirs[i].Name) < strings.ToLower(dirs[j].Name)
	})

	// ディレクトリを追加
	fc.folders = append(fc.folders, dirs...)

	fc.updateFolderList()
	fc.updatePathLabel()
}

func (fc *FolderChooser) updateFolderList() {
	var items []string
	for _, folder := range fc.folders {
		if folder.Name == ".." {
			items = append(items, "📁 "+folder.Name+" (Parent Directory)")
		} else {
			items = append(items, "📁 "+folder.Name)
		}
	}

	if listPresenter, ok := fc.folderList.TextPresenter.(*presenter.ListTextPresenter); ok {
		listPresenter.Items = items
		listPresenter.SelectedIndex = fc.selectedIndex
	}
}

func (fc *FolderChooser) updatePathLabel() {
	if labelPresenter, ok := fc.pathLabel.TextPresenter.(*presenter.LabelTextPresenter); ok {
		labelPresenter.Text = "Current Path: " + fc.currentPath
	}
}

func (fc *FolderChooser) Move(x, y int) {
	fc.x = x
	fc.y = y
	fc.chooserBox.Move(x, y)
}

func (fc *FolderChooser) Layout(width, height int) {
	fc.Width = width
	fc.Height = height
	fc.chooserBox.Layout(width, height)
}

func (fc *FolderChooser) MinimumSize(width, height int) (int, int) {
	return fc.chooserBox.MinimumSize(width, height)
}

func (fc *FolderChooser) Update() {
	fc.chooserBox.Update()
}

func (fc *FolderChooser) Draw(g *base.Graphics) {
	fc.chooserBox.Draw(g)
}

func (fc *FolderChooser) Handle(ev base.Event) {
	if keyEvent, ok := ev.GetSource().(*tcell.EventKey); ok {
		switch keyEvent.Key() {
		case tcell.KeyUp:
			if fc.selectedIndex > 0 {
				fc.selectedIndex--
				fc.updateFolderList()
			}
			return
		case tcell.KeyDown:
			if fc.selectedIndex < len(fc.folders)-1 {
				fc.selectedIndex++
				fc.updateFolderList()
			}
			return
		case tcell.KeyEnter:
			// Enterキーでフォルダを選択（決定）
			if len(fc.folders) > 0 && fc.selectedIndex < len(fc.folders) {
				selectedFolder := fc.folders[fc.selectedIndex]
				if fc.onFolderSelect != nil {
					fc.onFolderSelect(ev.GetRuntime(), selectedFolder.Path)
				}
			}
			return
		case tcell.KeyRight:
			// 右矢印キーでフォルダを展開（移動）
			if len(fc.folders) > 0 && fc.selectedIndex < len(fc.folders) {
				selectedFolder := fc.folders[fc.selectedIndex]
				fc.currentPath = selectedFolder.Path
				fc.loadDirectory()
			}
			return
		case tcell.KeyLeft, tcell.KeyBackspace, tcell.KeyBackspace2:
			// 左矢印キーまたはBackspaceで親ディレクトリに移動
			if fc.currentPath != "/" && fc.currentPath != "" {
				fc.currentPath = filepath.Dir(fc.currentPath)
				fc.loadDirectory()
			}
			return
		case tcell.KeyEscape:
			if fc.onCancel != nil {
				fc.onCancel(ev.GetRuntime())
			}
			return
		}

		// 文字キーでの処理
		switch keyEvent.Rune() {
		case ' ':
			// スペースキーでフォルダを展開（移動）
			if len(fc.folders) > 0 && fc.selectedIndex < len(fc.folders) {
				selectedFolder := fc.folders[fc.selectedIndex]
				fc.currentPath = selectedFolder.Path
				fc.loadDirectory()
			}
			return
		}
	}

	// デフォルトのイベント処理
	fc.folderList.Handle(ev)
}

func (fc *FolderChooser) Traverse(fm *tui.FocusManager) {
	if fc.IsFocusable() {
		fm.Register(fc)
	}
}

func (fc *FolderChooser) Focus(on bool) {
	// フォーカス状態の管理
}

func (fc *FolderChooser) SubFocusFirst() {
	// サブフォーカスの最初の要素
}

func (fc *FolderChooser) SubFocusPrev() bool {
	// サブフォーカスの前の要素
	return false
}

func (fc *FolderChooser) SubFocusNext() bool {
	// サブフォーカスの次の要素
	return false
}

func (fc *FolderChooser) SubFocusLast() {
	// サブフォーカスの最後の要素
}

func (fc *FolderChooser) IsFocusable() bool {
	return true
}

func (fc *FolderChooser) IsFlexibleWidth() bool {
	return fc.chooserBox.IsFlexibleWidth()
}

func (fc *FolderChooser) IsFlexibleHeight() bool {
	return fc.chooserBox.IsFlexibleHeight()
}

// GetCurrentPath returns the current directory path
func (fc *FolderChooser) GetCurrentPath() string {
	return fc.currentPath
}

// GetSelectedFolder returns the currently selected folder entry
func (fc *FolderChooser) GetSelectedFolder() *FolderEntry {
	if fc.selectedIndex >= 0 && fc.selectedIndex < len(fc.folders) {
		return &fc.folders[fc.selectedIndex]
	}
	return nil
}

// SetPath changes the current directory path
func (fc *FolderChooser) SetPath(path string) {
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		fc.currentPath = path
		fc.loadDirectory()
	}
}

// Refresh reloads the current directory
func (fc *FolderChooser) Refresh() {
	fc.loadDirectory()
}

// SelectCurrentFolder selects the current directory (useful for "select current folder" functionality)
func (fc *FolderChooser) SelectCurrentFolder() {
	if fc.onFolderSelect != nil {
		// Create a dummy runtime for this call
		// In real usage, you should pass the actual runtime
		fc.onFolderSelect(nil, fc.currentPath)
	}
}
