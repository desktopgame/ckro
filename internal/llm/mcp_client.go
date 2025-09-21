package llm

import (
	"context"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type McpClient struct {
	client  *mcp.Client
	session *mcp.ClientSession
}

func (m *McpClient) Init() {
	m.client = mcp.NewClient(&mcp.Implementation{Name: "ckro", Version: "v0.1.0"}, nil)
}

func (m *McpClient) Connect(ctx context.Context, name string, args ...string) error {
	if m.session == nil {
		transport := &mcp.CommandTransport{Command: exec.Command(name, args...)}
		session, err := m.client.Connect(ctx, transport, nil)

		if err == nil {
			m.session = session
			session.InitializeResult()
			return nil
		}
		return err
	}
	return nil
}

func (m *McpClient) ConnectLocal(ctx context.Context, server *mcp.Server) error {
	if m.session == nil {
		t1, t2 := mcp.NewInMemoryTransports()
		_, err := server.Connect(ctx, t1, nil)

		if err == nil {
			session, err := m.client.Connect(ctx, t2, nil)

			if err == nil {
				m.session = session
				session.InitializeResult()
				return nil
			}
			return err
		}
		return err
	}
	return nil
}

func (m *McpClient) Call(ctx context.Context, toolName string, args interface{}) (*mcp.CallToolResult, error) {
	params := &mcp.CallToolParams{
		Name:      toolName,
		Arguments: args,
	}
	return m.session.CallTool(ctx, params)
}

func (m *McpClient) Close() {
	if m.session != nil {
		m.session.Close()
		m.session = nil
	}
}
