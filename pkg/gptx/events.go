package gptx

import (
	"context"
)

// EventType identifies different types of model events.
type EventType string

const (
	// Core event types
	EventStart    EventType = "start"    // Model interaction started
	EventReply    EventType = "reply"    // Model produced a text response
	EventTool     EventType = "tool"     // Tool used or results available
	EventReason   EventType = "reason"   // Model reasoned a response
	EventComplete EventType = "complete" // Model interaction completed
	EventError    EventType = "error"    // Error occurred
)

// Payload represents an event payload in the model interaction lifecycle.
type Payload struct {
	Type EventType // Type of event
	Data any       // Event-specific data
}

// Event provides a simple system for event handling.
type Event struct {
	Channel chan Payload // Channel for events
}

// NewEvent creates a new event instance.
func NewEvent() *Event {
	return &Event{
		Channel: make(chan Payload, 10), // Buffer to prevent blocking
	}
}

// Emit sends an event through the channel.
// If the channel is full or closed, the event is dropped.
func (e *Event) Emit(ctx context.Context, typ EventType, payload any) {
	select {
	case e.Channel <- Payload{Type: typ, Data: payload}:
		// Successfully sent
	case <-ctx.Done():
		// Context cancelled
	default:
		// Channel full or closed, dropping event
	}
}

// Close closes the event channel.
func (e *Event) Close() {
	close(e.Channel)
}
