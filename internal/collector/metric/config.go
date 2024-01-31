package metric

import (
	"fmt"
	"github.com/turbolytics/latte/internal/collector/schedule"
	"github.com/turbolytics/latte/internal/collector/sink"
	configSource "github.com/turbolytics/latte/internal/collector/source"
	"github.com/turbolytics/latte/internal/collector/state"
	"github.com/turbolytics/latte/internal/collector/template"
	"github.com/turbolytics/latte/internal/metrics"
	"github.com/turbolytics/latte/internal/sinks"
	"github.com/turbolytics/latte/internal/source"
	"github.com/turbolytics/latte/internal/source/metric/postgres"
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
	Source     configSource.Config
	Sinks      map[string]sink.Config
	StateStore state.Config `yaml:"state_store"`

	logger *zap.Logger
	// validate will skip initializing network dependencies
	validate bool
}

func (c Config) GetMetric() *Metric {
	return &c.Metric
}

func (c Config) GetSource() configSource.Config {
	return c.Source
}

func (c Config) GetSchedule() schedule.Schedule {
	return c.Schedule
}

func (c Config) GetSinks() []sinks.Sinker {
	var ss []sinks.Sinker
	for _, s := range c.Sinks {
		ss = append(ss, s.Sinker)
	}
	return ss
}

func (c Config) CollectorName() string {
	return c.Name
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
	var s source.Sourcer
	var err error
	switch c.Source.Type {
	case configSource.TypePostgres:
		s, err = postgres.NewFromGenericConfig(
			c.Source.Config,
			c.validate,
		)
		/*
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

		*/
	default:
		return fmt.Errorf("source type: %q unknown", c.Source.Type)
	}

	if err != nil {
		return err
	}

	c.Source.Sourcer = s
	return nil
}

// initSinks initializes all the outputs
func initSinks(c *Config) error {
	for k, v := range c.Sinks {
		if err := v.Init(c.validate, c.logger); err != nil {
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
