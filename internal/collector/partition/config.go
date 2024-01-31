package partition

import (
	"github.com/turbolytics/latte/internal/collector/schedule"
	"github.com/turbolytics/latte/internal/collector/source"
	"github.com/turbolytics/latte/internal/collector/state"
	"github.com/turbolytics/latte/internal/collector/transform"
	"github.com/turbolytics/latte/internal/sink"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Name       string
	Schedule   schedule.Schedule
	StateStore state.Config `yaml:"state_store"`
	Source     source.Config
	Transform  transform.Config
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

func (c Config) GetSinks() []sink.Sinker {
	var ss []sink.Sinker
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
		if err := v.Init(c.validate, c.logger); err != nil {
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

	// TODO figure out how to maintain this as a template but still allow templatizing the rest of the config
	// We could go through and only parse the allowed fields.
	// Will end users need to provide templated values for the partition?
	// Not sure if we can escape it.
	/*
		bs, err := template.Parse(raw)
		if err != nil {
			return nil, err
		}
	*/

	if err := yaml.Unmarshal(raw, &conf); err != nil {
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

	if err := conf.Transform.Init(); err != nil {
		return nil, err
	}

	if err := initSinks(&conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
