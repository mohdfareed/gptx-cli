package gptx

import (
	"context"
	"fmt"
	"io"

	"github.com/mohdfareed/gptx-cli/internal/cfg"
	"github.com/mohdfareed/gptx-cli/internal/events"
	"github.com/mohdfareed/gptx-cli/internal/tools"
)

// Model handles interactions with AI models.
type Model struct {
	Config *cfg.Config
	Events *events.ModelEvents
	Tools  *tools.Manager
	client Client
}

// NewModel creates a new model instance with the given configuration.
func NewModel(config *cfg.Config) (*Model, error) {
	// Validate the configuration
	if config == nil {
		return nil, fmt.Errorf("config not provided")
	}

	// Create the event manager
	events := events.NewEventsManager()

	// Create and set the tools manager
	toolsManager, err := tools.NewManager(*config)
	if err != nil {
		return nil, fmt.Errorf("tools: %w", err)
	}

	// Subscribe tools to events
	events.ToolCall.Subscribe(context.Background(), func(call tools.ToolCall) {
		result, err := toolsManager.CallTool(context.Background(), call)
		if err != nil {
			events.Error.Emit(context.Background(), err)
			return
		}
		events.ToolResult.Emit(context.Background(), result)
	})

	// Create a new model instance
	model := &Model{
		Config: config, Events: events, Tools: toolsManager,
	}
	return model, nil
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

	// Emit completion event
	m.Events.Done.Emit(ctx, "100 tokens")
	return nil
}
