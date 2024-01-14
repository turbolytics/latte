package config

import (
	"fmt"
	"github.com/turbolytics/collector/internal/sources"
)

type TypeSourceStrategy string

const (
	TypeSourceStrategyHistoricWindow TypeSourceStrategy = "historic_window"
	TypeSourceStrategyTick           TypeSourceStrategy = "tick"
)

type Source struct {
	Type     sources.Type
	Strategy TypeSourceStrategy
	Sourcer  sources.Sourcer
	Config   map[string]any
}

func (s Source) Validate() error {
	vs := map[TypeSourceStrategy]struct{}{
		TypeSourceStrategyTick:           {},
		TypeSourceStrategyHistoricWindow: {},
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
