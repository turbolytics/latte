package collector

import (
	"github.com/turbolytics/collector/internal"
	"github.com/turbolytics/collector/internal/metrics"
)

type Collector struct {
	Config *internal.Config
}

func (c *Collector) Source() ([]metrics.Metric, error) {
	return nil, nil
}

func (c *Collector) Sink(metrics []metrics.Metric) error {
	return nil
}

func (c *Collector) Invoke() ([]metrics.Metric, error) {
	ms, err := c.Source()
	if err != nil {
		return nil, err
	}
	err = c.Sink(ms)
	return ms, err
}

func New(config *internal.Config) (*Collector, error) {
	return &Collector{
		Config: config,
	}, nil
}
