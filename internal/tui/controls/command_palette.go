package controls

import (
	"strings"

	"github.com/desktopgame/ckro/internal/tui"
	"github.com/desktopgame/ckro/internal/tui/base"
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/desktopgame/ckro/internal/tui/presenter"
	"github.com/gdamore/tcell/v2"
)

// CommandPalette is vscode like command palette control.
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

// NewCommandPalette returns CommandPalette
func NewCommandPalette(commands []Command) *CommandPalette {
	palette := &CommandPalette{}
	palette.Init(commands)
	return palette
}

// Init is initialize CommandPalette.
func (cp *CommandPalette) Init(commands []Command) {
	cp.allCommands = commands
	cp.filteredCommands = make([]Command, len(commands))
	copy(cp.filteredCommands, commands)
	cp.inputFocused = true

	cp.searchInput = tui.NewEditTile()
	cp.searchInput.FlexibleWidth = true
	cp.searchInput.MinimumHeight = 1

	cp.commandList = tui.NewListTile(cp.GetLabels())
	cp.commandList.FlexibleWidth = true
	cp.commandList.FlexibleHeight = true

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
				if listPresenter, ok := cp.commandList.TextPresenter.(*presenter.ListTextPresenter); ok {
					if listPresenter.SelectedIndex > 0 {
						listPresenter.SelectedIndex--
					}
				}
				return
			}
		case tcell.KeyDown:
			if !cp.inputFocused {
				if listPresenter, ok := cp.commandList.TextPresenter.(*presenter.ListTextPresenter); ok {
					if listPresenter.SelectedIndex < len(listPresenter.Items)-1 {
						listPresenter.SelectedIndex++
					}
				}
				return
			}
		case tcell.KeyEnter:
			if !cp.inputFocused {
				// execute command
				if listPresenter, ok := cp.commandList.TextPresenter.(*presenter.ListTextPresenter); ok {
					selectedIndex := listPresenter.GetSelectedIndex()
					if selectedIndex >= 0 {
						cp.filteredCommands[selectedIndex].Execute(ev.GetRuntime(), cp)
					}
				}
				return
			}
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			if cp.inputFocused {
				cp.removeLastChar()
				cp.filterCommands()
				return
			}
		case tcell.KeyEscape:
			ev.GetRuntime().Pop(-1)
		default:
			if cp.inputFocused && keyEvent.Rune() != 0 {
				cp.addChar(keyEvent.Rune())
				cp.filterCommands()
				return
			}
		}
	}

	if cp.inputFocused {
		cp.searchInput.Handle(ev)
	} else {
		cp.commandList.Handle(ev)
	}
}

func (cp *CommandPalette) addChar(r rune) {
	cp.searchInput.TextBox.InsertString(string(r))
}

func (cp *CommandPalette) removeLastChar() {
	cp.searchInput.TextBox.RemoveChar()
}

func (cp *CommandPalette) getSearchQuery() string {
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

	// update list
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
	if cp.IsFocusable() {
		fm.Register(cp)
	}
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
