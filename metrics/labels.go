package metrics

import (
	"sort"
	"strings"
)

// Labels is a sorted set of label key-value pairs.
type Labels struct {
	keys   []string
	values []string
}

// NewLabels creates labels from key-value pairs.
func NewLabels(pairs ...string) Labels {
	if len(pairs)%2 != 0 {
		pairs = pairs[:len(pairs)-1]
	}

	l := Labels{
		keys:   make([]string, 0, len(pairs)/2),
		values: make([]string, 0, len(pairs)/2),
	}

	for i := 0; i < len(pairs); i += 2 {
		l.keys = append(l.keys, pairs[i])
		l.values = append(l.values, pairs[i+1])
	}

	l.sort()
	return l
}

// LabelsFromMap creates labels from a map.
func LabelsFromMap(m map[string]string) Labels {
	l := Labels{
		keys:   make([]string, 0, len(m)),
		values: make([]string, 0, len(m)),
	}

	for k, v := range m {
		l.keys = append(l.keys, k)
		l.values = append(l.values, v)
	}

	l.sort()
	return l
}

func (l *Labels) sort() {
	if len(l.keys) <= 1 {
		return
	}

	indices := make([]int, len(l.keys))
	for i := range indices {
		indices[i] = i
	}

	sort.Slice(indices, func(i, j int) bool {
		return l.keys[indices[i]] < l.keys[indices[j]]
	})

	newKeys := make([]string, len(l.keys))
	newValues := make([]string, len(l.values))
	for i, idx := range indices {
		newKeys[i] = l.keys[idx]
		newValues[i] = l.values[idx]
	}

	l.keys = newKeys
	l.values = newValues
}

// Hash returns a unique string for the label set.
func (l Labels) Hash() string {
	if len(l.keys) == 0 {
		return ""
	}

	var sb strings.Builder
	for i, k := range l.keys {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(k)
		sb.WriteByte('=')
		sb.WriteString(l.values[i])
	}
	return sb.String()
}

// Keys returns label names.
func (l Labels) Keys() []string {
	return l.keys
}

// Values returns label values.
func (l Labels) Values() []string {
	return l.values
}

// Len returns the number of labels.
func (l Labels) Len() int {
	return len(l.keys)
}

// Get returns a label value.
func (l Labels) Get(key string) string {
	for i, k := range l.keys {
		if k == key {
			return l.values[i]
		}
	}
	return ""
}

// Merge combines two label sets.
func (l Labels) Merge(other Labels) Labels {
	merged := Labels{
		keys:   make([]string, 0, len(l.keys)+len(other.keys)),
		values: make([]string, 0, len(l.values)+len(other.values)),
	}

	merged.keys = append(merged.keys, l.keys...)
	merged.values = append(merged.values, l.values...)
	merged.keys = append(merged.keys, other.keys...)
	merged.values = append(merged.values, other.values...)

	merged.sort()
	return merged
}
