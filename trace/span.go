package trace

import (
	"sync"
	"sync/atomic"
	"time"
)

// SpanStatus represents span completion status.
type SpanStatus int

const (
	StatusUnset SpanStatus = iota
	StatusOK
	StatusError
)

func (s SpanStatus) String() string {
	switch s {
	case StatusOK:
		return "ok"
	case StatusError:
		return "error"
	default:
		return "unset"
	}
}

// Attribute is a key-value pair attached to a span.
type Attribute struct {
	Key   string
	Value any
}

// Event is a timestamped annotation.
type Event struct {
	Name       string
	Timestamp  time.Time
	Attributes []Attribute
}

// Span represents a unit of work.
type Span struct {
	tracer     *Tracer
	traceID    TraceID
	spanID     SpanID
	parentID   SpanID
	name       string
	startTime  time.Time
	endTime    time.Time
	status     SpanStatus
	statusMsg  string
	attributes []Attribute
	events     []Event
	sampled    bool
	noop       bool
	ended      atomic.Bool
	mu         sync.Mutex
}

// SpanOption configures span creation.
type SpanOption func(*Span)

// WithAttributes sets initial attributes.
func WithAttributes(attrs ...Attribute) SpanOption {
	return func(s *Span) {
		s.attributes = append(s.attributes, attrs...)
	}
}

// WithStartTime sets a custom start time.
func WithStartTime(t time.Time) SpanOption {
	return func(s *Span) {
		s.startTime = t
	}
}

// TraceID returns the trace ID.
func (s *Span) TraceID() TraceID {
	return s.traceID
}

// SpanID returns the span ID.
func (s *Span) SpanID() SpanID {
	return s.spanID
}

// ParentID returns the parent span ID.
func (s *Span) ParentID() SpanID {
	return s.parentID
}

// Name returns the span name.
func (s *Span) Name() string {
	return s.name
}

// StartTime returns the start time.
func (s *Span) StartTime() time.Time {
	return s.startTime
}

// EndTime returns the end time.
func (s *Span) EndTime() time.Time {
	return s.endTime
}

// Status returns the span status.
func (s *Span) Status() SpanStatus {
	return s.status
}

// StatusMessage returns the status message.
func (s *Span) StatusMessage() string {
	return s.statusMsg
}

// IsSampled returns whether the span is sampled.
func (s *Span) IsSampled() bool {
	return s.sampled
}

// SetAttribute adds an attribute.
func (s *Span) SetAttribute(key string, value any) {
	if s.noop || s.ended.Load() {
		return
	}
	s.mu.Lock()
	s.attributes = append(s.attributes, Attribute{Key: key, Value: value})
	s.mu.Unlock()
}

// SetAttributes adds multiple attributes.
func (s *Span) SetAttributes(attrs ...Attribute) {
	if s.noop || s.ended.Load() {
		return
	}
	s.mu.Lock()
	s.attributes = append(s.attributes, attrs...)
	s.mu.Unlock()
}

// AddEvent adds a timestamped event.
func (s *Span) AddEvent(name string, attrs ...Attribute) {
	if s.noop || s.ended.Load() {
		return
	}
	s.mu.Lock()
	s.events = append(s.events, Event{
		Name:       name,
		Timestamp:  time.Now(),
		Attributes: attrs,
	})
	s.mu.Unlock()
}

// SetStatus sets the span status.
func (s *Span) SetStatus(status SpanStatus, msg string) {
	if s.noop || s.ended.Load() {
		return
	}
	s.mu.Lock()
	s.status = status
	s.statusMsg = msg
	s.mu.Unlock()
}

// RecordError records an error as an event and sets error status.
func (s *Span) RecordError(err error) {
	if err == nil || s.noop || s.ended.Load() {
		return
	}
	s.AddEvent("exception", Attribute{Key: "exception.message", Value: err.Error()})
	s.SetStatus(StatusError, err.Error())
}

// End completes the span.
func (s *Span) End() {
	if s.noop || !s.ended.CompareAndSwap(false, true) {
		return
	}
	s.endTime = time.Now()

	if !s.sampled || s.tracer == nil {
		return
	}

	s.tracer.exportSpan(s)
}

// EndFunc returns a function suitable for defer.
func (s *Span) EndFunc(errPtr *error) func() {
	return func() {
		if errPtr != nil && *errPtr != nil {
			s.RecordError(*errPtr)
		}
		s.End()
	}
}

// Duration returns the span duration.
func (s *Span) Duration() time.Duration {
	if s.endTime.IsZero() {
		return time.Since(s.startTime)
	}
	return s.endTime.Sub(s.startTime)
}

// Attributes returns a copy of span attributes.
func (s *Span) Attributes() []Attribute {
	s.mu.Lock()
	defer s.mu.Unlock()
	attrs := make([]Attribute, len(s.attributes))
	copy(attrs, s.attributes)
	return attrs
}

// Events returns a copy of span events.
func (s *Span) Events() []Event {
	s.mu.Lock()
	defer s.mu.Unlock()
	events := make([]Event, len(s.events))
	copy(events, s.events)
	return events
}

func (s *Span) reset() {
	s.tracer = nil
	s.traceID = TraceID{}
	s.spanID = SpanID{}
	s.parentID = SpanID{}
	s.name = ""
	s.startTime = time.Time{}
	s.endTime = time.Time{}
	s.status = StatusUnset
	s.statusMsg = ""
	s.attributes = s.attributes[:0]
	s.events = s.events[:0]
	s.sampled = false
	s.noop = false
	s.ended.Store(false)
}
