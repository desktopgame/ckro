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

type FileChooser struct {
	x, y          int
	Width, Height int

	pathLabel  *tui.Tile
	fileList   *tui.Tile
	chooserBox *tui.Box

	currentPath   string
	entries       []FileEntry
	selectedIndex int
	onFileSelect  func(string)
	onCancel      func()
}

type FileEntry struct {
	Name  string
	Path  string
	IsDir bool
	Size  int64
}

func NewFileChooser(initialPath string, onFileSelect func(string), onCancel func()) *FileChooser {
	fc := &FileChooser{
		currentPath:  initialPath,
		onFileSelect: onFileSelect,
		onCancel:     onCancel,
	}
	fc.Init()
	return fc
}

func (fc *FileChooser) Init() {
	// パス表示ラベル
	fc.pathLabel = tui.NewLabelTile(fc.currentPath)
	fc.pathLabel.FlexibleWidth = true
	fc.pathLabel.MinimumHeight = 1

	// ファイルリスト
	fc.fileList = tui.NewListTile([]string{})
	fc.fileList.FlexibleWidth = true
	fc.fileList.FlexibleHeight = true

	// 垂直レイアウトで組み合わせ
	fc.chooserBox = tui.NewVBox(
		fc.pathLabel,
		tui.NewHorizontalSeparator(),
		fc.fileList,
	)

	// 初期ディレクトリを読み込み
	fc.loadDirectory()
}

func (fc *FileChooser) loadDirectory() {
	fc.entries = nil
	fc.selectedIndex = 0

	// 現在のディレクトリを読み込み
	entries, err := os.ReadDir(fc.currentPath)
	if err != nil {
		// エラーの場合は空のリストを表示
		fc.updateFileList()
		return
	}

	// 親ディレクトリへのエントリを追加（ルートディレクトリでない場合）
	if fc.currentPath != "/" && fc.currentPath != "" {
		fc.entries = append(fc.entries, FileEntry{
			Name:  "..",
			Path:  filepath.Dir(fc.currentPath),
			IsDir: true,
			Size:  0,
		})
	}

	// ディレクトリとファイルを分けて追加
	var dirs []FileEntry
	var files []FileEntry

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		fileEntry := FileEntry{
			Name:  entry.Name(),
			Path:  filepath.Join(fc.currentPath, entry.Name()),
			IsDir: entry.IsDir(),
			Size:  info.Size(),
		}

		if entry.IsDir() {
			dirs = append(dirs, fileEntry)
		} else {
			files = append(files, fileEntry)
		}
	}

	// ディレクトリとファイルをそれぞれソート
	sort.Slice(dirs, func(i, j int) bool {
		return strings.ToLower(dirs[i].Name) < strings.ToLower(dirs[j].Name)
	})
	sort.Slice(files, func(i, j int) bool {
		return strings.ToLower(files[i].Name) < strings.ToLower(files[j].Name)
	})

	// ディレクトリを先に、ファイルを後に追加
	fc.entries = append(fc.entries, dirs...)
	fc.entries = append(fc.entries, files...)

	fc.updateFileList()
	fc.updatePathLabel()
}

func (fc *FileChooser) updateFileList() {
	var items []string
	for _, entry := range fc.entries {
		if entry.IsDir {
			items = append(items, "[DIR] "+entry.Name)
		} else {
			items = append(items, entry.Name)
		}
	}

	if listPresenter, ok := fc.fileList.TextPresenter.(*presenter.ListTextPresenter); ok {
		listPresenter.Items = items
		listPresenter.SelectedIndex = fc.selectedIndex
	}
}

func (fc *FileChooser) updatePathLabel() {
	if labelPresenter, ok := fc.pathLabel.TextPresenter.(*presenter.LabelTextPresenter); ok {
		labelPresenter.Text = "Path: " + fc.currentPath
	}
}

func (fc *FileChooser) Move(x, y int) {
	fc.x = x
	fc.y = y
	fc.chooserBox.Move(x, y)
}

func (fc *FileChooser) Layout(width, height int) {
	fc.Width = width
	fc.Height = height
	fc.chooserBox.Layout(width, height)
}

func (fc *FileChooser) MinimumSize(width, height int) (int, int) {
	return fc.chooserBox.MinimumSize(width, height)
}

func (fc *FileChooser) Update() {
	fc.chooserBox.Update()
}

func (fc *FileChooser) Draw(screen tcell.Screen) {
	fc.chooserBox.Draw(screen)
}

func (fc *FileChooser) Handle(ev base.Event) {
	if keyEvent, ok := ev.GetSource().(*tcell.EventKey); ok {
		switch keyEvent.Key() {
		case tcell.KeyUp:
			if fc.selectedIndex > 0 {
				fc.selectedIndex--
				fc.updateFileList()
			}
			return
		case tcell.KeyDown:
			if fc.selectedIndex < len(fc.entries)-1 {
				fc.selectedIndex++
				fc.updateFileList()
			}
			return
		case tcell.KeyEnter:
			if len(fc.entries) > 0 && fc.selectedIndex < len(fc.entries) {
				selectedEntry := fc.entries[fc.selectedIndex]
				if selectedEntry.IsDir {
					// ディレクトリの場合は移動
					fc.currentPath = selectedEntry.Path
					fc.loadDirectory()
				} else {
					// ファイルの場合は選択
					if fc.onFileSelect != nil {
						fc.onFileSelect(selectedEntry.Path)
					}
				}
			}
			return
		case tcell.KeyEscape:
			if fc.onCancel != nil {
				fc.onCancel()
			}
			return
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			// 親ディレクトリに移動
			if fc.currentPath != "/" && fc.currentPath != "" {
				fc.currentPath = filepath.Dir(fc.currentPath)
				fc.loadDirectory()
			}
			return
		}
	}

	// デフォルトのイベント処理
	fc.fileList.Handle(ev)
}

func (fc *FileChooser) Traverse(fm *tui.FocusManager) {
	if fc.IsFocusable() {
		fm.Register(fc)
	}
}

func (fc *FileChooser) Focus(on bool) {
	// フォーカス状態の管理
}

func (fc *FileChooser) SubFocusFirst() {
	// サブフォーカスの最初の要素
}

func (fc *FileChooser) SubFocusPrev() bool {
	// サブフォーカスの前の要素
	return false
}

func (fc *FileChooser) SubFocusNext() bool {
	// サブフォーカスの次の要素
	return false
}

func (fc *FileChooser) SubFocusLast() {
	// サブフォーカスの最後の要素
}

func (fc *FileChooser) IsFocusable() bool {
	return true
}

func (fc *FileChooser) IsFlexibleWidth() bool {
	return fc.chooserBox.IsFlexibleWidth()
}

func (fc *FileChooser) IsFlexibleHeight() bool {
	return fc.chooserBox.IsFlexibleHeight()
}

// GetCurrentPath returns the current directory path
func (fc *FileChooser) GetCurrentPath() string {
	return fc.currentPath
}

// GetSelectedEntry returns the currently selected file entry
func (fc *FileChooser) GetSelectedEntry() *FileEntry {
	if fc.selectedIndex >= 0 && fc.selectedIndex < len(fc.entries) {
		return &fc.entries[fc.selectedIndex]
	}
	return nil
}

// SetPath changes the current directory path
func (fc *FileChooser) SetPath(path string) {
	if _, err := os.Stat(path); err == nil {
		fc.currentPath = path
		fc.loadDirectory()
	}
}

// Refresh reloads the current directory
func (fc *FileChooser) Refresh() {
	fc.loadDirectory()
}
