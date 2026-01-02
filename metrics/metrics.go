// Package metrics provides metrics collection with Prometheus and push support.
package metrics

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// Metric is the interface all metric types implement.
type Metric interface {
	Name() string
	Help() string
	Type() MetricType
	LabelNames() []string
	Collect() []Sample
}

// MetricType identifies the metric type.
type MetricType int

const (
	MetricTypeCounter MetricType = iota
	MetricTypeGauge
	MetricTypeHistogram
)

func (t MetricType) String() string {
	switch t {
	case MetricTypeCounter:
		return "counter"
	case MetricTypeGauge:
		return "gauge"
	case MetricTypeHistogram:
		return "histogram"
	default:
		return "unknown"
	}
}

// Sample is a single metric observation.
type Sample struct {
	Name      string
	Labels    Labels
	Value     float64
	Timestamp time.Time
}

// Registry manages metric registration and collection.
type Registry struct {
	opts       *Options
	metrics    sync.Map
	pushCancel context.CancelFunc
	pushWg     sync.WaitGroup
	closed     atomic.Bool
	closeOnce  sync.Once
}

// NewRegistry creates a new metrics registry.
func NewRegistry(opts *Options) *Registry {
	if opts == nil {
		opts = &Options{}
	}
	opts.applyDefaults()

	r := &Registry{opts: opts}

	if opts.PushInterval > 0 && opts.PushExporter != nil {
		ctx, cancel := context.WithCancel(context.Background())
		r.pushCancel = cancel
		r.pushWg.Add(1)
		go r.pushLoop(ctx)
	}

	return r
}

// Register adds a metric to the registry.
func (r *Registry) Register(m Metric) error {
	if r.closed.Load() {
		return ErrRegistryClosed
	}

	_, loaded := r.metrics.LoadOrStore(m.Name(), m)
	if loaded {
		return ErrMetricExists
	}

	return nil
}

// Unregister removes a metric from the registry.
func (r *Registry) Unregister(name string) {
	r.metrics.Delete(name)
}

// Get retrieves a metric by name.
func (r *Registry) Get(name string) (Metric, error) {
	if m, ok := r.metrics.Load(name); ok {
		return m.(Metric), nil
	}
	return nil, ErrMetricNotFound
}

// Collect gathers all metric samples.
func (r *Registry) Collect() []Sample {
	var samples []Sample

	r.metrics.Range(func(_, value any) bool {
		m := value.(Metric)
		samples = append(samples, m.Collect()...)
		return true
	})

	return samples
}

// Close shuts down the registry.
func (r *Registry) Close() error {
	r.closeOnce.Do(func() {
		r.closed.Store(true)
		if r.pushCancel != nil {
			r.pushCancel()
			r.pushWg.Wait()
		}
	})
	return nil
}

func (r *Registry) pushLoop(ctx context.Context) {
	defer r.pushWg.Done()

	ticker := time.NewTicker(r.opts.PushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			samples := r.Collect()
			r.opts.PushExporter.Export(samples)
		}
	}
}

// Counter creates and registers a counter.
func (r *Registry) Counter(name, help string, labelNames ...string) *Counter {
	c := NewCounter(name, help, labelNames...)
	r.Register(c)
	return c
}

// Gauge creates and registers a gauge.
func (r *Registry) Gauge(name, help string, labelNames ...string) *Gauge {
	g := NewGauge(name, help, labelNames...)
	r.Register(g)
	return g
}

// Histogram creates and registers a histogram.
func (r *Registry) Histogram(name, help string, buckets []float64, labelNames ...string) *Histogram {
	h := NewHistogram(name, help, buckets, labelNames...)
	r.Register(h)
	return h
}

// Default registry
var defaultRegistry = NewRegistry(nil)

// DefaultRegistry returns the default registry.
func DefaultRegistry() *Registry {
	return defaultRegistry
}

// SetDefaultRegistry sets the default registry.
func SetDefaultRegistry(r *Registry) {
	defaultRegistry = r
}
