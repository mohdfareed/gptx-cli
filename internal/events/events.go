package events

import (
	"context"
)

// EventType is the type of event
type EventType string

const (
	Start         EventType = "start"
	ToolCall      EventType = "tool-call"
	ToolResult    EventType = "tool-result"
	Reply         EventType = "reply"
	InternalReply EventType = "internal-reply"
	Done          EventType = "done"
	Error         EventType = "error"
)

// Manager manages events for the model.
type Manager struct {
	Events map[EventType]chan string
}

// New creates a new event manager
func New() *Manager {
	return &Manager{
		Events: make(map[EventType]chan string),
	}
}

// Subscribe returns a read‑only channel you can range over.
// Callers should cancel the context to stop listening.
func (m *Manager) Subscribe(
	ctx context.Context, event EventType, handler func(string),
) {
	go func() {
		for {
			select {
			case v := <-m.Events[event]:
				handler(v)
			case <-ctx.Done():
				return
			}
		}
	}()
}

// Emit sends data into the event channel (in a non‑blocking goroutine).
func (m *Manager) Emit(ctx context.Context, event EventType, data string) {
	go func() {
		select {
		case m.Events[event] <- data:
		case <-ctx.Done():
		}
	}()
}
