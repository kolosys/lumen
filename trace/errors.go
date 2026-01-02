package trace

import "errors"

var (
	ErrTracerClosed    = errors.New("trace: tracer is closed")
	ErrInvalidTraceID  = errors.New("trace: invalid trace ID")
	ErrInvalidSpanID   = errors.New("trace: invalid span ID")
	ErrInvalidContext  = errors.New("trace: invalid trace context")
	ErrSamplerRejected = errors.New("trace: span rejected by sampler")
	ErrExporterFailed  = errors.New("trace: exporter failed")
)
