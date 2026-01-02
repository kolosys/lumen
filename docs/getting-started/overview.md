# Overview

A developer-first observability library for Go. Zero dependencies, high performance.

## About lumen

This documentation provides comprehensive guidance for using lumen, a Go library designed to help you build better software.

## Project Information

- **Repository**: [https://github.com/kolosys/lumen](https://github.com/kolosys/lumen)
- **Import Path**: `github.com/kolosys/lumen`
- **License**: MIT
- **Version**: latest

## What You'll Find Here

This documentation is organized into several sections to help you find what you need:

- **[Getting Started](../getting-started/)** - Installation instructions and quick start guides
- **[Core Concepts](../core-concepts/)** - Fundamental concepts and architecture details
- **[Advanced Topics](../advanced/)** - Performance tuning and advanced usage patterns
- **[API Reference](../api-reference/)** - Complete API reference documentation
- **[Examples](../examples/)** - Working code examples and tutorials

## Project Features

lumen provides:
- **logs** - Package logs provides a high-performance, context-aware structured logging library.

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

- **metrics** - Package metrics provides metrics collection with Prometheus and push support.

- **trace** - Package trace provides distributed tracing with W3C and Kolosys format support.


## Quick Links

- [Installation Guide](installation.md)
- [Quick Start Guide](quick-start.md)
- [API Reference](../api-reference/)
- [Examples](../examples/README.md)

## Community & Support

- **GitHub Issues**: [https://github.com/kolosys/lumen/issues](https://github.com/kolosys/lumen/issues)
- **Discussions**: [https://github.com/kolosys/lumen/discussions](https://github.com/kolosys/lumen/discussions)
- **Repository Owner**: [kolosys](https://github.com/kolosys)

## Getting Help

If you encounter any issues or have questions:

1. Check the [API Reference](../api-reference/) for detailed documentation
2. Browse the [Examples](../examples/README.md) for common use cases
3. Search existing [GitHub Issues](https://github.com/kolosys/lumen/issues)
4. Open a new issue if you've found a bug or have a feature request

## Next Steps

Ready to get started? Head over to the [Installation Guide](installation.md) to begin using lumen.

