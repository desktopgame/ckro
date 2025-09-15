package presenter

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type TreeNode struct {
	Name       string
	Path       string
	IsDir      bool
	IsExpanded bool
	Children   []*TreeNode
	Parent     *TreeNode
	Level      int
}

type TreeTextPresenter struct {
	RootDirectory string
	rootNode      *TreeNode
	flatNodes     []*TreeNode
	selectedIndex int
}

func (t *TreeTextPresenter) Present(view View) {
	if t.rootNode == nil {
		t.buildTree()
	}
	t.updateFlatNodes()
	t.renderTree(view)
}

func (t *TreeTextPresenter) Handle(view View, ev tcell.Event) {
	switch e := ev.(type) {
	case *tcell.EventKey:
		switch e.Key() {
		case tcell.KeyUp:
			if t.selectedIndex > 0 {
				t.selectedIndex--
			}
		case tcell.KeyDown:
			if t.selectedIndex < len(t.flatNodes)-1 {
				t.selectedIndex++
			}
		case tcell.KeyEnter, tcell.KeyRight:
			// 展開
			if t.selectedIndex < len(t.flatNodes) {
				node := t.flatNodes[t.selectedIndex]
				if node.IsDir && !node.IsExpanded {
					node.IsExpanded = true
					t.loadChildren(node)
				}
			}
		case tcell.KeyLeft:
			// 折りたたみ
			if t.selectedIndex < len(t.flatNodes) {
				node := t.flatNodes[t.selectedIndex]
				if node.IsDir && node.IsExpanded {
					node.IsExpanded = false
					node.Children = nil
				}
			}
		}
		t.Present(view)
	}

	view.CursorUpdate()
}

func (t *TreeTextPresenter) ShowCursor() bool {
	return false
}

func (t *TreeTextPresenter) IsFocusable() bool {
	return true
}

func (t *TreeTextPresenter) buildTree() {
	if t.RootDirectory == "" {
		t.RootDirectory = "."
	}

	info, err := os.Stat(t.RootDirectory)
	if err != nil {
		return
	}

	t.rootNode = &TreeNode{
		Name:       filepath.Base(t.RootDirectory),
		Path:       t.RootDirectory,
		IsDir:      info.IsDir(),
		IsExpanded: true,
		Level:      0,
	}

	if t.rootNode.IsDir {
		t.loadChildren(t.rootNode)
	}
}

func (t *TreeTextPresenter) loadChildren(node *TreeNode) {
	entries, err := os.ReadDir(node.Path)
	if err != nil {
		return
	}

	// ディレクトリを先に、ファイルを後にソート
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir() != entries[j].IsDir() {
			return entries[i].IsDir()
		}
		return entries[i].Name() < entries[j].Name()
	})

	node.Children = make([]*TreeNode, 0, len(entries))
	for _, entry := range entries {
		// 隠しファイルをスキップ
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		child := &TreeNode{
			Name:   entry.Name(),
			Path:   filepath.Join(node.Path, entry.Name()),
			IsDir:  entry.IsDir(),
			Parent: node,
			Level:  node.Level + 1,
		}
		node.Children = append(node.Children, child)
	}
}

func (t *TreeTextPresenter) updateFlatNodes() {
	t.flatNodes = make([]*TreeNode, 0)
	if t.rootNode != nil {
		t.addNodeToFlat(t.rootNode)
	}
}

func (t *TreeTextPresenter) addNodeToFlat(node *TreeNode) {
	t.flatNodes = append(t.flatNodes, node)

	if node.IsExpanded && node.Children != nil {
		for _, child := range node.Children {
			t.addNodeToFlat(child)
		}
	}
}

func (t *TreeTextPresenter) renderTree(view View) {
	view.TextClear()
	doc := view.GetDocument()

	for i, node := range t.flatNodes {
		// インデント
		indent := strings.Repeat("  ", node.Level)

		// 展開/折りたたみアイコン
		var icon string
		if node.IsDir {
			if node.IsExpanded {
				icon = "▼ "
			} else {
				icon = "▶ "
			}
		} else {
			icon = "  "
		}

		// 選択インジケーター
		var prefix string
		if i == t.selectedIndex {
			prefix = "> "
		} else {
			prefix = "  "
		}

		// ファイル/ディレクトリアイコン
		var typeIcon string
		if node.IsDir {
			typeIcon = "📁 "
		} else {
			typeIcon = "📄 "
		}

		line := prefix + indent + icon + typeIcon + node.Name
		doc.InsertString(line)

		if i < len(t.flatNodes)-1 {
			doc.InsertLine()
		}
	}

	doc.MoveReset()
}

func (t *TreeTextPresenter) GetSelectedPath() string {
	if t.selectedIndex < len(t.flatNodes) {
		return t.flatNodes[t.selectedIndex].Path
	}
	return ""
}
