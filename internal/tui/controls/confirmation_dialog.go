package controls

import (
	"github.com/desktopgame/ckro/internal/tui"
	"github.com/desktopgame/ckro/internal/tui/base"
	"github.com/desktopgame/ckro/internal/tui/presenter"
	"github.com/gdamore/tcell/v2"
)

// ConfirmationDialog is choice a yes or no dialog.
type ConfirmationDialog struct {
	x, y          int
	Width, Height int

	titleLabel   *tui.Tile
	messageLabel *tui.Tile
	buttonBox    *tui.Box
	yesButton    *tui.Tile
	noButton     *tui.Tile
	dialogBox    *tui.Box

	selectedButton int // 0: Yes, 1: No
	onYes          func(base.Runtime)
	onNo           func(base.Runtime)
}

// NewConfirmationDialog returns ConfirmationDialog.
func NewConfirmationDialog(title, message string, onYes, onNo func(base.Runtime)) *ConfirmationDialog {
	cd := &ConfirmationDialog{
		selectedButton: 0,
		onYes:          onYes,
		onNo:           onNo,
	}
	cd.Init(title, message)
	return cd
}

// Init is initialize ConfirmationDialog.
func (cd *ConfirmationDialog) Init(title, message string) {
	cd.titleLabel = tui.NewCenteredLabelTile(title)
	cd.titleLabel.FlexibleWidth = true
	cd.titleLabel.MinimumHeight = 1

	cd.messageLabel = tui.NewCenteredLabelTile(message)
	cd.messageLabel.FlexibleWidth = true
	cd.messageLabel.MinimumHeight = 3

	cd.yesButton = tui.NewCenteredLabelTile("[ YES ]")
	cd.yesButton.MinimumWidth = 8
	cd.yesButton.MinimumHeight = 3

	cd.noButton = tui.NewCenteredLabelTile("[ NO ]")
	cd.noButton.MinimumWidth = 8
	cd.noButton.MinimumHeight = 3

	cd.buttonBox = tui.NewHBox(
		cd.yesButton,
		tui.NewFixedTile(&presenter.LabelTextPresenter{Text: "  "}, 2, 1),
		cd.noButton,
	)

	cd.dialogBox = tui.NewVBox(
		cd.titleLabel,
		tui.NewHorizontalSeparator(),
		cd.messageLabel,
		tui.NewHorizontalSeparator(),
		cd.buttonBox,
	)

	cd.updateButtonStyles()
}

func (cd *ConfirmationDialog) updateButtonStyles() {
	// highlight selected button
	if cd.selectedButton == 0 {
		if labelPresenter, ok := cd.yesButton.TextPresenter.(*presenter.LabelTextPresenter); ok {
			labelPresenter.Text = "> YES <"
		}
		if labelPresenter, ok := cd.noButton.TextPresenter.(*presenter.LabelTextPresenter); ok {
			labelPresenter.Text = "[ NO ]"
		}
	} else {
		if labelPresenter, ok := cd.yesButton.TextPresenter.(*presenter.LabelTextPresenter); ok {
			labelPresenter.Text = "[ YES ]"
		}
		if labelPresenter, ok := cd.noButton.TextPresenter.(*presenter.LabelTextPresenter); ok {
			labelPresenter.Text = "> NO <"
		}
	}
}

func (cd *ConfirmationDialog) Move(x, y int) {
	cd.x = x
	cd.y = y
	cd.dialogBox.Move(x, y)
}

func (cd *ConfirmationDialog) Layout(width, height int) {
	cd.Width = width
	cd.Height = height
	cd.dialogBox.Layout(width, height)
}

func (cd *ConfirmationDialog) MinimumSize(width, height int) (int, int) {
	return cd.dialogBox.MinimumSize(width, height)
}

func (cd *ConfirmationDialog) Update() {
	cd.dialogBox.Update()
}

func (cd *ConfirmationDialog) Draw(g *base.Graphics) {
	cd.dialogBox.Draw(g)
}

func (cd *ConfirmationDialog) Handle(ev base.Event) {
	if keyEvent, ok := ev.GetSource().(*tcell.EventKey); ok {
		switch keyEvent.Key() {
		case tcell.KeyEnter:
			// execute selected button
			if cd.selectedButton == 0 {
				if cd.onYes != nil {
					cd.onYes(ev.GetRuntime())
				}
			} else {
				if cd.onNo != nil {
					cd.onNo(ev.GetRuntime())
				}
			}
			return
		case tcell.KeyEscape:
			// cancel by escape
			if cd.onNo != nil {
				cd.onNo(ev.GetRuntime())
			}
			return
		}

		if keyEvent.Rune() == 'y' || keyEvent.Rune() == 'Y' {
			if cd.onYes != nil {
				cd.onYes(ev.GetRuntime())
			}
			return
		}
		if keyEvent.Rune() == 'n' || keyEvent.Rune() == 'N' {
			if cd.onNo != nil {
				cd.onNo(ev.GetRuntime())
			}
			return
		}
	}
}

func (cd *ConfirmationDialog) Traverse(fm *tui.FocusManager) {
	if cd.IsFocusable() {
		fm.Register(cd)
	}
}

func (cd *ConfirmationDialog) Focus(on bool) {
}

func (cd *ConfirmationDialog) SubFocusFirst() {
	cd.selectedButton = 0
	cd.updateButtonStyles()
}

func (cd *ConfirmationDialog) SubFocusPrev() bool {
	cd.selectedButton = 1 - cd.selectedButton
	cd.updateButtonStyles()
	return false
}

func (cd *ConfirmationDialog) SubFocusNext() bool {
	cd.selectedButton = 1 - cd.selectedButton
	cd.updateButtonStyles()
	return false
}

func (cd *ConfirmationDialog) SubFocusLast() {
	cd.selectedButton = 1
	cd.updateButtonStyles()
}

func (cd *ConfirmationDialog) IsFocusable() bool {
	return true
}

func (cd *ConfirmationDialog) IsFlexibleWidth() bool {
	return cd.dialogBox.IsFlexibleWidth()
}

func (cd *ConfirmationDialog) IsFlexibleHeight() bool {
	return cd.dialogBox.IsFlexibleHeight()
}

// SetTitle updates the dialog title
func (cd *ConfirmationDialog) SetTitle(title string) {
	if labelPresenter, ok := cd.titleLabel.TextPresenter.(*presenter.LabelTextPresenter); ok {
		labelPresenter.Text = title
	}
}

// SetMessage updates the dialog message
func (cd *ConfirmationDialog) SetMessage(message string) {
	if labelPresenter, ok := cd.messageLabel.TextPresenter.(*presenter.LabelTextPresenter); ok {
		labelPresenter.Text = message
	}
}

// GetSelectedButton returns the currently selected button (0: Yes, 1: No)
func (cd *ConfirmationDialog) GetSelectedButton() int {
	return cd.selectedButton
}

// SetSelectedButton sets the selected button (0: Yes, 1: No)
func (cd *ConfirmationDialog) SetSelectedButton(button int) {
	if button >= 0 && button <= 1 {
		cd.selectedButton = button
		cd.updateButtonStyles()
	}
}
