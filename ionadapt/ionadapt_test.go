package ionadapt_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/kolosys/ion/observe"
	"github.com/kolosys/lumen/ionadapt"
	"github.com/kolosys/lumen/logs"
	"github.com/kolosys/lumen/metrics"
	"github.com/kolosys/lumen/trace"
)

func TestLoggerRoundTrip(t *testing.T) {
	t.Parallel()

	boom := errors.New("boom")
	buf := &bytes.Buffer{}
	l := logs.New(&logs.Options{
		Output:    buf,
		Level:     logs.DebugLevel,
		Formatter: &logs.JSONFormatter{DisableTimestamp: true},
	})
	ol := ionadapt.Logger(l)

	ol.Debug("dbg", "k", "v")
	ol.Info("inf", "n", 7)
	ol.Error("fail", boom, "op", "write")

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("got %d lines:\n%s", len(lines), buf.String())
	}

	want := []map[string]any{
		{"level": "debug", "msg": "dbg", "k": "v"},
		{"level": "info", "msg": "inf", "n": "7"},
		{"level": "error", "msg": "fail", "op": "write", "error": "boom"},
	}
	for i, line := range lines {
		got := map[string]any{}
		if err := json.Unmarshal([]byte(line), &got); err != nil {
			t.Fatalf("line %d: %v\n%s", i, err, line)
		}
		for k, v := range want[i] {
			if got[k] != v {
				t.Errorf("line %d %q: got %#v want %#v", i, k, got[k], v)
			}
		}
	}
}

func TestLoggerNilAndOddKV(t *testing.T) {
	t.Parallel()

	ionadapt.Logger(nil).Info("discarded", "orphan")

	buf := &bytes.Buffer{}
	l := logs.New(&logs.Options{
		Output:    buf,
		Level:     logs.DebugLevel,
		Formatter: &logs.JSONFormatter{DisableTimestamp: true},
	})
	ionadapt.Logger(l).Warn("w", 1, "nonstr", "odd", "ok", "orphan")

	got := map[string]any{}
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got["msg"] != "w" || got["1"] != "nonstr" || got["odd"] != "ok" {
		t.Fatalf("got %#v", got)
	}
	if _, ok := got["orphan"]; ok {
		t.Fatalf("unpaired key should be dropped: %#v", got)
	}
}

func TestTracerStartEnd(t *testing.T) {
	t.Parallel()

	exp := &spanCapture{}
	ot := ionadapt.Tracer(trace.New(&trace.Options{Exporter: exp}))

	ctx, finish := ot.Start(context.Background(), "op", "user", "ada")
	if ctx == nil {
		t.Fatal("nil ctx")
	}
	finish(nil)

	boom := errors.New("nope")
	_, finishErr := ot.Start(context.Background(), "child")
	finishErr(boom)

	spans := exp.snapshot()
	if len(spans) != 2 {
		t.Fatalf("got %d spans, want 2", len(spans))
	}
	if spans[0].name != "op" || spans[0].status != trace.StatusUnset.String() {
		t.Errorf("first span: %+v", spans[0])
	}
	if got := attrVal(spans[0].attrs, "user"); got != "ada" {
		t.Errorf("user attr: got %#v", got)
	}
	if spans[1].name != "child" || spans[1].status != trace.StatusError.String() {
		t.Errorf("error span: %+v", spans[1])
	}
	if spans[1].statusMsg != boom.Error() {
		t.Errorf("status msg: %q", spans[1].statusMsg)
	}
}

func TestTracerNil(t *testing.T) {
	t.Parallel()
	ctx, finish := ionadapt.Tracer(nil).Start(context.Background(), "x")
	if ctx == nil {
		t.Fatal("nil ctx")
	}
	finish(errors.New("ignored"))
}

func TestMetricsNoop(t *testing.T) {
	t.Parallel()
	m := ionadapt.Metrics(metrics.NewRegistry(nil))
	m.Inc("n")
	m.Add("n", 1)
	m.Gauge("g", 2)
	m.Histogram("h", 3)
}

func TestSatisfiesObserve(t *testing.T) {
	t.Parallel()
	var (
		_ observe.Logger  = ionadapt.Logger(logs.New(nil))
		_ observe.Tracer  = ionadapt.Tracer(trace.New(nil))
		_ observe.Metrics = ionadapt.Metrics(nil)
	)
}

type captured struct {
	name      string
	status    string
	statusMsg string
	attrs     []trace.Attribute
}

type spanCapture struct {
	mu    sync.Mutex
	spans []captured
}

func (c *spanCapture) Export(span *trace.Span) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.spans = append(c.spans, captured{
		name:      span.Name(),
		status:    span.Status().String(),
		statusMsg: span.StatusMessage(),
		attrs:     span.Attributes(),
	})
}

func (c *spanCapture) Close() error { return nil }

func (c *spanCapture) snapshot() []captured {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]captured, len(c.spans))
	copy(out, c.spans)
	return out
}

func attrVal(attrs []trace.Attribute, key string) any {
	for _, a := range attrs {
		if a.Key == key {
			return a.Value
		}
	}
	return nil
}
