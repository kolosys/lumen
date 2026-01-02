# logs

**Import Path:** `github.com/kolosys/lumen/logs`

High-performance, context-aware structured logging with zero-allocation hot paths.

## Architecture

The logs package is built around these core components:

- **Logger** - Main logging interface with level-based methods
- **Entry** - Represents a log entry with message, level, time, and fields
- **Field** - Type-safe structured data attached to entries
- **Formatter** - Transforms entries into output bytes (Text, JSON, Pretty)
- **Hook** - Extensibility points for custom processing
- **Sampler** - Rate limiting for high-volume logs

## Log Levels

Levels from most to least severe:

| Level | Method | Description |
|-------|--------|-------------|
| Panic | `Panic()` | Logs and panics |
| Fatal | `Fatal()` | Logs and exits |
| Error | `Error()` | Errors requiring attention |
| Warn | `Warn()` | Non-critical issues |
| Info | `Info()` | General information |
| Debug | `Debug()` | Debugging details |
| Trace | `Trace()` | Fine-grained debugging |

```go
log.SetLevel(logs.DebugLevel)  // Enable debug and above
```

## Named Loggers

Create hierarchical loggers for different components:

```go
log := logs.NewNamed("gateway")
log.Info("connected")  // INFO [gateway] connected

shard := log.Named("shard.0")
shard.Info("ready")    // INFO [gateway.shard.0] ready
```

Each named logger can have independent configuration:

```go
shard.SetLevel(logs.DebugLevel)  // Per-instance level
```

## Structured Fields

Type-safe field constructors avoid interface{} boxing:

```go
log.Info("request",
    logs.String("method", "GET"),
    logs.Int("status", 200),
    logs.Int64("bytes", 1024),
    logs.Float64("latency_ms", 12.5),
    logs.Bool("cached", true),
    logs.Duration("elapsed", time.Since(start)),
    logs.Time("timestamp", time.Now()),
    logs.Err(err),
    logs.Any("data", complexStruct),
)
```

## Context Integration

Attach fields to context for automatic inclusion:

```go
ctx := logs.WithContextFields(ctx,
    logs.String("request_id", reqID),
    logs.String("user_id", userID),
)

log.InfoContext(ctx, "processing")  // Includes request_id, user_id
```

Helper functions for common fields:

```go
ctx = logs.WithRequestID(ctx, "req-123")
ctx = logs.WithTraceID(ctx, "trace-456")
ctx = logs.WithUserID(ctx, "user-789")
```

## Formatters

### TextFormatter

Default human-readable format:

```go
log := logs.New(&logs.Options{
    Formatter: &logs.TextFormatter{
        TimestampFormat: time.RFC3339,
        DisableColors:   false,
    },
})
```

### JSONFormatter

Machine-readable JSON output:

```go
log := logs.New(&logs.Options{
    Formatter: &logs.JSONFormatter{
        TimestampKey: "ts",
        MessageKey:   "message",
        LevelKey:     "severity",
    },
})
```

### PrettyFormatter

Development-friendly with colors and emojis:

```go
log := logs.New(&logs.Options{
    Formatter: &logs.PrettyFormatter{
        ShowTimestamp: true,
        ShowCaller:    true,
    },
})
```

## Hooks

Extend logging behavior with hooks:

```go
// Write errors to file
fileHook, _ := logs.NewFileHook("/var/log/errors.log",
    &logs.JSONFormatter{},
    logs.ErrorLevel, logs.FatalLevel, logs.PanicLevel,
)
log.AddHook(fileHook)

// Track metrics
metricsHook := logs.NewMetricsHook()
log.AddHook(metricsHook)
fmt.Println(metricsHook.Count(logs.ErrorLevel))

// Custom function hook
log.AddHook(logs.NewFuncHook(func(e *logs.Entry) {
    // Custom processing
}, logs.ErrorLevel))
```

## Sampling

Reduce log volume in high-throughput scenarios:

```go
// Log every 100th occurrence
log := logs.New(&logs.Options{
    Sampler: logs.NewCountSampler(100),
})

// Rate limit: max 10 per second per message
log := logs.New(&logs.Options{
    Sampler: logs.NewRateSampler(10, time.Second),
})

// Log first 5 occurrences only
log := logs.New(&logs.Options{
    Sampler: logs.NewFirstNSampler(5),
})

// Per-level sampling
sampler := logs.NewLevelSampler(nil).
    WithLevel(logs.DebugLevel, logs.NewCountSampler(10))
```

## Async Logging

Enable non-blocking logging for high throughput:

```go
log := logs.New(&logs.Options{
    AsyncBufferSize: 10000,  // Enable with buffer size
})
defer log.Close()  // Flush pending entries
```

## Builder API

Fluent interface for building entries:

```go
log.Build().
    Str("user", "alice").
    Int("items", 5).
    WithError(err).
    Info("order placed")

// With context
log.Ctx(ctx).
    Str("action", "update").
    Info("completed")

// Key-value pairs
log.F("user", "bob", "age", 30).Info("created")
```

## Error Helpers

Convenient error logging patterns:

```go
// Conditional logging
log.IfErr(err).With("user", id).Error("failed to create user")

// Log and wrap error
return log.WrapErr(err, "database query failed", logs.String("table", "users"))

// Check and log
if log.CheckErr(err, "operation failed") {
    return
}
```

## Configuration Reference

```go
opts := &logs.Options{
    Output:          os.Stdout,           // io.Writer for output
    Level:           logs.InfoLevel,      // Minimum log level
    Formatter:       &logs.TextFormatter{}, // Output format
    AddCaller:       true,                // Include file:line
    CallerDepth:     2,                   // Caller stack depth
    AddStack:        false,               // Stack traces for errors
    AsyncBufferSize: 0,                   // 0 = sync, >0 = async
    Hooks:           []logs.Hook{},       // Processing hooks
    Fields:          []logs.Field{},      // Default fields
    Sampler:         nil,                 // Rate limiting
}
```

## Further Reading

- [API Reference](../api-reference/logs.md) - Complete API documentation
- [Best Practices](../advanced/best-practices.md) - Production patterns
- [Performance Tuning](../advanced/performance-tuning.md) - Optimization guide
