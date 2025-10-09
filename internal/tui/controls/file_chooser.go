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

// FileChooser is file choose dialog.
type FileChooser struct {
	x, y          int
	Width, Height int

	pathLabel  *tui.Tile
	fileList   *tui.Tile
	chooserBox *tui.Box

	currentPath   string
	entries       []FileEntry
	selectedIndex int
	onFileSelect  func(base.Runtime, string)
	onCancel      func(base.Runtime)
}

type FileEntry struct {
	Name  string
	Path  string
	IsDir bool
	Size  int64
}

// NewFileChooser returns FileChooser.
func NewFileChooser(initialPath string, onFileSelect func(base.Runtime, string), onCancel func(base.Runtime)) *FileChooser {
	fc := &FileChooser{
		currentPath:  initialPath,
		onFileSelect: onFileSelect,
		onCancel:     onCancel,
	}
	fc.Init()
	return fc
}

// Init is initialize FileChooser.
func (fc *FileChooser) Init() {
	fc.pathLabel = tui.NewLabelTile(fc.currentPath)
	fc.pathLabel.FlexibleWidth = true
	fc.pathLabel.MinimumHeight = 1

	fc.fileList = tui.NewListTile([]string{})
	fc.fileList.FlexibleWidth = true
	fc.fileList.FlexibleHeight = true

	fc.chooserBox = tui.NewVBox(
		fc.pathLabel,
		tui.NewHorizontalSeparator(),
		fc.fileList,
	)

	fc.loadDirectory()
}

func (fc *FileChooser) loadDirectory() {
	fc.entries = nil
	fc.selectedIndex = 0

	// load directory
	entries, err := os.ReadDir(fc.currentPath)
	if err != nil {
		fc.updateFileList()
		return
	}

	// parent directory
	if fc.currentPath != "/" && fc.currentPath != "" {
		fc.entries = append(fc.entries, FileEntry{
			Name:  "..",
			Path:  filepath.Dir(fc.currentPath),
			IsDir: true,
			Size:  0,
		})
	}

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

	sort.Slice(dirs, func(i, j int) bool {
		return strings.ToLower(dirs[i].Name) < strings.ToLower(dirs[j].Name)
	})
	sort.Slice(files, func(i, j int) bool {
		return strings.ToLower(files[i].Name) < strings.ToLower(files[j].Name)
	})

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

func (fc *FileChooser) Draw(g *base.Graphics) {
	fc.chooserBox.Draw(g)
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
					// move to directory
					fc.currentPath = selectedEntry.Path
					fc.loadDirectory()
				} else {
					// select file
					if fc.onFileSelect != nil {
						fc.onFileSelect(ev.GetRuntime(), selectedEntry.Path)
					}
				}
			}
			return
		case tcell.KeyEscape:
			if fc.onCancel != nil {
				fc.onCancel(ev.GetRuntime())
			}
			return
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			// move to parent directory
			if fc.currentPath != "/" && fc.currentPath != "" {
				fc.currentPath = filepath.Dir(fc.currentPath)
				fc.loadDirectory()
			}
			return
		}
	}

	fc.fileList.Handle(ev)
}

func (fc *FileChooser) Traverse(fm *tui.FocusManager) {
	if fc.IsFocusable() {
		fm.Register(fc)
	}
}

func (fc *FileChooser) Focus(on bool) {
}

func (fc *FileChooser) SubFocusFirst() {
}

func (fc *FileChooser) SubFocusPrev() bool {
	return false
}

func (fc *FileChooser) SubFocusNext() bool {
	return false
}

func (fc *FileChooser) SubFocusLast() {
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
