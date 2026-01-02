# trace

**Import Path:** `github.com/kolosys/lumen/trace`

Distributed tracing with W3C Trace Context and custom header propagation.

## Architecture

The trace package provides:

- **Tracer** - Creates and manages spans
- **Span** - Represents a unit of work with timing and metadata
- **TraceID/SpanID** - Unique identifiers for correlation
- **Sampler** - Determines which traces to record
- **Exporter** - Outputs completed spans
- **Propagator** - Injects/extracts context across boundaries

## Creating Spans

```go
tracer := trace.New(nil)

ctx, span := tracer.Start(ctx, "operation-name")
defer span.End()

// Add attributes
span.SetAttribute("user.id", userID)
span.SetAttribute("order.total", 99.99)
```

Using the default tracer:

```go
ctx, span := trace.Start(ctx, "operation")
defer span.End()
```

## Span Hierarchy

Spans automatically form parent-child relationships through context:

```go
func HandleRequest(ctx context.Context) {
    ctx, span := trace.Start(ctx, "handle-request")
    defer span.End()

    validateInput(ctx)   // Child span
    processData(ctx)     // Child span
    saveResult(ctx)      // Child span
}

func validateInput(ctx context.Context) {
    ctx, span := trace.Start(ctx, "validate-input")
    defer span.End()
    // Automatically linked to parent
}
```

## Attributes

Add metadata to spans:

```go
span.SetAttribute("http.method", "GET")
span.SetAttribute("http.status_code", 200)
span.SetAttribute("db.statement", query)

// Multiple at once
span.SetAttributes(
    trace.Attribute{Key: "user.id", Value: userID},
    trace.Attribute{Key: "user.role", Value: "admin"},
)
```

## Events

Record timestamped annotations:

```go
span.AddEvent("cache.miss", trace.Attribute{Key: "key", Value: cacheKey})
span.AddEvent("retry.attempt", trace.Attribute{Key: "attempt", Value: 2})
span.AddEvent("query.executed")
```

## Error Handling

Record errors and set status:

```go
ctx, span := trace.Start(ctx, "database-query")
defer span.End()

result, err := db.Query(ctx, query)
if err != nil {
    span.RecordError(err)  // Adds exception event and sets error status
    return nil, err
}

span.SetStatus(trace.StatusOK, "")
```

Status values: `StatusUnset`, `StatusOK`, `StatusError`

## Span Options

Configure spans at creation:

```go
ctx, span := trace.Start(ctx, "operation",
    trace.WithAttributes(
        trace.Attribute{Key: "component", Value: "api"},
    ),
    trace.WithStartTime(customTime),
)
```

## Sampling

Control which traces are recorded:

```go
// Always sample (default)
tracer := trace.New(&trace.Options{
    Sampler: trace.AlwaysSample(),
})

// Never sample
tracer := trace.New(&trace.Options{
    Sampler: trace.NeverSample(),
})

// Sample 10% of traces (counter-based)
tracer := trace.New(&trace.Options{
    Sampler: trace.RatioSample(0.1),
})

// Sample based on trace ID (deterministic)
tracer := trace.New(&trace.Options{
    Sampler: trace.TraceIDRatioSample(0.1),
})

// Follow parent's sampling decision
tracer := trace.New(&trace.Options{
    Sampler: trace.ParentBasedSample(trace.RatioSample(0.1)),
})
```

## Exporters

Output completed spans:

### Writer Exporter

```go
exporter := trace.NewWriterExporter(os.Stdout)
tracer := trace.New(&trace.Options{
    Exporter: exporter,
})
```

Outputs JSON lines:

```json
{"trace_id":"abc123","span_id":"def456","name":"operation","duration_ns":1500000}
```

### In-Memory Exporter (Testing)

```go
exporter := trace.NewInMemoryExporter()
tracer := trace.New(&trace.Options{
    Exporter: exporter,
})

// After running code
spans := exporter.Spans()
assert.Equal(t, 3, exporter.Len())
exporter.Clear()
```

### Custom Exporter

```go
type OTLPExporter struct{}

func (e *OTLPExporter) Export(span *trace.Span) {
    // Send to OTLP endpoint
}

func (e *OTLPExporter) Close() error {
    return nil
}
```

## Context Propagation

Propagate trace context across service boundaries.

### W3C Trace Context

```go
propagator := &trace.W3CPropagator{}

// Inject into outgoing request
propagator.Inject(ctx, trace.MapCarrier(headers))
// Sets: traceparent: 00-{trace_id}-{span_id}-{flags}

// Extract from incoming request
ctx = propagator.Extract(ctx, trace.MapCarrier(headers))
```

### Custom Headers

```go
propagator := &trace.HeaderPropagator{
    TraceIDHeader: "X-Trace-ID",
    SpanIDHeader:  "X-Span-ID",
}
```

### Composite Propagator

```go
propagator := trace.NewCompositePropagator(
    &trace.W3CPropagator{},
    &trace.HeaderPropagator{},
)

// Or use default (supports both)
propagator := trace.DefaultPropagator()
```

### HTTP Integration

```go
func Handler(w http.ResponseWriter, r *http.Request) {
    propagator := trace.DefaultPropagator()
    ctx := propagator.Extract(r.Context(), headerCarrier(r.Header))

    ctx, span := trace.Start(ctx, "handle-request")
    defer span.End()

    // Make outgoing request
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    propagator.Inject(ctx, headerCarrier(req.Header))
}
```

## Async Export

Non-blocking span export for high throughput:

```go
tracer := trace.New(&trace.Options{
    AsyncExport:     true,
    AsyncBufferSize: 1024,
    Exporter:        exporter,
})
defer tracer.Close()  // Flush pending spans
```

## Configuration Reference

```go
opts := &trace.Options{
    ServiceName:       "my-service",         // Service identifier
    Sampler:           trace.AlwaysSample(), // Sampling strategy
    Exporter:          trace.NopExporter{},  // Span destination
    PropagationFormat: "both",               // "w3c", "kolosys", "both"
    AsyncExport:       false,                // Async export mode
    AsyncBufferSize:   1024,                 // Async buffer size
    MaxSpansPerSecond: 0,                    // Rate limit (0 = unlimited)
}
```

## Further Reading

- [API Reference](../api-reference/trace.md) - Complete API documentation
- [Best Practices](../advanced/best-practices.md) - Production patterns
- [Performance Tuning](../advanced/performance-tuning.md) - Optimization guide
