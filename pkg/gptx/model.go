package gptx

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/mohdfareed/gptx-cli/internal/cfg"
	"github.com/mohdfareed/gptx-cli/internal/events"
	"github.com/mohdfareed/gptx-cli/internal/files"
	"github.com/mohdfareed/gptx-cli/internal/tools"
)

// Message represents a chat message
type Message struct {
	Role    string // "user", "assistant", "system"
	Content string // Message content
}

// Request represents a request to the model
type Request struct {
	Model        string          // Model name to use
	SystemPrompt string          // System prompt/instructions
	Messages     []Message       // Conversation history
	Files        []string        // Files to attach
	Tools        []tools.ToolDef // Tools to enable
	Temperature  float32         // Temperature (0-1)
	MaxTokens    int             // Max tokens to generate
	User         string          // End-user identifier
}

// ResponseStream represents a stream of model responses
type ResponseStream interface {
	// Close closes the stream
	Close()
	// HasNext returns true if there are more events
	HasNext() bool
	// Next returns the next event
	Next() ResponseEvent
	// Err returns any error that occurred during streaming
	Err() error
	// SubmitToolOutputs submits tool outputs to the model
	SubmitToolOutputs(ctx context.Context, outputs []ToolOutput) error
}

// ResponseEvent represents a streaming response event
type ResponseEvent interface {
	// GetType returns the event type
	GetType() string
	// GetContent returns the text content if available
	GetContent() string
	// GetToolCall returns tool call info if available
	GetToolCall() *ToolCall
}

// ToolCall represents a tool call from the model
type ToolCall struct {
	ID        string          // Tool call ID
	Name      string          // Tool name
	Arguments json.RawMessage // Tool arguments
}

// ToolOutput represents a tool output to be submitted
type ToolOutput struct {
	ID     string // Tool call ID
	Name   string // Tool name
	Output string // Tool output
}

// ToolHandler is a function that handles tool calls
type ToolHandler func(ctx context.Context, toolName string, args json.RawMessage) (json.RawMessage, error)

// Client defines an interface for model API operations
type Client interface {
	// Generate sends a request to the model and returns a stream of events
	Generate(ctx context.Context, req *Request) (ResponseStream, error)

	// ProcessStream processes a response stream, handling text and tool calls
	ProcessStream(
		ctx context.Context,
		stream ResponseStream,
		textHandler func(string),
		toolHandlers map[string]ToolHandler,
	) error
}

// Model handles interactions with AI models.
type Model struct {
	Config *cfg.Config
	Events *events.Manager
	Tools  *tools.Manager
	client Client
}

// NewModel creates a new model instance with the given configuration.
func NewModel(config *cfg.Config) (*Model, error) {
	if config == nil {
		config = &cfg.Config{}
	}

	// Create the event manager
	eventManager := events.New()

	// Create and set the tools manager
	toolsManager, err := tools.NewManager(*config, eventManager)
	if err != nil {
		return nil, fmt.Errorf("tools manager: %w", err)
	}

	// Create a new model instance
	model := &Model{
		Config: config,
		Events: eventManager,
		Tools:  toolsManager,
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

	// Process prompt tags if present
	tagPrefix := "@file" // Default tag prefix
	processedPrompt, err := files.ProcessTags(prompt, tagPrefix)
	if err != nil {
		return fmt.Errorf("process prompt: %w", err)
	}

	// Create context with cancellation
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Emit start event
	m.Events.EventStart.Emit(ctx, struct{}{})

	// Create request with messages
	req := &Request{
		Model:        m.Config.Model,
		SystemPrompt: m.Config.SysPrompt,
		Messages: []Message{
			{Role: "user", Content: processedPrompt},
		},
		Temperature: float32(m.Config.Temp) / 100.0,
	}

	// Add user identifier if present
	if m.Config.UserID != "" {
		req.User = m.Config.UserID
	}

	// Add max tokens if specified
	if m.Config.Tokens != nil {
		req.MaxTokens = *m.Config.Tokens
	}

	// Add files from the processed tags and config
	allFiles := append(attachments, m.Config.Files...)
	if len(allFiles) > 0 {
		req.Files = allFiles
	}

	// Add tools from the tool manager if tools are enabled
	if len(m.Config.Tools) > 0 {
		req.Tools = m.Tools.GetToolDefinitions()
	}

	// Send request to model via client
	stream, err := m.client.Generate(ctx, req)
	if err != nil {
		m.Events.EventError.Emit(ctx, fmt.Errorf("generate: %w", err))
		return fmt.Errorf("generate: %w", err)
	}
	defer stream.Close()

	// Set up handlers for stream events
	textHandler := func(text string) {
		// Emit reply event
		m.Events.EventReply.Emit(ctx, text)

		// Write to output
		fmt.Fprint(output, text)
	}

	// Process the stream
	err = m.client.ProcessStream(ctx, stream, textHandler, m.handleTools(ctx))
	if err != nil {
		m.Events.EventError.Emit(ctx, fmt.Errorf("process stream: %w", err))
		return fmt.Errorf("process stream: %w", err)
	}

	// Emit completion event
	m.Events.EventDone.Emit(ctx, struct{}{})
	return nil
}

// handleTools returns a map of tool handlers for the client
func (m *Model) handleTools(ctx context.Context) map[string]ToolHandler {
	handlers := make(map[string]ToolHandler)

	for _, toolName := range m.Config.Tools {
		// Create a closure that delegates to the tools manager
		// Capture toolName to avoid loop variable issues
		name := toolName
		handlers[name] = func(ctx context.Context, toolName string, args json.RawMessage) (json.RawMessage, error) {
			// Emit tool call event
			m.Events.EventToolCall.Emit(ctx, string(args))

			// Execute tool using the tools manager
			result, err := m.Tools.HandleToolCall(ctx, name, args)
			if err != nil {
				return nil, err
			}

			// Emit tool result event
			m.Events.EventToolResult.Emit(ctx, string(result))

			return result, nil
		}
	}

	return handlers
}
