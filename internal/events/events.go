package events

import (
	"context"

	"github.com/mohdfareed/gptx-cli/internal/cfg"
	"github.com/mohdfareed/gptx-cli/internal/tools"
)

// Event is the type of event that can be emitted by the model.
type Event[P any] struct {
	Name string // Name of the event
	ch   chan P // Channel to send data
}

// ModelEvents is the manager of events for the model.
type ModelEvents struct {
	Start Event[cfg.Config]
	Error Event[error]
	Done  Event[string]

	Reply     Event[string]
	Reasoning Event[string]

	ToolCall   Event[tools.ToolCall]
	ToolResult Event[string]
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
