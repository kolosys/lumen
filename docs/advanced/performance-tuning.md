# Performance Tuning

Optimization techniques for high-throughput applications.

## Zero-Allocation Logging

The logs package uses several techniques to minimize allocations:

### Use Typed Field Constructors

```go
// Zero allocation: uses stack-allocated Field struct
log.Info("request", logs.Int("status", 200))

// Allocates: interface{} boxing
log.Info("request", logs.Any("status", 200))
```

### Pre-allocate Fields

```go
// Reusable fields for hot paths
var (
    methodField = logs.String("method", "")
    statusField = logs.Int("status", 0)
)

// Update values in hot path (still allocates new Field)
log.Info("request", logs.String("method", r.Method), logs.Int("status", code))
```

### Use Logger.With() for Common Fields

```go
// Create once at startup
requestLog := log.With(
    logs.String("service", "api"),
    logs.String("version", "1.0"),
)

// Hot path uses pre-allocated fields
requestLog.Info("request", logs.String("path", path))
```

## Async Logging

Enable async mode for non-blocking logging:

```go
log := logs.New(&logs.Options{
    AsyncBufferSize: 10000,  // Buffer 10k entries
})
defer log.Close()
```

**Trade-offs:**
- Faster: main goroutine doesn't wait for I/O
- Risk: entries may be lost on crash
- Memory: buffer consumes memory

### Buffer Sizing

```text
Rule of thumb:
- High throughput (>10k logs/sec): 50,000 - 100,000
- Medium throughput (1k-10k/sec): 10,000 - 50,000
- Low throughput (<1k/sec): 1,000 - 10,000
```

## Sampling

Reduce log volume in high-throughput scenarios:

### Rate Sampling

```go
// Max 100 logs per second per unique message
log := logs.New(&logs.Options{
    Sampler: logs.NewRateSampler(100, time.Second),
})
```

### Count Sampling

```go
// Log every 1000th occurrence
log := logs.New(&logs.Options{
    Sampler: logs.NewCountSampler(1000),
})
```

### Per-Level Sampling

```go
// Full fidelity for errors, sample debug heavily
sampler := logs.NewLevelSampler(&logs.AlwaysSampler{}).
    WithLevel(logs.DebugLevel, logs.NewCountSampler(100)).
    WithLevel(logs.TraceLevel, logs.NewCountSampler(1000))

log := logs.New(&logs.Options{Sampler: sampler})
```

## Metrics Optimization

### Atomic Operations

All metric types use lock-free atomic operations:

```go
// Counter: uses atomic.Uint64 with fixed-point math (6 decimal places)
counter.Inc("method", "status")  // No locks

// Gauge: uses atomic.Uint64 for float64 bits
gauge.Set(value, "label")  // No locks

// Histogram: per-bucket atomic counters
histogram.Observe(0.5, "endpoint")  // Minimal contention
```

### Label Cardinality

High cardinality labels impact performance and storage:

```go
// Good: bounded cardinality (<100 combinations)
counter.Inc("GET", "200")
counter.Inc("POST", "201")

// Bad: unbounded cardinality (millions of combinations)
counter.Inc(userID, requestID)
```

### Pre-compute Labels

```go
// Pre-compute label hash for repeated use
labels := metrics.NewLabels("method", "GET", "status", "200")
hash := labels.Hash()  // Reuse this hash
```

## Trace Optimization

### Sampling Strategies

Choose sampling based on your needs:

```go
// Development: 100% sampling
trace.AlwaysSample()

// Production: ratio-based
trace.RatioSample(0.1)  // 10% of traces

// Consistent: trace ID-based (same trace always sampled/not)
trace.TraceIDRatioSample(0.1)

// Cascade: follow parent decision
trace.ParentBasedSample(trace.RatioSample(0.1))
```

### Async Export

Enable async span export for non-blocking:

```go
tracer := trace.New(&trace.Options{
    AsyncExport:     true,
    AsyncBufferSize: 1024,
})
```

### Minimal Attributes

Add only necessary attributes:

```go
// Good: essential attributes
span.SetAttribute("user.id", userID)
span.SetAttribute("http.status_code", status)

// Avoid: large payloads
span.SetAttribute("request.body", largeBody)  // Memory intensive
```

## Benchmarking

### Run Benchmarks

```bash
go test -bench=. -benchmem ./...
```

### Example Benchmark

```go
func BenchmarkLogInfo(b *testing.B) {
    log := logs.New(&logs.Options{
        Formatter: &logs.NoopFormatter{},
    })

    b.ResetTimer()
    for range b.N {
        log.Info("test message", logs.Int("count", 42))
    }
}

func BenchmarkCounterInc(b *testing.B) {
    counter := metrics.NewCounter("test", "test counter", "label")

    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            counter.Inc("value")
        }
    })
}
```

### Expected Performance

```text
Logging (no-op formatter):
  BenchmarkLogInfo-8          10000000    100 ns/op    0 B/op    0 allocs/op

Metrics:
  BenchmarkCounterInc-8       50000000     25 ns/op    0 B/op    0 allocs/op
  BenchmarkGaugeSet-8         50000000     30 ns/op    0 B/op    0 allocs/op
  BenchmarkHistogramObserve-8 20000000     80 ns/op    0 B/op    0 allocs/op

Tracing:
  BenchmarkSpanStart-8        5000000     250 ns/op   48 B/op    1 allocs/op
```

## Profiling

### CPU Profiling

```bash
go test -cpuprofile=cpu.prof -bench=. ./...
go tool pprof cpu.prof
```

### Memory Profiling

```bash
go test -memprofile=mem.prof -bench=. ./...
go tool pprof mem.prof
```

### Allocation Tracking

```bash
go test -benchmem -bench=. ./...
```

## Production Checklist

- [ ] Enable async logging for high-throughput paths
- [ ] Configure appropriate sampling rates
- [ ] Use JSON formatter in production
- [ ] Bound label cardinality for metrics
- [ ] Enable trace sampling (10-100% based on volume)
- [ ] Set appropriate buffer sizes for async operations
- [ ] Close all resources on shutdown
- [ ] Monitor observability system overhead

## Further Reading

- [Best Practices](best-practices.md) - Production patterns
- [API Reference](../api-reference/) - Complete API documentation
