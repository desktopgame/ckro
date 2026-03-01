package main

import (
	"context"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func NewTimeMcp() *mcp.Server {
	type TimeArgs struct {
	}

	server := mcp.NewServer(&mcp.Implementation{Name: "time-server", Version: "v0.1.0"}, nil)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_current_time",
		Description: "Return current time string",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in TimeArgs) (*mcp.CallToolResult, any, error) {
		t := time.Now()

		res := mcp.CallToolResult{}
		res.Content = []mcp.Content{
			&mcp.TextContent{Text: t.String()},
		}
		return &res, nil, nil
	})
	return server
}
