# metrics

**Import Path:** `github.com/kolosys/lumen/metrics`

Metrics collection with Prometheus export support and push-based delivery.

## Architecture

The metrics package provides:

- **Registry** - Manages metric registration and collection
- **Counter** - Cumulative values that only increase
- **Gauge** - Values that can go up or down
- **Histogram** - Observations in configurable buckets
- **Labels** - Dimension key-value pairs
- **Exporter** - Outputs metrics (Prometheus, push)

## Counter

Counters track cumulative values that only increase (or reset to zero).

```go
registry := metrics.NewRegistry(nil)
requests := registry.Counter(
    "http_requests_total",
    "Total HTTP requests",
    "method", "status",  // Label names
)

requests.Inc("GET", "200")      // Increment by 1
requests.Add(5, "POST", "201")  // Add arbitrary value

value := requests.Value("GET", "200")  // Read current value
```

Use cases: request counts, bytes transferred, errors.

## Gauge

Gauges track values that can increase or decrease.

```go
connections := registry.Gauge(
    "active_connections",
    "Current active connections",
    "pool",
)

connections.Set(10, "primary")  // Set absolute value
connections.Inc("primary")       // Increment by 1
connections.Dec("primary")       // Decrement by 1
connections.Add(-5, "primary")   // Add (negative allowed)
```

Use cases: active connections, queue size, temperature.

## Histogram

Histograms sample observations and count them in buckets.

```go
latency := registry.Histogram(
    "request_duration_seconds",
    "Request duration distribution",
    metrics.DefaultHistogramBuckets(),  // Default buckets
    "endpoint",
)

latency.Observe(0.25, "/api/users")
latency.Observe(0.05, "/api/health")
```

Generates `_bucket`, `_sum`, and `_count` series for Prometheus.

### Custom Buckets

```go
// Linear: start=0, width=0.1, count=5 => [0, 0.1, 0.2, 0.3, 0.4]
buckets := metrics.LinearBuckets(0, 0.1, 5)

// Exponential: start=0.001, factor=10, count=5 => [0.001, 0.01, 0.1, 1, 10]
buckets := metrics.ExponentialBuckets(0.001, 10, 5)

latency := registry.Histogram("duration", "Duration", buckets, "method")
```

## Labels

Labels add dimensions to metrics:

```go
// Create labels from key-value pairs
labels := metrics.NewLabels("method", "GET", "status", "200")

// Create from map
labels := metrics.LabelsFromMap(map[string]string{
    "method": "GET",
    "status": "200",
})

// Access
labels.Get("method")  // "GET"
labels.Keys()         // ["method", "status"]
labels.Len()          // 2
```

Labels are automatically sorted for consistent hashing.

## Registry

The registry manages metric lifecycle:

```go
registry := metrics.NewRegistry(&metrics.Options{
    Prefix:       "myapp_",
    DefaultLabels: map[string]string{"env": "prod"},
})

// Register metrics
counter := metrics.NewCounter("requests", "Total requests")
registry.Register(counter)

// Or use convenience methods
counter := registry.Counter("requests", "Total requests")

// Retrieve
m, err := registry.Get("requests")

// Collect all samples
samples := registry.Collect()

// Cleanup
registry.Close()
```

## Prometheus Export

Expose metrics via HTTP endpoint:

```go
import "net/http"

// Using default registry
http.Handle("/metrics", metrics.DefaultHTTPHandler())

// Using custom registry
http.Handle("/metrics", metrics.HTTPHandler(registry))

http.ListenAndServe(":9090", nil)
```

Output format:

```
http_requests_total{method="GET",status="200"} 1523
http_requests_total{method="POST",status="201"} 842
request_duration_seconds_bucket{endpoint="/api/users",le="0.1"} 450
request_duration_seconds_bucket{endpoint="/api/users",le="0.25"} 890
request_duration_seconds_bucket{endpoint="/api/users",le="+Inf"} 1000
request_duration_seconds_sum{endpoint="/api/users"} 125.5
request_duration_seconds_count{endpoint="/api/users"} 1000
```

## Push-Based Export

Push metrics periodically to an exporter:

```go
registry := metrics.NewRegistry(&metrics.Options{
    PushInterval: 10 * time.Second,
    PushExporter: myExporter,  // Implements Exporter interface
})
defer registry.Close()
```

Custom exporter:

```go
type MyExporter struct{}

func (e *MyExporter) Export(samples []metrics.Sample) {
    for _, s := range samples {
        fmt.Printf("%s{%s} = %f\n", s.Name, s.Labels.Hash(), s.Value)
    }
}
```

## Thread Safety

All metric operations are thread-safe using atomic operations:

- Counter uses `atomic.Uint64` with fixed-point math
- Gauge uses `atomic.Uint64` for float64 bits
- Histogram uses atomic counters per bucket
- Registry uses `sync.Map` for metric storage

## Configuration Reference

```go
opts := &metrics.Options{
    Prefix:           "myapp_",                    // Metric name prefix
    DefaultLabels:    map[string]string{},         // Added to all metrics
    HistogramBuckets: metrics.DefaultHistogramBuckets(),
    PushInterval:     0,                           // 0 = disabled
    PushExporter:     nil,                         // Push destination
}
```

## Further Reading

- [API Reference](../api-reference/metrics.md) - Complete API documentation
- [Best Practices](../advanced/best-practices.md) - Production patterns
- [Performance Tuning](../advanced/performance-tuning.md) - Optimization guide
