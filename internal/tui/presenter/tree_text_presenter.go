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
	OnFileOpen    func(filePath string)
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
		case tcell.KeyEnter:
			// open file or expand tree
			if t.selectedIndex < len(t.flatNodes) {
				node := t.flatNodes[t.selectedIndex]
				if node.IsDir {
					// expand or collapse
					if node.IsExpanded {
						node.IsExpanded = false
						node.Children = nil
					} else {
						node.IsExpanded = true
						t.loadChildren(node)
					}
				} else {
					if t.OnFileOpen != nil {
						t.OnFileOpen(node.Path)
					}
				}
			}
		case tcell.KeyRight:
			// expand
			if t.selectedIndex < len(t.flatNodes) {
				node := t.flatNodes[t.selectedIndex]
				if node.IsDir && !node.IsExpanded {
					node.IsExpanded = true
					t.loadChildren(node)
				}
			}
		case tcell.KeyLeft:
			// collapse
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

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir() != entries[j].IsDir() {
			return entries[i].IsDir()
		}
		return entries[i].Name() < entries[j].Name()
	})

	node.Children = make([]*TreeNode, 0, len(entries))
	for _, entry := range entries {
		// skip the hidden files.
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
		indent := strings.Repeat("  ", node.Level)

		// icon for expand or collapse
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

		// cursor
		var prefix string
		if i == t.selectedIndex {
			prefix = "> "
		} else {
			prefix = "  "
		}

		// icon for file or directory
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

	t.moveToSelectedItem(view)
}

// moveToSelectedItem is moves the document cursor to the selected item
func (t *TreeTextPresenter) moveToSelectedItem(view View) {
	// TODO: impl
	//doc := view.GetDocument()
	//if t.selectedIndex >= 0 && t.selectedIndex < len(t.flatNodes) {
	//	doc.MoveReset()
	//
	//	for i := 0; i < t.selectedIndex; i++ {
	//		doc.MoveDown()
	//	}
	//
	//	for doc.GetCursorColumn() > 0 {
	//		doc.MoveLeft()
	//	}
	//}
}

func (t *TreeTextPresenter) GetSelectedPath() string {
	if t.selectedIndex < len(t.flatNodes) {
		return t.flatNodes[t.selectedIndex].Path
	}
	return ""
}

func (t *TreeTextPresenter) Reset(rootDirectory string) {
	t.rootNode = nil
	t.RootDirectory = rootDirectory
}

// Reload reloads the children of all expanded nodes in the tree
func (t *TreeTextPresenter) Reload() {
	if t.rootNode != nil {
		t.reloadNode(t.rootNode)
	}
}

// ReloadSelected reloads the children of the currently selected node if it's a directory
func (t *TreeTextPresenter) ReloadSelected() {
	if t.selectedIndex < len(t.flatNodes) {
		node := t.flatNodes[t.selectedIndex]
		if node.IsDir {
			t.reloadNode(node)
		}
	}
}

// ReloadPath reloads the children of the node at the specified path if it exists and is expanded
func (t *TreeTextPresenter) ReloadPath(path string) {
	if t.rootNode != nil {
		node := t.findNodeByPath(t.rootNode, path)
		if node != nil && node.IsDir {
			t.reloadNode(node)
		}
	}
}

// reloadNode reloads the children of a specific node if it's expanded
func (t *TreeTextPresenter) reloadNode(node *TreeNode) {
	if !node.IsDir || !node.IsExpanded {
		return
	}

	// Store the expanded state of current children
	expandedPaths := make(map[string]bool)
	if node.Children != nil {
		for _, child := range node.Children {
			if child.IsExpanded {
				expandedPaths[child.Path] = true
			}
		}
	}

	// Reload children from filesystem
	t.loadChildren(node)

	// Restore expanded state and recursively reload expanded children
	if node.Children != nil {
		for _, child := range node.Children {
			if expandedPaths[child.Path] {
				child.IsExpanded = true
				t.reloadNode(child) // Recursively reload expanded children
			}
		}
	}
}

// findNodeByPath finds a node by its path in the tree
func (t *TreeTextPresenter) findNodeByPath(node *TreeNode, path string) *TreeNode {
	if node.Path == path {
		return node
	}

	if node.Children != nil {
		for _, child := range node.Children {
			if found := t.findNodeByPath(child, path); found != nil {
				return found
			}
		}
	}

	return nil
}
