package gptx

import (
	"context"
	"encoding/json"

	"github.com/mohdfareed/gptx-cli/internal/tools"
)

// Message represents a chat message
type Message struct {
	Role    string
	Content string
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
