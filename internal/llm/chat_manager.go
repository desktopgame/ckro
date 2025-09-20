package llm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/openai/openai-go/v2"
)

type ChatManager struct {
	client          *openai.Client
	model           string
	systemPrompt    string
	inputList       []openai.ChatCompletionMessageParamUnion
	backgroundToken chan int
	backgroundDone  bool

	tools2client map[string]*McpClient
	tools        []string
	allowTools   []string
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

func (cm *ChatManager) toolUse(ctx context.Context, response *mcp.CallToolResult, toolCall openai.ChatCompletionMessageToolCallUnion, input chan Status) ChatEvent {
	sb := strings.Builder{}
	for i := 0; i < len(response.Content); i++ {
		c := response.Content[i]
		switch v := c.(type) {
		case *mcp.TextContent:
			sb.WriteString(v.Text)
			sb.WriteString("\n")
		default:
		}
	}
	cm.inputList = append(cm.inputList, openai.ToolMessage(sb.String(), toolCall.ID))
	return ChatEvent{
		client:    cm.client,
		inputList: cm.inputList,
		model:     cm.model,
		EventBase: EventBase[*openai.ChatCompletion]{
			ch: input,
		},
	}
}

func (cm *ChatManager) turn(ctx context.Context, response *openai.ChatCompletion, input chan Status, output chan Event) {
	cm.inputList = append(cm.inputList, response.Choices[0].Message.ToParam())
	toolCalls := response.Choices[0].Message.ToolCalls
	if len(toolCalls) > 0 {
		for _, toolCall := range toolCalls {
			fn := toolCall.Function

			var args interface{}
			err := json.Unmarshal([]byte(fn.Arguments), &args)
			if err == nil {
				if client, ok := cm.tools2client[fn.Name]; ok {
					confirmEvent := ConfirmEvent{
						EventBase: EventBase[struct{}]{
							ch: input,
						},
					}
					output <- &confirmEvent

					status := <-input
					if status == Complete {
						if confirmEvent.Approve {
							toolEvent := ToolEvent{
								mcpClient: client,
								name:      fn.Name,
								args:      args,
								EventBase: EventBase[*mcp.CallToolResult]{
									ch: input,
								},
							}
							output <- &toolEvent

							status = <-input
							if status == Complete {
								chatEvent := cm.toolUse(ctx, toolEvent.result, toolCall, input)
								output <- &chatEvent

								status = <-input
								if status == Complete {
									cm.turn(ctx, chatEvent.result, input, output)
								}
							}
						} else {
							message := fmt.Sprintf("user rejected a %s execute", fn.Name)
							cm.inputList = append(cm.inputList, openai.UserMessage(message))

							chatEvent := ChatEvent{
								client:    cm.client,
								inputList: cm.inputList,
								model:     cm.model,
								EventBase: EventBase[*openai.ChatCompletion]{
									ch: input,
								},
							}
							output <- &chatEvent

							status = <-input
							if status == Complete {
								cm.turn(ctx, chatEvent.result, input, output)
							}
						}
					}
				} else {
					message := fmt.Sprintf("%s is not found.", fn.Name)
					cm.inputList = append(cm.inputList, openai.SystemMessage(message))

					chatEvent := ChatEvent{
						client:    cm.client,
						inputList: cm.inputList,
						model:     cm.model,
						EventBase: EventBase[*openai.ChatCompletion]{
							ch: input,
						},
					}
					output <- &chatEvent

					status := <-input
					if status == Complete {
						cm.turn(ctx, chatEvent.result, input, output)
					}
				}
			}
		}
	} else {
		messageEvent := MessageEvent{
			EventBase: EventBase[*openai.ChatCompletion]{
				result: response,
			},
		}
		output <- &messageEvent
	}
}

func (cm *ChatManager) Post(ctx context.Context, message string, output chan Event) (string, error) {
	if !cm.backgroundDone {
		<-cm.backgroundToken
		cm.backgroundDone = true
	}
	cm.inputList = append(cm.inputList, openai.UserMessage(message))

	input := make(chan Status)
	defer close(input)

	chatEvent := ChatEvent{
		client:    cm.client,
		inputList: cm.inputList,
		model:     cm.model,
		EventBase: EventBase[*openai.ChatCompletion]{
			ch: input,
		},
	}
	output <- &chatEvent

	status := <-input
	if status == Cancel {
		return "", errors.New("chat is cancelled")
	}
	cm.turn(ctx, chatEvent.result, input, output)
	return "", nil
}
