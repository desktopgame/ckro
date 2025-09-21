package llm

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ToolEvent struct {
	EventBase[*mcp.CallToolResult]
	mcpClient *McpClient
	name      string
	args      interface{}
	e         error
}

func (t *ToolEvent) Consume(ctx context.Context) {
	response, err := t.mcpClient.Call(ctx, t.name, t.args)
	if err == nil {
		t.result = response
		t.ch <- Complete
		return
	}
	t.e = err
	t.ch <- Error
}

func (t *ToolEvent) Cancel(ctx context.Context) {
	t.ch <- Cancel
}

func (t *ToolEvent) GetError() error {
	return t.e
}

func (t *ToolEvent) GetResult() *mcp.CallToolResult {
	return t.result
}
