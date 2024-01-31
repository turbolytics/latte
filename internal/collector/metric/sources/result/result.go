package result

import (
	"github.com/turbolytics/latte/internal/metrics"
	"github.com/turbolytics/latte/internal/source"
)

//

type Metrics struct {
	metrics []*metrics.Metric
}

func (m *Metrics) Transform() error {
	/*
		metricConfig := c.(metric.Config)

		for _, metric := range m.metrics {
			metric.Name = metricConfig.Metric.Name
			metric.Type = metricConfig.Metric.Type

			// enrich with tags
			// should these be copied?
			for _, t := range metricConfig.Metric.Tags {
				metric.Tags[t.Key] = t.Value
			}
		}

	*/

	return nil
}

func (m *Metrics) Records() []source.Record {
	var records []source.Record
	for _, m := range m.metrics {
		records = append(records, m)
	}
	return records
}

func NewMetricsResult(ms []*metrics.Metric) *Metrics {
	metrics := &Metrics{
		metrics: ms,
	}

	return metrics
}
