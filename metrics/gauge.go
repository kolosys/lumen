package metrics

import (
	"math"
	"sync"
	"sync/atomic"
	"time"
)

// Gauge is a metric that can go up and down.
type Gauge struct {
	name       string
	help       string
	labelNames []string
	values     sync.Map
}

type gaugeValue struct {
	labels Labels
	bits   atomic.Uint64
}

// NewGauge creates a new gauge.
func NewGauge(name, help string, labelNames ...string) *Gauge {
	return &Gauge{
		name:       name,
		help:       help,
		labelNames: labelNames,
	}
}

func (g *Gauge) Name() string       { return g.name }
func (g *Gauge) Help() string       { return g.help }
func (g *Gauge) Type() MetricType   { return MetricTypeGauge }
func (g *Gauge) LabelNames() []string { return g.labelNames }

// Set sets the gauge to a value.
func (g *Gauge) Set(value float64, labelValues ...string) {
	labels := g.makeLabels(labelValues)
	hash := labels.Hash()

	val, _ := g.values.LoadOrStore(hash, &gaugeValue{labels: labels})
	gv := val.(*gaugeValue)
	gv.bits.Store(math.Float64bits(value))
}

// Inc increments by 1.
func (g *Gauge) Inc(labelValues ...string) {
	g.Add(1, labelValues...)
}

// Dec decrements by 1.
func (g *Gauge) Dec(labelValues ...string) {
	g.Add(-1, labelValues...)
}

// Add adds a delta.
func (g *Gauge) Add(delta float64, labelValues ...string) {
	labels := g.makeLabels(labelValues)
	hash := labels.Hash()

	val, _ := g.values.LoadOrStore(hash, &gaugeValue{labels: labels})
	gv := val.(*gaugeValue)

	for {
		oldBits := gv.bits.Load()
		newVal := math.Float64frombits(oldBits) + delta
		if gv.bits.CompareAndSwap(oldBits, math.Float64bits(newVal)) {
			return
		}
	}
}

// Value returns the current value.
func (g *Gauge) Value(labelValues ...string) float64 {
	labels := g.makeLabels(labelValues)
	hash := labels.Hash()

	if val, ok := g.values.Load(hash); ok {
		return math.Float64frombits(val.(*gaugeValue).bits.Load())
	}
	return 0
}

// Collect returns all samples.
func (g *Gauge) Collect() []Sample {
	var samples []Sample
	now := time.Now()

	g.values.Range(func(_, value any) bool {
		gv := value.(*gaugeValue)
		samples = append(samples, Sample{
			Name:      g.name,
			Labels:    gv.labels,
			Value:     math.Float64frombits(gv.bits.Load()),
			Timestamp: now,
		})
		return true
	})

	return samples
}

// Reset resets all gauge values.
func (g *Gauge) Reset() {
	g.values.Range(func(key, _ any) bool {
		g.values.Delete(key)
		return true
	})
}

func (g *Gauge) makeLabels(values []string) Labels {
	if len(g.labelNames) == 0 {
		return Labels{}
	}

	if len(values) != len(g.labelNames) {
		if len(values) < len(g.labelNames) {
			padded := make([]string, len(g.labelNames))
			copy(padded, values)
			values = padded
		} else {
			values = values[:len(g.labelNames)]
		}
	}

	pairs := make([]string, 0, len(g.labelNames)*2)
	for i, name := range g.labelNames {
		pairs = append(pairs, name, values[i])
	}
	return NewLabels(pairs...)
}
