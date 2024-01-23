package source

import (
	"context"
	"fmt"
	"github.com/turbolytics/collector/internal/metrics"
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
	Source(ctx context.Context) error
}

type Source struct {
	Strategy TypeStrategy
	Config   map[string]any
	Type     Type

	MetricSourcer    MetricSourcer
	PartitionSourcer PartitionSourcer
}

func (s Source) Validate() error {
	vs := map[TypeStrategy]struct{}{
		TypeStrategyTick:                   {},
		TypeStrategyHistoricTumblingWindow: {},
		TypeStrategyIncremental:            {},
	}

	if _, ok := vs[s.Strategy]; !ok {
		return fmt.Errorf("unknown strategy: %q", s.Strategy)
	}
	return nil
}

func (s *Source) SetDefaults() {
	if s.Strategy == "" {
		s.Strategy = TypeStrategyTick
	}
}
