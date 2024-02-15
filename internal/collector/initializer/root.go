package initializer

import (
	"fmt"
	"github.com/turbolytics/latte/internal/collector/config"
	"github.com/turbolytics/latte/internal/invoker"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type RootOption func(*config.Root)

func WithJustValidation(validate bool) RootOption {
	return func(c *config.Root) {
		c.Validate = validate
	}
}

func RootWithLogger(l *zap.Logger) RootOption {
	return func(c *config.Root) {
		c.Logger = l
	}
}

func NewCollectorsFromGlob(glob string, opts ...RootOption) ([]invoker.Collector, error) {
	files, err := filepath.Glob(glob)
	if err != nil {
		return nil, err
	}
	var collectors []invoker.Collector

	for _, fName := range files {
		c, err := NewCollectorFromFile(fName, opts...)
		if err != nil {
			return nil, err
		}
		collectors = append(collectors, c)
	}
	return collectors, nil
}

func NewCollectorFromFile(fpath string, opts ...RootOption) (invoker.Collector, error) {
	fmt.Printf("loading config from file: %q\n", fpath)

	bs, err := os.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	return NewCollector(bs, opts...)
}

func NewCollector(bs []byte, opts ...RootOption) (invoker.Collector, error) {
	var conf config.Root

	for _, opt := range opts {
		opt(&conf)
	}

	if err := yaml.Unmarshal(bs, &conf); err != nil {
		return nil, err
	}

	var coll invoker.Collector
	var err error

	switch conf.Collector.Type {
	case config.TypeMetric:
		coll, err = NewMetricCollectorFromConfig(
			bs,
			conf.Validate,
			conf.Logger,
		)
	case config.TypePartition:
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
