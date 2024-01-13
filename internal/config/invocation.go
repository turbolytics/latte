package config

import (
	"fmt"
	"github.com/turbolytics/collector/internal/invocation"
)

type Invocation struct {
	Strategy invocation.TypeStrategy
}

func (i Invocation) Validate() error {
	vs := map[invocation.TypeStrategy]struct{}{
		invocation.TypeStrategyTick:           {},
		invocation.TypeStrategyHistoricWindow: {},
	}

	if _, ok := vs[i.Strategy]; !ok {
		return fmt.Errorf("unknown strategy: %q", i.Strategy)
	}
	return nil
}

func (i *Invocation) SetDefaults() {
	if i.Strategy == "" {
		i.Strategy = invocation.TypeStrategyTick
	}
}
