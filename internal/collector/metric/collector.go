package metric

import (
	"fmt"
	"github.com/turbolytics/latte/internal/invoker"
	"github.com/turbolytics/latte/internal/metric"
	"github.com/turbolytics/latte/internal/record"
	"github.com/turbolytics/latte/internal/schedule"
	"go.uber.org/zap"
)

type Transformer struct {
	wrapped invoker.Transformer
	config  *config
}

func (t Transformer) Transform(r record.Result) error {
	// generic invoker passes control back to concrete metrics
	// collector
	metricResult, ok := r.(*metric.Metrics)
	if !ok {
		return fmt.Errorf("cannot convert %v to *metric.Metrics result", r)
	}

	for _, m := range metricResult.Metrics {
		m.Name = t.config.Metric.Name
		m.Type = t.config.Metric.Type

		// enrich with tags
		// should these be copied?
		for _, t := range t.config.Metric.Tags {
			m.Tags[t.Key] = t.Value
		}
	}

	return t.wrapped.Transform(r)
}

type Collector struct {
	config      *config
	logger      *zap.Logger
	schedule    schedule.Schedule
	sinks       map[string]invoker.Sinker
	stateStore  invoker.Storer
	sourcer     invoker.Sourcer
	transformer invoker.Transformer
	validate    bool
}

func (c *Collector) Name() string {
	return c.config.Name
}

func (c *Collector) InvocationStrategy() invoker.TypeStrategy {
	return c.config.Collector.InvocationStrategy
}

func (c *Collector) Sinks() []invoker.Sinker {
	var sinks []invoker.Sinker
	for _, sink := range c.sinks {
		sinks = append(sinks, sink)
	}
	return sinks
}

func (c *Collector) Schedule() invoker.Schedule {
	return c.schedule
}

func (c *Collector) Sourcer() invoker.Sourcer {
	return c.sourcer
}

func (c *Collector) Storer() invoker.Storer {
	return c.stateStore
}

func (c *Collector) Transformer() invoker.Transformer {
	// this is the configuration defined transform.
	// Metrics also need to be enriched with default behavior.
	return c.transformer
}

type Option func(*Collector)

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

func WithSourcer(s invoker.Sourcer) Option {
	return func(c *Collector) {
		c.sourcer = s
	}
}

func WithStateStore(ss invoker.Storer) Option {
	return func(c *Collector) {
		c.stateStore = ss
	}
}

func WithTransformer(t invoker.Transformer) Option {
	return func(c *Collector) {
		c.transformer = Transformer{
			config:  c.config,
			wrapped: t,
		}
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
