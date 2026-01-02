# Installation

## Prerequisites

- **Go 1.21** or later
- A Go module initialized in your project (`go mod init`)

## Install

```bash
go get github.com/kolosys/lumen
```

## Import Packages

Import the packages you need:

```go
import (
    "github.com/kolosys/lumen/logs"    // Structured logging
    "github.com/kolosys/lumen/metrics" // Metrics collection
    "github.com/kolosys/lumen/trace"   // Distributed tracing
)
```

## Verify Installation

```go
package main

import (
    "github.com/kolosys/lumen/logs"
)

func main() {
    log := logs.New(nil)
    log.Info("lumen installed successfully")
}
```

```bash
go run main.go
```

## Update

Update to the latest version:

```bash
go get -u github.com/kolosys/lumen
```

Update to a specific version:

```bash
go get github.com/kolosys/lumen@v1.2.3
```

## Development Setup

Clone for local development:

```bash
git clone https://github.com/kolosys/lumen.git
cd lumen
go mod download
go test -race ./...
```

## Troubleshooting

### Module Not Found

1. Verify `GOPATH` is set correctly
2. Check network access to GitHub
3. Run `go clean -modcache` and reinstall

### Private Repository

Configure Git for SSH or set GOPRIVATE:

```bash
export GOPRIVATE=github.com/kolosys/lumen
```

## Next Steps

- [Quick Start Guide](quick-start.md) - Basic usage examples
- [Core Concepts](../core-concepts/) - Detailed documentation
- [API Reference](../api-reference/) - Complete API docs
