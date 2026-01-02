# Quick Start

This guide covers the essential usage patterns for each lumen package.

## Logging

### Basic Usage

```go
package main

import "github.com/kolosys/lumen/logs"

func main() {
    log := logs.New(nil)
    log.Info("application started", logs.String("version", "1.0.0"))
    log.Debug("debug message", logs.Int("count", 42))
    log.Warn("warning message")
    log.Error("error occurred", logs.Err(err))
}
```

### Named Loggers

Create hierarchical loggers for different components:

```go
gateway := logs.NewNamed("gateway")
gateway.Info("connected")  // Output: INFO [gateway] connected

shard := gateway.Named("shard.0")
shard.SetLevel(logs.DebugLevel)
shard.Debug("heartbeat")  // Output: DEBG [gateway.shard.0] heartbeat
```

### Structured Fields

Type-safe field builders:

```go
log.Info("request completed",
    logs.String("method", "GET"),
    logs.Int("status", 200),
    logs.Duration("latency", time.Since(start)),
    logs.Bool("cached", true),
    logs.Err(err),
)
```

### Context-Aware Logging

```go
ctx := logs.WithRequestID(ctx, "req-123")
ctx = logs.WithUserID(ctx, "user-456")

log.InfoContext(ctx, "processing request")  // Includes request_id and user_id
```

### Fluent Builder API

```go
log.Build().
    Str("user", "alice").
    Int("items", 5).
    Info("order placed")
```

### Configuration

```go
log := logs.New(&logs.Options{
    Level:           logs.DebugLevel,
    Formatter:       &logs.JSONFormatter{},
    AddCaller:       true,
    AsyncBufferSize: 1000,  // Enable async logging
})
defer log.Close()
```

## Metrics

### Counter

Tracks cumulative values that only increase:

```go
registry := metrics.NewRegistry(nil)
requests := registry.Counter("http_requests_total", "Total HTTP requests", "method", "status")

requests.Inc("GET", "200")
requests.Add(5, "POST", "201")
```

### Gauge

Tracks values that can go up or down:

```go
connections := registry.Gauge("active_connections", "Active connections", "pool")

connections.Set(10, "primary")
connections.Inc("primary")
connections.Dec("primary")
```

### Histogram

Records observations in configurable buckets:

```go
latency := registry.Histogram(
    "request_duration_seconds",
    "Request duration",
    metrics.DefaultHistogramBuckets(),
    "endpoint",
)

latency.Observe(0.25, "/api/users")
```

### Prometheus Export

```go
import "net/http"

http.Handle("/metrics", metrics.DefaultHTTPHandler())
http.ListenAndServe(":9090", nil)
```

### Custom Buckets

```go
// Linear: 0, 0.1, 0.2, 0.3, 0.4
buckets := metrics.LinearBuckets(0, 0.1, 5)

// Exponential: 0.001, 0.01, 0.1, 1, 10
buckets := metrics.ExponentialBuckets(0.001, 10, 5)
```

## Tracing

### Basic Span

```go
ctx, span := trace.Start(ctx, "process-order")
defer span.End()

span.SetAttribute("order.id", orderID)
span.SetAttribute("order.total", total)
```

### Nested Spans

```go
func HandleRequest(ctx context.Context) {
    ctx, span := trace.Start(ctx, "handle-request")
    defer span.End()

    processData(ctx)  // Creates child span
}

func processData(ctx context.Context) {
    ctx, span := trace.Start(ctx, "process-data")
    defer span.End()
    // Parent span is automatically linked
}
```

### Error Recording

```go
ctx, span := trace.Start(ctx, "database-query")
defer span.End()

result, err := db.Query(ctx, query)
if err != nil {
    span.RecordError(err)  // Records error and sets status
    return err
}
span.SetStatus(trace.StatusOK, "")
```

### Span Events

```go
span.AddEvent("cache.miss", trace.Attribute{Key: "key", Value: cacheKey})
span.AddEvent("retry.attempt", trace.Attribute{Key: "attempt", Value: 2})
```

### Tracer Configuration

```go
tracer := trace.New(&trace.Options{
    ServiceName: "my-service",
    Sampler:     trace.RatioSample(0.1),  // Sample 10%
    Exporter:    trace.NewWriterExporter(os.Stdout),
    AsyncExport: true,
})
defer tracer.Close()

trace.SetDefault(tracer)
```

### HTTP Propagation

```go
propagator := trace.DefaultPropagator()

// Inject into outgoing request
propagator.Inject(ctx, trace.MapCarrier(headers))

// Extract from incoming request
ctx = propagator.Extract(ctx, trace.MapCarrier(headers))
```

## Next Steps

- [Core Concepts](../core-concepts/) - Detailed package architecture
- [Best Practices](../advanced/best-practices.md) - Production patterns
- [Performance Tuning](../advanced/performance-tuning.md) - Optimization guide
- [API Reference](../api-reference/) - Complete API documentation
