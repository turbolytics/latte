package collector

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	latteMetric "github.com/turbolytics/latte/internal/collector/metric"
	"github.com/turbolytics/latte/internal/collector/partition"
	"github.com/turbolytics/latte/internal/collector/schedule"
	"github.com/turbolytics/latte/internal/obs"
	"github.com/turbolytics/latte/internal/sink"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"time"
)

var meter = otel.Meter("latte-collector")

type Config interface {
	CollectorName() string
	GetSinks() []sink.Sinker
	GetSchedule() schedule.Schedule
}

type Collector interface {
	Invoke(context.Context) error
}

type InvokerOption func(*Invoker)

func InvokerWithLogger(l *zap.Logger) InvokerOption {
	return func(c *Invoker) {
		c.logger = l
	}
}

type Invoker struct {
	Collector Collector
	Config    Config

	logger *zap.Logger
	now    func() time.Time
}

func (i *Invoker) Close() error {
	ss := i.Config.GetSinks()
	for _, s := range ss {
		s.Close()
	}
	return nil
}

// InvokeHandleError will log any Invoke errors and not return them.
// Useful for async scheduling.
func (i *Invoker) InvokeHandleError(ctx context.Context) {
	if err := i.Invoke(ctx); err != nil {
		i.logger.Error(err.Error())
	}
}

func (i *Invoker) Invoke(ctx context.Context) (err error) {
	start := time.Now().UTC()

	histogram, _ := meter.Float64Histogram(
		"collector.invoke.duration",
		metric.WithUnit("s"),
	)

	counter, _ := meter.Int64Counter(
		"collector.invoke.count",
	)

	defer func() {
		duration := time.Since(start)

		counter.Add(ctx, 1, metric.WithAttributeSet(
			attribute.NewSet(
				attribute.String("result.status_code", obs.ErrToStatus(err)),
				attribute.String("collector.name", i.Config.CollectorName()),
			),
		))

		histogram.Record(ctx, duration.Seconds(), metric.WithAttributeSet(
			attribute.NewSet(
				attribute.String("result.status_code", obs.ErrToStatus(err)),
				attribute.String("collector.name", i.Config.CollectorName()),
			),
		))
	}()

	id := uuid.New()
	i.logger.Info(
		"collector.Invoke",
		zap.String("id", id.String()),
		zap.String("name", i.Config.CollectorName()),
	)
	ctx = context.WithValue(ctx, "id", id)
	return i.Collector.Invoke(ctx)
}

func NewFromGlob(glob string, opts ...RootOption) ([]*Invoker, error) {
	files, err := filepath.Glob(glob)
	if err != nil {
		return nil, err
	}
	var invokers []*Invoker

	for _, fName := range files {
		c, err := NewFromFile(fName, opts...)
		if err != nil {
			return nil, err
		}
		invokers = append(invokers, c)
	}
	return invokers, nil
}

func NewFromFile(fpath string, opts ...RootOption) (*Invoker, error) {
	fmt.Printf("loading config from file: %q\n", fpath)

	bs, err := os.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	return New(bs, opts...)
}

func New(bs []byte, opts ...RootOption) (*Invoker, error) {
	var conf RootConfig

	for _, opt := range opts {
		opt(&conf)
	}

	if err := yaml.Unmarshal(bs, &conf); err != nil {
		return nil, err
	}

	var coll Collector
	var collConfig Config

	var err error
	switch conf.Collector.Type {
	case TypeMetric:
		mc, err := latteMetric.NewConfig(
			bs,
			latteMetric.ConfigWithJustValidation(conf.validate),
			latteMetric.ConfigWithLogger(conf.logger),
		)
		if err != nil {
			return nil, err
		}
		coll, err = latteMetric.NewCollector(
			mc,
			latteMetric.CollectorWithLogger(conf.logger),
		)
		collConfig = mc
	case TypePartition:
		pc, err := partition.NewConfig(
			bs,
			partition.ConfigWithJustValidation(conf.validate),
			partition.ConfigWithLogger(conf.logger),
		)
		if err != nil {
			return nil, err
		}
		coll, err = partition.NewCollector(
			pc,
			partition.CollectorWithLogger(conf.logger),
		)
		collConfig = pc

	default:
		return nil, fmt.Errorf("collector type: %v not supported", conf.Collector.Type)
	}

	i := &Invoker{
		Config:    collConfig,
		Collector: coll,

		logger: conf.logger,
	}
	return i, err
}
