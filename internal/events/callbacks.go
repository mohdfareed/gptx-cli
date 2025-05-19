// Package events provides event handling for model interactions.
package events

import (
	"github.com/mohdfareed/gptx-cli/internal/cfg"
	"github.com/mohdfareed/gptx-cli/internal/tools"
)

// Callbacks defines handlers for various model events.
// Each field is a function that will be called when the corresponding event occurs.
// This simpler approach replaces the pub/sub channel-based system.
type Callbacks struct {
	// Model lifecycle events
	OnStart func(cfg.Config) // When model starts
	OnError func(error)      // Error handling
	OnDone  func(string)     // When model completes (with usage stats)

	// Model output events
	OnReply     func(string) // Text from the model
	OnReasoning func(string) // Reasoning steps (when available)

	// Tool-related events
	OnToolCall   func(tools.ToolCall) // Request to execute a tool
	OnToolResult func(string)         // Results from tool execution
}

// Builder provides a fluent API for constructing callbacks.
// This pattern makes it easy to register only the callbacks you need.
type Builder struct {
	callbacks Callbacks
}

// NewCallbacks creates a new event callbacks builder.
func NewCallbacks() *Builder {
	return &Builder{
		callbacks: Callbacks{
			// Default no-op handlers
			OnStart:      func(cfg.Config) {},
			OnError:      func(error) {},
			OnDone:       func(string) {},
			OnReply:      func(string) {},
			OnReasoning:  func(string) {},
			OnToolCall:   func(tools.ToolCall) {},
			OnToolResult: func(string) {},
		},
	}
}

// Build returns the configured callbacks.
func (b *Builder) Build() Callbacks {
	return b.callbacks
}

// WithStartHandler sets the handler for the start event.
func (b *Builder) WithStartHandler(handler func(cfg.Config)) *Builder {
	b.callbacks.OnStart = handler
	return b
}

// WithErrorHandler sets the handler for errors.
func (b *Builder) WithErrorHandler(handler func(error)) *Builder {
	b.callbacks.OnError = handler
	return b
}

// WithDoneHandler sets the handler for the done event.
func (b *Builder) WithDoneHandler(handler func(string)) *Builder {
	b.callbacks.OnDone = handler
	return b
}

// WithReplyHandler sets the handler for model replies.
func (b *Builder) WithReplyHandler(handler func(string)) *Builder {
	b.callbacks.OnReply = handler
	return b
}

// WithReasoningHandler sets the handler for reasoning events.
func (b *Builder) WithReasoningHandler(handler func(string)) *Builder {
	b.callbacks.OnReasoning = handler
	return b
}

// WithToolCallHandler sets the handler for tool calls.
func (b *Builder) WithToolCallHandler(handler func(tools.ToolCall)) *Builder {
	b.callbacks.OnToolCall = handler
	return b
}

// WithToolResultHandler sets the handler for tool results.
func (b *Builder) WithToolResultHandler(handler func(string)) *Builder {
	b.callbacks.OnToolResult = handler
	return b
}
