package metric

import (
	"github.com/turbolytics/latte/internal/invoker"
	"github.com/turbolytics/latte/internal/schedule"
	"go.uber.org/zap"
)

type Collector struct {
	config     *config
	logger     *zap.Logger
	schedule   schedule.Schedule
	sinks      map[string]invoker.Sinker
	stateStore invoker.Storer
	sourcer    invoker.Sourcer
	validate   bool
}

func (c *Collector) Name() string {
	return c.config.Name
}

func (c *Collector) InvocationStrategy() invoker.TypeStrategy {
	return c.config.Collector.InvocationStrategy
}

func (c *Collector) Sinks() []invoker.Sinker {
	return nil
}

func (c *Collector) Schedule() invoker.Schedule {
	return c.schedule
}

func (c *Collector) Sourcer() invoker.Sourcer {
	return nil
}

func (c *Collector) Storer() invoker.Storer {
	return c.stateStore
}

func (c *Collector) Transformer() invoker.Transformer {
	return nil
}

type Option func(*Collector)

func WithStateStore(ss invoker.Storer) Option {
	return func(c *Collector) {
		c.stateStore = ss
	}
}

func WithSchedule(sch schedule.Schedule) Option {
	return func(c *Collector) {
		c.schedule = sch
	}
}

func WithSinks(ss map[string]invoker.Sinker) Option {
	return func(c *Collector) {
		c.sinks = ss
	}
}

func WithValidation(validate bool) Option {
	return func(c *Collector) {
		c.validate = validate
	}
}

func WithLogger(l *zap.Logger) Option {
	return func(c *Collector) {
		c.logger = l
	}
}

func NewCollector(conf *config, opts ...Option) (*Collector, error) {
	c := &Collector{
		config: conf,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}
