// Package tools manages model tool capabilities like shell commands and web search.
package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mohdfareed/gptx-cli/internal/cfg"
)

// Tool is a function that executes with parameters and returns a result.
type Tool func(ctx context.Context, params map[string]any) (string, error)

// ToolDef defines a tool's metadata and implementation.
type ToolDef struct {
	Name    string         // Tool identifier
	Desc    string         // Description
	Params  map[string]any // Parameter schema
	Handler Tool           // Implementation
}

// ToolCall represents a request from the model to use a tool.
type ToolCall struct {
	Name   string // Tool name
	Params string // JSON parameters
}

// WebSearchToolDef is the identifier for the web search tool.
// This tool is handled directly by the OpenAI API client rather than locally.
const WebSearchToolDef = "web-search"

// Tools manages the registration and execution of tools available to the model.
// It serves as a registry for available tools and dispatches tool calls to
// the appropriate handler based on the tool name.
type Tools struct {
	tools map[string]ToolDef // Registry of available tools indexed by name
}

// NewTools creates a tool manager and registers tools based on the configuration.
// It dynamically enables tools based on the provided configuration, allowing
// different tool sets to be available in different sessions or environments.
//
// Currently supported tools:
// - Shell commands (when config.Shell is set)
// - Web search (when config.WebSearch is set, handled by OpenAI)
func NewTools(config cfg.Config) *Tools {
	m := &Tools{
		tools: make(map[string]ToolDef),
	}

	// Register shell tool if configured
	if config.Shell != nil {
		m.tools[ShellToolDef] = NewShellTool(config)
	}

	// Register web search tool if configured
	// Note: Web search is actually handled by the OpenAI client,
	// but we register it here to make it visible to the model
	if config.WebSearch != nil {
		m.tools[WebSearchToolDef] = ToolDef{}
	}

	return m
}

// CallTool executes a tool based on the provided ToolCall.
// It handles parameter parsing, tool lookup, execution, and error formatting.
//
// The flow is:
// 1. Look up the requested tool by name
// 2. Skip execution for special tools (e.g., web search handled by OpenAI)
// 3. Parse the JSON parameters
// 4. Execute the tool with the parsed parameters
// 5. Return the result or a formatted error
func (m *Tools) CallTool(ctx context.Context, call ToolCall) (string, error) {
	// Look up the tool definition
	tool, ok := m.tools[call.Name]
	if !ok {
		return "", fmt.Errorf("unknown tool: %s", call.Name)
	}

	// Special handling for web search (handled by OpenAI client)
	if call.Name == WebSearchToolDef {
		return "", nil // Skip web search tool (handled by client)
	}

	// Parse the JSON parameters
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
