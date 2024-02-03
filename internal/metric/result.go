package metric

import (
	"github.com/turbolytics/latte/internal/record"
)

type Metrics struct {
	metrics []*Metric
}

func (m *Metrics) Records() []record.Record {
	var records []record.Record
	for _, m := range m.metrics {
		records = append(records, m)
	}
	return records
}

func NewMetricsResult(ms []*Metric) *Metrics {
	metrics := &Metrics{
		metrics: ms,
	}

	return metrics
}
