package collector

import (
	"fmt"
	"github.com/turbolytics/collector/internal/collector/metric"
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

type Config struct {
	Collector collector

	validate bool
	logger   *zap.Logger
}

type Option func(*Config)

func WithJustValidation(validate bool) Option {
	return func(c *Config) {
		c.validate = validate
	}
}

func WithLogger(l *zap.Logger) Option {
	return func(c *Config) {
		c.logger = l
	}
}

func NewFromGlob(glob string, opts ...Option) ([]Collector, error) {
	files, err := filepath.Glob(glob)
	if err != nil {
		return nil, err
	}
	var collectors []Collector

	for _, fName := range files {
		c, err := NewFromFile(fName, opts...)
		if err != nil {
			return nil, err
		}
		collectors = append(collectors, c)
	}
	return collectors, nil
}

func NewFromFile(fpath string, opts ...Option) (Collector, error) {
	fmt.Printf("loading config from file: %q\n", fpath)

	bs, err := os.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	return New(bs, opts...)
}

func New(bs []byte, opts ...Option) (Collector, error) {
	var conf Config

	for _, opt := range opts {
		opt(&conf)
	}

	if err := yaml.Unmarshal(bs, &conf); err != nil {
		return nil, err
	}

	var c Collector
	var err error
	switch conf.Collector.Type {
	case TypeMetric:
		mc, err := metric.NewConfig(
			bs,
			metric.ConfigWithJustValidation(conf.validate),
			metric.ConfigWithLogger(conf.logger),
		)
		if err != nil {
			return nil, err
		}
		c, err = metric.NewCollector(
			mc,
			metric.CollectorWithLogger(conf.logger),
		)

	default:
		return nil, fmt.Errorf("collector type: %v not supported", conf.Collector.Type)
	}

	return c, err
}
