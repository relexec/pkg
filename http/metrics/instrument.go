package metrics

import (
	"go.opentelemetry.io/otel/metric"
)

const (
	InstrumentNameRequests = "http_requests"
	InstrumentDescRequests = "Count of HTTP requests. Labels: 'method', 'status.code', 'route'."

	InstrumentNameRequestsCurrent = "http_requests_current"
	InstrumentDescRequestsCurrent = "Count of HTTP requests currently being processed. Labels: none."

	InstrumentNameRequestsDuration = "http_requests_duration"
	InstrumentDescRequestsDuration = "Histogram of HTTP request duration in seconds. Labels: 'method', 'status.code', 'route'."
)

var (
	InstrumentRequests         metric.Int64Counter
	InstrumentRequestsCurrent  metric.Int64UpDownCounter
	InstrumentRequestsDuration metric.Float64Histogram
)
