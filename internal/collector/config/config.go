package config

import (
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

// Root contains top level collector configuration
// required for initializing more concrete collectors.
type Root struct {
	Collector Collector

	Validate bool
	Logger   *zap.Logger
}
