package trace

// Options configures a Tracer.
type Options struct {
	// ServiceName identifies the service in traces.
	ServiceName string

	// Sampler determines which spans to record.
	Sampler Sampler

	// Exporter receives completed spans.
	Exporter Exporter

	// MaxSpansPerSecond limits span creation rate (0 = unlimited).
	MaxSpansPerSecond int

	// PropagationFormat sets the context propagation format.
	// Supports: "w3c", "kolosys", "both" (default: "both")
	PropagationFormat string

	// AsyncExport enables asynchronous span export.
	AsyncExport bool

	// AsyncBufferSize sets the async export buffer size.
	AsyncBufferSize int
}

func (o *Options) applyDefaults() {
	if o.ServiceName == "" {
		o.ServiceName = "unknown"
	}
	if o.Sampler == nil {
		o.Sampler = AlwaysSample()
	}
	if o.Exporter == nil {
		o.Exporter = NopExporter{}
	}
	if o.PropagationFormat == "" {
		o.PropagationFormat = "both"
	}
	if o.AsyncBufferSize == 0 {
		o.AsyncBufferSize = 1024
	}
}
