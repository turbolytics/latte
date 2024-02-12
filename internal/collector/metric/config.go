package metric

import (
	collconf "github.com/turbolytics/latte/internal/collector/config"
	"github.com/turbolytics/latte/internal/collector/template"
	"github.com/turbolytics/latte/internal/metric"
	"github.com/turbolytics/latte/internal/schedule"
	"github.com/turbolytics/latte/internal/sink"
	"github.com/turbolytics/latte/internal/source"
	"github.com/turbolytics/latte/internal/state"
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
	Type metric.Type
	Tags []Tag
}

type config struct {
	Collector  collconf.Collector
	Name       string
	Metric     Metric
	Schedule   schedule.Config
	Source     source.Config
	Sinks      map[string]sink.Config
	StateStore state.Config `yaml:"state_store"`

	logger *zap.Logger
	// validate will skip initializing network dependencies
	validate bool
}

func defaults(c *config) error {
	(&c.Source).SetDefaults()

	return nil
}

func validate(c config) error {
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
func NewConfig(raw []byte) (*config, error) {
	var conf config

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

	/*
		if err := conf.Source.Init(conf.logger, conf.validate); err != nil {
			return nil, err
		}
	*/
	return &conf, nil
}
