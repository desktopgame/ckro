package controls

type Command interface {
	Execute()
	GetLabel() string
}

type DelegateCommand struct {
	Label string
	Func  func()
}

func (dc DelegateCommand) Execute() {
	dc.Func()
}

func (dc DelegateCommand) GetLabel() string {
	return dc.Label
}
