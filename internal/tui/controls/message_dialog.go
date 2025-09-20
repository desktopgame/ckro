package controls

import (
	"github.com/desktopgame/ckro/internal/tui"
	"github.com/desktopgame/ckro/internal/tui/base"
	"github.com/desktopgame/ckro/internal/tui/presenter"
	"github.com/gdamore/tcell/v2"
)

type MessageDialog struct {
	x, y          int
	Width, Height int

	titleLabel   *tui.Tile
	messageLabel *tui.Tile
	okButton     *tui.Tile
	dialogBox    *tui.Box

	onOK func(base.Runtime)
}

func NewMessageDialog(title, message string, onOK func(base.Runtime)) *MessageDialog {
	md := &MessageDialog{
		onOK: onOK,
	}
	md.Init(title, message)
	return md
}

func (md *MessageDialog) Init(title, message string) {
	// タイトルラベル
	md.titleLabel = tui.NewCenteredLabelTile(title)
	md.titleLabel.FlexibleWidth = true
	md.titleLabel.MinimumHeight = 1

	// メッセージラベル
	md.messageLabel = tui.NewCenteredLabelTile(message)
	md.messageLabel.FlexibleWidth = true
	md.messageLabel.MinimumHeight = 3

	// OKボタン
	md.okButton = tui.NewCenteredLabelTile("> OK <")
	md.okButton.FlexibleWidth = true
	md.okButton.MinimumHeight = 3

	// 全体を垂直に配置
	md.dialogBox = tui.NewVBox(
		md.titleLabel,
		tui.NewHorizontalSeparator(),
		md.messageLabel,
		tui.NewHorizontalSeparator(),
		md.okButton,
	)
}

func (md *MessageDialog) Move(x, y int) {
	md.x = x
	md.y = y
	md.dialogBox.Move(x, y)
}

func (md *MessageDialog) Layout(width, height int) {
	md.Width = width
	md.Height = height
	md.dialogBox.Layout(width, height)
}

func (md *MessageDialog) MinimumSize(width, height int) (int, int) {
	return md.dialogBox.MinimumSize(width, height)
}

func (md *MessageDialog) Update() {
	md.dialogBox.Update()
}

func (md *MessageDialog) Draw(g *base.Graphics) {
	md.dialogBox.Draw(g)
}

func (md *MessageDialog) Handle(ev base.Event) {
	if keyEvent, ok := ev.GetSource().(*tcell.EventKey); ok {
		switch keyEvent.Key() {
		case tcell.KeyEnter, tcell.KeyEscape:
			// EnterキーまたはEscapeキーでOKボタンを実行
			if md.onOK != nil {
				md.onOK(ev.GetRuntime())
			}
			return
		}

		// スペースキーでもOKボタンを実行
		if keyEvent.Rune() == ' ' {
			if md.onOK != nil {
				md.onOK(ev.GetRuntime())
			}
			return
		}
	}
}

func (md *MessageDialog) Traverse(fm *tui.FocusManager) {
	if md.IsFocusable() {
		fm.Register(md)
	}
}

func (md *MessageDialog) Focus(on bool) {
	// フォーカス状態の管理
}

func (md *MessageDialog) SubFocusFirst() {
	// OKボタンにフォーカス
}

func (md *MessageDialog) SubFocusPrev() bool {
	// 単一ボタンなので移動なし
	return false
}

func (md *MessageDialog) SubFocusNext() bool {
	// 単一ボタンなので移動なし
	return false
}

func (md *MessageDialog) SubFocusLast() {
	// OKボタンにフォーカス
}

func (md *MessageDialog) IsFocusable() bool {
	return true
}

func (md *MessageDialog) IsFlexibleWidth() bool {
	return md.dialogBox.IsFlexibleWidth()
}

func (md *MessageDialog) IsFlexibleHeight() bool {
	return md.dialogBox.IsFlexibleHeight()
}

// SetTitle updates the dialog title
func (md *MessageDialog) SetTitle(title string) {
	if labelPresenter, ok := md.titleLabel.TextPresenter.(*presenter.LabelTextPresenter); ok {
		labelPresenter.Text = title
	}
}

// SetMessage updates the dialog message
func (md *MessageDialog) SetMessage(message string) {
	if labelPresenter, ok := md.messageLabel.TextPresenter.(*presenter.LabelTextPresenter); ok {
		labelPresenter.Text = message
	}
}

// SetButtonText updates the OK button text
func (md *MessageDialog) SetButtonText(text string) {
	if labelPresenter, ok := md.okButton.TextPresenter.(*presenter.LabelTextPresenter); ok {
		labelPresenter.Text = "> " + text + " <"
	}
}

// ExecuteOK manually executes the OK callback
func (md *MessageDialog) ExecuteOK(runtime base.Runtime) {
	if md.onOK != nil {
		md.onOK(runtime)
	}
}
