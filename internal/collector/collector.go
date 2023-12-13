package collector

import (
	"context"
	"encoding/json"
	"github.com/turbolytics/collector/internal"
	"github.com/turbolytics/collector/internal/metrics"
)

type Collector struct {
	Config *internal.Config
}

func (c *Collector) Source(ctx context.Context) ([]*metrics.Metric, error) {
	ms, err := c.Config.Source.Sourcer.Source(ctx)
	if err != nil {
		return nil, err
	}

	for _, m := range ms {
		m.Name = c.Config.Metric.Name
		m.Type = c.Config.Metric.Type
	}

	return ms, nil
}

func (c *Collector) Sink(metrics []*metrics.Metric) error {
	// need to add a serializer
	bs, err := json.Marshal(metrics)
	if err != nil {
		return err
	}
	for _, s := range c.Config.Sinks {
		_, err := s.Sinker.Write(bs)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Collector) Invoke(ctx context.Context) ([]*metrics.Metric, error) {
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
