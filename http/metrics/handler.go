package metrics

import (
	"go.opentelemetry.io/otel/sdk/metric"
)

// Handler handles OTEL Meters for an HTTP server
type Handler struct {
	// name is the name of the HTTP server -- i.e. the name of the group of
	// meters being provided.
	name string
	// mp is the OTEL metric.MeterProvider
	mp *metric.MeterProvider
	// exporter is the OTEL metric.Exporter
	exporter metric.Exporter
	// reader is the OTEL metric.Reader
	reader metric.Reader
	// cfg contains configuration options for the metrics Handler.
	cfg Config
}

// Name returns the name of the group of metrics being provided.
func (h Handler) Name() string {
	return h.name
}

// MeterProvider returns the OTEL metric.MeterProvider
func (h Handler) MeterProvider() *metric.MeterProvider {
	return h.mp
}

// Exporter returns the OTEL metric.Exporter.
func (h Handler) Exporter() metric.Exporter {
	return h.exporter
}

// Reader returns the OTEL metric.Reader.
func (h Handler) Reader() metric.Reader {
	return h.reader
}
