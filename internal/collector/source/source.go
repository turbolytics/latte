package source

import (
	"fmt"
	"github.com/turbolytics/collector/internal/collector/metric/sources"
)

type Type string

const (
	TypePostgres   Type = "postgres"
	TypeMongoDB    Type = "mongodb"
	TypePrometheus Type = "prometheus"
)

type TypeStrategy string

const (
	TypeStrategyHistoricTumblingWindow TypeStrategy = "historic_tumbling_window"
	TypeStrategyTick                   TypeStrategy = "tick"
)

// TODO need to clean all of this up
type Source struct {
	Strategy TypeStrategy
	Config   map[string]any
	Type     Type

	MetricSourcer sources.MetricSourcer
}

func (s Source) Validate() error {
	vs := map[TypeStrategy]struct{}{
		TypeStrategyTick:                   {},
		TypeStrategyHistoricTumblingWindow: {},
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
