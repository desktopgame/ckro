package controls

import "github.com/desktopgame/ckro/internal/tui/base"

type Command interface {
	Execute(runtime base.Runtime, cp *CommandPalette)
	GetLabel() string
}

type DelegateCommand struct {
	Label string
	Func  func(runtime base.Runtime, cp *CommandPalette)
}

func (dc DelegateCommand) Execute(runtime base.Runtime, cp *CommandPalette) {
	dc.Func(runtime, cp)
}

func (dc DelegateCommand) GetLabel() string {
	return dc.Label
}
