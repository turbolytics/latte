package metric

import (
	"github.com/turbolytics/latte/internal/collector/template"
	"github.com/turbolytics/latte/internal/invoker"
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

func (c config) GetSchedule() schedule.Config {
	return c.Schedule
}

func (c config) GetSource() source.Config {
	return c.Source
}

/*
func (c Collector) GetSinks() []_type.Sinker {
	var ss []_type.Sinker
	for _, s := range c.Sinks {
		ss = append(ss, s.Sinker)
	}
	return ss
}
*/

func (c config) CollectorName() string {
	return c.Name
}

/*
// initSinks initializes all the outputs
func initSinks(c *Collector) error {
	for k, v := range c.Sinks {
		if err := v.Init(c.validate, c.logger); err != nil {
			return err
		}
		c.Sinks[k] = v
	}
	return nil
}
*/

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

	/*
		for _, opt := range opts {
			opt(&conf)
		}
	*/

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
		if err := conf.StateStore.Init(); err != nil {
			return nil, err
		}

			if err := conf.Source.Init(conf.logger, conf.validate); err != nil {
				return nil, err
			}

			if err := initSinks(&conf); err != nil {
				return nil, err
			}

	*/
	return &conf, nil
}

type Collector struct {
	config   *config
	logger   *zap.Logger
	validate bool
}

func (c *Collector) Name() string {
	return ""
}

func (c *Collector) InvocationStrategy() invoker.TypeStrategy {
	return invoker.TypeStrategyTick
}

func (c *Collector) Sinks() []invoker.Sinker {
	return nil
}

func (c *Collector) Schedule() invoker.Schedule {
	return nil
}

func (c *Collector) Sourcer() invoker.Sourcer {
	return nil
}

func (c *Collector) Storer() state.Storer {
	return nil
}

func (c *Collector) Transformer() invoker.Transformer {
	return nil
}

type Option func(*Collector)

func WithJustValidation(validate bool) Option {
	return func(c *Collector) {
		c.validate = validate
	}
}

func WithLogger(l *zap.Logger) Option {
	return func(c *Collector) {
		c.logger = l
	}
}
func NewCollectorFromConfig(bs []byte, opts ...Option) (*Collector, error) {
	conf, err := NewConfig(bs)
	if err != nil {
		return nil, err
	}

	collector := &Collector{
		config: conf,
	}

	for _, opt := range opts {
		opt(collector)
	}

	return collector, nil
}
