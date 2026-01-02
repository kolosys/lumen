package metrics

import "time"

// Options configures a Registry.
type Options struct {
	// Prefix is prepended to all metric names.
	Prefix string

	// DefaultLabels are added to all metrics.
	DefaultLabels map[string]string

	// HistogramBuckets defines default histogram bucket boundaries.
	HistogramBuckets []float64

	// PushInterval sets the interval for push exporters (0 = disabled).
	PushInterval time.Duration

	// PushExporter is the exporter for push-based metrics.
	PushExporter Exporter
}

func (o *Options) applyDefaults() {
	if o.HistogramBuckets == nil {
		o.HistogramBuckets = DefaultHistogramBuckets()
	}
}

// DefaultHistogramBuckets returns commonly used bucket boundaries.
func DefaultHistogramBuckets() []float64 {
	return []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}
}

// LinearBuckets creates n buckets of equal width.
func LinearBuckets(start, width float64, count int) []float64 {
	buckets := make([]float64, count)
	for i := range count {
		buckets[i] = start + width*float64(i)
	}
	return buckets
}

// ExponentialBuckets creates exponentially growing buckets.
func ExponentialBuckets(start, factor float64, count int) []float64 {
	buckets := make([]float64, count)
	buckets[0] = start
	for i := 1; i < count; i++ {
		buckets[i] = buckets[i-1] * factor
	}
	return buckets
}
