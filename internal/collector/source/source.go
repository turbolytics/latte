package source

import (
	"fmt"
	"github.com/turbolytics/latte/internal/collector/partition/sources/s3"
	"github.com/turbolytics/latte/internal/source"
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

type Config struct {
	Strategy TypeStrategy
	Config   map[string]any
	Type     Type
	Sourcer  source.Sourcer
}

func (c Config) Window() *time.Duration {
	return c.Sourcer.Window()
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
	var s source.Sourcer
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
	c.Sourcer = s
	return nil
}
