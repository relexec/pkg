package metrics

import (
	"context"

	"go.opentelemetry.io/otel/metric"

	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

type WithOption func(*Handler)

// New returns a new Handler and instantiates all instruments published by
// the metrics provider.
func New(
	ctx context.Context,
	name string,
	opts ...WithOption,
) (*Handler, error) {
	var err error
	h := Handler{
		name: name,
	}
	for _, opt := range opts {
		opt(&h)
	}
	if h.exporter == nil {
		exp, err := stdoutmetric.New(stdoutmetric.WithPrettyPrint())
		if err != nil {
			return nil, err
		}
		h.exporter = exp
	}
	if h.reader == nil {
		if !h.cfg.HTTP.Enabled {
			h.reader = sdkmetric.NewPeriodicReader(h.exporter)
		} else {
			h.reader = sdkmetric.NewManualReader()
		}
	}
	if h.mp == nil {
		mp := sdkmetric.NewMeterProvider(
			sdkmetric.WithView(Views...),
			sdkmetric.WithReader(h.reader),
		)
		h.mp = mp
	}
	p := h.MeterProvider()
	m := p.Meter(h.name)

	InstrumentRequests, err = m.Int64Counter(
		InstrumentNameRequests,
		metric.WithDescription(InstrumentDescRequests),
		metric.WithUnit("{call}"),
	)
	if err != nil {
		return nil, err
	}

	InstrumentRequestsCurrent, err = m.Int64UpDownCounter(
		InstrumentNameRequestsCurrent,
		metric.WithDescription(InstrumentDescRequestsCurrent),
		metric.WithUnit("{call}"),
	)
	if err != nil {
		return nil, err
	}

	InstrumentRequestsDuration, err = m.Float64Histogram(
		InstrumentNameRequestsDuration,
		metric.WithDescription(InstrumentDescRequestsDuration),
		metric.WithUnit("{seconds}"),
	)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

// WithConfig sets the Metrics handler's Config.
func WithConfig(cfg Config) WithOption {
	return func(h *Handler) {
		h.cfg = cfg
	}
}

// WithMeterProvider sets the Metrics handler's MeterProvider.
func WithMeterProvider(mp *sdkmetric.MeterProvider) WithOption {
	return func(h *Handler) {
		h.mp = mp
	}
}

// WithExporter sets the Metrics handler's Exporter.
func WithExporter(exp sdkmetric.Exporter) WithOption {
	return func(h *Handler) {
		h.exporter = exp
	}
}

// WithReader sets the Metrics handler's Reader.
func WithReader(reader sdkmetric.Reader) WithOption {
	return func(h *Handler) {
		h.reader = reader
	}
}
