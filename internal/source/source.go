package source

import (
	"fmt"
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

/*
func (c *Config) Init(l *zap.Logger, validate bool) error {
	var err error
	var s Sourcer
	switch c.Type {
			case TypePartitionS3:
				s, err = s3.NewFromGenericConfig(
					c.Config,
				)
			}
		case TypePostgres:
			s, err = postgres.NewFromGenericConfig(
				c.Config,
				validate,
			)
	case TypeMongoDB:
		s, err = mongodb.NewFromGenericConfig(
			context.TODO(),
			c.Config,
			validate,
		)
			case TypePrometheus:
				s, err = prometheus.NewFromGenericConfig(
					c.Config,
					prometheus.WithLogger(l),
				)
	default:
		return fmt.Errorf("source type: %q unknown", c.Type)
	}

	if err != nil {
		return err
	}
	c.Sourcer = s
	return nil
}
*/
