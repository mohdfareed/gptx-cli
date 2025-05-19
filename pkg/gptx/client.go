// Package gptx provides core model interaction logic.
package gptx

import (
	"context"

	"github.com/mohdfareed/gptx-cli/internal/cfg"
	"github.com/mohdfareed/gptx-cli/internal/tools"
)

// Client defines the minimal interface for model API operations.
// This keeps the interface simple while allowing different implementations.
type Client interface {
	// SendRequest sends a single request with messages to the model and returns the response.
	// The client should stream events through callbacks but doesn't need to handle the loop.
	SendRequest(ctx context.Context, request Request) (Response, error)
}

// Request contains all the information needed for a model request.
// This decouples the model from the client implementation.
type Request struct {
	Config      cfg.Config      // Configuration
	Messages    []Message       // Conversation history
	ToolHandler ToolHandler     // Function to handle tool calls
	Callbacks   ModelCallbacks  // Event callbacks
	ToolDefs    []tools.ToolDef // Tool definitions from registry
}

// Message represents a message in the conversation.
type Message struct {
	Role    string // Role of the message sender (user, assistant, tool)
	Content string // Content of the message
	Name    string // Optional name for tool messages
}

// Response contains the model's response data
type Response struct {
	Messages     []Message // Messages from the model's response
	Usage        string    // Usage information as a JSON string
	HasToolCalls bool      // Whether the response contains tool calls
}

// ToolCall represents a tool call from the model
type ToolCall struct {
	Name      string // Name of the tool
	Arguments string // Arguments as a JSON string
}

// ModelCallbacks defines handlers for model interaction events.
type ModelCallbacks struct {
	OnStart     func(cfg.Config) // Called when model starts
	OnReply     func(string)     // Called when model returns text
	OnReasoning func(string)     // Called when model exposes reasoning
	OnWebSearch func()           // Called when model initiates a web search
	OnError     func(error)      // Called when an error occurs
	OnDone      func(string)     // Called when model completes
}

// ToolHandler is a function that handles tool calls.
// It takes a tool name and parameters, and returns a result or an error.
type ToolHandler func(ctx context.Context, name string, params string) (string, error)
