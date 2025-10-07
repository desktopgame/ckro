package llm

import (
	"context"
	"encoding/json"
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

	mcpClients  map[string]*McpClient
	tool2client map[string]*McpClient
	tools       []string
	toolParams  []openai.ChatCompletionToolUnionParam
	allowTools  []string
}

func (cm *ChatManager) Init(client *openai.Client, model string, systemPrompt string, mcpClients map[string]*McpClient) {
	cm.client = client
	cm.model = model
	cm.systemPrompt = systemPrompt
	cm.mcpClients = mcpClients
	cm.backgroundToken = make(chan int)
	cm.backgroundDone = false
}

func (cm *ChatManager) background(ctx context.Context) error {
	// register tools
	for _, mcpClient := range cm.mcpClients {
		if mcpClient.session == nil {
			continue
		}
		tools, err := mcpClient.Tools(ctx)
		if err == nil {
			for _, tool := range tools {
				cm.tools = append(cm.tools, tool.Name)
				cm.tool2client[tool.Name] = mcpClient

				t := openai.ChatCompletionFunctionToolParam{
					Function: openai.FunctionDefinitionParam{
						Name:        tool.Name,
						Description: openai.String(tool.Description),
						Parameters:  tool.Parameters,
					},
				}
				tUnion := openai.ChatCompletionToolUnionParam{
					OfFunction: &t,
				}
				cm.toolParams = append(cm.toolParams, tUnion)
			}
		}
	}

	// initialize conversation
	cm.inputList = nil
	cm.backgroundToken <- 0
	return nil
}

func (cm *ChatManager) Setup() {
	cm.tools = nil
	cm.tool2client = map[string]*McpClient{}
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
		client:     cm.client,
		inputList:  cm.inputList,
		model:      cm.model,
		toolParams: cm.toolParams,
		EventBase: EventBase[*openai.ChatCompletion]{
			ch: input,
		},
	}
}

func (cm *ChatManager) turn(ctx context.Context, response *openai.ChatCompletion, input chan Status, output chan Event) error {
	cm.inputList = append(cm.inputList, response.Choices[0].Message.ToParam())
	toolCalls := response.Choices[0].Message.ToolCalls
	if len(toolCalls) > 0 {
		for _, toolCall := range toolCalls {
			fn := toolCall.Function

			var args interface{}
			err := json.Unmarshal([]byte(fn.Arguments), &args)
			if err == nil {
				if client, ok := cm.tool2client[fn.Name]; ok {
					confirmEvent := ConfirmEvent{
						ToolName: fn.Name,
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
							switch status {
							case Complete:
								chatEvent := cm.toolUse(ctx, toolEvent.result, toolCall, input)
								output <- &chatEvent

								status = <-input
								switch status {
								case Complete:
									return cm.turn(ctx, chatEvent.result, input, output)
								case Error:
									return chatEvent.GetError()
								}
							case Error:
								return toolEvent.GetError()
							}
						} else {
							message := fmt.Sprintf("user rejected a %s execute", fn.Name)
							cm.inputList = append(cm.inputList, openai.UserMessage(message))

							chatEvent := ChatEvent{
								client:     cm.client,
								inputList:  cm.inputList,
								model:      cm.model,
								toolParams: cm.toolParams,
								EventBase: EventBase[*openai.ChatCompletion]{
									ch: input,
								},
							}
							output <- &chatEvent

							status = <-input
							switch status {
							case Complete:
								return cm.turn(ctx, chatEvent.result, input, output)
							case Error:
								return chatEvent.GetError()
							}
						}
					}
				} else {
					message := fmt.Sprintf("%s is not found.", fn.Name)
					cm.inputList = append(cm.inputList, openai.SystemMessage(message))

					chatEvent := ChatEvent{
						client:     cm.client,
						inputList:  cm.inputList,
						model:      cm.model,
						toolParams: cm.toolParams,
						EventBase: EventBase[*openai.ChatCompletion]{
							ch: input,
						},
					}
					output <- &chatEvent

					status := <-input
					switch status {
					case Complete:
						return cm.turn(ctx, chatEvent.result, input, output)
					case Error:
						return chatEvent.GetError()
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
	return nil
}

func (cm *ChatManager) Post(ctx context.Context, message string, output chan Event) error {
	if !cm.backgroundDone {
		<-cm.backgroundToken
		cm.backgroundDone = true
	}
	cm.inputList = append(cm.inputList, openai.UserMessage(message))

	input := make(chan Status)
	defer close(input)

	chatEvent := ChatEvent{
		client:     cm.client,
		inputList:  cm.inputList,
		model:      cm.model,
		toolParams: cm.toolParams,
		EventBase: EventBase[*openai.ChatCompletion]{
			ch: input,
		},
	}
	output <- &chatEvent

	status := <-input
	if status == Cancel {
		return nil
	}

	err := cm.turn(ctx, chatEvent.result, input, output)
	if err != nil {
		output <- &ErrorEvent{e: err}
	}
	return err
}
