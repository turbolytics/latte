package invoker

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/turbolytics/latte/internal/obs"
	"github.com/turbolytics/latte/internal/record"
	"github.com/turbolytics/latte/internal/sink"
	"github.com/turbolytics/latte/internal/source"
	"github.com/turbolytics/latte/internal/timeseries"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
	"io"
	"time"
)

var meter = otel.Meter("latte-collector")

type TypeStrategy string

const (
	TypeStrategyHistoricTumblingWindow TypeStrategy = "historic_tumbling_window"
	TypeStrategyIncremental            TypeStrategy = "incremental"
	TypeStrategyTick                   TypeStrategy = "tick"
)

type Invocation struct {
	CollectorName string
	Time          time.Time
	Window        *timeseries.Window
}

func (i Invocation) End() *time.Time {
	if i.Window != nil {
		return &i.Window.End
	}
	return nil
}

type Sourcer interface {
	Source(ctx context.Context) (record.Result, error)
	WindowDuration() *time.Duration
	Type() source.Type
}

// Sinker is responsible for sinking
// TODO - Starting with an io.Writer for right now.
type Sinker interface {
	Write(context.Context, record.Record) (int, error)
	Close() error
	Type() sink.Type
	Flush(context.Context) error
}

type Schedule interface {
	Interval() *time.Duration
	Cron() *string
}

type Storer interface {
	io.Closer

	MostRecentInvocation(ctx context.Context, collectorName string) (*Invocation, error)
	SaveInvocation(invocation *Invocation) error
}

type Transformer interface {
	Transform(record.Result) error
}

type Collector interface {
	Name() string
	InvocationStrategy() TypeStrategy
	Sinks() []Sinker
	Schedule() Schedule
	Sourcer() Sourcer
	Storer() Storer
	Transformer() Transformer
}

type Option func(*Invoker)

func WithLogger(l *zap.Logger) Option {
	return func(i *Invoker) {
		i.logger = l
	}
}

func WithStartTime(t time.Time) Option {
	return func(i *Invoker) {
		i.now = func() time.Time {
			return t
		}
	}
}

type Invoker struct {
	Collector Collector

	logger *zap.Logger
	now    func() time.Time
}

func (i *Invoker) Close() error {
	ss := i.Collector.Sinks()
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

	s := i.Collector.Sourcer()

	defer func() {
		duration := time.Since(start)

		histogram.Record(ctx, duration.Seconds(), metric.WithAttributeSet(
			attribute.NewSet(
				attribute.String("result.status_code", obs.ErrToStatus(err)),
				attribute.String("source.type", string(s.Type())),
			),
		))

		meter.Int64ObservableGauge(
			"collector.source.metrics.total",
			metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
				o.Observe(int64(len(sr.Records())), metric.WithAttributeSet(
					attribute.NewSet(
						attribute.String("source.type", string(s.Type())),
					),
				))
				return nil
			}),
		)
	}()

	sr, err = s.Source(ctx)
	if err != nil {
		return nil, err
	}

	var numRecords int
	if sr != nil {
		numRecords = len(sr.Records())
	}

	i.logger.Debug("collector.Source",
		zap.String("collector.invocation_strategy", string(i.Collector.InvocationStrategy())),
		zap.String("id", id.String()),
		zap.String("name", i.Collector.Name()),
		zap.Int("results.count", numRecords),
	)
	return sr, err
}

func (i *Invoker) invokeWindowSourceAndSave(ctx context.Context, window timeseries.Window) error {
	id := ctx.Value("id").(uuid.UUID)
	i.logger.Info(
		"invoker.invokeWindowSourceAndSave",
		zap.String("msg", "invoking for window"),
		zap.String("window.start", window.Start.String()),
		zap.String("window.end", window.End.String()),
		zap.String("id", id.String()),
		zap.String("name", i.Collector.Name()),
	)
	ctx = context.WithValue(ctx, "window.start", window.Start)
	ctx = context.WithValue(ctx, "window.end", window.End)

	r, err := i.Source(ctx)
	if err != nil {
		return err
	}

	tr := i.Collector.Transformer()
	if err := tr.Transform(r); err != nil {
		return err
	}

	if err = i.Sink(ctx, r); err != nil {
		return err
	}

	state := i.Collector.Storer()
	err = state.SaveInvocation(&Invocation{
		CollectorName: i.Collector.Name(),
		Time:          i.now(),
		Window:        &window,
	})

	return err
}

func (i *Invoker) invokeTick(ctx context.Context) error {
	id := ctx.Value("id").(uuid.UUID)
	start := ctx.Value("invocation.start").(time.Time)
	ctx = context.WithValue(ctx, "window.start", start)

	i.logger.Debug("collector.invokeTick",
		zap.String("id", id.String()),
		zap.String("name", i.Collector.Name()),
	)

	sr, err := i.Source(ctx)
	if err != nil {
		return err
	}

	if sr == nil || len(sr.Records()) == 0 {
		i.logger.Warn(
			"collector.Invoke",
			zap.String("msg", "no results found"),
			zap.String("id", id.String()),
			zap.String("name", i.Collector.Name()),
		)
	}

	tr := i.Collector.Transformer()
	if err := tr.Transform(sr); err != nil {
		return err
	}

	if err = i.Sink(ctx, sr); err != nil {
		return err
	}

	return nil
}

func (i *Invoker) invokeHistoricTumblingWindow(ctx context.Context) error {
	id := ctx.Value("id").(uuid.UUID)
	i.logger.Debug("collector.invokeHistoricTumblingWindow",
		zap.String("id", id.String()),
		zap.String("name", i.Collector.Name()),
	)

	storer := i.Collector.Storer()
	inv, err := storer.MostRecentInvocation(ctx, i.Collector.Name())
	if err != nil {
		return err
	}

	var lastWindowEnd *time.Time
	if inv != nil {
		lastWindowEnd = inv.End()
	}

	hw := timeseries.NewHistoricTumblingWindower(
		timeseries.WithHistoricTumblingWindowerNow(i.now),
	)

	s := i.Collector.Sourcer()
	windows, err := hw.FullWindowsSince(
		lastWindowEnd,
		*(s.WindowDuration()),
	)
	if err != nil {
		return err
	}

	switch len(windows) {
	case 0:
		// no full windows have passed, just return
		return nil
	case 1:
		// a single window is available, collect it
		return i.invokeWindowSourceAndSave(ctx, windows[0])
	default:
		// multiple windows have been found, currently
		// not supported
		i.logger.Error(
			"collector.invokeHistoricTumblingWindow",
			zap.String("msg", "multiple windows detected"),
			zap.Int("windows", len(windows)),
			zap.String("id", id.String()),
			zap.String("name", i.Collector.Name()),
		)
		return fmt.Errorf("backfilling multiple windows not yet supported: %v", windows)
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

	sinks := i.Collector.Sinks()

	// need to add a serializer
	for _, s := range sinks {
		start := time.Now().UTC()
		var err error

		for _, r := range rs {

			_, err = s.Write(ctx, r)

			if err != nil {
				break
			}

		}

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

		if err := s.Flush(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (i *Invoker) Invoke(ctx context.Context) (err error) {
	start := i.now()

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
				attribute.String("collector.name", i.Collector.Name()),
			),
		))

		histogram.Record(ctx, duration.Seconds(), metric.WithAttributeSet(
			attribute.NewSet(
				attribute.String("result.status_code", obs.ErrToStatus(err)),
				attribute.String("collector.name", i.Collector.Name()),
			),
		))
	}()

	id := uuid.New()
	i.logger.Info(
		"collector.Invoke",
		zap.String("invocation.start", start.String()),
		zap.String("id", id.String()),
		zap.String("name", i.Collector.Name()),
	)
	ctx = context.WithValue(ctx, "id", id)
	ctx = context.WithValue(ctx, "invocation.start", start)

	strat := i.Collector.InvocationStrategy()
	switch strat {
	case TypeStrategyHistoricTumblingWindow:
		return i.invokeHistoricTumblingWindow(ctx)
	case TypeStrategyTick:
		return i.invokeTick(ctx)
	default:
		return fmt.Errorf("strategy: %q not supported", strat)
	}

	return nil
}

func New(collector Collector, opts ...Option) (*Invoker, error) {
	i := &Invoker{
		Collector: collector,
		now: func() time.Time {
			return time.Now().UTC()
		},
	}

	for _, opt := range opts {
		opt(i)
	}

	return i, nil
}
