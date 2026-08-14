# ionadapt API

Complete API documentation for the ionadapt package.

**Import Path:** `github.com/kolosys/lumen/ionadapt`

## Package Documentation

Package ionadapt adapts lumen logs and traces to [observe] interfaces.

Ion stays zero-dependency. This package is the optional bridge:

	circuit.WithLogger(ionadapt.Logger(logs.NewNamed("circuit")))


## Functions

### Logger
Logger wraps l as an [observe.Logger]. Nil l becomes [observe.NopLogger]. Alternating kv pairs become [logs.Field] values: keys that are not strings are fmt.Sprint'd; values are always fmt.Sprint'd. [observe.Logger.Error] adds a logs "error" field from the error argument.

```go
func Logger(l *logs.Logger) observe.Logger
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `l` | `*logs.Logger` | |

**Returns:**
| Type | Description |
|------|-------------|
| `observe.Logger` | |

**Example:**

```go
// Example usage of Logger
result := Logger(/* parameters */)
```

### Metrics
Metrics returns a no-op [observe.Metrics]. Lumen's Registry requires pre-registered names and a fixed label set. ion/observe.Metrics is fire-and-forget with dynamic names and kv pairs. A thin wrap would invent unbounded series or drop labels, so this adapter does not pretend. The registry argument is reserved for a future honest map.

```go
func Metrics(_ *metrics.Registry) observe.Metrics
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `_` | `*metrics.Registry` | |

**Returns:**
| Type | Description |
|------|-------------|
| `observe.Metrics` | |

**Example:**

```go
// Example usage of Metrics
result := Metrics(/* parameters */)
```

### Tracer
Tracer wraps t as an [observe.Tracer]. Nil t becomes [observe.NopTracer]. Start creates a lumen span with kv pairs as attributes. The finish func records a non-nil error on the span, then calls End.

```go
func Tracer(t *trace.Tracer) observe.Tracer
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*trace.Tracer` | |

**Returns:**
| Type | Description |
|------|-------------|
| `observe.Tracer` | |

**Example:**

```go
// Example usage of Tracer
result := Tracer(/* parameters */)
```

## External Links

- [Package Overview](../packages/ionadapt.md)
- [pkg.go.dev Documentation](https://pkg.go.dev/github.com/kolosys/lumen/ionadapt)
- [Source Code](https://github.com/kolosys/lumen/tree/main/ionadapt)
