// Package events provides pub/sub communication using Go channels and generics.
package events

import (
	"context"

	"github.com/mohdfareed/gptx-cli/internal/cfg"
	"github.com/mohdfareed/gptx-cli/internal/tools"
)

// Event is a typed event channel with a name.
type Event[P any] struct {
	Name string // Event name
	ch   chan P // Data channel
}

// ModelEvents manages all event types for model interactions.
type ModelEvents struct {
	// Events from model
	Start      Event[cfg.Config] // Config loaded
	ToolResult Event[string]     // Tool results
	Error      Event[error]      // Errors

	// Events from client
	Reply     Event[string]          // Model text
	Reasoning Event[string]          // Reasoning steps from model (when available)
	ToolCall  Event[tools.ToolCall]  // Request to execute a tool
	Done      Event[string]          // Emitted when model interaction completes
}

// NewEventsManager creates a new ModelEvent with the given name.
func NewEventsManager() *ModelEvents {
	return &ModelEvents{
		Start: Event[cfg.Config]{Name: "start", ch: make(chan cfg.Config)},
		Error: Event[error]{Name: "error", ch: make(chan error)},
		Done:  Event[string]{Name: "done", ch: make(chan string)},

		Reply:     Event[string]{Name: "reply", ch: make(chan string)},
		Reasoning: Event[string]{Name: "reasoning", ch: make(chan string)},

		ToolCall: Event[tools.ToolCall]{
			Name: "tool-call", ch: make(chan tools.ToolCall),
		},
		ToolResult: Event[string]{Name: "tool-result", ch: make(chan string)},
	}
}

// Subscribe returns a read‑only channel you can range over.
// Callers should cancel the context to stop listening.
func (e *Event[P]) Subscribe(ctx context.Context, handler func(P)) {
	go func() {
		for {
			select {
			case v := <-e.ch:
				handler(v)
			case <-ctx.Done():
				return
			}
		}
	}()
}

// Emit sends data into the event channel (in a non‑blocking goroutine).
func (e *Event[P]) Emit(ctx context.Context, data P) {
	go func() {
		select {
		case e.ch <- data:
		case <-ctx.Done():
		}
	}()
}
