// Package gptx provides core model interaction logic, acting as the controller
// layer between the CLI interface, configuration, events, tools, and API clients.
package gptx

import (
	"context"
	"fmt"

	"github.com/mohdfareed/gptx-cli/internal/cfg"
	"github.com/mohdfareed/gptx-cli/internal/events"
	"github.com/mohdfareed/gptx-cli/internal/tools"
)

// Model handles interactions with AI models.
// It serves as the central component that:
// 1. Manages configuration
// 2. Coordinates tool execution
// 3. Handles event callbacks
// 4. Delegates to the API client
type Model struct {
	client       Client           // API client
	config       cfg.Config       // Configuration
	toolRegistry *tools.Registry  // Tool registry
	callbacks    events.Callbacks // Event callbacks
}

// ModelOption is a function that configures a Model.
// This follows the functional options pattern for clean configuration.
type ModelOption func(*Model)

// NewModel creates a new model with the given configuration and options.
func NewModel(config cfg.Config, options ...ModelOption) *Model {
	// Create a new model with default configuration
	model := &Model{
		config:       config,
		toolRegistry: tools.NewRegistry(),
	}

	// Apply all options
	for _, option := range options {
		option(model)
	}

	return model
}

// WithClient is an option that sets the client for the model.
func WithClient(client Client) ModelOption {
	return func(m *Model) {
		m.client = client
	}
}

// WithCallbacks is an option that sets the event callbacks.
func WithCallbacks(callbacks events.Callbacks) ModelOption {
	return func(m *Model) {
		m.callbacks = callbacks
	}
}

// RegisterTool adds a tool to the model's registry.
// This makes it easy to add custom tools or extensions.
func (m *Model) RegisterTool(tool tools.ToolDef) {
	m.toolRegistry.Register(tool)
}

// Config returns the model's configuration.
func (m *Model) Config() cfg.Config {
	return m.config
}

// Tools returns all registered tool definitions.
func (m *Model) Tools() []tools.ToolDef {
	return m.toolRegistry.GetDefinitions()
}

// Message sends a message to the model and processes the response through callbacks.
func (m *Model) Message(ctx context.Context, prompt string) error {
	if m.client == nil {
		return fmt.Errorf("no client set, use WithClient option")
	}

	// Create the tool handler function
	toolHandler := func(ctx context.Context, name string, params string) (string, error) {
		result, err := m.toolRegistry.Execute(ctx, tools.ToolCall{
			Name: name, Params: params,
		})
		if err != nil && m.callbacks.OnError != nil {
			m.callbacks.OnError(err)
		}
		if m.callbacks.OnToolResult != nil {
			m.callbacks.OnToolResult(result)
		}
		return result, err
	}

	// Create callbacks for the client
	clientCallbacks := ModelCallbacks{
		OnStart:     m.callbacks.OnStart,
		OnReply:     m.callbacks.OnReply,
		OnReasoning: m.callbacks.OnReasoning,
		OnError:     m.callbacks.OnError,
		OnDone:      m.callbacks.OnDone,
	}

	// Send the request to the client
	request := Request{
		Config:      m.config,
		Prompt:      prompt,
		ToolHandler: toolHandler,
		Callbacks:   clientCallbacks,
	}

	// Let the client process the request
	if err := m.client.Generate(ctx, request); err != nil {
		return fmt.Errorf("generate: %w", err)
	}
	return nil
}
