package llm

import (
	"context"

	"github.com/openai/openai-go/v2"
)

type ChatManager struct {
	client          *openai.Client
	model           string
	systemPrompt    string
	inputList       []openai.ChatCompletionMessageParamUnion
	backgroundToken chan int
	backgroundDone  bool
}

func (cm *ChatManager) Init(client *openai.Client, model string, systemPrompt string) {
	cm.client = client
	cm.model = model
	cm.systemPrompt = systemPrompt
	cm.backgroundToken = make(chan int)
	cm.backgroundDone = false
}

func (cm *ChatManager) background(ctx context.Context) error {
	cm.inputList = nil
	cm.inputList = append(cm.inputList, openai.SystemMessage(cm.systemPrompt))

	chatCompletion, err := cm.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: cm.inputList,
		Model:    cm.model,
	})
	if err == nil {
		cm.inputList = append(cm.inputList, chatCompletion.Choices[0].Message.ToParam())
		cm.backgroundToken <- 0
		return nil
	}
	return err
}

func (cm *ChatManager) Setup() {
	cm.backgroundDone = false
	go cm.background(context.TODO())
}

func (cm *ChatManager) Post(ctx context.Context, message string) (string, error) {
	if !cm.backgroundDone {
		<-cm.backgroundToken
		cm.backgroundDone = true
	}
	cm.inputList = append(cm.inputList, openai.UserMessage(message))
	chatCompletion, err := cm.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: cm.inputList,
		Model:    cm.model,
	})
	if err == nil {
		cm.inputList = append(cm.inputList, chatCompletion.Choices[0].Message.ToParam())
		return chatCompletion.Choices[0].Message.Content, nil
	}
	return "", err
}
