# trace API

Complete API documentation for the trace package.

**Import Path:** `github.com/kolosys/lumen/trace`

## Package Documentation

Package trace provides distributed tracing with W3C and Kolosys format support.


## Constants

**W3CTraceparentHeader, W3CTracestateHeader**

W3C Trace Context header names.


```go
const W3CTraceparentHeader = "traceparent"
const W3CTracestateHeader = "tracestate"
```

**TraceIDHeader, SpanIDHeader**

Simple header names for trace propagation.


```go
const TraceIDHeader = "X-Trace-ID"
const SpanIDHeader = "X-Span-ID"
```

## Variables

**ErrTracerClosed, ErrInvalidTraceID, ErrInvalidSpanID, ErrInvalidContext, ErrSamplerRejected, ErrExporterFailed**



```go
var ErrTracerClosed = errors.New("trace: tracer is closed")
var ErrInvalidTraceID = errors.New("trace: invalid trace ID")
var ErrInvalidSpanID = errors.New("trace: invalid span ID")
var ErrInvalidContext = errors.New("trace: invalid trace context")
var ErrSamplerRejected = errors.New("trace: span rejected by sampler")
var ErrExporterFailed = errors.New("trace: exporter failed")
```

## Types

### AlwaysSampler
AlwaysSampler always samples.

#### Example Usage

```go
// Create a new AlwaysSampler
alwayssampler := AlwaysSampler{

}
```

#### Type Definition

```go
type AlwaysSampler struct {
}
```

### Constructor Functions

### AlwaysSample

AlwaysSample returns a sampler that always samples.

```go
func AlwaysSample() *AlwaysSampler
```

**Parameters:**
  None

**Returns:**
- *AlwaysSampler

## Methods

### ShouldSample



```go
func (*ParentBasedSampler) ShouldSample(params SamplingParams) bool
```

**Parameters:**
- `params` (SamplingParams)

**Returns:**
- bool

### Attribute
Attribute is a key-value pair attached to a span.

#### Example Usage

```go
// Create a new Attribute
attribute := Attribute{
    Key: "example",
    Value: any{},
}
```

#### Type Definition

```go
type Attribute struct {
    Key string
    Value any
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Key | `string` |  |
| Value | `any` |  |

### Carrier
Carrier is an interface for propagation carriers (e.g., HTTP headers).

#### Example Usage

```go
// Example implementation of Carrier
type MyCarrier struct {
    // Add your fields here
}

func (m MyCarrier) Get(param1 string) string {
    // Implement your logic here
    return
}

func (m MyCarrier) Set(param1 string)  {
    // Implement your logic here
    return
}

func (m MyCarrier) Keys() []string {
    // Implement your logic here
    return
}


```

#### Type Definition

```go
type Carrier interface {
    Get(key string) string
    Set(key, value string)
    Keys() []string
}
```

## Methods

| Method | Description |
| ------ | ----------- |

### CompositePropagator
CompositePropagator combines multiple propagators.

#### Example Usage

```go
// Create a new CompositePropagator
compositepropagator := CompositePropagator{

}
```

#### Type Definition

```go
type CompositePropagator struct {
}
```

### Constructor Functions

### NewCompositePropagator

NewCompositePropagator creates a composite propagator.

```go
func NewCompositePropagator(propagators ...Propagator) *CompositePropagator
```

**Parameters:**
- `propagators` (...Propagator)

**Returns:**
- *CompositePropagator

## Methods

### Extract



```go
func (*CompositePropagator) Extract(ctx context.Context, carrier Carrier) context.Context
```

**Parameters:**
- `ctx` (context.Context)
- `carrier` (Carrier)

**Returns:**
- context.Context

### Inject



```go
func (*CompositePropagator) Inject(ctx context.Context, carrier Carrier)
```

**Parameters:**
- `ctx` (context.Context)
- `carrier` (Carrier)

**Returns:**
  None

### Event
Event is a timestamped annotation.

#### Example Usage

```go
// Create a new Event
event := Event{
    Name: "example",
    Timestamp: /* value */,
    Attributes: [],
}
```

#### Type Definition

```go
type Event struct {
    Name string
    Timestamp time.Time
    Attributes []Attribute
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Name | `string` |  |
| Timestamp | `time.Time` |  |
| Attributes | `[]Attribute` |  |

### Exporter
Exporter receives completed spans.

#### Example Usage

```go
// Example implementation of Exporter
type MyExporter struct {
    // Add your fields here
}

func (m MyExporter) Export(param1 *Span)  {
    // Implement your logic here
    return
}

func (m MyExporter) Close() error {
    // Implement your logic here
    return
}


```

#### Type Definition

```go
type Exporter interface {
    Export(span *Span)
    Close() error
}
```

## Methods

| Method | Description |
| ------ | ----------- |

### HeaderPropagator
HeaderPropagator implements customizable header-based propagation.

#### Example Usage

```go
// Create a new HeaderPropagator
headerpropagator := HeaderPropagator{
    TraceIDHeader: "example",
    SpanIDHeader: "example",
}
```

#### Type Definition

```go
type HeaderPropagator struct {
    TraceIDHeader string
    SpanIDHeader string
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| TraceIDHeader | `string` | TraceIDHeader is the header name for trace ID. Default: "X-Trace-ID" |
| SpanIDHeader | `string` | SpanIDHeader is the header name for span ID. Default: "X-Span-ID" |

## Methods

### Extract



```go
func (*CompositePropagator) Extract(ctx context.Context, carrier Carrier) context.Context
```

**Parameters:**
- `ctx` (context.Context)
- `carrier` (Carrier)

**Returns:**
- context.Context

### Inject



```go
func (*CompositePropagator) Inject(ctx context.Context, carrier Carrier)
```

**Parameters:**
- `ctx` (context.Context)
- `carrier` (Carrier)

**Returns:**
  None

### InMemoryExporter
InMemoryExporter collects spans in memory for testing.

#### Example Usage

```go
// Create a new InMemoryExporter
inmemoryexporter := InMemoryExporter{

}
```

#### Type Definition

```go
type InMemoryExporter struct {
}
```

### Constructor Functions

### NewInMemoryExporter

NewInMemoryExporter creates an in-memory exporter.

```go
func NewInMemoryExporter() *InMemoryExporter
```

**Parameters:**
  None

**Returns:**
- *InMemoryExporter

## Methods

### Clear

Clear removes all collected spans.

```go
func (*InMemoryExporter) Clear()
```

**Parameters:**
  None

**Returns:**
  None

### Close



```go
func (*InMemoryExporter) Close() error
```

**Parameters:**
  None

**Returns:**
- error

### Export



```go
func (*InMemoryExporter) Export(span *Span)
```

**Parameters:**
- `span` (*Span)

**Returns:**
  None

### Len

Len returns the number of collected spans.

```go
func (*InMemoryExporter) Len() int
```

**Parameters:**
  None

**Returns:**
- int

### Spans

Spans returns collected spans.

```go
func (*InMemoryExporter) Spans() []*spanData
```

**Parameters:**
  None

**Returns:**
- []*spanData

### MapCarrier
MapCarrier is a map-based carrier.

#### Example Usage

```go
// Example usage of MapCarrier
var value MapCarrier
// Initialize with appropriate value
```

#### Type Definition

```go
type MapCarrier map[string]string
```

## Methods

### Get



```go
func (MapCarrier) Get(key string) string
```

**Parameters:**
- `key` (string)

**Returns:**
- string

### Keys



```go
func (MapCarrier) Keys() []string
```

**Parameters:**
  None

**Returns:**
- []string

### Set



```go
func (MapCarrier) Set(key, value string)
```

**Parameters:**
- `key` (string)
- `value` (string)

**Returns:**
  None

### NeverSampler
NeverSampler never samples.

#### Example Usage

```go
// Create a new NeverSampler
neversampler := NeverSampler{

}
```

#### Type Definition

```go
type NeverSampler struct {
}
```

### Constructor Functions

### NeverSample

NeverSample returns a sampler that never samples.

```go
func NeverSample() *NeverSampler
```

**Parameters:**
  None

**Returns:**
- *NeverSampler

## Methods

### ShouldSample



```go
func (*ParentBasedSampler) ShouldSample(params SamplingParams) bool
```

**Parameters:**
- `params` (SamplingParams)

**Returns:**
- bool

### NopExporter
NopExporter discards all spans.

#### Example Usage

```go
// Create a new NopExporter
nopexporter := NopExporter{

}
```

#### Type Definition

```go
type NopExporter struct {
}
```

## Methods

### Close



```go
func (*InMemoryExporter) Close() error
```

**Parameters:**
  None

**Returns:**
- error

### Export



```go
func (*InMemoryExporter) Export(span *Span)
```

**Parameters:**
- `span` (*Span)

**Returns:**
  None

### Options
Options configures a Tracer.

#### Example Usage

```go
// Create a new Options
options := Options{
    ServiceName: "example",
    Sampler: Sampler{},
    Exporter: Exporter{},
    MaxSpansPerSecond: 42,
    PropagationFormat: "example",
    AsyncExport: true,
    AsyncBufferSize: 42,
}
```

#### Type Definition

```go
type Options struct {
    ServiceName string
    Sampler Sampler
    Exporter Exporter
    MaxSpansPerSecond int
    PropagationFormat string
    AsyncExport bool
    AsyncBufferSize int
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| ServiceName | `string` | ServiceName identifies the service in traces. |
| Sampler | `Sampler` | Sampler determines which spans to record. |
| Exporter | `Exporter` | Exporter receives completed spans. |
| MaxSpansPerSecond | `int` | MaxSpansPerSecond limits span creation rate (0 = unlimited). |
| PropagationFormat | `string` | PropagationFormat sets the context propagation format. Supports: "w3c", "kolosys", "both" (default: "both") |
| AsyncExport | `bool` | AsyncExport enables asynchronous span export. |
| AsyncBufferSize | `int` | AsyncBufferSize sets the async export buffer size. |

### ParentBasedSampler
ParentBasedSampler follows parent sampling decision.

#### Example Usage

```go
// Create a new ParentBasedSampler
parentbasedsampler := ParentBasedSampler{

}
```

#### Type Definition

```go
type ParentBasedSampler struct {
}
```

### Constructor Functions

### ParentBasedSample

ParentBasedSample returns a sampler that follows parent decisions.

```go
func ParentBasedSample(root Sampler) *ParentBasedSampler
```

**Parameters:**
- `root` (Sampler)

**Returns:**
- *ParentBasedSampler

## Methods

### ShouldSample



```go
func (*ParentBasedSampler) ShouldSample(params SamplingParams) bool
```

**Parameters:**
- `params` (SamplingParams)

**Returns:**
- bool

### Propagator
Propagator handles trace context injection and extraction.

#### Example Usage

```go
// Example implementation of Propagator
type MyPropagator struct {
    // Add your fields here
}

func (m MyPropagator) Inject(param1 context.Context, param2 Carrier)  {
    // Implement your logic here
    return
}

func (m MyPropagator) Extract(param1 context.Context, param2 Carrier) context.Context {
    // Implement your logic here
    return
}


```

#### Type Definition

```go
type Propagator interface {
    Inject(ctx context.Context, carrier Carrier)
    Extract(ctx context.Context, carrier Carrier) context.Context
}
```

## Methods

| Method | Description |
| ------ | ----------- |

### Constructor Functions

### DefaultPropagator

DefaultPropagator returns a propagator supporting both W3C and header formats.

```go
func DefaultPropagator() Propagator
```

**Parameters:**
  None

**Returns:**
- Propagator

### RatioSampler
RatioSampler samples a percentage of traces.

#### Example Usage

```go
// Create a new RatioSampler
ratiosampler := RatioSampler{

}
```

#### Type Definition

```go
type RatioSampler struct {
}
```

### Constructor Functions

### RatioSample

RatioSample returns a sampler that samples the given ratio of traces.

```go
func RatioSample(ratio float64) *RatioSampler
```

**Parameters:**
- `ratio` (float64)

**Returns:**
- *RatioSampler

## Methods

### ShouldSample



```go
func (*ParentBasedSampler) ShouldSample(params SamplingParams) bool
```

**Parameters:**
- `params` (SamplingParams)

**Returns:**
- bool

### Sampler
Sampler determines whether a span should be recorded.

#### Example Usage

```go
// Example implementation of Sampler
type MySampler struct {
    // Add your fields here
}

func (m MySampler) ShouldSample(param1 SamplingParams) bool {
    // Implement your logic here
    return
}


```

#### Type Definition

```go
type Sampler interface {
    ShouldSample(params SamplingParams) bool
}
```

## Methods

| Method | Description |
| ------ | ----------- |

### SamplingParams
SamplingParams provides data for sampling decisions.

#### Example Usage

```go
// Create a new SamplingParams
samplingparams := SamplingParams{
    TraceID: TraceID{},
    Name: "example",
    ParentID: SpanID{},
}
```

#### Type Definition

```go
type SamplingParams struct {
    TraceID TraceID
    Name string
    ParentID SpanID
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| TraceID | `TraceID` |  |
| Name | `string` |  |
| ParentID | `SpanID` |  |

### Span
Span represents a unit of work.

#### Example Usage

```go
// Create a new Span
span := Span{

}
```

#### Type Definition

```go
type Span struct {
}
```

### Constructor Functions

### SpanFromContext

SpanFromContext retrieves a span from context.

```go
func SpanFromContext(ctx context.Context) *Span
```

**Parameters:**
- `ctx` (context.Context)

**Returns:**
- *Span

### Start

Start creates a span using the default tracer.

```go
func Start(ctx context.Context, name string, opts ...SpanOption) (context.Context, *Span)
```

**Parameters:**
- `ctx` (context.Context)
- `name` (string)
- `opts` (...SpanOption)

**Returns:**
- context.Context
- *Span

## Methods

### AddEvent

AddEvent adds a timestamped event.

```go
func (*Span) AddEvent(name string, attrs ...Attribute)
```

**Parameters:**
- `name` (string)
- `attrs` (...Attribute)

**Returns:**
  None

### Attributes

Attributes returns a copy of span attributes.

```go
func (*Span) Attributes() []Attribute
```

**Parameters:**
  None

**Returns:**
- []Attribute

### Duration

Duration returns the span duration.

```go
func (*Span) Duration() time.Duration
```

**Parameters:**
  None

**Returns:**
- time.Duration

### End

End completes the span.

```go
func (*Span) End()
```

**Parameters:**
  None

**Returns:**
  None

### EndFunc

EndFunc returns a function suitable for defer.

```go
func (*Span) EndFunc(errPtr *error) func()
```

**Parameters:**
- `errPtr` (*error)

**Returns:**
- func()

### EndTime

EndTime returns the end time.

```go
func (*Span) EndTime() time.Time
```

**Parameters:**
  None

**Returns:**
- time.Time

### Events

Events returns a copy of span events.

```go
func (*Span) Events() []Event
```

**Parameters:**
  None

**Returns:**
- []Event

### IsSampled

IsSampled returns whether the span is sampled.

```go
func (*Span) IsSampled() bool
```

**Parameters:**
  None

**Returns:**
- bool

### Name

Name returns the span name.

```go
func (*Span) Name() string
```

**Parameters:**
  None

**Returns:**
- string

### ParentID

ParentID returns the parent span ID.

```go
func (*Span) ParentID() SpanID
```

**Parameters:**
  None

**Returns:**
- SpanID

### RecordError

RecordError records an error as an event and sets error status.

```go
func (*Span) RecordError(err error)
```

**Parameters:**
- `err` (error)

**Returns:**
  None

### SetAttribute

SetAttribute adds an attribute.

```go
func (*Span) SetAttribute(key string, value any)
```

**Parameters:**
- `key` (string)
- `value` (any)

**Returns:**
  None

### SetAttributes

SetAttributes adds multiple attributes.

```go
func (*Span) SetAttributes(attrs ...Attribute)
```

**Parameters:**
- `attrs` (...Attribute)

**Returns:**
  None

### SetStatus

SetStatus sets the span status.

```go
func (*Span) SetStatus(status SpanStatus, msg string)
```

**Parameters:**
- `status` (SpanStatus)
- `msg` (string)

**Returns:**
  None

### SpanID

SpanID returns the span ID.

```go
func (*Span) SpanID() SpanID
```

**Parameters:**
  None

**Returns:**
- SpanID

### StartTime

StartTime returns the start time.

```go
func (*Span) StartTime() time.Time
```

**Parameters:**
  None

**Returns:**
- time.Time

### Status

Status returns the span status.

```go
func (*Span) Status() SpanStatus
```

**Parameters:**
  None

**Returns:**
- SpanStatus

### StatusMessage

StatusMessage returns the status message.

```go
func (*Span) StatusMessage() string
```

**Parameters:**
  None

**Returns:**
- string

### TraceID

TraceID returns the trace ID.

```go
func (*Span) TraceID() TraceID
```

**Parameters:**
  None

**Returns:**
- TraceID

### SpanID
SpanID is an 8-byte span identifier.

#### Example Usage

```go
// Example usage of SpanID
var value SpanID
// Initialize with appropriate value
```

#### Type Definition

```go
type SpanID [*ast.BasicLit]byte
```

## Methods

### IsValid



```go
func (SpanID) IsValid() bool
```

**Parameters:**
  None

**Returns:**
- bool

### String



```go
func (SpanStatus) String() string
```

**Parameters:**
  None

**Returns:**
- string

### SpanOption
SpanOption configures span creation.

#### Example Usage

```go
// Example usage of SpanOption
var value SpanOption
// Initialize with appropriate value
```

#### Type Definition

```go
type SpanOption func(*Span)
```

### Constructor Functions

### WithAttributes

WithAttributes sets initial attributes.

```go
func WithAttributes(attrs ...Attribute) SpanOption
```

**Parameters:**
- `attrs` (...Attribute)

**Returns:**
- SpanOption

### WithStartTime

WithStartTime sets a custom start time.

```go
func WithStartTime(t time.Time) SpanOption
```

**Parameters:**
- `t` (time.Time)

**Returns:**
- SpanOption

### SpanStatus
SpanStatus represents span completion status.

#### Example Usage

```go
// Example usage of SpanStatus
var value SpanStatus
// Initialize with appropriate value
```

#### Type Definition

```go
type SpanStatus int
```

## Methods

### String



```go
func (SpanStatus) String() string
```

**Parameters:**
  None

**Returns:**
- string

### TraceContext
TraceContext holds W3C Trace Context data.

#### Example Usage

```go
// Create a new TraceContext
tracecontext := TraceContext{
    TraceID: TraceID{},
    SpanID: SpanID{},
    TraceFlags: byte{},
    TraceState: "example",
}
```

#### Type Definition

```go
type TraceContext struct {
    TraceID TraceID
    SpanID SpanID
    TraceFlags byte
    TraceState string
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| TraceID | `TraceID` |  |
| SpanID | `SpanID` |  |
| TraceFlags | `byte` |  |
| TraceState | `string` |  |

### Constructor Functions

### ParseHeaders

ParseHeaders parses X-Trace-ID and X-Span-ID headers.

```go
func ParseHeaders(traceID, spanID string) (*TraceContext, error)
```

**Parameters:**
- `traceID` (string)
- `spanID` (string)

**Returns:**
- *TraceContext
- error

### ParseW3CTraceparent

ParseW3CTraceparent parses a W3C traceparent header.

```go
func ParseW3CTraceparent(header string) (*TraceContext, error)
```

**Parameters:**
- `header` (string)

**Returns:**
- *TraceContext
- error

### TraceContextFromContext

TraceContextFromContext retrieves parsed trace context.

```go
func TraceContextFromContext(ctx context.Context) *TraceContext
```

**Parameters:**
- `ctx` (context.Context)

**Returns:**
- *TraceContext

## Methods

### FormatW3CTraceparent

FormatW3CTraceparent formats a TraceContext as W3C traceparent.

```go
func (*TraceContext) FormatW3CTraceparent() string
```

**Parameters:**
  None

**Returns:**
- string

### IsSampled

IsSampled returns whether the trace is sampled.

```go
func (*TraceContext) IsSampled() bool
```

**Parameters:**
  None

**Returns:**
- bool

### SetSampled

SetSampled sets the sampled flag.

```go
func (*TraceContext) SetSampled(sampled bool)
```

**Parameters:**
- `sampled` (bool)

**Returns:**
  None

### TraceID
TraceID is a 16-byte trace identifier.

#### Example Usage

```go
// Example usage of TraceID
var value TraceID
// Initialize with appropriate value
```

#### Type Definition

```go
type TraceID [*ast.BasicLit]byte
```

## Methods

### IsValid



```go
func (SpanID) IsValid() bool
```

**Parameters:**
  None

**Returns:**
- bool

### String



```go
func (SpanStatus) String() string
```

**Parameters:**
  None

**Returns:**
- string

### TraceIDRatioSampler
TraceIDRatioSampler samples based on trace ID for consistency.

#### Example Usage

```go
// Create a new TraceIDRatioSampler
traceidratiosampler := TraceIDRatioSampler{

}
```

#### Type Definition

```go
type TraceIDRatioSampler struct {
}
```

### Constructor Functions

### TraceIDRatioSample

TraceIDRatioSample returns a sampler based on trace ID.

```go
func TraceIDRatioSample(ratio float64) *TraceIDRatioSampler
```

**Parameters:**
- `ratio` (float64)

**Returns:**
- *TraceIDRatioSampler

## Methods

### ShouldSample



```go
func (*ParentBasedSampler) ShouldSample(params SamplingParams) bool
```

**Parameters:**
- `params` (SamplingParams)

**Returns:**
- bool

### Tracer
Tracer creates and manages spans.

#### Example Usage

```go
// Create a new Tracer
tracer := Tracer{

}
```

#### Type Definition

```go
type Tracer struct {
}
```

### Constructor Functions

### Default

Default returns the default tracer.

```go
func Default() *Tracer
```

**Parameters:**
  None

**Returns:**
- *Tracer

### New

New creates a new Tracer.

```go
func New(opts *Options) *Tracer
```

**Parameters:**
- `opts` (*Options)

**Returns:**
- *Tracer

## Methods

### Close

Close shuts down the tracer.

```go
func (*Tracer) Close() error
```

**Parameters:**
  None

**Returns:**
- error

### Start

Start creates a new span.

```go
func Start(ctx context.Context, name string, opts ...SpanOption) (context.Context, *Span)
```

**Parameters:**
- `ctx` (context.Context)
- `name` (string)
- `opts` (...SpanOption)

**Returns:**
- context.Context
- *Span

### W3CPropagator
W3CPropagator implements W3C Trace Context propagation.

#### Example Usage

```go
// Create a new W3CPropagator
w3cpropagator := W3CPropagator{

}
```

#### Type Definition

```go
type W3CPropagator struct {
}
```

## Methods

### Extract



```go
func (*CompositePropagator) Extract(ctx context.Context, carrier Carrier) context.Context
```

**Parameters:**
- `ctx` (context.Context)
- `carrier` (Carrier)

**Returns:**
- context.Context

### Inject



```go
func (*CompositePropagator) Inject(ctx context.Context, carrier Carrier)
```

**Parameters:**
- `ctx` (context.Context)
- `carrier` (Carrier)

**Returns:**
  None

### WriterExporter
WriterExporter writes spans as JSON to an io.Writer.

#### Example Usage

```go
// Create a new WriterExporter
writerexporter := WriterExporter{

}
```

#### Type Definition

```go
type WriterExporter struct {
}
```

### Constructor Functions

### NewWriterExporter

NewWriterExporter creates an exporter that writes to w.

```go
func NewWriterExporter(w io.Writer) *WriterExporter
```

**Parameters:**
- `w` (io.Writer)

**Returns:**
- *WriterExporter

## Methods

### Close



```go
func (*Tracer) Close() error
```

**Parameters:**
  None

**Returns:**
- error

### Export



```go
func (*InMemoryExporter) Export(span *Span)
```

**Parameters:**
- `span` (*Span)

**Returns:**
  None

## Functions

### ContextWithSpan
ContextWithSpan returns a context with the span attached.

```go
func ContextWithSpan(ctx context.Context, span *Span) context.Context
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `ctx` | `context.Context` | |
| `span` | `*Span` | |

**Returns:**
| Type | Description |
|------|-------------|
| `context.Context` | |

**Example:**

```go
// Example usage of ContextWithSpan
result := ContextWithSpan(/* parameters */)
```

### ContextWithTraceContext
ContextWithTraceContext returns a context with trace context attached.

```go
func ContextWithTraceContext(ctx context.Context, tc *TraceContext) context.Context
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `ctx` | `context.Context` | |
| `tc` | `*TraceContext` | |

**Returns:**
| Type | Description |
|------|-------------|
| `context.Context` | |

**Example:**

```go
// Example usage of ContextWithTraceContext
result := ContextWithTraceContext(/* parameters */)
```

### SetDefault
SetDefault sets the default tracer.

```go
func SetDefault(t *Tracer)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*Tracer` | |

**Returns:**
None

**Example:**

```go
// Example usage of SetDefault
result := SetDefault(/* parameters */)
```

## External Links

- [Package Overview](../packages/trace.md)
- [pkg.go.dev Documentation](https://pkg.go.dev/github.com/kolosys/lumen/trace)
- [Source Code](https://github.com/kolosys/lumen/tree/main/trace)
