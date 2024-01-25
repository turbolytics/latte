package collector

import (
	"go.uber.org/zap"
)

type Type string

const (
	TypeMetric    Type = "metric"
	TypePartition Type = "partition"
	TypeEntity    Type = "entity"
)

type collector struct {
	Type Type
}

// RootConfig contains top level collector configuration
// required for initializing more concrete collectors.
type RootConfig struct {
	Collector collector

	validate bool
	logger   *zap.Logger
}

type RootOption func(*RootConfig)

func WithJustValidation(validate bool) RootOption {
	return func(c *RootConfig) {
		c.validate = validate
	}
}

func RootWithLogger(l *zap.Logger) RootOption {
	return func(c *RootConfig) {
		c.logger = l
	}
}
