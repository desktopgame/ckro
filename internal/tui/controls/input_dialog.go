package controls

import (
	"github.com/desktopgame/ckro/internal/tui"
	"github.com/desktopgame/ckro/internal/tui/base"
	"github.com/desktopgame/ckro/internal/tui/presenter"
	"github.com/gdamore/tcell/v2"
)

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
	onOK           func(base.Stackable, string)
	onCancel       func(base.Stackable)
	initialValue   string
}

func NewInputDialog(title, prompt, initialValue string, onOK func(base.Stackable, string), onCancel func(base.Stackable)) *InputDialog {
	id := &InputDialog{
		selectedButton: 0, // デフォルトでOKを選択
		onOK:           onOK,
		onCancel:       onCancel,
		initialValue:   initialValue,
	}
	id.Init(title, prompt)
	return id
}

func (id *InputDialog) Init(title, prompt string) {
	// タイトルラベル
	id.titleLabel = tui.NewCenteredLabelTile(title)
	id.titleLabel.FlexibleWidth = true
	id.titleLabel.MinimumHeight = 1

	// プロンプトラベル
	id.promptLabel = tui.NewLabelTile(prompt)
	id.promptLabel.FlexibleWidth = true
	id.promptLabel.MinimumHeight = 1

	// 入力フィールド
	id.inputField = tui.NewEditTile()
	id.inputField.FlexibleWidth = true
	id.inputField.MinimumHeight = 1
	id.inputField.TextBox.ShowCursor = true

	// 初期値を設定
	if id.initialValue != "" {
		doc := id.inputField.TextBox.GetDocument()
		doc.InsertString(id.initialValue)
	}

	// ボタン
	id.okButton = tui.NewCenteredLabelTile("[ OK ]")
	id.okButton.MinimumWidth = 8
	id.okButton.MinimumHeight = 1

	id.cancelButton = tui.NewCenteredLabelTile("[ Cancel ]")
	id.cancelButton.MinimumWidth = 12
	id.cancelButton.MinimumHeight = 1

	// ボタンを水平に配置
	id.buttonBox = tui.NewHBox(
		id.okButton,
		tui.NewFixedTile(&presenter.LabelTextPresenter{Text: "  "}, 2, 1), // スペーサー
		id.cancelButton,
	)

	// 全体を垂直に配置
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
	// 選択されたボタンをハイライト表示
	if id.selectedButton == 0 {
		// OKボタンを選択状態に
		if labelPresenter, ok := id.okButton.TextPresenter.(*presenter.LabelTextPresenter); ok {
			labelPresenter.Text = "> OK <"
		}
		if labelPresenter, ok := id.cancelButton.TextPresenter.(*presenter.LabelTextPresenter); ok {
			labelPresenter.Text = "[ Cancel ]"
		}
	} else {
		// Cancelボタンを選択状態に
		if labelPresenter, ok := id.okButton.TextPresenter.(*presenter.LabelTextPresenter); ok {
			labelPresenter.Text = "[ OK ]"
		}
		if labelPresenter, ok := id.cancelButton.TextPresenter.(*presenter.LabelTextPresenter); ok {
			labelPresenter.Text = "> Cancel <"
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
				// 入力フィールドフォーカス時はOKボタンと同じ動作
				if id.onOK != nil {
					inputValue := id.getInputValue()
					id.onOK(ev.GetStackable(), inputValue)
				}
			} else {
				// ボタンフォーカス時は選択されたボタンを実行
				if id.selectedButton == 0 {
					// OKボタン
					if id.onOK != nil {
						inputValue := id.getInputValue()
						id.onOK(ev.GetStackable(), inputValue)
					}
				} else {
					// Cancelボタン
					if id.onCancel != nil {
						id.onCancel(ev.GetStackable())
					}
				}
			}
			return
		case tcell.KeyEscape:
			// Escapeキーでキャンセル
			if id.onCancel != nil {
				id.onCancel(ev.GetStackable())
			}
			return
		case tcell.KeyLeft, tcell.KeyRight:
			// 入力フィールドフォーカス時は入力フィールドに転送
			if id.inputField.TextBox.ShowCursor {
				id.inputField.Handle(ev)
				return
			}
		default:
			// 入力フィールドフォーカス時は文字入力を転送
			if id.inputField.TextBox.ShowCursor {
				id.inputField.Handle(ev)
			}
			return
		}
	}
}

func (id *InputDialog) getInputValue() string {
	// 入力フィールドの内容を取得
	doc := id.inputField.TextBox.GetDocument()
	buffer := doc.GetBuffer()
	if buffer.GetLineCount() > 0 {
		return buffer.GetLineAt(0).GetContent()
	}
	return ""
}

func (id *InputDialog) Traverse(fm *tui.FocusManager) {
	if id.IsFocusable() {
		fm.Register(id)
	}
}

func (id *InputDialog) Focus(on bool) {
	// フォーカス状態の管理
	if on {
		id.inputField.TextBox.ShowCursor = true
	}
}

func (id *InputDialog) SubFocusFirst() {
	id.inputField.TextBox.ShowCursor = true
}

func (id *InputDialog) SubFocusPrev() bool {
	if id.inputField.TextBox.ShowCursor {
		id.inputField.TextBox.ShowCursor = false
		id.selectedButton = 1 // Cancelボタン
		id.updateButtonStyles()
	} else {
		if id.selectedButton == 0 {
			id.inputField.TextBox.ShowCursor = true
		} else {
			id.selectedButton = 0
			id.updateButtonStyles()
		}
	}
	return false
}

func (id *InputDialog) SubFocusNext() bool {
	if id.inputField.TextBox.ShowCursor {
		id.inputField.TextBox.ShowCursor = false
		id.selectedButton = 0 // OKボタン
		id.updateButtonStyles()
	} else {
		if id.selectedButton == 1 {
			id.inputField.TextBox.ShowCursor = true
		} else {
			id.selectedButton = 1
			id.updateButtonStyles()
		}
	}
	return false
}

func (id *InputDialog) SubFocusLast() {
	id.inputField.TextBox.ShowCursor = false
	id.selectedButton = 1 // Cancelボタン
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
	doc := id.inputField.TextBox.GetDocument()
	doc.Init()
	doc.InsertString(value)
}

// GetInputValue returns the current input field value
func (id *InputDialog) GetInputValue() string {
	return id.getInputValue()
}

// ClearInput clears the input field
func (id *InputDialog) ClearInput() {
	doc := id.inputField.TextBox.GetDocument()
	doc.Init()
}
