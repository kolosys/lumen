// Package trace provides distributed tracing with W3C and Kolosys format support.
package trace

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"sync"
	"sync/atomic"
	"time"
)

// TraceID is a 16-byte trace identifier.
type TraceID [16]byte

func (t TraceID) String() string {
	return hex.EncodeToString(t[:])
}

func (t TraceID) IsValid() bool {
	return t != TraceID{}
}

// SpanID is an 8-byte span identifier.
type SpanID [8]byte

func (s SpanID) String() string {
	return hex.EncodeToString(s[:])
}

func (s SpanID) IsValid() bool {
	return s != SpanID{}
}

// Tracer creates and manages spans.
type Tracer struct {
	opts      *Options
	spanPool  *sync.Pool
	closed    atomic.Bool
	asyncCh   chan *Span
	asyncWg   sync.WaitGroup
	closeOnce sync.Once
}

// New creates a new Tracer.
func New(opts *Options) *Tracer {
	if opts == nil {
		opts = &Options{}
	}
	opts.applyDefaults()

	t := &Tracer{
		opts: opts,
		spanPool: &sync.Pool{
			New: func() any {
				return &Span{
					attributes: make([]Attribute, 0, 8),
					events:     make([]Event, 0, 4),
				}
			},
		},
	}

	if opts.AsyncExport {
		t.asyncCh = make(chan *Span, opts.AsyncBufferSize)
		t.asyncWg.Add(1)
		go t.asyncWorker()
	}

	return t
}

// Start creates a new span.
func (t *Tracer) Start(ctx context.Context, name string, opts ...SpanOption) (context.Context, *Span) {
	if t.closed.Load() {
		return ctx, &Span{noop: true}
	}

	parent := SpanFromContext(ctx)
	tc := TraceContextFromContext(ctx)

	span := t.getSpan()
	span.tracer = t
	span.name = name
	span.startTime = time.Now()

	if parent != nil && parent.traceID.IsValid() {
		span.traceID = parent.traceID
		span.parentID = parent.spanID
	} else if tc != nil && tc.TraceID.IsValid() {
		span.traceID = tc.TraceID
		span.parentID = tc.SpanID
	} else {
		span.traceID = generateTraceID()
	}
	span.spanID = generateSpanID()

	for _, opt := range opts {
		opt(span)
	}

	if !t.opts.Sampler.ShouldSample(SamplingParams{
		TraceID:  span.traceID,
		Name:     name,
		ParentID: span.parentID,
	}) {
		span.sampled = false
	} else {
		span.sampled = true
	}

	return ContextWithSpan(ctx, span), span
}

// Close shuts down the tracer.
func (t *Tracer) Close() error {
	t.closeOnce.Do(func() {
		t.closed.Store(true)
		if t.asyncCh != nil {
			close(t.asyncCh)
			t.asyncWg.Wait()
		}
	})
	return nil
}

func (t *Tracer) getSpan() *Span {
	s := t.spanPool.Get().(*Span)
	s.reset()
	return s
}

func (t *Tracer) releaseSpan(s *Span) {
	s.reset()
	t.spanPool.Put(s)
}

func (t *Tracer) asyncWorker() {
	defer t.asyncWg.Done()
	for span := range t.asyncCh {
		t.opts.Exporter.Export(span)
		t.releaseSpan(span)
	}
}

func (t *Tracer) exportSpan(span *Span) {
	if t.asyncCh != nil && !t.closed.Load() {
		select {
		case t.asyncCh <- span:
			return
		default:
		}
	}
	t.opts.Exporter.Export(span)
	t.releaseSpan(span)
}

func generateTraceID() TraceID {
	var id TraceID
	rand.Read(id[:])
	return id
}

func generateSpanID() SpanID {
	var id SpanID
	rand.Read(id[:])
	return id
}

// Default tracer
var defaultTracer = New(nil)

// SetDefault sets the default tracer.
func SetDefault(t *Tracer) {
	defaultTracer = t
}

// Default returns the default tracer.
func Default() *Tracer {
	return defaultTracer
}

// Start creates a span using the default tracer.
func Start(ctx context.Context, name string, opts ...SpanOption) (context.Context, *Span) {
	return defaultTracer.Start(ctx, name, opts...)
}
