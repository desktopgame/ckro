package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func NewFileSystemMcp() *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: "file-system-server", Version: "v0.1.0"}, nil)

	type ListArgs struct {
		RelativePath string `json:"relativePath" jsonschema:"Relative path from current working directory; set '.' for root"`
	}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_entries",
		Description: "Return entries at specified location",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in ListArgs) (*mcp.CallToolResult, any, error) {
		wd, err := os.Getwd()
		if err != nil {
			return nil, nil, err
		}

		log.Println(in.RelativePath)

		dir := wd
		if in.RelativePath != "." {
			dir = filepath.Join(wd, in.RelativePath)
		}

		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil, nil, err
		}

		sb := strings.Builder{}
		for _, entry := range entries {
			if entry.IsDir() {
				sb.WriteString("[DIR] ")
				sb.WriteString(entry.Name())
				sb.WriteString("\n")
			}
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				sb.WriteString("[FILE] ")
				sb.WriteString(entry.Name())
				sb.WriteString("\n")
			}
		}

		res := mcp.CallToolResult{}
		res.Content = []mcp.Content{
			&mcp.TextContent{Text: sb.String()},
		}
		return &res, nil, nil
	})
	return server
}
