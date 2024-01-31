package sources

import (
	"github.com/turbolytics/latte/internal/metrics"
	"github.com/turbolytics/latte/internal/source"
)

type Metrics struct {
	Ms []source.Record
}

func (m *Metrics) Transform() error {
	return nil
}

func (m *Metrics) Records() []source.Record {
	return m.Ms
}

func NewMetricsResult(ms []*metrics.Metric) *Metrics {
	metrics := &Metrics{}

	for _, m := range ms {
		metrics.Ms = append(metrics.Ms, m)
	}

	return metrics
}
