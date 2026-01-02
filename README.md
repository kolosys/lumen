# Lumen

A developer-first observability library for Go. Zero dependencies, high performance.

## Packages

| Package | Description |
|---------|-------------|
| `logs` | Structured logging with named instances |
| `trace` | Distributed tracing with W3C support |
| `metrics` | Prometheus-compatible metrics |

## Installation

```bash
go get github.com/kolosys/lumen
```

## Logs

High-performance structured logging with named logger instances.

```go
import "github.com/kolosys/lumen/logs"

// Simple usage
logs.Info("server started", logs.Int("port", 8080))

// Named loggers (micro-instances)
gateway := logs.NewNamed("gateway")
gateway.Info("connected")  // Output: INFO [gateway] connected

shard := gateway.Named("shard.0")
shard.SetLevel(logs.DebugLevel)
shard.Debug("heartbeat")   // Output: DEBG [gateway.shard.0] heartbeat

// With fields
log := logs.New(nil).With(logs.String("service", "api"))
log.Info("request handled", logs.Duration("latency", latency))
```

### Formatters

```go
// Text (default)
logs.New(&logs.Options{Formatter: &logs.TextFormatter{}})

// JSON
logs.New(&logs.Options{Formatter: &logs.JSONFormatter{}})

// Pretty (development)
logs.New(&logs.Options{Formatter: &logs.PrettyFormatter{}})
```

## Trace

Distributed tracing with W3C Trace Context and custom header support.

```go
import "github.com/kolosys/lumen/trace"

tracer := trace.New(nil)

ctx, span := tracer.Start(ctx, "handle-request")
defer span.End()

span.SetAttribute("user.id", userID)
span.AddEvent("validated input")

if err != nil {
    span.RecordError(err)
}
```

### Propagation

```go
// W3C Trace Context (traceparent header)
propagator := &trace.W3CPropagator{}

// Custom headers (X-Trace-ID, X-Span-ID)
propagator := &trace.HeaderPropagator{}

// Both formats
propagator := trace.DefaultPropagator()

// Inject into outgoing request
propagator.Inject(ctx, trace.MapCarrier(headers))

// Extract from incoming request
ctx = propagator.Extract(ctx, trace.MapCarrier(headers))
```

## Metrics

Prometheus-compatible metrics with labels.

```go
import "github.com/kolosys/lumen/metrics"

registry := metrics.New(nil)

// Counter
requests := metrics.NewCounter("http_requests_total", "Total requests", "method", "path")
registry.Register(requests)
requests.Inc("GET", "/api/users")

// Gauge
connections := metrics.NewGauge("active_connections", "Active connections")
registry.Register(connections)
connections.Set(42)

// Histogram
latency := metrics.NewHistogram("request_duration_seconds", "Request latency", nil, "method")
registry.Register(latency)
latency.Observe(0.025, "GET")

// Prometheus endpoint
http.Handle("/metrics", metrics.HTTPHandler(registry))
```

## Design Principles

- **Zero dependencies** - stdlib only
- **Zero-allocation hot paths** - sync.Pool, atomics
- **Developer-first** - simple API, sensible defaults
- **Production-ready** - context-aware, thread-safe

## License

MIT
