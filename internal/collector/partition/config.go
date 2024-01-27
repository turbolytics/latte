package partition

import (
	"github.com/turbolytics/collector/internal/collector/schedule"
	"github.com/turbolytics/collector/internal/collector/sink"
	"github.com/turbolytics/collector/internal/collector/source"
	"github.com/turbolytics/collector/internal/collector/state"
	"github.com/turbolytics/collector/internal/collector/template"
	"github.com/turbolytics/collector/internal/sinks"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Name       string
	Schedule   schedule.Schedule
	StateStore state.Config `yaml:"state_store"`
	Source     source.Config
	Sinks      map[string]sink.Config

	logger *zap.Logger
	// validate will skip initializing network dependencies
	validate bool
}

func (c Config) CollectorName() string {
	return c.Name
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

func defaults(c *Config) error {
	(&c.Source).SetDefaults()

	return nil
}

type validater interface {
	Validate() error
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

func initSinks(c *Config) error {
	for k, v := range c.Sinks {
		if err := v.Init(c.validate); err != nil {
			return err
		}
		c.Sinks[k] = v
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

	if err := conf.Source.Init(); err != nil {
		return nil, err
	}

	if err := initSinks(&conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
