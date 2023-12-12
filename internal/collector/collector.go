package collector

import (
	"context"
	"github.com/turbolytics/collector/internal"
	"github.com/turbolytics/collector/internal/metrics"
)

type Collector struct {
	Config *internal.Config
}

func (c *Collector) Source(ctx context.Context) ([]metrics.Metric, error) {
	return c.Config.Source.Sourcer.Source(ctx)
	return nil, nil
}

func (c *Collector) Sink(metrics []metrics.Metric) error {
	return nil
}

func (c *Collector) Invoke(ctx context.Context) ([]metrics.Metric, error) {
	ms, err := c.Source(ctx)
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
