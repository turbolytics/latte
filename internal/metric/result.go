package metric

import (
	"github.com/turbolytics/latte/internal/record"
)

type Metrics struct {
	Metrics []*Metric
}

func (m *Metrics) Records() []record.Record {
	var records []record.Record
	for _, m := range m.Metrics {
		records = append(records, m)
	}
	return records
}

func NewMetricsResult(ms []*Metric) *Metrics {
	metrics := &Metrics{
		Metrics: ms,
	}

	return metrics
}
