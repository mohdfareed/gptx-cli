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
	// Generate processes a prompt and returns a model response.
	// The client should call the appropriate handlers in ModelCallbacks
	// as events occur during generation.
	Generate(ctx context.Context, request Request) error
}

// Request contains all the information needed for a model request.
// This decouples the model from the client implementation.
type Request struct {
	Config      cfg.Config      // Configuration
	Prompt      string          // User input
	ToolHandler ToolHandler     // Function to handle tool calls
	Callbacks   ModelCallbacks  // Event callbacks
	ToolDefs    []tools.ToolDef // Tool definitions from registry
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
