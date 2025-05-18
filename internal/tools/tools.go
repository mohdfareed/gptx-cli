package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mohdfareed/gptx-cli/internal/cfg"
)

type Tool func(ctx context.Context, params map[string]any) (string, error)

type ToolDef struct {
	Name    string
	Desc    string
	Params  map[string]any
	Handler Tool
}

type ToolCall struct {
	Name   string
	Params string
}

const WebSearchToolDef = "web-search"

// Tools wires up tool calls to handlers and emits results.
type Tools struct {
	tools map[string]ToolDef
}

// NewTools creates a tool manager and registers the requested tools.
func NewTools(config cfg.Config) *Tools {
	m := &Tools{
		tools: make(map[string]ToolDef),
	}

	// Register requested tools
	if config.Shell != nil {
		m.tools[ShellToolDef] = NewShellTool(config)
	}

	// Register requested tools
	if config.WebSearch != nil {
		m.tools[WebSearchToolDef] = ToolDef{}
	}

	return m
}

func (m *Tools) CallTool(ctx context.Context, call ToolCall) (string, error) {
	tool, ok := m.tools[call.Name]
	if !ok {
		return "", fmt.Errorf("unknown tool: %s", call.Name)
	}

	if call.Name == WebSearchToolDef {
		return "", nil // Skip web search tool (handled by client)
	}

	var params map[string]any
	err := json.Unmarshal([]byte(call.Params), &params)
	if err != nil {
		return "", fmt.Errorf("tool %s: %w", call.Name, err)
	}

	// Call the tool handler
	result, err := tool.Handler(ctx, params)
	if err != nil {
		return "", fmt.Errorf("tool %s: %w", call.Name, err)
	}
	return result, nil
}

func (m *Tools) GetTools() []ToolDef {
	tools := make([]ToolDef, 0, len(m.tools))
	for _, tool := range m.tools {
		if tool.Name == WebSearchToolDef {
			continue
		} // Skip web search tool (handled by client)
		tools = append(tools, tool)
	}
	return tools
}
