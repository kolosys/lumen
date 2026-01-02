package metrics

import (
	"sync"
	"sync/atomic"
	"time"
)

// Counter is a cumulative metric that only increases.
type Counter struct {
	name       string
	help       string
	labelNames []string
	values     sync.Map
}

type counterValue struct {
	labels Labels
	value  atomic.Uint64
}

// NewCounter creates a new counter.
func NewCounter(name, help string, labelNames ...string) *Counter {
	return &Counter{
		name:       name,
		help:       help,
		labelNames: labelNames,
	}
}

func (c *Counter) Name() string       { return c.name }
func (c *Counter) Help() string       { return c.help }
func (c *Counter) Type() MetricType   { return MetricTypeCounter }
func (c *Counter) LabelNames() []string { return c.labelNames }

// Inc increments by 1.
func (c *Counter) Inc(labelValues ...string) {
	c.Add(1, labelValues...)
}

// Add increments by the given value.
func (c *Counter) Add(delta float64, labelValues ...string) {
	if delta < 0 {
		return
	}

	labels := c.makeLabels(labelValues)
	hash := labels.Hash()

	val, _ := c.values.LoadOrStore(hash, &counterValue{labels: labels})
	cv := val.(*counterValue)

	cv.value.Add(uint64(delta * 1000000))
}

// Value returns the current value for the given labels.
func (c *Counter) Value(labelValues ...string) float64 {
	labels := c.makeLabels(labelValues)
	hash := labels.Hash()

	if val, ok := c.values.Load(hash); ok {
		return float64(val.(*counterValue).value.Load()) / 1000000
	}
	return 0
}

// Collect returns all samples.
func (c *Counter) Collect() []Sample {
	var samples []Sample
	now := time.Now()

	c.values.Range(func(_, value any) bool {
		cv := value.(*counterValue)
		samples = append(samples, Sample{
			Name:      c.name,
			Labels:    cv.labels,
			Value:     float64(cv.value.Load()) / 1000000,
			Timestamp: now,
		})
		return true
	})

	return samples
}

// Reset resets all counter values.
func (c *Counter) Reset() {
	c.values.Range(func(key, _ any) bool {
		c.values.Delete(key)
		return true
	})
}

func (c *Counter) makeLabels(values []string) Labels {
	if len(c.labelNames) == 0 {
		return Labels{}
	}

	if len(values) != len(c.labelNames) {
		if len(values) < len(c.labelNames) {
			padded := make([]string, len(c.labelNames))
			copy(padded, values)
			values = padded
		} else {
			values = values[:len(c.labelNames)]
		}
	}

	pairs := make([]string, 0, len(c.labelNames)*2)
	for i, name := range c.labelNames {
		pairs = append(pairs, name, values[i])
	}
	return NewLabels(pairs...)
}
