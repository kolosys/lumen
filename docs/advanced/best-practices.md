# Best Practices

Production patterns and recommendations for using lumen effectively.

## Logging

### Use Named Loggers

Create dedicated loggers for each component:

```go
var (
    httpLog = logs.NewNamed("http")
    dbLog   = logs.NewNamed("db")
    cacheLog = logs.NewNamed("cache")
)

func HandleRequest(w http.ResponseWriter, r *http.Request) {
    httpLog.Info("request", logs.String("path", r.URL.Path))
}
```

### Prefer Structured Fields Over Printf

```go
// Good: structured, searchable, zero-allocation
log.Info("user created", logs.String("user_id", id), logs.Int("age", 25))

// Avoid: unstructured, allocates
log.Infof("user %s created, age %d", id, 25)
```

### Include Context in Logs

```go
func ProcessOrder(ctx context.Context, orderID string) error {
    ctx = logs.WithContextFields(ctx, logs.String("order_id", orderID))

    log.InfoContext(ctx, "processing started")
    // All subsequent logs include order_id

    if err := validateOrder(ctx); err != nil {
        log.ErrorContext(ctx, "validation failed", logs.Err(err))
        return err
    }
    return nil
}
```

### Use Appropriate Log Levels

```go
log.Trace("detailed debugging info")     // Development only
log.Debug("debugging information")        // Development/debugging
log.Info("normal operations")             // Default production level
log.Warn("unexpected but handled")        // Warning signs
log.Error("operation failed", logs.Err(err))  // Errors
```

### Configure Level by Environment

```go
level := logs.ParseLevel(os.Getenv("LOG_LEVEL"))
log := logs.New(&logs.Options{Level: level})
```

### Use JSON in Production

```go
var formatter logs.Formatter
if os.Getenv("ENV") == "production" {
    formatter = &logs.JSONFormatter{}
} else {
    formatter = &logs.PrettyFormatter{ShowCaller: true}
}
log := logs.New(&logs.Options{Formatter: formatter})
```

## Metrics

### Follow Naming Conventions

```go
// Use snake_case
// Include unit as suffix (_seconds, _bytes, _total)
// Use _total suffix for counters

registry.Counter("http_requests_total", "Total HTTP requests")
registry.Histogram("http_request_duration_seconds", "Request duration")
registry.Gauge("db_connections_active", "Active database connections")
```

### Use Labels Judiciously

```go
// Good: bounded cardinality
counter.Inc("GET", "200")  // method, status

// Avoid: unbounded cardinality (user IDs, timestamps, etc.)
counter.Inc(userID, timestamp)  // Creates millions of time series
```

### Pre-Register Metrics

```go
var (
    requestCounter = metrics.NewCounter(
        "http_requests_total", "Total requests", "method", "status")
    requestDuration = metrics.NewHistogram(
        "http_request_duration_seconds", "Duration",
        metrics.DefaultHistogramBuckets(), "endpoint")
)

func init() {
    registry := metrics.DefaultRegistry()
    registry.Register(requestCounter)
    registry.Register(requestDuration)
}
```

### Choose Appropriate Buckets

```go
// HTTP latency (most requests < 1s)
httpBuckets := []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}

// Database queries (wider range)
dbBuckets := metrics.ExponentialBuckets(0.001, 2, 15)  // 1ms to ~16s

// Batch processing (longer durations)
batchBuckets := metrics.LinearBuckets(0, 60, 10)  // 0 to 10 minutes
```

## Tracing

### Use Meaningful Span Names

```go
// Good: descriptive, hierarchical
trace.Start(ctx, "http.server.request")
trace.Start(ctx, "db.query.users")
trace.Start(ctx, "cache.redis.get")

// Avoid: generic, uninformative
trace.Start(ctx, "operation")
trace.Start(ctx, "process")
```

### Add Relevant Attributes

```go
ctx, span := trace.Start(ctx, "http.request")
defer span.End()

span.SetAttribute("http.method", r.Method)
span.SetAttribute("http.url", r.URL.Path)
span.SetAttribute("http.status_code", statusCode)
span.SetAttribute("http.request_content_length", r.ContentLength)
```

### Record Errors Consistently

```go
func ProcessItem(ctx context.Context, id string) error {
    ctx, span := trace.Start(ctx, "process-item")
    defer span.End()

    span.SetAttribute("item.id", id)

    result, err := doWork(ctx, id)
    if err != nil {
        span.RecordError(err)  // Always record errors
        return err
    }

    span.SetStatus(trace.StatusOK, "")
    return nil
}
```

### Use Parent-Based Sampling

```go
// Production: sample based on parent, 10% for root spans
tracer := trace.New(&trace.Options{
    Sampler: trace.ParentBasedSample(trace.RatioSample(0.1)),
})
```

### Propagate Context Across Boundaries

```go
// HTTP client
func CallService(ctx context.Context, url string) error {
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    trace.DefaultPropagator().Inject(ctx, headerCarrier(req.Header))
    return client.Do(req)
}

// HTTP server
func Handler(w http.ResponseWriter, r *http.Request) {
    ctx := trace.DefaultPropagator().Extract(r.Context(), headerCarrier(r.Header))
    ctx, span := trace.Start(ctx, "handle-request")
    defer span.End()
}
```

## Integration Patterns

### Correlation Across Packages

```go
func HandleRequest(ctx context.Context, r *http.Request) {
    // Extract trace context
    ctx = trace.DefaultPropagator().Extract(ctx, headerCarrier(r.Header))

    // Start span
    ctx, span := trace.Start(ctx, "http.request")
    defer span.End()

    // Add trace ID to logs
    ctx = logs.WithTraceID(ctx, span.TraceID().String())

    // Now logs and traces are correlated
    log.InfoContext(ctx, "request started")
}
```

### HTTP Middleware

```go
func ObservabilityMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        ctx := trace.DefaultPropagator().Extract(r.Context(), headerCarrier(r.Header))
        ctx, span := trace.Start(ctx, "http.request")
        defer span.End()

        ctx = logs.WithContextFields(ctx,
            logs.String("trace_id", span.TraceID().String()),
            logs.String("method", r.Method),
            logs.String("path", r.URL.Path),
        )

        rw := &responseWriter{ResponseWriter: w}
        next.ServeHTTP(rw, r.WithContext(ctx))

        duration := time.Since(start)

        span.SetAttribute("http.status_code", rw.status)
        requestDuration.Observe(duration.Seconds(), r.URL.Path)
        requestCounter.Inc(r.Method, strconv.Itoa(rw.status))

        log.InfoContext(ctx, "request completed",
            logs.Int("status", rw.status),
            logs.Duration("duration", duration),
        )
    })
}
```

## Error Handling

### Wrap Errors with Context

```go
if err := db.Query(ctx, query); err != nil {
    return log.WrapErr(err, "database query failed",
        logs.String("query", query),
        logs.String("table", "users"),
    )
}
```

### Use CheckErr for Control Flow

```go
if log.CheckErr(err, "failed to connect") {
    return
}
```

### Conditional Error Logging

```go
log.IfErr(err).
    With("user_id", userID).
    With("action", "create").
    Error("user creation failed")
```

## Resource Management

### Close Async Loggers

```go
log := logs.New(&logs.Options{AsyncBufferSize: 1000})
defer log.Close()  // Flush pending entries
```

### Close Tracers

```go
tracer := trace.New(&trace.Options{AsyncExport: true})
defer tracer.Close()  // Flush pending spans
```

### Close Registries

```go
registry := metrics.NewRegistry(&metrics.Options{PushInterval: 10 * time.Second})
defer registry.Close()  // Stop push loop
```

## Further Reading

- [Performance Tuning](performance-tuning.md) - Optimization techniques
- [API Reference](../api-reference/) - Complete API documentation
