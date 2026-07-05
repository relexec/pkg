package metrics

type SeriesSummary struct {
	Min float64 `json:"min" doc:"Minimum value."`
	Max float64 `json:"max" doc:"Maximum value."`
	Avg float64 `json:"avg" doc:"Average value."`
}

type RequestsMetrics struct {
	Current         int64                     `json:"current" doc:"Count of in-flight HTTP requests."`
	Count           int64                     `json:"count" doc:"Aggregate count of all HTTP requests."`
	CountByRoute    map[string]int64          `json:"count_by_route" doc:"Counts of HTTP requests by route/path."`
	CountByStatus   map[int64]int64           `json:"count_by_status" doc:"Counts of HTTP requests by status code."`
	Duration        *SeriesSummary            `json:"duration" doc:"min/max/avg/mean of HTTP request duration in seconds."`
	DurationByRoute map[string]*SeriesSummary `json:"duration_by_route" doc:"min/max/avg/means of HTTP request duration by route/path."`
}

type ListMetricsOutputBody struct {
	Requests RequestsMetrics `json:"requests" doc:"HTTP request metrics."`
}

type ListMetricsOutput struct {
	Body ListMetricsOutputBody
}
