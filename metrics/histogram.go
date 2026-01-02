package metrics

import (
	"math"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// Histogram samples observations and counts them in buckets.
type Histogram struct {
	name       string
	help       string
	labelNames []string
	buckets    []float64
	values     sync.Map
}

type histogramValue struct {
	labels     Labels
	buckets    []float64
	counts     []atomic.Uint64
	countTotal atomic.Uint64
	sumBits    atomic.Uint64
}

// NewHistogram creates a new histogram.
func NewHistogram(name, help string, buckets []float64, labelNames ...string) *Histogram {
	if len(buckets) == 0 {
		buckets = DefaultHistogramBuckets()
	}

	sort.Float64s(buckets)
	deduped := make([]float64, 0, len(buckets))
	for i, b := range buckets {
		if i == 0 || b != buckets[i-1] {
			deduped = append(deduped, b)
		}
	}

	return &Histogram{
		name:       name,
		help:       help,
		labelNames: labelNames,
		buckets:    deduped,
	}
}

func (h *Histogram) Name() string       { return h.name }
func (h *Histogram) Help() string       { return h.help }
func (h *Histogram) Type() MetricType   { return MetricTypeHistogram }
func (h *Histogram) LabelNames() []string { return h.labelNames }

// Observe adds an observation.
func (h *Histogram) Observe(value float64, labelValues ...string) {
	labels := h.makeLabels(labelValues)
	hash := labels.Hash()

	val, _ := h.values.LoadOrStore(hash, h.newHistogramValue(labels))
	hv := val.(*histogramValue)

	for i, bucket := range h.buckets {
		if value <= bucket {
			hv.counts[i].Add(1)
		}
	}

	hv.countTotal.Add(1)

	for {
		oldBits := hv.sumBits.Load()
		newSum := math.Float64frombits(oldBits) + value
		if hv.sumBits.CompareAndSwap(oldBits, math.Float64bits(newSum)) {
			break
		}
	}
}

func (h *Histogram) newHistogramValue(labels Labels) *histogramValue {
	return &histogramValue{
		labels:  labels,
		buckets: h.buckets,
		counts:  make([]atomic.Uint64, len(h.buckets)),
	}
}

// Collect returns all samples.
func (h *Histogram) Collect() []Sample {
	var samples []Sample
	now := time.Now()

	h.values.Range(func(_, value any) bool {
		hv := value.(*histogramValue)

		for i, bucket := range h.buckets {
			count := hv.counts[i].Load()

			bucketLabels := hv.labels.Merge(NewLabels("le", formatFloat(bucket)))
			samples = append(samples, Sample{
				Name:      h.name + "_bucket",
				Labels:    bucketLabels,
				Value:     float64(count),
				Timestamp: now,
			})
		}

		infLabels := hv.labels.Merge(NewLabels("le", "+Inf"))
		samples = append(samples, Sample{
			Name:      h.name + "_bucket",
			Labels:    infLabels,
			Value:     float64(hv.countTotal.Load()),
			Timestamp: now,
		})

		samples = append(samples, Sample{
			Name:      h.name + "_sum",
			Labels:    hv.labels,
			Value:     math.Float64frombits(hv.sumBits.Load()),
			Timestamp: now,
		})

		samples = append(samples, Sample{
			Name:      h.name + "_count",
			Labels:    hv.labels,
			Value:     float64(hv.countTotal.Load()),
			Timestamp: now,
		})

		return true
	})

	return samples
}

// Reset resets all histogram values.
func (h *Histogram) Reset() {
	h.values.Range(func(key, _ any) bool {
		h.values.Delete(key)
		return true
	})
}

func (h *Histogram) makeLabels(values []string) Labels {
	if len(h.labelNames) == 0 {
		return Labels{}
	}

	if len(values) != len(h.labelNames) {
		if len(values) < len(h.labelNames) {
			padded := make([]string, len(h.labelNames))
			copy(padded, values)
			values = padded
		} else {
			values = values[:len(h.labelNames)]
		}
	}

	pairs := make([]string, 0, len(h.labelNames)*2)
	for i, name := range h.labelNames {
		pairs = append(pairs, name, values[i])
	}
	return NewLabels(pairs...)
}

func formatFloat(f float64) string {
	if f == math.Inf(1) {
		return "+Inf"
	}
	return strconv.FormatFloat(f, 'g', -1, 64)
}
