package collector

import (
	"fmt"
	"github.com/turbolytics/latte/internal/collector/metric"
	"github.com/turbolytics/latte/internal/invoker"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
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

func NewFromFile(fpath string, opts ...RootOption) (invoker.Collector, error) {
	fmt.Printf("loading config from file: %q\n", fpath)

	bs, err := os.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	return New(bs, opts...)
}

func New(bs []byte, opts ...RootOption) (invoker.Collector, error) {
	var conf RootConfig

	for _, opt := range opts {
		opt(&conf)
	}

	if err := yaml.Unmarshal(bs, &conf); err != nil {
		return nil, err
	}

	var coll invoker.Collector
	var err error

	switch conf.Collector.Type {
	case TypeMetric:
		coll, err = metric.NewCollectorFromConfig(
			bs,
			metric.WithJustValidation(conf.validate),
			metric.WithLogger(conf.logger),
		)
	case TypePartition:
		/*
			collConfig, err := partition.NewConfig(
				bs,
				partition.ConfigWithJustValidation(conf.validate),
				partition.ConfigWithLogger(conf.logger),
			)
		*/
	default:
		return nil, fmt.Errorf("collector type: %v not supported", conf.Collector.Type)
	}

	if err != nil {
		return nil, err
	}

	return coll, err
}

func NewFromGlob(glob string, opts ...RootOption) ([]invoker.Collector, error) {
	files, err := filepath.Glob(glob)
	if err != nil {
		return nil, err
	}
	var collectors []invoker.Collector

	for _, fName := range files {
		c, err := NewFromFile(fName, opts...)
		if err != nil {
			return nil, err
		}
		collectors = append(collectors, c)
	}
	return collectors, nil
}
