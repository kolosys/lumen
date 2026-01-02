# metrics API

Complete API documentation for the metrics package.

**Import Path:** `github.com/kolosys/lumen/metrics`

## Package Documentation

Package metrics provides metrics collection with Prometheus and push support.


## Variables

**ErrRegistryClosed, ErrMetricNotFound, ErrMetricExists, ErrInvalidMetricName, ErrInvalidLabelName, ErrLabelMismatch, ErrExporterFailed**



```go
var ErrRegistryClosed = errors.New("metrics: registry is closed")
var ErrMetricNotFound = errors.New("metrics: metric not found")
var ErrMetricExists = errors.New("metrics: metric already exists")
var ErrInvalidMetricName = errors.New("metrics: invalid metric name")
var ErrInvalidLabelName = errors.New("metrics: invalid label name")
var ErrLabelMismatch = errors.New("metrics: label names do not match")
var ErrExporterFailed = errors.New("metrics: exporter failed")
```

## Types

### Counter
Counter is a cumulative metric that only increases.

#### Example Usage

```go
// Create a new Counter
counter := Counter{

}
```

#### Type Definition

```go
type Counter struct {
}
```

### Constructor Functions

### NewCounter

NewCounter creates a new counter.

```go
func NewCounter(name, help string, labelNames ...string) *Counter
```

**Parameters:**
- `name` (string)
- `help` (string)
- `labelNames` (...string)

**Returns:**
- *Counter

## Methods

### Add

Add increments by the given value.

```go
func (*Gauge) Add(delta float64, labelValues ...string)
```

**Parameters:**
- `delta` (float64)
- `labelValues` (...string)

**Returns:**
  None

### Collect

Collect returns all samples.

```go
func (*Gauge) Collect() []Sample
```

**Parameters:**
  None

**Returns:**
- []Sample

### Help



```go
func (*Histogram) Help() string
```

**Parameters:**
  None

**Returns:**
- string

### Inc

Inc increments by 1.

```go
func (*Gauge) Inc(labelValues ...string)
```

**Parameters:**
- `labelValues` (...string)

**Returns:**
  None

### LabelNames



```go
func (*Gauge) LabelNames() []string
```

**Parameters:**
  None

**Returns:**
- []string

### Name



```go
func (*Gauge) Name() string
```

**Parameters:**
  None

**Returns:**
- string

### Reset

Reset resets all counter values.

```go
func (*Counter) Reset()
```

**Parameters:**
  None

**Returns:**
  None

### Type



```go
func (*Counter) Type() MetricType
```

**Parameters:**
  None

**Returns:**
- MetricType

### Value

Value returns the current value for the given labels.

```go
func (*Gauge) Value(labelValues ...string) float64
```

**Parameters:**
- `labelValues` (...string)

**Returns:**
- float64

### Exporter
Exporter exports metrics.

#### Example Usage

```go
// Example implementation of Exporter
type MyExporter struct {
    // Add your fields here
}

func (m MyExporter) Export(param1 []Sample)  {
    // Implement your logic here
    return
}


```

#### Type Definition

```go
type Exporter interface {
    Export(samples []Sample)
}
```

## Methods

| Method | Description |
| ------ | ----------- |

### Gauge
Gauge is a metric that can go up and down.

#### Example Usage

```go
// Create a new Gauge
gauge := Gauge{

}
```

#### Type Definition

```go
type Gauge struct {
}
```

### Constructor Functions

### NewGauge

NewGauge creates a new gauge.

```go
func NewGauge(name, help string, labelNames ...string) *Gauge
```

**Parameters:**
- `name` (string)
- `help` (string)
- `labelNames` (...string)

**Returns:**
- *Gauge

## Methods

### Add

Add adds a delta.

```go
func (*Counter) Add(delta float64, labelValues ...string)
```

**Parameters:**
- `delta` (float64)
- `labelValues` (...string)

**Returns:**
  None

### Collect

Collect returns all samples.

```go
func (*Histogram) Collect() []Sample
```

**Parameters:**
  None

**Returns:**
- []Sample

### Dec

Dec decrements by 1.

```go
func (*Gauge) Dec(labelValues ...string)
```

**Parameters:**
- `labelValues` (...string)

**Returns:**
  None

### Help



```go
func (*Counter) Help() string
```

**Parameters:**
  None

**Returns:**
- string

### Inc

Inc increments by 1.

```go
func (*Gauge) Inc(labelValues ...string)
```

**Parameters:**
- `labelValues` (...string)

**Returns:**
  None

### LabelNames



```go
func (*Gauge) LabelNames() []string
```

**Parameters:**
  None

**Returns:**
- []string

### Name



```go
func (*Histogram) Name() string
```

**Parameters:**
  None

**Returns:**
- string

### Reset

Reset resets all gauge values.

```go
func (*Gauge) Reset()
```

**Parameters:**
  None

**Returns:**
  None

### Set

Set sets the gauge to a value.

```go
func (*Gauge) Set(value float64, labelValues ...string)
```

**Parameters:**
- `value` (float64)
- `labelValues` (...string)

**Returns:**
  None

### Type



```go
func (*Gauge) Type() MetricType
```

**Parameters:**
  None

**Returns:**
- MetricType

### Value

Value returns the current value.

```go
func (*Gauge) Value(labelValues ...string) float64
```

**Parameters:**
- `labelValues` (...string)

**Returns:**
- float64

### Histogram
Histogram samples observations and counts them in buckets.

#### Example Usage

```go
// Create a new Histogram
histogram := Histogram{

}
```

#### Type Definition

```go
type Histogram struct {
}
```

### Constructor Functions

### NewHistogram

NewHistogram creates a new histogram.

```go
func NewHistogram(name, help string, buckets []float64, labelNames ...string) *Histogram
```

**Parameters:**
- `name` (string)
- `help` (string)
- `buckets` ([]float64)
- `labelNames` (...string)

**Returns:**
- *Histogram

## Methods

### Collect

Collect returns all samples.

```go
func (*Histogram) Collect() []Sample
```

**Parameters:**
  None

**Returns:**
- []Sample

### Help



```go
func (*Counter) Help() string
```

**Parameters:**
  None

**Returns:**
- string

### LabelNames



```go
func (*Histogram) LabelNames() []string
```

**Parameters:**
  None

**Returns:**
- []string

### Name



```go
func (*Gauge) Name() string
```

**Parameters:**
  None

**Returns:**
- string

### Observe

Observe adds an observation.

```go
func (*Histogram) Observe(value float64, labelValues ...string)
```

**Parameters:**
- `value` (float64)
- `labelValues` (...string)

**Returns:**
  None

### Reset

Reset resets all histogram values.

```go
func (*Counter) Reset()
```

**Parameters:**
  None

**Returns:**
  None

### Type



```go
func (*Gauge) Type() MetricType
```

**Parameters:**
  None

**Returns:**
- MetricType

### Labels
Labels is a sorted set of label key-value pairs.

#### Example Usage

```go
// Create a new Labels
labels := Labels{

}
```

#### Type Definition

```go
type Labels struct {
}
```

### Constructor Functions

### LabelsFromMap

LabelsFromMap creates labels from a map.

```go
func LabelsFromMap(m map[string]string) Labels
```

**Parameters:**
- `m` (map[string]string)

**Returns:**
- Labels

### NewLabels

NewLabels creates labels from key-value pairs.

```go
func NewLabels(pairs ...string) Labels
```

**Parameters:**
- `pairs` (...string)

**Returns:**
- Labels

## Methods

### Get

Get returns a label value.

```go
func (Labels) Get(key string) string
```

**Parameters:**
- `key` (string)

**Returns:**
- string

### Hash

Hash returns a unique string for the label set.

```go
func (Labels) Hash() string
```

**Parameters:**
  None

**Returns:**
- string

### Keys

Keys returns label names.

```go
func (Labels) Keys() []string
```

**Parameters:**
  None

**Returns:**
- []string

### Len

Len returns the number of labels.

```go
func (Labels) Len() int
```

**Parameters:**
  None

**Returns:**
- int

### Merge

Merge combines two label sets.

```go
func (Labels) Merge(other Labels) Labels
```

**Parameters:**
- `other` (Labels)

**Returns:**
- Labels

### Values

Values returns label values.

```go
func (Labels) Values() []string
```

**Parameters:**
  None

**Returns:**
- []string

### Metric
Metric is the interface all metric types implement.

#### Example Usage

```go
// Example implementation of Metric
type MyMetric struct {
    // Add your fields here
}

func (m MyMetric) Name() string {
    // Implement your logic here
    return
}

func (m MyMetric) Help() string {
    // Implement your logic here
    return
}

func (m MyMetric) Type() MetricType {
    // Implement your logic here
    return
}

func (m MyMetric) LabelNames() []string {
    // Implement your logic here
    return
}

func (m MyMetric) Collect() []Sample {
    // Implement your logic here
    return
}


```

#### Type Definition

```go
type Metric interface {
    Name() string
    Help() string
    Type() MetricType
    LabelNames() []string
    Collect() []Sample
}
```

## Methods

| Method | Description |
| ------ | ----------- |

### MetricType
MetricType identifies the metric type.

#### Example Usage

```go
// Example usage of MetricType
var value MetricType
// Initialize with appropriate value
```

#### Type Definition

```go
type MetricType int
```

## Methods

### String



```go
func (MetricType) String() string
```

**Parameters:**
  None

**Returns:**
- string

### NopExporter
NopExporter discards metrics.

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

### Export



```go
func (NopExporter) Export([]Sample)
```

**Parameters:**
- `` ([]Sample)

**Returns:**
  None

### Options
Options configures a Registry.

#### Example Usage

```go
// Create a new Options
options := Options{
    Prefix: "example",
    DefaultLabels: map[],
    HistogramBuckets: [],
    PushInterval: /* value */,
    PushExporter: Exporter{},
}
```

#### Type Definition

```go
type Options struct {
    Prefix string
    DefaultLabels map[string]string
    HistogramBuckets []float64
    PushInterval time.Duration
    PushExporter Exporter
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Prefix | `string` | Prefix is prepended to all metric names. |
| DefaultLabels | `map[string]string` | DefaultLabels are added to all metrics. |
| HistogramBuckets | `[]float64` | HistogramBuckets defines default histogram bucket boundaries. |
| PushInterval | `time.Duration` | PushInterval sets the interval for push exporters (0 = disabled). |
| PushExporter | `Exporter` | PushExporter is the exporter for push-based metrics. |

### Registry
Registry manages metric registration and collection.

#### Example Usage

```go
// Create a new Registry
registry := Registry{

}
```

#### Type Definition

```go
type Registry struct {
}
```

### Constructor Functions

### DefaultRegistry

DefaultRegistry returns the default registry.

```go
func DefaultRegistry() *Registry
```

**Parameters:**
  None

**Returns:**
- *Registry

### NewRegistry

NewRegistry creates a new metrics registry.

```go
func NewRegistry(opts *Options) *Registry
```

**Parameters:**
- `opts` (*Options)

**Returns:**
- *Registry

## Methods

### Close

Close shuts down the registry.

```go
func (*Registry) Close() error
```

**Parameters:**
  None

**Returns:**
- error

### Collect

Collect gathers all metric samples.

```go
func (*Counter) Collect() []Sample
```

**Parameters:**
  None

**Returns:**
- []Sample

### Counter

Counter creates and registers a counter.

```go
func (*Registry) Counter(name, help string, labelNames ...string) *Counter
```

**Parameters:**
- `name` (string)
- `help` (string)
- `labelNames` (...string)

**Returns:**
- *Counter

### Gauge

Gauge creates and registers a gauge.

```go
func (*Registry) Gauge(name, help string, labelNames ...string) *Gauge
```

**Parameters:**
- `name` (string)
- `help` (string)
- `labelNames` (...string)

**Returns:**
- *Gauge

### Get

Get retrieves a metric by name.

```go
func (*Registry) Get(name string) (Metric, error)
```

**Parameters:**
- `name` (string)

**Returns:**
- Metric
- error

### Histogram

Histogram creates and registers a histogram.

```go
func (*Registry) Histogram(name, help string, buckets []float64, labelNames ...string) *Histogram
```

**Parameters:**
- `name` (string)
- `help` (string)
- `buckets` ([]float64)
- `labelNames` (...string)

**Returns:**
- *Histogram

### Register

Register adds a metric to the registry.

```go
func (*Registry) Register(m Metric) error
```

**Parameters:**
- `m` (Metric)

**Returns:**
- error

### Unregister

Unregister removes a metric from the registry.

```go
func (*Registry) Unregister(name string)
```

**Parameters:**
- `name` (string)

**Returns:**
  None

### Sample
Sample is a single metric observation.

#### Example Usage

```go
// Create a new Sample
sample := Sample{
    Name: "example",
    Labels: Labels{},
    Value: 3.14,
    Timestamp: /* value */,
}
```

#### Type Definition

```go
type Sample struct {
    Name string
    Labels Labels
    Value float64
    Timestamp time.Time
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Name | `string` |  |
| Labels | `Labels` |  |
| Value | `float64` |  |
| Timestamp | `time.Time` |  |

## Functions

### DefaultHTTPHandler
DefaultHTTPHandler returns an http.Handler using the default registry.

```go
func DefaultHTTPHandler() http.Handler
```

**Parameters:**
None

**Returns:**
| Type | Description |
|------|-------------|
| `http.Handler` | |

**Example:**

```go
// Example usage of DefaultHTTPHandler
result := DefaultHTTPHandler(/* parameters */)
```

### DefaultHistogramBuckets
DefaultHistogramBuckets returns commonly used bucket boundaries.

```go
func DefaultHistogramBuckets() []float64
```

**Parameters:**
None

**Returns:**
| Type | Description |
|------|-------------|
| `[]float64` | |

**Example:**

```go
// Example usage of DefaultHistogramBuckets
result := DefaultHistogramBuckets(/* parameters */)
```

### ExponentialBuckets
ExponentialBuckets creates exponentially growing buckets.

```go
func ExponentialBuckets(start, factor float64, count int) []float64
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `start` | `float64` | |
| `factor` | `float64` | |
| `count` | `int` | |

**Returns:**
| Type | Description |
|------|-------------|
| `[]float64` | |

**Example:**

```go
// Example usage of ExponentialBuckets
result := ExponentialBuckets(/* parameters */)
```

### HTTPHandler
HTTPHandler returns an http.Handler for the Prometheus endpoint.

```go
func HTTPHandler(registry *Registry) http.Handler
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `registry` | `*Registry` | |

**Returns:**
| Type | Description |
|------|-------------|
| `http.Handler` | |

**Example:**

```go
// Example usage of HTTPHandler
result := HTTPHandler(/* parameters */)
```

### LinearBuckets
LinearBuckets creates n buckets of equal width.

```go
func LinearBuckets(start, width float64, count int) []float64
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `start` | `float64` | |
| `width` | `float64` | |
| `count` | `int` | |

**Returns:**
| Type | Description |
|------|-------------|
| `[]float64` | |

**Example:**

```go
// Example usage of LinearBuckets
result := LinearBuckets(/* parameters */)
```

### SetDefaultRegistry
SetDefaultRegistry sets the default registry.

```go
func SetDefaultRegistry(r *Registry)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `r` | `*Registry` | |

**Returns:**
None

**Example:**

```go
// Example usage of SetDefaultRegistry
result := SetDefaultRegistry(/* parameters */)
```

### WritePrometheus
WritePrometheus writes samples in Prometheus text format.

```go
func WritePrometheus(w io.Writer, samples []Sample)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `w` | `io.Writer` | |
| `samples` | `[]Sample` | |

**Returns:**
None

**Example:**

```go
// Example usage of WritePrometheus
result := WritePrometheus(/* parameters */)
```

## External Links

- [Package Overview](../packages/metrics.md)
- [pkg.go.dev Documentation](https://pkg.go.dev/github.com/kolosys/lumen/metrics)
- [Source Code](https://github.com/kolosys/lumen/tree/main/metrics)
