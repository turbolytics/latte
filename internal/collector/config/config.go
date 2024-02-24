package config

import (
	"fmt"
	"github.com/turbolytics/latte/internal/invoker"
	"go.uber.org/zap"
)

type Type string

const (
	TypeMetric    Type = "metric"
	TypePartition Type = "partition"
	TypeEntity    Type = "entity"
)

type Collector struct {
	Type               Type
	InvocationStrategy invoker.TypeStrategy `yaml:"invocation_strategy"`
}

func (c Collector) Validate() error {
	vs := map[invoker.TypeStrategy]struct{}{
		invoker.TypeStrategyTick:                   {},
		invoker.TypeStrategyHistoricTumblingWindow: {},
		invoker.TypeStrategyIncremental:            {},
	}

	if _, ok := vs[c.InvocationStrategy]; !ok {
		return fmt.Errorf("unknown strategy: %q", c.InvocationStrategy)
	}
	return nil
}

func (c *Collector) SetDefaults() {
	if c.InvocationStrategy == "" {
		c.InvocationStrategy = invoker.TypeStrategyTick
	}
}

// Root contains top level collector configuration
// required for initializing more concrete collectors.
type Root struct {
	Collector Collector

	Validate bool
	Logger   *zap.Logger
}
