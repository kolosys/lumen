package ionadapt

import (
	"context"

	"github.com/kolosys/lumen/trace"
)

type tracer struct{ t *trace.Tracer }

func (a tracer) Start(ctx context.Context, name string, kv ...any) (context.Context, func(err error)) {
	ctx, span := a.t.Start(ctx, name, trace.WithAttributes(spanAttrs(kv)...))
	return ctx, func(err error) {
		if err != nil {
			span.RecordError(err)
		}
		span.End()
	}
}

func spanAttrs(kv []any) []trace.Attribute {
	n := len(kv) / 2
	if n == 0 {
		return nil
	}
	attrs := make([]trace.Attribute, 0, n)
	eachKV(kv, func(key string, val any) {
		attrs = append(attrs, trace.Attribute{Key: key, Value: val})
	})
	return attrs
}
