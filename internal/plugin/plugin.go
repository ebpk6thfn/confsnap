// Package plugin provides a simple hook-based plugin system for confsnap,
// allowing external or internal extensions to react to snapshot events.
package plugin

import (
	"fmt"
	"sync"
)

// EventType represents the kind of event emitted during a snapshot run.
type EventType string

const (
	EventSnapshotTaken  EventType = "snapshot_taken"
	EventDiffDetected   EventType = "diff_detected"
	EventRunComplete    EventType = "run_complete"
)

// Event carries contextual data for a plugin hook invocation.
type Event struct {
	Type    EventType
	Host    string
	Path    string
	Meta    map[string]string
}

// Handler is a function that processes a plugin event.
type Handler func(e Event) error

// Registry holds named handlers registered for each event type.
type Registry struct {
	mu       sync.RWMutex
	handlers map[EventType][]namedHandler
}

type namedHandler struct {
	name    string
	handle  Handler
}

// New returns an initialised Registry.
func New() *Registry {
	return &Registry{
		handlers: make(map[EventType][]namedHandler),
	}
}

// Register adds a named handler for the given event type.
// Returns an error if a handler with the same name is already registered
// for that event type.
func (r *Registry) Register(event EventType, name string, h Handler) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, nh := range r.handlers[event] {
		if nh.name == name {
			return fmt.Errorf("plugin: handler %q already registered for event %q", name, event)
		}
	}
	r.handlers[event] = append(r.handlers[event], namedHandler{name: name, handle: h})
	return nil
}

// Emit dispatches the event to all registered handlers in registration order.
// All handlers are called; errors are collected and returned as a combined error.
func (r *Registry) Emit(e Event) error {
	r.mu.RLock()
	handlers := make([]namedHandler, len(r.handlers[e.Type]))
	copy(handlers, r.handlers[e.Type])
	r.mu.RUnlock()

	var errs []error
	for _, nh := range handlers {
		if err := nh.handle(e); err != nil {
			errs = append(errs, fmt.Errorf("plugin %q: %w", nh.name, err))
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("plugin emit errors: %v", errs)
}

// Names returns the registered handler names for the given event type.
func (r *Registry) Names(event EventType) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var names []string
	for _, nh := range r.handlers[event] {
		names = append(names, nh.name)
	}
	return names
}
