// Package tools provides the tool registry and execution system.
package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
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

// Registry manages tool definitions and execution.
// It provides a unified interface for:
// - Registering built-in tools
// - Registering user-defined tools
// - Looking up tool definitions
// - Executing tools
type Registry struct {
	definitions map[string]ToolDef // Tool definitions indexed by name
	mu          sync.RWMutex       // Protects concurrent access to definitions
}

// NewRegistry creates a new tool registry.
func NewRegistry() *Registry {
	return &Registry{
		definitions: make(map[string]ToolDef),
	}
}

// Register adds a tool to the registry.
// This is the primary method for adding both built-in and user-defined tools.
func (r *Registry) Register(tool ToolDef) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.definitions[tool.Name] = tool
}

// GetDefinitions returns all registered tool definitions.
// This is useful for generating tool descriptions for the model.
func (r *Registry) GetDefinitions() []ToolDef {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tools := make([]ToolDef, 0, len(r.definitions))
	for _, def := range r.definitions {
		tools = append(tools, def)
	}
	return tools
}

// Get returns a specific tool definition by name.
func (r *Registry) Get(name string) (ToolDef, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	def, ok := r.definitions[name]
	return def, ok
}

// Execute runs a tool with the provided parameters.
// It handles parameter parsing, validation, and tool execution.
func (r *Registry) Execute(ctx context.Context, call ToolCall) (string, error) {
	r.mu.RLock()
	tool, ok := r.definitions[call.Name]
	r.mu.RUnlock()

	if !ok {
		return "", fmt.Errorf("unknown tool: %s", call.Name)
	}

	// Special handling for tools that don't need execution
	if tool.Handler == nil {
		return "", nil
	}

	// Parse parameters
	var params map[string]any
	if err := json.Unmarshal([]byte(call.Params), &params); err != nil {
		return "", fmt.Errorf("tool %s: invalid parameters: %w", call.Name, err)
	}

	// Execute the tool
	result, err := tool.Handler(ctx, params)
	if err != nil {
		return "", fmt.Errorf("tool %s: %w", call.Name, err)
	}

	return result, nil
}
