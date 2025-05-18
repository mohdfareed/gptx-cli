package tools

import (
	"context"
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
	Params map[string]any
}

// Manager wires up tool calls to handlers and emits results.
type Manager struct {
	tools map[string]ToolDef
}

// NewManager creates a tool manager and registers the requested tools.
func NewManager(config cfg.Config) (*Manager, error) {
	m := &Manager{
		tools: make(map[string]ToolDef),
	}

	// Register requested tools
	if config.Shell != nil {
		m.tools[ShellToolDef] = NewShellTool(config)
	}

	return m, nil
}

func (m *Manager) CallTool(ctx context.Context, call ToolCall) (string, error) {
	tool, ok := m.tools[call.Name]
	if !ok {
		return "", fmt.Errorf("unknown tool: %s", call.Name)
	}

	// Call the tool handler
	result, err := tool.Handler(ctx, call.Params)
	if err != nil {
		return "", fmt.Errorf("tool %s: %w", call.Name, err)
	}
	return result, nil
}
