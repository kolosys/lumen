package trace

import (
	"encoding/json"
	"io"
	"sync"
)

// Exporter receives completed spans.
type Exporter interface {
	Export(span *Span)
	Close() error
}

// NopExporter discards all spans.
type NopExporter struct{}

func (NopExporter) Export(*Span) {}
func (NopExporter) Close() error { return nil }

// WriterExporter writes spans as JSON to an io.Writer.
type WriterExporter struct {
	writer io.Writer
	mu     sync.Mutex
}

// NewWriterExporter creates an exporter that writes to w.
func NewWriterExporter(w io.Writer) *WriterExporter {
	return &WriterExporter{writer: w}
}

type spanData struct {
	TraceID    string      `json:"trace_id"`
	SpanID     string      `json:"span_id"`
	ParentID   string      `json:"parent_id,omitempty"`
	Name       string      `json:"name"`
	StartTime  int64       `json:"start_time_ns"`
	EndTime    int64       `json:"end_time_ns"`
	Duration   int64       `json:"duration_ns"`
	Status     string      `json:"status"`
	StatusMsg  string      `json:"status_message,omitempty"`
	Attributes []attrData  `json:"attributes,omitempty"`
	Events     []eventData `json:"events,omitempty"`
}

type attrData struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

type eventData struct {
	Name       string     `json:"name"`
	Timestamp  int64      `json:"timestamp_ns"`
	Attributes []attrData `json:"attributes,omitempty"`
}

func (e *WriterExporter) Export(span *Span) {
	data := spanData{
		TraceID:   span.traceID.String(),
		SpanID:    span.spanID.String(),
		Name:      span.name,
		StartTime: span.startTime.UnixNano(),
		EndTime:   span.endTime.UnixNano(),
		Duration:  span.Duration().Nanoseconds(),
		Status:    span.status.String(),
		StatusMsg: span.statusMsg,
	}

	if span.parentID.IsValid() {
		data.ParentID = span.parentID.String()
	}

	for _, attr := range span.attributes {
		data.Attributes = append(data.Attributes, attrData{
			Key:   attr.Key,
			Value: attr.Value,
		})
	}

	for _, event := range span.events {
		ev := eventData{
			Name:      event.Name,
			Timestamp: event.Timestamp.UnixNano(),
		}
		for _, attr := range event.Attributes {
			ev.Attributes = append(ev.Attributes, attrData{
				Key:   attr.Key,
				Value: attr.Value,
			})
		}
		data.Events = append(data.Events, ev)
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	encoded, _ := json.Marshal(data)
	e.writer.Write(encoded)
	e.writer.Write([]byte("\n"))
}

func (e *WriterExporter) Close() error { return nil }

// InMemoryExporter collects spans in memory for testing.
type InMemoryExporter struct {
	spans []*spanData
	mu    sync.Mutex
}

// NewInMemoryExporter creates an in-memory exporter.
func NewInMemoryExporter() *InMemoryExporter {
	return &InMemoryExporter{
		spans: make([]*spanData, 0),
	}
}

func (e *InMemoryExporter) Export(span *Span) {
	data := &spanData{
		TraceID:   span.traceID.String(),
		SpanID:    span.spanID.String(),
		Name:      span.name,
		StartTime: span.startTime.UnixNano(),
		EndTime:   span.endTime.UnixNano(),
		Duration:  span.Duration().Nanoseconds(),
		Status:    span.status.String(),
	}

	if span.parentID.IsValid() {
		data.ParentID = span.parentID.String()
	}

	e.mu.Lock()
	e.spans = append(e.spans, data)
	e.mu.Unlock()
}

// Spans returns collected spans.
func (e *InMemoryExporter) Spans() []*spanData {
	e.mu.Lock()
	defer e.mu.Unlock()
	result := make([]*spanData, len(e.spans))
	copy(result, e.spans)
	return result
}

// Len returns the number of collected spans.
func (e *InMemoryExporter) Len() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return len(e.spans)
}

// Clear removes all collected spans.
func (e *InMemoryExporter) Clear() {
	e.mu.Lock()
	e.spans = e.spans[:0]
	e.mu.Unlock()
}

func (e *InMemoryExporter) Close() error { return nil }
