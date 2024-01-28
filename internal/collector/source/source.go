package source

import (
	"context"
	"fmt"
	"github.com/turbolytics/collector/internal/collector/partition/sources/s3"
	"github.com/turbolytics/collector/internal/metrics"
	"github.com/turbolytics/collector/internal/partition"
	"time"
)

type Type string

const (
	TypeMongoDB     Type = "mongodb"
	TypePartitionS3 Type = "partition.s3"
	TypePostgres    Type = "postgres"
	TypePrometheus  Type = "prometheus"
)

type TypeStrategy string

const (
	TypeStrategyHistoricTumblingWindow TypeStrategy = "historic_tumbling_window"
	TypeStrategyIncremental            TypeStrategy = "incremental"
	TypeStrategyTick                   TypeStrategy = "tick"
)

type MetricSourcer interface {
	Source(ctx context.Context) ([]*metrics.Metric, error)
	Window() *time.Duration
}

type PartitionSourcer interface {
	Source(ctx context.Context) (*partition.Partition, error)
	Window() *time.Duration
}

type Config struct {
	Strategy TypeStrategy
	Config   map[string]any
	Type     Type

	MetricSourcer    MetricSourcer
	PartitionSourcer PartitionSourcer
}

func (c Config) Window() *time.Duration {
	if c.MetricSourcer != nil {
		return c.MetricSourcer.Window()
	} else if c.PartitionSourcer != nil {
		return c.PartitionSourcer.Window()
	}
	return nil
}

func (c Config) Validate() error {
	vs := map[TypeStrategy]struct{}{
		TypeStrategyTick:                   {},
		TypeStrategyHistoricTumblingWindow: {},
		TypeStrategyIncremental:            {},
	}

	if _, ok := vs[c.Strategy]; !ok {
		return fmt.Errorf("unknown strategy: %q", c.Strategy)
	}
	return nil
}

func (c *Config) SetDefaults() {
	if c.Strategy == "" {
		c.Strategy = TypeStrategyTick
	}
}

func (c *Config) Init() error {
	var s PartitionSourcer
	var err error
	switch c.Type {
	case TypePartitionS3:
		s, err = s3.NewFromGenericConfig(
			c.Config,
		)
	}
	if err != nil {
		return err
	}
	c.PartitionSourcer = s
	return nil
}
