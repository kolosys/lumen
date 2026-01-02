package trace

import (
	"hash/fnv"
	"sync/atomic"
)

// SamplingParams provides data for sampling decisions.
type SamplingParams struct {
	TraceID  TraceID
	Name     string
	ParentID SpanID
}

// Sampler determines whether a span should be recorded.
type Sampler interface {
	ShouldSample(params SamplingParams) bool
}

// AlwaysSampler always samples.
type AlwaysSampler struct{}

// AlwaysSample returns a sampler that always samples.
func AlwaysSample() *AlwaysSampler { return &AlwaysSampler{} }

func (*AlwaysSampler) ShouldSample(SamplingParams) bool { return true }

// NeverSampler never samples.
type NeverSampler struct{}

// NeverSample returns a sampler that never samples.
func NeverSample() *NeverSampler { return &NeverSampler{} }

func (*NeverSampler) ShouldSample(SamplingParams) bool { return false }

// RatioSampler samples a percentage of traces.
type RatioSampler struct {
	ratio     float64
	threshold uint32
	counter   atomic.Uint64
}

// RatioSample returns a sampler that samples the given ratio of traces.
func RatioSample(ratio float64) *RatioSampler {
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	return &RatioSampler{
		ratio:     ratio,
		threshold: uint32(ratio * 100),
	}
}

func (s *RatioSampler) ShouldSample(params SamplingParams) bool {
	count := s.counter.Add(1)
	return uint32(count%100) < s.threshold
}

// TraceIDRatioSampler samples based on trace ID for consistency.
type TraceIDRatioSampler struct {
	threshold uint32
}

// TraceIDRatioSample returns a sampler based on trace ID.
func TraceIDRatioSample(ratio float64) *TraceIDRatioSampler {
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	return &TraceIDRatioSampler{
		threshold: uint32(ratio * 0xFFFFFFFF),
	}
}

func (s *TraceIDRatioSampler) ShouldSample(params SamplingParams) bool {
	h := fnv.New32a()
	h.Write(params.TraceID[:])
	return h.Sum32() < s.threshold
}

// ParentBasedSampler follows parent sampling decision.
type ParentBasedSampler struct {
	root Sampler
}

// ParentBasedSample returns a sampler that follows parent decisions.
func ParentBasedSample(root Sampler) *ParentBasedSampler {
	return &ParentBasedSampler{root: root}
}

func (s *ParentBasedSampler) ShouldSample(params SamplingParams) bool {
	if params.ParentID.IsValid() {
		return true
	}
	return s.root.ShouldSample(params)
}
