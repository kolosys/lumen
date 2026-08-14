// Package ionadapt adapts lumen logs and traces to [observe] interfaces.
//
// Ion stays zero-dependency. This package is the optional bridge:
//
//	circuit.WithLogger(ionadapt.Logger(logs.NewNamed("circuit")))
package ionadapt

import (
	"fmt"

	"github.com/kolosys/ion/observe"
	"github.com/kolosys/lumen/logs"
	"github.com/kolosys/lumen/metrics"
	"github.com/kolosys/lumen/trace"
)

// Logger wraps l as an [observe.Logger]. Nil l becomes [observe.NopLogger].
//
// Alternating kv pairs become [logs.Field] values: keys that are not strings
// are fmt.Sprint'd; values are always fmt.Sprint'd. [observe.Logger.Error]
// adds a logs "error" field from the error argument.
func Logger(l *logs.Logger) observe.Logger {
	if l == nil {
		return observe.NopLogger{}
	}
	return logger{l: l}
}

// Tracer wraps t as an [observe.Tracer]. Nil t becomes [observe.NopTracer].
//
// Start creates a lumen span with kv pairs as attributes. The finish func
// records a non-nil error on the span, then calls End.
func Tracer(t *trace.Tracer) observe.Tracer {
	if t == nil {
		return observe.NopTracer{}
	}
	return tracer{t: t}
}

// Metrics returns a no-op [observe.Metrics].
//
// Lumen's Registry requires pre-registered names and a fixed label set.
// ion/observe.Metrics is fire-and-forget with dynamic names and kv pairs.
// A thin wrap would invent unbounded series or drop labels, so this adapter
// does not pretend. The registry argument is reserved for a future honest map.
func Metrics(_ *metrics.Registry) observe.Metrics {
	return observe.NopMetrics{}
}

var (
	_ observe.Logger = logger{}
	_ observe.Tracer = tracer{}
)

func eachKV(kv []any, fn func(key string, val any)) {
	for i := 0; i+1 < len(kv); i += 2 {
		key, ok := kv[i].(string)
		if !ok {
			key = fmt.Sprint(kv[i])
		}
		fn(key, kv[i+1])
	}
}
