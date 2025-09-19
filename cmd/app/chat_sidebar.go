package main

import (
	"strings"

	"github.com/desktopgame/ckro/internal/tui"
	"github.com/desktopgame/ckro/internal/tui/base"
	"github.com/desktopgame/ckro/internal/tui/presenter"
	"github.com/gdamore/tcell/v2"
)

type ChatSidebar struct {
	x, y          int
	Width, Height int

	historyArea *tui.Tile
	inputArea   *tui.Tile
	sidebarBox  *tui.Box

	onMessageSend func(message string)
	messages      []ChatMessageEntry
}

type ChatMessageEntry struct {
	Role    string // "user" or "assistant"
	Content string
}

func NewChatSidebar(onMessageSend func(message string)) *ChatSidebar {
	cs := &ChatSidebar{
		onMessageSend: onMessageSend,
		messages:      make([]ChatMessageEntry, 0),
	}
	cs.Init()
	return cs
}

func (cs *ChatSidebar) Init() {
	// 会話履歴表示エリア
	cs.historyArea = &tui.Tile{}
	cs.historyArea.Init()
	cs.historyArea.MinimumWidth = 30
	cs.historyArea.FlexibleHeight = true
	cs.historyArea.TextPresenter = &presenter.EditTextPresenter{}

	// 入力エリア
	cs.inputArea = &tui.Tile{}
	cs.inputArea.Init()
	cs.inputArea.MinimumWidth = 30
	cs.inputArea.MinimumHeight = 3
	cs.inputArea.TextBox.ShowCursor = true
	cs.inputArea.TextPresenter = &presenter.EditTextPresenter{}

	// 垂直レイアウトで組み合わせ
	cs.sidebarBox = tui.NewVBox(
		tui.WithFrame(cs.historyArea),
		tui.WithFrame(cs.inputArea),
	)

	cs.updateHistoryDisplay()
}

func (cs *ChatSidebar) Move(x, y int) {
	cs.x = x
	cs.y = y
	cs.sidebarBox.Move(x, y)
}

func (cs *ChatSidebar) Layout(width, height int) {
	cs.Width = width
	cs.Height = height
	cs.sidebarBox.Layout(width, height)
}

func (cs *ChatSidebar) MinimumSize(width, height int) (int, int) {
	return cs.sidebarBox.MinimumSize(width, height)
}

func (cs *ChatSidebar) Update() {
	cs.sidebarBox.Update()
}

func (cs *ChatSidebar) Draw(g *base.Graphics) {
	cs.sidebarBox.Draw(g)
}

func (cs *ChatSidebar) Handle(ev base.Event) {
	if keyEvent, ok := ev.GetSource().(*tcell.EventKey); ok {
		// 入力エリアにフォーカスがある場合の特別処理
		if cs.inputArea.TextBox.ShowCursor {
			switch keyEvent.Key() {
			case tcell.KeyEnter:
				// Enterキーでメッセージ送信
				message := cs.getInputMessage()
				if message != "" {
					cs.sendMessage(message)
					cs.clearInput()
				}
				return
			}
		}
	}

	// 入力エリアにフォーカスがある場合は入力エリアにイベントを転送
	if cs.inputArea.TextBox.ShowCursor {
		cs.inputArea.Handle(ev)
	} else if cs.historyArea.TextBox.ShowCursor {
		cs.historyArea.Handle(ev)
	}
}

func (cs *ChatSidebar) getInputMessage() string {
	doc := cs.inputArea.TextBox.GetDocument()
	buffer := doc.GetBuffer()

	var lines []string
	for i := 0; i < buffer.GetLineCount(); i++ {
		lines = append(lines, buffer.GetLineAt(i).GetContent())
	}

	// 空行を除去して結合
	var nonEmptyLines []string
	for _, line := range lines {
		if line != "" {
			nonEmptyLines = append(nonEmptyLines, line)
		}
	}

	if len(nonEmptyLines) == 0 {
		return ""
	}

	return strings.Join(nonEmptyLines, "\n")
}

func (cs *ChatSidebar) clearInput() {
	doc := cs.inputArea.TextBox.GetDocument()
	doc.Init()
	cs.inputArea.TextBox.CursorReset()
}

func (cs *ChatSidebar) sendMessage(message string) {
	// ユーザーメッセージを履歴に追加
	cs.AddMessage("user", message)

	// コールバック関数を呼び出し
	if cs.onMessageSend != nil {
		cs.onMessageSend(message)
	}
}

func (cs *ChatSidebar) updateHistoryDisplay() {
	doc := cs.historyArea.TextBox.GetDocument()
	doc.Init()

	for _, msg := range cs.messages {
		var prefix string
		switch msg.Role {
		case "user":
			prefix = "You: "
		case "assistant":
			prefix = "AI: "
		default:
			prefix = msg.Role + ": "
		}

		doc.InsertString(prefix + msg.Content + "\n\n")
	}

	// 最新メッセージが見えるようにスクロール
	cs.historyArea.TextBox.CursorReset()
	// カーソルを最後に移動
	buffer := doc.GetBuffer()
	if buffer.GetLineCount() > 0 {
		// 最後の行に移動
		for i := 0; i < buffer.GetLineCount()-1; i++ {
			doc.MoveDown()
		}
		// 行末に移動
		lastLine := buffer.GetLineAt(buffer.GetLineCount() - 1)
		for j := 0; j < len(lastLine.GetContent()); j++ {
			doc.MoveRight()
		}
	}
}

func (cs *ChatSidebar) Traverse(fm *tui.FocusManager) {
	cs.sidebarBox.Traverse(fm)
}

func (cs *ChatSidebar) IsFlexibleWidth() bool {
	return cs.sidebarBox.IsFlexibleWidth()
}

func (cs *ChatSidebar) IsFlexibleHeight() bool {
	return cs.sidebarBox.IsFlexibleHeight()
}

// AddMessage adds a new message to the chat history
func (cs *ChatSidebar) AddMessage(role, content string) {
	cs.messages = append(cs.messages, ChatMessageEntry{
		Role:    role,
		Content: content,
	})
	cs.updateHistoryDisplay()
}

// ClearHistory clears all chat messages
func (cs *ChatSidebar) ClearHistory() {
	cs.messages = make([]ChatMessageEntry, 0)
	cs.updateHistoryDisplay()
}

// GetMessages returns all chat messages
func (cs *ChatSidebar) GetMessages() []ChatMessageEntry {
	return cs.messages
}

// SetMinimumWidth sets the minimum width for the sidebar
func (cs *ChatSidebar) SetMinimumWidth(width int) {
	cs.historyArea.MinimumWidth = width
	cs.inputArea.MinimumWidth = width
}

// SetInputHeight sets the height for the input area
func (cs *ChatSidebar) SetInputHeight(height int) {
	cs.inputArea.MinimumHeight = height
}

// FocusInput focuses the input area
func (cs *ChatSidebar) FocusInput() {
	cs.historyArea.TextBox.ShowCursor = false
	cs.inputArea.TextBox.ShowCursor = true
}

// FocusHistory focuses the history area
func (cs *ChatSidebar) FocusHistory() {
	cs.inputArea.TextBox.ShowCursor = false
	cs.historyArea.TextBox.ShowCursor = true
}
