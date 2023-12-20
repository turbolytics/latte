package collector

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/turbolytics/collector/internal"
	"github.com/turbolytics/collector/internal/metrics"
	"go.uber.org/zap"
)

type Collector struct {
	logger *zap.Logger
	Config *internal.Config
}

func (c *Collector) Close() error {
	for _, s := range c.Config.Sinks {
		s.Sinker.Close()
	}
	return nil
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

func (c *Collector) Sink(ctx context.Context, metrics []*metrics.Metric) error {
	// need to add a serializer
	for _, metric := range metrics {
		bs, err := json.Marshal(metric)
		if err != nil {
			return err
		}
		for _, s := range c.Config.Sinks {
			_, err := s.Sinker.Write(bs)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// InvokeHandleError will log any Invoke errors and not return them.
// Useful for async scheduling.
func (c *Collector) InvokeHandleError(ctx context.Context) {
	_, err := c.Invoke(ctx)
	if err != nil {
		c.logger.Error(err.Error())
	}
}

func (c *Collector) Invoke(ctx context.Context) ([]*metrics.Metric, error) {
	id := uuid.New()
	c.logger.Info(
		"collector.Invoke",
		zap.String("id", id.String()),
		zap.String("name", c.Config.Name),
	)
	ctx = context.WithValue(ctx, "id", id)
	ms, err := c.Source(ctx)
	if err != nil {
		return nil, err
	}

	// only sink if metrics are present:
	if len(ms) > 0 {
		err = c.Sink(ctx, ms)
	} else {
		c.logger.Warn(
			"collector.Invoke",
			zap.String("msg", "no metrics found"),
			zap.String("id", id.String()),
			zap.String("name", c.Config.Name),
		)
	}
	return ms, err
}

type Option func(*Collector)

func WithLogger(l *zap.Logger) Option {
	return func(c *Collector) {
		c.logger = l
	}
}

func New(config *internal.Config, opts ...Option) (*Collector, error) {
	c := &Collector{
		Config: config,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

func NewFromConfigs(configs []*internal.Config, opts ...Option) ([]*Collector, error) {
	var cs []*Collector
	for _, config := range configs {
		coll, err := New(config, opts...)
		if err != nil {
			return nil, err
		}
		cs = append(cs, coll)
	}
	return cs, nil
}
