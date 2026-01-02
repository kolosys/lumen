package metrics

import (
	"net/http"
)

// Exporter exports metrics.
type Exporter interface {
	Export(samples []Sample)
}

// NopExporter discards metrics.
type NopExporter struct{}

func (NopExporter) Export([]Sample) {}

// HTTPHandler returns an http.Handler for the Prometheus endpoint.
func HTTPHandler(registry *Registry) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		samples := registry.Collect()
		w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
		WritePrometheus(w, samples)
	})
}

// DefaultHTTPHandler returns an http.Handler using the default registry.
func DefaultHTTPHandler() http.Handler {
	return HTTPHandler(defaultRegistry)
}
