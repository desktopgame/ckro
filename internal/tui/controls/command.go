package controls

import "github.com/desktopgame/ckro/internal/tui/base"

type Command interface {
	Execute(cp *CommandPalette, stackable base.Stackable)
	GetLabel() string
}

type DelegateCommand struct {
	Label string
	Func  func(cp *CommandPalette, stackable base.Stackable)
}

func (dc DelegateCommand) Execute(cp *CommandPalette, stackable base.Stackable) {
	dc.Func(cp, stackable)
}

func (dc DelegateCommand) GetLabel() string {
	return dc.Label
}
