// Package gptx provides core model interaction logic, acting as the controller
// layer between the CLI interface, configuration, events, tools, and API clients.
package gptx

import (
	"context"
	"fmt"
	"io"

	"github.com/mohdfareed/gptx-cli/internal/cfg"
	"github.com/mohdfareed/gptx-cli/internal/events"
	"github.com/mohdfareed/gptx-cli/internal/tools"
)

// Client defines a minimal interface for model API operations.
type Client interface {
	// Generate starts a conversation with the model and emits events.
	Generate(ctx context.Context, config Model, prompt string) error
}

// Model handles interactions with AI models.
type Model struct {
	Config cfg.Config           // Configuration
	Events *events.ModelEvents  // Event system
	Tools  *tools.Tools         // Tool manager
	client Client               // API client
}

// NewModel creates a new model with the given configuration.
// Sets up events and tools with automatic tool call handling.
func NewModel(config cfg.Config) *Model {
	// Create the event manager
	events := events.NewEventsManager()

	// Create tools manager
	toolsManager := tools.NewTools(config)

	// Wire up tool execution flow via events
	events.ToolCall.Subscribe(context.Background(), func(call tools.ToolCall) {
		result, err := toolsManager.CallTool(context.Background(), call)
		if err != nil {
			events.Error.Emit(context.Background(), err)
			return
		}
		events.ToolResult.Emit(context.Background(), result)
	})

	// Create a new model instance with all components wired together
	model := &Model{
		Config: config, Events: events, Tools: toolsManager,
	}
	return model
}

// WithClient sets a custom client for the model and returns the model.
func (m *Model) WithClient(client Client) *Model {
	m.client = client
	return m
}

// Message sends a message to the model and streams the response.
func (m *Model) Message(ctx context.Context, prompt string, output io.Writer) error {
	if m.client == nil {
		return fmt.Errorf("no client set, use WithClient to set a client")
	}

	err := m.client.Generate(ctx, *m, prompt)
	if err != nil {
		return fmt.Errorf("generate: %w", err)
	}
	return nil
}
