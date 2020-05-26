package monitoring

import (
	"time"
)

// Events can be used to report events to the underlying event collector.
var Events EventReporter

// EventReporter is an interface that defines the methods expected to be implemented by event reporting clients.
type EventReporter interface {
	// Emit sends the provided Event.
	Emit(e *Event)

	// EmitSimple sends an event with the provided title and description.
	EmitSimple(title, description string)
}

// Event encapsulates the data representing a single event.
type Event struct {
	// Title of the event.
	Title string
	// Description of the event.
	Description string
	// Timestamp is the time the event occurred.
	Timestamp time.Time
	// AggregationKey groups this event with others of the same key.
	AggregationKey string
	// Priority of the event.
	Priority EventPriority
	// Source is the name of the event source.
	Source string
	// EventType is the type or severity level of the event.
	Type EventType
	// Tags for the event.
	Tags []Tag
}

// EventPriority is the event priority for events.
type EventPriority int

const (
	// Low is the "low" priority for events.
	Low EventPriority = iota
	// Normal is the "normal" priority for events.
	Normal
)

// EventType is the type or severity level of the event.
type EventType int

const (
	// Info indicates the event is informational in nature.
	Info EventType = iota
	// Warning indicates the event represents a warning.
	Warning
	// Error indicates the event represents an error.
	Error
	// Success indicates the event represents a successful operation.
	Success
)
