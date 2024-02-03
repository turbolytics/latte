package collector

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	latteMetric "github.com/turbolytics/latte/internal/collector/metric"
	"github.com/turbolytics/latte/internal/collector/partition"
	"github.com/turbolytics/latte/internal/obs"
	"github.com/turbolytics/latte/internal/record"
	"github.com/turbolytics/latte/internal/schedule"
	"github.com/turbolytics/latte/internal/sink/type"
	"github.com/turbolytics/latte/internal/source"
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
	GetSinks() []_type.Sinker
	GetSchedule() schedule.Config
	GetSource() source.Config
}

type Invoker struct {
	Config Config

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

func (i *Invoker) Source(ctx context.Context) (sr record.Result, err error) {
	id := ctx.Value("id").(uuid.UUID)
	start := time.Now().UTC()

	histogram, _ := meter.Float64Histogram(
		"collector.source.duration",
		metric.WithUnit("s"),
	)

	s := i.Config.GetSource()

	defer func() {
		duration := time.Since(start)

		histogram.Record(ctx, duration.Seconds(), metric.WithAttributeSet(
			attribute.NewSet(
				attribute.String("result.status_code", obs.ErrToStatus(err)),
				attribute.String("source.type", string(s.Type)),
			),
		))

		meter.Int64ObservableGauge(
			"collector.source.metrics.total",
			metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
				o.Observe(int64(len(sr.Records())), metric.WithAttributeSet(
					attribute.NewSet(
						attribute.String("source.type", string(s.Type)),
					),
				))
				return nil
			}),
		)
	}()

	sr, err = s.Sourcer.Source(ctx)

	i.logger.Debug("collector.Source",
		zap.String("source.strategy", string(i.Config.GetSource().Strategy)),
		zap.String("id", id.String()),
		zap.String("name", i.Config.CollectorName()),
		zap.Int("results.count", len(sr.Records())),
	)
	return sr, err
}

func (i *Invoker) invokeTick(ctx context.Context) error {
	id := ctx.Value("id").(uuid.UUID)

	i.logger.Debug("collector.invokeTick",
		zap.String("id", id.String()),
		zap.String("name", i.Config.CollectorName()),
	)

	sr, err := i.Source(ctx)
	if err != nil {
		return err
	}

	if len(sr.Records()) == 0 {
		i.logger.Warn(
			"collector.Invoke",
			zap.String("msg", "no results found"),
			zap.String("id", id.String()),
			zap.String("name", i.Config.CollectorName()),
		)
	}

	// how to get additional context to transform function?
	/*
		if err := sr.Transform(); err != nil {
			return err
		}
	*/

	if err = i.Sink(ctx, sr); err != nil {
		return err
	}

	return nil
}

func (i *Invoker) Sink(ctx context.Context, res record.Result) error {

	histogram, _ := meter.Float64Histogram(
		"collector.sink.duration",
		metric.WithUnit("s"),
	)
	// add tags from config
	rs := res.Records()

	sinks := i.Config.GetSinks()

	// need to add a serializer
	for _, r := range rs {
		bs, err := r.Bytes()
		if err != nil {
			return err
		}
		for _, s := range sinks {
			start := time.Now().UTC()

			_, err := s.Write(bs)

			duration := time.Since(start)
			histogram.Record(ctx, duration.Seconds(), metric.WithAttributeSet(
				attribute.NewSet(
					attribute.String("result.status_code", obs.ErrToStatus(err)),
					attribute.String("sink.name", string(s.Type())),
				),
			))

			if err != nil {
				return err
			}

		}
	}
	return nil
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

	s := i.Config.GetSource()
	switch s.Strategy {
	case source.TypeStrategyHistoricTumblingWindow:
		fmt.Println("tumbling window")
	case source.TypeStrategyTick:
		i.invokeTick(ctx)
	default:
		return fmt.Errorf("strategy: %q not supported", s.Strategy)
	}

	return nil
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

	var collConfig Config
	var err error

	switch conf.Collector.Type {
	case TypeMetric:
		collConfig, err = latteMetric.NewConfig(
			bs,
			latteMetric.ConfigWithJustValidation(conf.validate),
			latteMetric.ConfigWithLogger(conf.logger),
		)
	case TypePartition:
		collConfig, err = partition.NewConfig(
			bs,
			partition.ConfigWithJustValidation(conf.validate),
			partition.ConfigWithLogger(conf.logger),
		)
	default:
		return nil, fmt.Errorf("collector type: %v not supported", conf.Collector.Type)
	}

	if err != nil {
		return nil, err
	}

	i := &Invoker{
		Config: collConfig,

		logger: conf.logger,
	}
	return i, err
}
