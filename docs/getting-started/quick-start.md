# Quick Start

This guide will help you get started with lumen quickly with a basic example.

## Basic Usage

Here's a simple example to get you started:

```go
package main

import (
    "fmt"
    "log"
    "github.com/kolosys/lumen/logs"
    "github.com/kolosys/lumen/metrics"
    "github.com/kolosys/lumen/trace"
)

func main() {
    // Basic usage example
    fmt.Println("Welcome to lumen!")
    
    // TODO: Add your code here
}
```

## Common Use Cases

### Using logs

**Import Path:** `github.com/kolosys/lumen/logs`

Package logs provides a high-performance, context-aware structured logging library.

Features:
  - Zero-allocation hot paths using sync.Pool
  - Context-aware logging with context.Context
  - Type-safe field builders
  - Multiple output formats (text, JSON, pretty)
  - Named logger instances with automatic prefix display
  - Per-instance log level configuration
  - Sampling for high-volume logs
  - Async logging option
  - Hook system for extensibility
  - Built-in caller information
  - Chained/fluent API

Basic usage:

	log := logs.New(nil)
	log.Info("server started", logs.Int("port", 8080))

Named loggers (micro-instances):

	gateway := logs.NewNamed("gateway")
	gateway.Info("connected")           // Output: INFO [gateway] connected

	shard := gateway.Named("shard.0")
	shard.SetLevel(logs.DebugLevel)     // Per-instance level
	shard.Debug("heartbeat received")   // Output: DEBG [gateway.shard.0] heartbeat received

With context:

	log.InfoContext(ctx, "request processed", logs.Duration("latency", time.Since(start)))


```go
package main

import (
    "fmt"
    "github.com/kolosys/lumen/logs"
)

func main() {
    // Example usage of logs
    fmt.Println("Using logs package")
}
```

#### Available Types
- **AlwaysSampler** - AlwaysSampler always allows logging.
- **Builder** - Builder provides a fluent/chainable API for building log entries. It accumulates fields and then emits a log entry when a level method is called.
- **CompositeSampler** - CompositeSampler combines multiple samplers with AND logic.
- **CountSampler** - CountSampler logs every Nth occurrence.
- **Entry** - Entry represents a log entry.
- **ErrorBuilder** - ErrorBuilder provides a fluent API for logging errors.
- **ErrorHook** - ErrorHook collects errors for inspection.
- **Field** - Field represents a structured log field.
- **FieldType** - FieldType represents the type of a field value.
- **FileHook** - FileHook writes entries to a file.
- **FilterHook** - FilterHook conditionally fires another hook.
- **FirstNSampler** - FirstNSampler logs only the first N occurrences.
- **Formatter** - Formatter formats log entries.
- **FuncHook** - FuncHook wraps a function as a hook.
- **Hook** - Hook is called when a log entry is written.
- **JSONFormatter** - JSONFormatter formats logs as JSON.
- **Level** - Level represents a log level.
- **LevelHook** - LevelHook fires only for specific levels.
- **LevelSampler** - LevelSampler applies different samplers per level.
- **Logger** - Logger is the main logging interface.
- **MetricsHook** - MetricsHook tracks log counts by level.
- **NamedFormatter** - NamedFormatter is a formatter wrapper for custom name formatting. Note: The built-in formatters (TextFormatter, JSONFormatter, PrettyFormatter) now natively display logger names in [brackets]. This wrapper is only needed if you want custom brackets or separators.
- **NeverSampler** - NeverSampler never allows logging.
- **NoopFormatter** - NoopFormatter discards all output.
- **OncePerSampler** - OncePerSampler logs a message only once per duration.
- **Options** - Options configures a Logger.
- **PrettyFormatter** - PrettyFormatter formats logs with colors and alignment for development.
- **RandomSampler** - RandomSampler samples a percentage of logs.
- **RateSampler** - RateSampler limits logs to a certain rate per message.
- **Sampler** - Sampler determines if a log entry should be emitted.
- **TextFormatter** - TextFormatter formats logs as text.
- **WriterHook** - WriterHook writes entries to an io.Writer.

#### Available Functions
- **CtxDebug** - CtxDebug logs at debug level using the logger from context.
- **CtxError** - CtxError logs at error level using the logger from context.
- **CtxInfo** - CtxInfo logs at info level using the logger from context.
- **CtxTrace** - CtxTrace logs at trace level using the logger from context.
- **CtxWarn** - CtxWarn logs at warn level using the logger from context.
- **Debug** - Debug logs at debug level using the default logger.
- **Debugf** - Debugf logs a formatted message at debug level.
- **Error** - Error logs at error level using the default logger.
- **Errorf** - Errorf logs a formatted message at error level.
- **Fatal** - Fatal logs at fatal level using the default logger and exits.
- **Fatalf** - Fatalf logs a formatted message at fatal level and exits.
- **Info** - Info logs at info level using the default logger.
- **Infof** - Infof logs a formatted message at info level.
- **Must** - Must logs and panics if error is not nil. Useful for initialization code. db := log.Must(sql.Open("postgres", dsn))
- **Panic** - Panic logs at panic level using the default logger and panics.
- **Panicf** - Panicf logs a formatted message at panic level and panics.
- **Print** - Print logs a message at info level.
- **Printf** - Printf logs a formatted message at info level.
- **Println** - Println logs a message at info level.
- **SetDefault** - SetDefault sets the default logger.
- **SetDefaultFormatter** - SetDefaultFormatter sets the default formatter.
- **SetDefaultLevel** - SetDefaultLevel sets the default level.
- **Trace** - Trace logs at trace level using the default logger.
- **Tracef** - Tracef logs a formatted message at trace level.
- **Warn** - Warn logs at warn level using the default logger.
- **Warnf** - Warnf logs a formatted message at warn level.
- **WithContextFields** - WithFields adds fields to the context that will be included in all logs.
- **WithLogger** - WithLogger attaches a logger to the context.
- **WithRequestID** - WithRequestID adds a request ID to the context.
- **WithTraceID** - WithTraceID adds a trace ID to the context.
- **WithUserID** - WithUserID adds a user ID to the context.

For detailed API documentation, see the [logs API Reference](../api-reference/logs.md).

### Using metrics

**Import Path:** `github.com/kolosys/lumen/metrics`

Package metrics provides metrics collection with Prometheus and push support.


```go
package main

import (
    "fmt"
    "github.com/kolosys/lumen/metrics"
)

func main() {
    // Example usage of metrics
    fmt.Println("Using metrics package")
}
```

#### Available Types
- **Counter** - Counter is a cumulative metric that only increases.
- **Exporter** - Exporter exports metrics.
- **Gauge** - Gauge is a metric that can go up and down.
- **Histogram** - Histogram samples observations and counts them in buckets.
- **Labels** - Labels is a sorted set of label key-value pairs.
- **Metric** - Metric is the interface all metric types implement.
- **MetricType** - MetricType identifies the metric type.
- **NopExporter** - NopExporter discards metrics.
- **Options** - Options configures a Registry.
- **Registry** - Registry manages metric registration and collection.
- **Sample** - Sample is a single metric observation.

#### Available Functions
- **DefaultHTTPHandler** - DefaultHTTPHandler returns an http.Handler using the default registry.
- **DefaultHistogramBuckets** - DefaultHistogramBuckets returns commonly used bucket boundaries.
- **ExponentialBuckets** - ExponentialBuckets creates exponentially growing buckets.
- **HTTPHandler** - HTTPHandler returns an http.Handler for the Prometheus endpoint.
- **LinearBuckets** - LinearBuckets creates n buckets of equal width.
- **SetDefaultRegistry** - SetDefaultRegistry sets the default registry.
- **WritePrometheus** - WritePrometheus writes samples in Prometheus text format.

For detailed API documentation, see the [metrics API Reference](../api-reference/metrics.md).

### Using trace

**Import Path:** `github.com/kolosys/lumen/trace`

Package trace provides distributed tracing with W3C and Kolosys format support.


```go
package main

import (
    "fmt"
    "github.com/kolosys/lumen/trace"
)

func main() {
    // Example usage of trace
    fmt.Println("Using trace package")
}
```

#### Available Types
- **AlwaysSampler** - AlwaysSampler always samples.
- **Attribute** - Attribute is a key-value pair attached to a span.
- **Carrier** - Carrier is an interface for propagation carriers (e.g., HTTP headers).
- **CompositePropagator** - CompositePropagator combines multiple propagators.
- **Event** - Event is a timestamped annotation.
- **Exporter** - Exporter receives completed spans.
- **HeaderPropagator** - HeaderPropagator implements customizable header-based propagation.
- **InMemoryExporter** - InMemoryExporter collects spans in memory for testing.
- **MapCarrier** - MapCarrier is a map-based carrier.
- **NeverSampler** - NeverSampler never samples.
- **NopExporter** - NopExporter discards all spans.
- **Options** - Options configures a Tracer.
- **ParentBasedSampler** - ParentBasedSampler follows parent sampling decision.
- **Propagator** - Propagator handles trace context injection and extraction.
- **RatioSampler** - RatioSampler samples a percentage of traces.
- **Sampler** - Sampler determines whether a span should be recorded.
- **SamplingParams** - SamplingParams provides data for sampling decisions.
- **Span** - Span represents a unit of work.
- **SpanID** - SpanID is an 8-byte span identifier.
- **SpanOption** - SpanOption configures span creation.
- **SpanStatus** - SpanStatus represents span completion status.
- **TraceContext** - TraceContext holds W3C Trace Context data.
- **TraceID** - TraceID is a 16-byte trace identifier.
- **TraceIDRatioSampler** - TraceIDRatioSampler samples based on trace ID for consistency.
- **Tracer** - Tracer creates and manages spans.
- **W3CPropagator** - W3CPropagator implements W3C Trace Context propagation.
- **WriterExporter** - WriterExporter writes spans as JSON to an io.Writer.

#### Available Functions
- **ContextWithSpan** - ContextWithSpan returns a context with the span attached.
- **ContextWithTraceContext** - ContextWithTraceContext returns a context with trace context attached.
- **SetDefault** - SetDefault sets the default tracer.

For detailed API documentation, see the [trace API Reference](../api-reference/trace.md).

## Step-by-Step Tutorial

### Step 1: Import the Package

First, import the necessary packages in your Go file:

```go
import (
    "fmt"
    "github.com/kolosys/lumen/logs"
    "github.com/kolosys/lumen/metrics"
    "github.com/kolosys/lumen/trace"
)
```

### Step 2: Initialize

Set up the basic configuration:

```go
func main() {
    // Initialize your application
    fmt.Println("Initializing lumen...")
}
```

### Step 3: Use the Library

Implement your specific use case:

```go
func main() {
    // Your implementation here
}
```

## Running Your Code

To run your Go program:

```bash
go run main.go
```

To build an executable:

```bash
go build -o myapp
./myapp
```

## Configuration Options

lumen can be configured to suit your needs. Check the [Core Concepts](../core-concepts/) section for detailed information about configuration options.

## Error Handling

Always handle errors appropriately:

```go
result, err := someFunction()
if err != nil {
    log.Fatalf("Error: %v", err)
}
```

## Best Practices

- Always handle errors returned by library functions
- Check the API documentation for detailed parameter information
- Use meaningful variable and function names
- Add comments to document your code

## Complete Example

Here's a complete working example:

```go
package main

import (
    "fmt"
    "log"
    "github.com/kolosys/lumen/logs"
    "github.com/kolosys/lumen/metrics"
    "github.com/kolosys/lumen/trace"
)

func main() {
    fmt.Println("Starting lumen application...")
    
    // Add your implementation here
    
    fmt.Println("Application completed successfully!")
}
```

## Next Steps

Now that you've seen the basics, explore:

- **[Core Concepts](../core-concepts/)** - Understanding the library architecture
- **[API Reference](../api-reference/)** - Complete API documentation
- **[Examples](../examples/README.md)** - More detailed examples
- **[Advanced Topics](../advanced/)** - Performance tuning and advanced patterns

## Getting Help

If you run into issues:

1. Check the [API Reference](../api-reference/)
2. Browse the [Examples](../examples/README.md)
3. Visit the [GitHub Issues](https://github.com/kolosys/lumen/issues) page

