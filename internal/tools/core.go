package tools

import (
	"context"
	"fmt"

	"github.com/mohdfareed/gptx-cli/internal/cfg"
	"github.com/mohdfareed/gptx-cli/internal/events"
)

type Tool func(ctx context.Context, params map[string]any) (string, error)

type ToolDef struct {
	Name    string
	Desc    string
	Params  map[string]any
	Handler Tool
}

// Manager wires up tool calls to handlers and emits results.
type Manager struct {
	tools map[string]ToolDef
}

// NewManager creates a tool manager and registers the requested tools.
func NewManager(config cfg.Config, ev *events.Manager) (*Manager, error) {
	m := &Manager{
		tools: make(map[string]ToolDef),
	}

	// Register requested tools
	for _, name := range config.Tools {
		switch name {
		case ShellToolName:
			m.tools[ShellToolName] = ShellTool(config)
		case RepoToolName:
			m.tools[RepoToolName] = RepoTool(config)
		default:
			return nil, fmt.Errorf("unknown tool: %s", name)
		}
	}
	return m, nil
}

func (m *Manager) CallTool(
	ctx context.Context, name string, params map[string]any,
) (string, error) {
	tool, ok := m.tools[name]
	if !ok {
		return "", fmt.Errorf("unknown tool: %s", name)
	}

	// Call the tool handler
	result, err := tool.Handler(ctx, params)
	if err != nil {
		return "", fmt.Errorf("tool %s: %w", name, err)
	}
	return result, nil
}
