package metrics

import (
	"fmt"
	"math"
	"strings"

	"github.com/samber/lo"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// Sum returns the int64 sum of the named instrument within the supplied
// metricsdata.ResourceMetrics.
//
// If the named instrument does not exist in the supplied ResourceMetrics,
// returns an error.
//
// If the named instrument exists in the supplied ResourceMetrics but is not an
// int64 sum, returns an error.
func Sum(
	data *metricdata.ResourceMetrics,
	instrument string,
) (int64, error) {
	for _, sm := range data.ScopeMetrics {
		m, ok := lo.Find(
			sm.Metrics,
			func(subject metricdata.Metrics) bool {
				return strings.EqualFold(instrument, subject.Name)
			},
		)
		if !ok {
			continue
		}
		sum, ok := m.Data.(metricdata.Sum[int64])
		if !ok {
			return -1, fmt.Errorf(
				"instrument %q is a %T, not a metricsdata.Sum[int64]",
				instrument, m.Data,
			)
		}
		return lo.SumBy(
			sum.DataPoints,
			func(dp metricdata.DataPoint[int64]) int64 {
				return dp.Value
			},
		), nil
	}
	return -1, fmt.Errorf(
		"instrument %q not in supplied ResourceMetrics", instrument,
	)
}

// SumByAttrString returns a map, keyed by unique values of the supplied
// string-valued attribute, of the int64 sums of the named instrument within
// the supplied metricsdata.ResourceMetrics.
//
// If the named instrument does not exist in the supplied ResourceMetrics,
// returns an error.
//
// If the named instrument exists in the supplied ResourceMetrics but is not an
// int64 sum, returns an error.
func SumByAttrString(
	data *metricdata.ResourceMetrics,
	instrument string,
	attributeKey string,
) (map[string]int64, error) {
	out := map[string]int64{}
	for _, sm := range data.ScopeMetrics {
		m, ok := lo.Find(
			sm.Metrics,
			func(subject metricdata.Metrics) bool {
				return strings.EqualFold(instrument, subject.Name)
			},
		)
		if !ok {
			continue
		}
		sum, ok := m.Data.(metricdata.Sum[int64])
		if !ok {
			return nil, fmt.Errorf(
				"instrument %q is a %T, not a metricsdata.Sum[int64]",
				instrument, m.Data,
			)
		}
		attrVals := map[string]int64{}
		for _, dp := range sum.DataPoints {
			for _, kv := range dp.Attributes.ToSlice() {
				if string(kv.Key) != attributeKey {
					continue
				}
				attrVals[kv.Value.AsString()] = dp.Value
				break
			}
		}
		out = lo.Assign(out, attrVals)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf(
			"instrument %q not in supplied ResourceMetrics", instrument,
		)
	}
	return out, nil
}

// SumByAttrInt64 returns a map, keyed by unique values of the supplied
// int64-valuedZ attribute, of the int64 sums of the named instrument within
// the supplied metricsdata.ResourceMetrics.
//
// If the named instrument does not exist in the supplied ResourceMetrics,
// returns an error.
//
// If the named instrument exists in the supplied ResourceMetrics but is not an
// int64 sum, returns an error.
func SumByAttrInt64(
	data *metricdata.ResourceMetrics,
	instrument string,
	attributeKey string,
) (map[int64]int64, error) {
	out := map[int64]int64{}
	for _, sm := range data.ScopeMetrics {
		m, ok := lo.Find(
			sm.Metrics,
			func(subject metricdata.Metrics) bool {
				return strings.EqualFold(instrument, subject.Name)
			},
		)
		if !ok {
			continue
		}
		sum, ok := m.Data.(metricdata.Sum[int64])
		if !ok {
			return nil, fmt.Errorf(
				"instrument %q is a %T, not a metricsdata.Sum[int64]",
				instrument, m.Data,
			)
		}
		attrVals := map[int64]int64{}
		for _, dp := range sum.DataPoints {
			for _, kv := range dp.Attributes.ToSlice() {
				if string(kv.Key) != attributeKey {
					continue
				}
				attrVals[kv.Value.AsInt64()] = dp.Value
				break
			}
		}
		out = lo.Assign(out, attrVals)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf(
			"instrument %q not in supplied ResourceMetrics", instrument,
		)
	}
	return out, nil
}

// Value returns the single data point value of the named instrument in the
// supplied metricsdata.ResourceMetrics.
//
// If the named instrument does not exist in the supplied ResourceMetrics,
// returns an error.
//
// If the named instrument exists in the supplied ResourceMetrics but is not an
// int64 counter with a single data point, returns an error.
func Value(
	data *metricdata.ResourceMetrics,
	instrument string,
) (int64, error) {
	for _, sm := range data.ScopeMetrics {
		m, ok := lo.Find(
			sm.Metrics,
			func(subject metricdata.Metrics) bool {
				return strings.EqualFold(instrument, subject.Name)
			},
		)
		if !ok {
			continue
		}
		sum, ok := m.Data.(metricdata.Sum[int64])
		if !ok {
			return -1, fmt.Errorf(
				"instrument %q is a %T, not a metricsdata.Sum[int64]",
				instrument, m.Data,
			)
		}
		if len(sum.DataPoints) != 1 {
			return -1, fmt.Errorf(
				"expected instrument %q to have 1 data point, but found %d",
				instrument, len(sum.DataPoints),
			)
		}
		return sum.DataPoints[0].Value, nil
	}
	return -1, fmt.Errorf(
		"instrument %q not in supplied ResourceMetrics", instrument,
	)
}

// SummarizeHistogram returns a SeriesSummary of the named instrument within
// the supplied metricsdata.ResourceMetrics.
//
// If the named instrument does not exist in the supplied ResourceMetrics,
// returns an error.
//
// If the named instrument exists in the supplied ResourceMetrics but is not a
// histogram instrument, returns an error.
func SummarizeHistogram(
	data *metricdata.ResourceMetrics,
	instrument string,
) (*SeriesSummary, error) {
	out := SeriesSummary{}
	for _, sm := range data.ScopeMetrics {
		m, ok := lo.Find(
			sm.Metrics,
			func(subject metricdata.Metrics) bool {
				return strings.EqualFold(instrument, subject.Name)
			},
		)
		if !ok {
			continue
		}
		histData, ok := m.Data.(metricdata.Histogram[float64])
		if !ok {
			return nil, fmt.Errorf(
				"instrument %q is a %T, not a metricsdata.Sum[int64]",
				instrument, m.Data,
			)
		}
		count := uint64(0)
		sumVal := 0.0
		minVal := math.MaxFloat64
		maxVal := 0.0
		for _, dp := range histData.DataPoints {
			if dp.Count == 0 {
				continue
			}
			count += dp.Count
			sumVal += dp.Sum
			minDP, _ := dp.Min.Value()
			minVal = min(minVal, minDP)
			maxDP, _ := dp.Max.Value()
			maxVal = max(maxVal, maxDP)
		}

		avgVal := sumVal / float64(count)

		out.Min = minVal
		out.Max = maxVal
		out.Avg = avgVal
		return &out, nil
	}
	return nil, fmt.Errorf(
		"instrument %q not in supplied ResourceMetrics", instrument,
	)
}

// SummarizeHistogramByAttrString returns a map, keyed by unique values of the
// supplied string-valued attribute, of SeriesSummary of the named instrument
// within the supplied metricsdata.ResourceMetrics.
//
// If the named instrument does not exist in the supplied ResourceMetrics,
// returns an error.
//
// If the named instrument exists in the supplied ResourceMetrics but is not a
// histogram instrument, returns an error.
func SummarizeHistogramByAttrString(
	data *metricdata.ResourceMetrics,
	instrument string,
	attributeKey string,
) (map[string]*SeriesSummary, error) {
	dpsByAttrValue := map[string][]metricdata.HistogramDataPoint[float64]{}
	for _, sm := range data.ScopeMetrics {
		m, ok := lo.Find(
			sm.Metrics,
			func(subject metricdata.Metrics) bool {
				return strings.EqualFold(instrument, subject.Name)
			},
		)
		if !ok {
			continue
		}
		histData, ok := m.Data.(metricdata.Histogram[float64])
		if !ok {
			return nil, fmt.Errorf(
				"instrument %q is a %T, not a metricsdata.Sum[int64]",
				instrument, m.Data,
			)
		}
		for _, dp := range histData.DataPoints {
			for _, kv := range dp.Attributes.ToSlice() {
				if string(kv.Key) != attributeKey {
					continue
				}
				attrVal := kv.Value.AsString()
				dps, ok := dpsByAttrValue[attrVal]
				if !ok {
					dps = []metricdata.HistogramDataPoint[float64]{}
				}
				dps = append(dps, dp)
				dpsByAttrValue[attrVal] = dps
			}
		}
	}
	if len(dpsByAttrValue) == 0 {
		return nil, fmt.Errorf(
			"instrument %q not in supplied ResourceMetrics", instrument,
		)
	}

	out := map[string]*SeriesSummary{}
	// Now that we've collected the histogram's data points into groups by
	// attribute value for the attribute key we're looking for, we
	// construct SeriesSummary structs for each of those groups.
	for attrVal, dps := range dpsByAttrValue {
		ss := &SeriesSummary{}
		count := uint64(0)
		sumVal := 0.0
		minVal := math.MaxFloat64
		maxVal := 0.0
		for _, dp := range dps {
			if dp.Count == 0 {
				continue
			}
			count += dp.Count
			sumVal += dp.Sum
			minDP, _ := dp.Min.Value()
			minVal = min(minVal, minDP)
			maxDP, _ := dp.Max.Value()
			maxVal = max(maxVal, maxDP)
		}

		avgVal := sumVal / float64(count)

		ss.Min = minVal
		ss.Max = maxVal
		ss.Avg = avgVal
		out[attrVal] = ss
	}
	return out, nil
}

// InResourceMetrics returns true if the named instrument is within
// the supplied metricdata.ResourceMetrics.
func InResourceMetrics(
	data *metricdata.ResourceMetrics,
	instrument string,
) bool {
	for _, sm := range data.ScopeMetrics {
		if lo.ContainsBy(
			sm.Metrics,
			func(subject metricdata.Metrics) bool {
				return strings.EqualFold(instrument, subject.Name)
			},
		) {
			return true
		}
	}
	return false
}
