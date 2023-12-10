package internal

import (
	"fmt"
	"github.com/turbolytics/collector/internal/metrics"
	"github.com/turbolytics/collector/internal/sinks"
	"github.com/turbolytics/collector/internal/sources"
	"github.com/turbolytics/collector/internal/sources/postgres"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Metric struct {
	Name string
	Type metrics.Type
}

type Schedule struct {
	Interval time.Duration
	Cron     string
}

type Collector struct {
	Type   string // enum
	Config map[string]any
}

type Sink struct {
	Type   string
	Sinker sinks.Sinker
	Config map[string]any
}

type Source struct {
	Type    sources.Type
	Sourcer sources.Sourcer
	Config  map[string]any
}

type Config struct {
	Name     string
	Metric   Metric
	Schedule Schedule
	Source   Source
	Sinks    map[string]Sink
}

func NewConfigFromFile(name string) (*Config, error) {
	bs, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return NewConfig(bs)
}

// initSource initializes the correct source.
func initSource(c *Config) error {
	return nil
}

// NewConfig initializes a config from yaml bytes.
// NewConfig initializes all subtypes as well.
func NewConfig(bs []byte) (*Config, error) {
	var conf Config
	if err := yaml.Unmarshal(bs, &conf); err != nil {
		return nil, err
	}

	var s sources.Sourcer
	var err error
	switch conf.Source.Type {
	case sources.TypePostgres:
		s, err = postgres.NewFromGenericConfig(conf.Source.Config)
	}

	if err != nil {
		return nil, err
	}

	conf.Source.Sourcer = s

	fmt.Printf("%+v\n", conf)
	return &conf, nil
}
