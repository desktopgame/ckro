package controls

import (
	"github.com/desktopgame/ckro/internal/tui"
	"github.com/desktopgame/ckro/internal/tui/base"
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/desktopgame/ckro/internal/tui/presenter"
	"github.com/gdamore/tcell/v2"
)

// InputDialog is dialog with input form.
type InputDialog struct {
	x, y          int
	Width, Height int

	titleLabel   *tui.Tile
	promptLabel  *tui.Tile
	inputField   *tui.Tile
	buttonBox    *tui.Box
	okButton     *tui.Tile
	cancelButton *tui.Tile
	dialogBox    *tui.Box

	selectedButton int // 0: OK, 1: Cancel
	onOK           func(base.Runtime, string)
	onCancel       func(base.Runtime)
	initialValue   string
}

// NewInputDialog returns InputDialog.
func NewInputDialog(title, prompt, initialValue string, onOK func(base.Runtime, string), onCancel func(base.Runtime)) *InputDialog {
	id := &InputDialog{
		selectedButton: 0,
		onOK:           onOK,
		onCancel:       onCancel,
		initialValue:   initialValue,
	}
	id.Init(title, prompt)
	return id
}

// Init is initialize InputDialog.
func (id *InputDialog) Init(title, prompt string) {
	id.titleLabel = tui.NewCenteredLabelTile(title)
	id.titleLabel.FlexibleWidth = true
	id.titleLabel.MinimumHeight = 1

	id.promptLabel = tui.NewLabelTile(prompt)
	id.promptLabel.FlexibleWidth = true
	id.promptLabel.MinimumHeight = 1

	id.inputField = tui.NewEditTile()
	id.inputField.FlexibleWidth = true
	id.inputField.MinimumHeight = 1
	id.inputField.TextBox.ShowCursor = true

	if id.initialValue != "" {
		id.inputField.TextBox.InsertString(id.initialValue)
	}

	id.okButton = tui.NewCenteredLabelTile("[ OK ]")
	id.okButton.MinimumWidth = 8
	id.okButton.MinimumHeight = 1

	id.cancelButton = tui.NewCenteredLabelTile("[ Cancel ]")
	id.cancelButton.MinimumWidth = 12
	id.cancelButton.MinimumHeight = 1

	id.buttonBox = tui.NewHBox(
		id.okButton,
		tui.NewFixedTile(&presenter.LabelTextPresenter{Text: "  "}, 2, 1), // スペーサー
		id.cancelButton,
	)

	id.dialogBox = tui.NewVBox(
		id.titleLabel,
		tui.NewHorizontalSeparator(),
		id.promptLabel,
		id.inputField,
		tui.NewHorizontalSeparator(),
		id.buttonBox,
	)

	id.updateButtonStyles()
}

func (id *InputDialog) updateButtonStyles() {
	// highlight selected button
	if !id.inputField.TextBox.ShowCursor {
		if id.selectedButton == 0 {
			if labelPresenter, ok := id.okButton.TextPresenter.(*presenter.LabelTextPresenter); ok {
				labelPresenter.Text = "> OK <"
			}
			if labelPresenter, ok := id.cancelButton.TextPresenter.(*presenter.LabelTextPresenter); ok {
				labelPresenter.Text = "[ Cancel ]"
			}
		} else {
			if labelPresenter, ok := id.okButton.TextPresenter.(*presenter.LabelTextPresenter); ok {
				labelPresenter.Text = "[ OK ]"
			}
			if labelPresenter, ok := id.cancelButton.TextPresenter.(*presenter.LabelTextPresenter); ok {
				labelPresenter.Text = "> Cancel <"
			}
		}
	} else {
		// normal shown if focus on input field
		if labelPresenter, ok := id.okButton.TextPresenter.(*presenter.LabelTextPresenter); ok {
			labelPresenter.Text = "[ OK ]"
		}
		if labelPresenter, ok := id.cancelButton.TextPresenter.(*presenter.LabelTextPresenter); ok {
			labelPresenter.Text = "[ Cancel ]"
		}
	}
}

func (id *InputDialog) Move(x, y int) {
	id.x = x
	id.y = y
	id.dialogBox.Move(x, y)
}

func (id *InputDialog) Layout(width, height int) {
	id.Width = width
	id.Height = height
	id.dialogBox.Layout(width, height)
}

func (id *InputDialog) MinimumSize(width, height int) (int, int) {
	return id.dialogBox.MinimumSize(width, height)
}

func (id *InputDialog) Update() {
	id.dialogBox.Update()
}

func (id *InputDialog) Draw(g *base.Graphics) {
	id.dialogBox.Draw(g)
}

func (id *InputDialog) Handle(ev base.Event) {
	if keyEvent, ok := ev.GetSource().(*tcell.EventKey); ok {
		switch keyEvent.Key() {
		case tcell.KeyEnter:
			if id.inputField.TextBox.ShowCursor {
				// same operation to ok button, when enter on input field
				if id.onOK != nil {
					inputValue := id.getInputValue()
					id.onOK(ev.GetRuntime(), inputValue)
				}
			} else {
				// execute button
				if id.selectedButton == 0 {
					if id.onOK != nil {
						inputValue := id.getInputValue()
						id.onOK(ev.GetRuntime(), inputValue)
					}
				} else {
					if id.onCancel != nil {
						id.onCancel(ev.GetRuntime())
					}
				}
			}
			return
		case tcell.KeyEscape:
			// cancel by escape
			if id.onCancel != nil {
				id.onCancel(ev.GetRuntime())
			}
			return
		case tcell.KeyLeft, tcell.KeyRight:
			if id.inputField.TextBox.ShowCursor {
				id.inputField.Handle(ev)
				return
			}
		default:
			if id.inputField.TextBox.ShowCursor {
				id.inputField.Handle(ev)
			}
			return
		}
	}
}

func (id *InputDialog) getInputValue() string {
	tb := id.inputField.TextBox
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

func (id *InputDialog) Traverse(fm *tui.FocusManager) {
	if id.IsFocusable() {
		fm.Register(id)
	}
}

func (id *InputDialog) Focus(on bool) {
	if on {
		id.inputField.TextBox.ShowCursor = true
	}
}

func (id *InputDialog) SubFocusFirst() {
	id.inputField.TextBox.ShowCursor = true
}

func (id *InputDialog) SubFocusPrev() bool {
	if id.inputField.TextBox.ShowCursor {
		// focus to cancel button
		id.inputField.TextBox.ShowCursor = false
		id.selectedButton = 1
		id.updateButtonStyles()
		return true
	} else {
		if id.selectedButton == 1 {
			id.selectedButton = 0
			id.updateButtonStyles()
			return true
		} else {
			id.inputField.TextBox.ShowCursor = true
			id.updateButtonStyles()
			return true
		}
	}
}

func (id *InputDialog) SubFocusNext() bool {
	if id.inputField.TextBox.ShowCursor {
		// focus to ok button
		id.inputField.TextBox.ShowCursor = false
		id.selectedButton = 0
		id.updateButtonStyles()
		return true
	} else {
		if id.selectedButton == 0 {
			id.selectedButton = 1
			id.updateButtonStyles()
			return true
		} else {
			id.inputField.TextBox.ShowCursor = true
			id.updateButtonStyles()
			return true
		}
	}
}

func (id *InputDialog) SubFocusLast() {
	id.inputField.TextBox.ShowCursor = false
	id.selectedButton = 1
	id.updateButtonStyles()
}

func (id *InputDialog) IsFocusable() bool {
	return true
}

func (id *InputDialog) IsFlexibleWidth() bool {
	return id.dialogBox.IsFlexibleWidth()
}

func (id *InputDialog) IsFlexibleHeight() bool {
	return id.dialogBox.IsFlexibleHeight()
}

// SetTitle updates the dialog title
func (id *InputDialog) SetTitle(title string) {
	if labelPresenter, ok := id.titleLabel.TextPresenter.(*presenter.LabelTextPresenter); ok {
		labelPresenter.Text = title
	}
}

// SetPrompt updates the dialog prompt
func (id *InputDialog) SetPrompt(prompt string) {
	if labelPresenter, ok := id.promptLabel.TextPresenter.(*presenter.LabelTextPresenter); ok {
		labelPresenter.Text = prompt
	}
}

// SetInputValue sets the input field value
func (id *InputDialog) SetInputValue(value string) {
	tb := id.inputField.TextBox
	doc := tb.GetDocument()
	doc.Clear()
	doc.InsertString(0, 0, value)
	tb.MoveReset()
}

// GetInputValue returns the current input field value
func (id *InputDialog) GetInputValue() string {
	return id.getInputValue()
}

// ClearInput clears the input field
func (id *InputDialog) ClearInput() {
	doc := id.inputField.TextBox.GetDocument()
	doc.Clear()
}
