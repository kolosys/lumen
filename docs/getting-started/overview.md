# Overview

A developer-first observability library for Go. Zero dependencies, high performance.

## About lumen

lumen is a comprehensive observability library providing structured logging, metrics collection, and distributed tracing. Built for production workloads with zero-allocation hot paths and full context propagation support.

## Project Information

- **Repository**: [https://github.com/kolosys/lumen](https://github.com/kolosys/lumen)
- **Import Path**: `github.com/kolosys/lumen`
- **License**: MIT
- **Go Version**: 1.21+

## Packages

### logs

High-performance, context-aware structured logging.

```go
import "github.com/kolosys/lumen/logs"

log := logs.New(nil)
log.Info("server started", logs.Int("port", 8080))
```

**Key Features:**
- Zero-allocation hot paths using sync.Pool
- Multiple formatters (Text, JSON, Pretty)
- Named loggers with hierarchical naming
- Sampling for high-volume logs
- Async logging support
- Hook system for extensibility
- Fluent builder API

### metrics

Metrics collection with Prometheus export support.

```go
import "github.com/kolosys/lumen/metrics"

registry := metrics.NewRegistry(nil)
counter := registry.Counter("http_requests_total", "Total HTTP requests", "method", "status")
counter.Inc("GET", "200")
```

**Key Features:**
- Counter, Gauge, and Histogram types
- Prometheus text format export
- Label support with efficient hashing
- Push-based export support
- Thread-safe atomic operations

### trace

Distributed tracing with W3C Trace Context support.

```go
import "github.com/kolosys/lumen/trace"

ctx, span := trace.Start(ctx, "operation")
defer span.End()
span.SetAttribute("user.id", userID)
```

**Key Features:**
- W3C Trace Context propagation
- Configurable sampling strategies
- Async span export
- Span events and attributes
- Error recording

## What You'll Find Here

- **[Getting Started](../getting-started/)** - Installation and quick start guides
- **[Core Concepts](../core-concepts/)** - Detailed package documentation
- **[Advanced Topics](../advanced/)** - Performance tuning and best practices
- **[API Reference](../api-reference/)** - Complete API documentation

## Next Steps

Ready to get started? Head over to the [Installation Guide](installation.md) to begin using lumen.
