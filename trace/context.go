package trace

import (
	"context"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

type spanContextKey struct{}
type traceContextKey struct{}

// ContextWithSpan returns a context with the span attached.
func ContextWithSpan(ctx context.Context, span *Span) context.Context {
	return context.WithValue(ctx, spanContextKey{}, span)
}

// SpanFromContext retrieves a span from context.
func SpanFromContext(ctx context.Context) *Span {
	if ctx == nil {
		return nil
	}
	if span, ok := ctx.Value(spanContextKey{}).(*Span); ok {
		return span
	}
	return nil
}

// TraceContext holds W3C Trace Context data.
type TraceContext struct {
	TraceID    TraceID
	SpanID     SpanID
	TraceFlags byte
	TraceState string
}

// W3C Trace Context header names.
const (
	W3CTraceparentHeader = "traceparent"
	W3CTracestateHeader  = "tracestate"
)

// Simple header names for trace propagation.
const (
	TraceIDHeader = "X-Trace-ID"
	SpanIDHeader  = "X-Span-ID"
)

var traceparentRegex = regexp.MustCompile(`^([0-9a-f]{2})-([0-9a-f]{32})-([0-9a-f]{16})-([0-9a-f]{2})$`)

// ParseW3CTraceparent parses a W3C traceparent header.
func ParseW3CTraceparent(header string) (*TraceContext, error) {
	header = strings.TrimSpace(strings.ToLower(header))
	matches := traceparentRegex.FindStringSubmatch(header)
	if matches == nil {
		return nil, ErrInvalidContext
	}

	version := matches[1]
	if version != "00" {
		return nil, ErrInvalidContext
	}

	var tc TraceContext

	traceIDBytes, err := hex.DecodeString(matches[2])
	if err != nil || len(traceIDBytes) != 16 {
		return nil, ErrInvalidTraceID
	}
	copy(tc.TraceID[:], traceIDBytes)

	spanIDBytes, err := hex.DecodeString(matches[3])
	if err != nil || len(spanIDBytes) != 8 {
		return nil, ErrInvalidSpanID
	}
	copy(tc.SpanID[:], spanIDBytes)

	flags, err := hex.DecodeString(matches[4])
	if err != nil || len(flags) != 1 {
		return nil, ErrInvalidContext
	}
	tc.TraceFlags = flags[0]

	return &tc, nil
}

// FormatW3CTraceparent formats a TraceContext as W3C traceparent.
func (tc *TraceContext) FormatW3CTraceparent() string {
	return fmt.Sprintf("00-%s-%s-%02x",
		tc.TraceID.String(),
		tc.SpanID.String(),
		tc.TraceFlags)
}

// ParseHeaders parses X-Trace-ID and X-Span-ID headers.
func ParseHeaders(traceID, spanID string) (*TraceContext, error) {
	var tc TraceContext

	traceIDBytes, err := hex.DecodeString(traceID)
	if err != nil || len(traceIDBytes) != 16 {
		return nil, ErrInvalidTraceID
	}
	copy(tc.TraceID[:], traceIDBytes)

	if spanID != "" {
		spanIDBytes, err := hex.DecodeString(spanID)
		if err != nil || len(spanIDBytes) != 8 {
			return nil, ErrInvalidSpanID
		}
		copy(tc.SpanID[:], spanIDBytes)
	}

	return &tc, nil
}

// IsSampled returns whether the trace is sampled.
func (tc *TraceContext) IsSampled() bool {
	return tc.TraceFlags&0x01 != 0
}

// SetSampled sets the sampled flag.
func (tc *TraceContext) SetSampled(sampled bool) {
	if sampled {
		tc.TraceFlags |= 0x01
	} else {
		tc.TraceFlags &^= 0x01
	}
}

// ContextWithTraceContext returns a context with trace context attached.
func ContextWithTraceContext(ctx context.Context, tc *TraceContext) context.Context {
	return context.WithValue(ctx, traceContextKey{}, tc)
}

// TraceContextFromContext retrieves parsed trace context.
func TraceContextFromContext(ctx context.Context) *TraceContext {
	if ctx == nil {
		return nil
	}
	if tc, ok := ctx.Value(traceContextKey{}).(*TraceContext); ok {
		return tc
	}
	return nil
}
