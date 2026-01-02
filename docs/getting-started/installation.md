# Installation

This guide will help you install and set up lumen in your Go project.

## Prerequisites

Before installing lumen, ensure you have:

- **Go 1.21** or later installed
- A Go module initialized in your project (run `go mod init` if needed)
- Access to the GitHub repository (for private repositories)

## Installation Steps

### Step 1: Install the Package

Use `go get` to install lumen:

```bash
go get github.com/kolosys/lumen
```

This will download the package and add it to your `go.mod` file.

### Step 2: Import in Your Code

Import the package in your Go source files:

```go
import "github.com/kolosys/lumen"
```

### Multiple Packages

lumen includes several packages. Import the ones you need:

```go
// Package logs provides a high-performance, context-aware structured logging library.

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

import "github.com/kolosys/lumen/logs"
```

```go
// Package metrics provides metrics collection with Prometheus and push support.

import "github.com/kolosys/lumen/metrics"
```

```go
// Package trace provides distributed tracing with W3C and Kolosys format support.

import "github.com/kolosys/lumen/trace"
```

### Step 3: Verify Installation

Create a simple test file to verify the installation:

```go
package main

import (
    "fmt"
    "github.com/kolosys/lumen"
)

func main() {
    fmt.Println("lumen installed successfully!")
}
```

Run the test:

```bash
go run main.go
```

## Updating the Package

To update to the latest version:

```bash
go get -u github.com/kolosys/lumen
```

To update to a specific version:

```bash
go get github.com/kolosys/lumen@v1.2.3
```

## Installing a Specific Version

To install a specific version of the package:

```bash
go get github.com/kolosys/lumen@v1.0.0
```

Check available versions on the [GitHub releases page](https://github.com/kolosys/lumen/releases).

## Development Setup

If you want to contribute or modify the library:

### Clone the Repository

```bash
git clone https://github.com/kolosys/lumen.git
cd lumen
```

### Install Dependencies

```bash
go mod download
```

### Run Tests

```bash
go test ./...
```

## Troubleshooting

### Module Not Found

If you encounter a "module not found" error:

1. Ensure your `GOPATH` is set correctly
2. Check that you have network access to GitHub
3. Try running `go clean -modcache` and reinstall

### Private Repository Access

For private repositories, configure Git to use SSH or a personal access token:

```bash
git config --global url."git@github.com:".insteadOf "https://github.com/"
```

Or set up GOPRIVATE:

```bash
export GOPRIVATE=github.com/kolosys/lumen
```

## Next Steps

Now that you have lumen installed, check out the [Quick Start Guide](quick-start.md) to learn how to use it.

## Additional Resources

- [Quick Start Guide](quick-start.md)
- [API Reference](../api-reference/)
- [Examples](../examples/README.md)

