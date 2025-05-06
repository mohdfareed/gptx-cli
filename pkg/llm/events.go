package llm

import (
	"reflect"
	"slices"
	"sync"
)

// EventHandler is any function that takes a payload.
type EventHandler[P any] func(payload P)

// EventPayload is an event payload.
type EventPayload any

// Event is a named event with a payload.
type Event[P EventPayload] struct {
	subs []EventHandler[P]
	mu   sync.RWMutex
}

// MARK: Invocation
// ============================================================================

// Publish invokes all handlers for an event, passing payload.
func (e *Event[P]) publish(payload P) *sync.WaitGroup {
	e.mu.RLock()
	defer e.mu.RUnlock()
	var wg sync.WaitGroup
	for _, h := range e.subs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			h(payload)
		}()
	}
	return &wg
}

// Wait blocks asynchronously until the event is published.
func (e *Event[P]) Wait() {
	ch := make(chan struct{})
	handler := func(_ P) {
		close(ch)
	}
	e.Subscribe(handler)
	<-ch
	e.Unsubscribe(handler)
}

// MARK: Subscription
// ============================================================================

// Subscribe registers sub to the event.
func (e *Event[P]) Subscribe(sub EventHandler[P]) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.subs = append(e.subs, sub)
}

// Unsubscribe removes sub from the event.
func (e *Event[P]) Unsubscribe(sub EventHandler[P]) {
	e.mu.Lock()
	defer e.mu.Unlock()
	subs := e.subs

	for i, fn := range subs {
		if reflect.ValueOf(fn).Pointer() == reflect.ValueOf(sub).Pointer() {
			e.subs = slices.Delete(e.subs, i, i+1)
			break
		}
	}
}

// MARK: Payloads
// ============================================================================

type P1[A1 any] struct {
	A A1 `json:"a1"`
}
type P2[A1 any, A2 any] struct {
	A A1 `json:"a1"`
	B A2 `json:"a2"`
}
type P3[A1 any, A2 any, A3 any] struct {
	A A1 `json:"a1"`
	B A2 `json:"a2"`
	C A3 `json:"a3"`
}
