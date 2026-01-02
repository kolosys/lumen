package trace

import "context"

// Carrier is an interface for propagation carriers (e.g., HTTP headers).
type Carrier interface {
	Get(key string) string
	Set(key, value string)
	Keys() []string
}

// MapCarrier is a map-based carrier.
type MapCarrier map[string]string

func (m MapCarrier) Get(key string) string  { return m[key] }
func (m MapCarrier) Set(key, value string)  { m[key] = value }
func (m MapCarrier) Keys() []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Propagator handles trace context injection and extraction.
type Propagator interface {
	Inject(ctx context.Context, carrier Carrier)
	Extract(ctx context.Context, carrier Carrier) context.Context
}

// W3CPropagator implements W3C Trace Context propagation.
type W3CPropagator struct{}

func (p *W3CPropagator) Inject(ctx context.Context, carrier Carrier) {
	span := SpanFromContext(ctx)
	if span == nil || !span.traceID.IsValid() {
		return
	}

	tc := &TraceContext{
		TraceID: span.traceID,
		SpanID:  span.spanID,
	}
	if span.sampled {
		tc.SetSampled(true)
	}

	carrier.Set(W3CTraceparentHeader, tc.FormatW3CTraceparent())
}

func (p *W3CPropagator) Extract(ctx context.Context, carrier Carrier) context.Context {
	traceparent := carrier.Get(W3CTraceparentHeader)
	if traceparent == "" {
		return ctx
	}

	tc, err := ParseW3CTraceparent(traceparent)
	if err != nil {
		return ctx
	}

	tc.TraceState = carrier.Get(W3CTracestateHeader)
	return ContextWithTraceContext(ctx, tc)
}

// HeaderPropagator implements customizable header-based propagation.
type HeaderPropagator struct {
	// TraceIDHeader is the header name for trace ID. Default: "X-Trace-ID"
	TraceIDHeader string
	// SpanIDHeader is the header name for span ID. Default: "X-Span-ID"
	SpanIDHeader string
}

func (p *HeaderPropagator) traceHeader() string {
	if p.TraceIDHeader == "" {
		return TraceIDHeader
	}
	return p.TraceIDHeader
}

func (p *HeaderPropagator) spanHeader() string {
	if p.SpanIDHeader == "" {
		return SpanIDHeader
	}
	return p.SpanIDHeader
}

func (p *HeaderPropagator) Inject(ctx context.Context, carrier Carrier) {
	span := SpanFromContext(ctx)
	if span == nil || !span.traceID.IsValid() {
		return
	}

	carrier.Set(p.traceHeader(), span.traceID.String())
	carrier.Set(p.spanHeader(), span.spanID.String())
}

func (p *HeaderPropagator) Extract(ctx context.Context, carrier Carrier) context.Context {
	traceID := carrier.Get(p.traceHeader())
	if traceID == "" {
		return ctx
	}

	tc, err := ParseHeaders(traceID, carrier.Get(p.spanHeader()))
	if err != nil {
		return ctx
	}

	return ContextWithTraceContext(ctx, tc)
}

// CompositePropagator combines multiple propagators.
type CompositePropagator struct {
	propagators []Propagator
}

// NewCompositePropagator creates a composite propagator.
func NewCompositePropagator(propagators ...Propagator) *CompositePropagator {
	return &CompositePropagator{propagators: propagators}
}

func (p *CompositePropagator) Inject(ctx context.Context, carrier Carrier) {
	for _, prop := range p.propagators {
		prop.Inject(ctx, carrier)
	}
}

func (p *CompositePropagator) Extract(ctx context.Context, carrier Carrier) context.Context {
	for _, prop := range p.propagators {
		ctx = prop.Extract(ctx, carrier)
	}
	return ctx
}

// DefaultPropagator returns a propagator supporting both W3C and header formats.
func DefaultPropagator() Propagator {
	return NewCompositePropagator(&W3CPropagator{}, &HeaderPropagator{})
}
