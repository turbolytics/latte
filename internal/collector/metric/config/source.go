package config

import (
	"fmt"
	"github.com/turbolytics/collector/internal/collector/metric/sources"
)

type TypeSourceStrategy string

const (
	TypeSourceStrategyHistoricTumblingWindow TypeSourceStrategy = "historic_tumbling_window"
	TypeSourceStrategyTick                   TypeSourceStrategy = "tick"
)

type Source struct {
	Type     sources.Type
	Strategy TypeSourceStrategy
	Sourcer  sources.MetricSourcer
	Config   map[string]any
}

func (s Source) Validate() error {
	vs := map[TypeSourceStrategy]struct{}{
		TypeSourceStrategyTick:                   {},
		TypeSourceStrategyHistoricTumblingWindow: {},
	}

	if _, ok := vs[s.Strategy]; !ok {
		return fmt.Errorf("unknown strategy: %q", s.Strategy)
	}
	return nil
}

func (s *Source) SetDefaults() {
	if s.Strategy == "" {
		s.Strategy = TypeSourceStrategyTick
	}
}
