package metric

import (
	"context"
	"fmt"
	"github.com/turbolytics/collector/internal/collector/metric/sources/mongodb"
	"github.com/turbolytics/collector/internal/collector/metric/sources/postgres"
	"github.com/turbolytics/collector/internal/collector/metric/sources/prometheus"
	"github.com/turbolytics/collector/internal/collector/schedule"
	"github.com/turbolytics/collector/internal/collector/sink"
	"github.com/turbolytics/collector/internal/collector/source"
	"github.com/turbolytics/collector/internal/collector/state"
	"github.com/turbolytics/collector/internal/collector/template"
	"github.com/turbolytics/collector/internal/metrics"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type validater interface {
	Validate() error
}

type Tag struct {
	Key   string
	Value string
}

type Metric struct {
	Name string
	Type metrics.Type
	Tags []Tag
}

type Config struct {
	Name       string
	Metric     Metric
	Schedule   schedule.Schedule
	Source     source.Config
	Sinks      map[string]sink.Config
	StateStore state.Config `yaml:"state_store"`

	logger *zap.Logger
	// validate will skip initializing network dependencies
	validate bool
}

type ConfigOption func(*Config)

func ConfigWithJustValidation(validate bool) ConfigOption {
	return func(c *Config) {
		c.validate = validate
	}
}

func ConfigWithLogger(l *zap.Logger) ConfigOption {
	return func(c *Config) {
		c.logger = l
	}
}

// initSource initializes the correct source.
func initSource(c *Config) error {
	var s source.MetricSourcer
	var err error
	switch c.Source.Type {
	case source.TypePostgres:
		s, err = postgres.NewFromGenericConfig(
			c.Source.Config,
			c.validate,
		)
	case source.TypeMongoDB:
		s, err = mongodb.NewFromGenericConfig(
			context.TODO(),
			c.Source.Config,
			c.validate,
		)
	case source.TypePrometheus:
		s, err = prometheus.NewFromGenericConfig(
			c.Source.Config,
			prometheus.WithLogger(c.logger),
		)
	default:
		return fmt.Errorf("source type: %q unknown", c.Source.Type)
	}

	if err != nil {
		return err
	}

	c.Source.MetricSourcer = s
	return nil
}

// initSinks initializes all the outputs
func initSinks(c *Config) error {
	for k, v := range c.Sinks {
		if err := v.Init(c.validate); err != nil {
			return err
		}
		c.Sinks[k] = v
	}
	return nil
}

func defaults(c *Config) error {
	(&c.Source).SetDefaults()

	return nil
}

func validate(c Config) error {
	validaters := []validater{
		c.Schedule,
		c.Source,
		c.StateStore,
	}

	for _, v := range validaters {
		if err := v.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// NewConfig initializes a config from yaml bytes.
// NewConfig initializes all subtypes as well.
func NewConfig(raw []byte, opts ...ConfigOption) (*Config, error) {
	var conf Config

	for _, opt := range opts {
		opt(&conf)
	}

	bs, err := template.Parse(raw)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(bs, &conf); err != nil {
		return nil, err
	}

	if err := defaults(&conf); err != nil {
		return nil, err
	}

	if err := validate(conf); err != nil {
		return nil, err
	}

	if err := conf.StateStore.Init(); err != nil {
		return nil, err
	}

	if err := initSource(&conf); err != nil {
		return nil, err
	}

	if err := initSinks(&conf); err != nil {
		return nil, err
	}
	return &conf, nil
}
