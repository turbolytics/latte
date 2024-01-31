package partition

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/turbolytics/latte/internal/collector/state"
	"github.com/turbolytics/latte/internal/obs"
	"github.com/turbolytics/latte/internal/partition"
	"github.com/turbolytics/latte/internal/source"
	"github.com/turbolytics/latte/internal/timeseries"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
	"time"
)

var meter = otel.Meter("latte-collector")

type Collector struct {
	config *Config
	logger *zap.Logger
	now    func() time.Time
}

func (c *Collector) invokeHistoricTumblingWindow(ctx context.Context) ([]*partition.Partition, error) {
	i, err := c.config.StateStore.Storer.MostRecentInvocation(
		ctx,
		c.config.Name,
	)

	if err != nil {
		return nil, err
	}
	var lastWindowEnd *time.Time
	if i != nil {
		lastWindowEnd = i.End()
	}

	hw := timeseries.NewHistoricTumblingWindower(
		timeseries.WithHistoricTumblingWindowerNow(c.now),
	)

	windows, err := hw.FullWindowsSince(
		lastWindowEnd,
		*(c.config.Source.Window()),
	)

	if err != nil {
		return nil, err
	}

	switch len(windows) {
	case 0:
		// no full windows have passed, just return
		return nil, nil
	case 1:
		// a single window is available, collect it
		p, err := c.invokeWindowSourceAndSave(ctx, windows[0])
		if err != nil {
			return []*partition.Partition{p}, err
		}

		if err := c.Sink(ctx, p); err != nil {
			return []*partition.Partition{p}, err
		}

		return []*partition.Partition{p}, nil
	}

	return nil, nil
}

func (c *Collector) Sink(ctx context.Context, p *partition.Partition) error {
	histogram, _ := meter.Float64Histogram(
		"collector.sink.duration",
		metric.WithUnit("s"),
	)

	out, err := c.config.Transform.Transformer.Transform(p)
	if err != nil {
		return err
	}

	for _, s := range c.config.Sinks {
		start := time.Now().UTC()

		_, err := s.Sinker.Write(out)

		duration := time.Since(start)
		histogram.Record(ctx, duration.Seconds(), metric.WithAttributeSet(
			attribute.NewSet(
				attribute.String("result.status_code", obs.ErrToStatus(err)),
				attribute.String("sink.name", string(s.Type)),
			),
		))

		if err != nil {
			return err
		}

	}
	return nil
}

/*
type source.Result interface {
	Transform() error
	Rows() ([]any, error)
}

- Source becomes generic
Source(ctx) (source.Result, error) {
}

- Sink becomes generic
Sink(ctx, sr source.Result) error {
}



*/

func (c *Collector) Source(ctx context.Context) (p *partition.Partition, err error) {
	start := time.Now().UTC()

	histogram, _ := meter.Float64Histogram(
		"collector.source.duration",
		metric.WithUnit("s"),
	)

	defer func() {
		duration := time.Since(start)

		histogram.Record(ctx, duration.Seconds(), metric.WithAttributeSet(
			attribute.NewSet(
				attribute.String("result.status_code", obs.ErrToStatus(err)),
				attribute.String("source.type", string(c.config.Source.Type)),
			),
		))

		meter.Int64ObservableGauge(
			"collector.source.metrics.total",
			metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
				o.Observe(int64(1), metric.WithAttributeSet(
					attribute.NewSet(
						attribute.String("source.type", string(c.config.Source.Type)),
					),
				))
				return nil
			}),
		)

	}()

	p, err = c.config.Source.PartitionSourcer.Source(ctx)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (c *Collector) invokeWindowSourceAndSave(ctx context.Context, window timeseries.Window) (*partition.Partition, error) {
	id := ctx.Value("id").(uuid.UUID)

	c.logger.Info(
		"collector.invokeWindow",
		zap.String("msg", "invoking for window"),
		zap.String("window.start", window.Start.String()),
		zap.String("window.end", window.End.String()),
		zap.String("id", id.String()),
		zap.String("name", c.config.CollectorName()),
	)

	ctx = context.WithValue(ctx, "window.start", window.Start)
	ctx = context.WithValue(ctx, "window.end", window.End)
	p, err := c.Source(ctx)
	if err != nil {
		return p, err
	}

	err = c.config.StateStore.Storer.SaveInvocation(&state.Invocation{
		CollectorName: c.config.CollectorName(),
		Time:          c.now(),
		Window:        &window,
	})
	return p, err
}

func (c *Collector) Invoke(ctx context.Context) error {
	// get the last windows that completed
	switch c.config.Source.Strategy {
	case source.TypeStrategyHistoricTumblingWindow:
		_, err := c.invokeHistoricTumblingWindow(ctx)
		return err
	default:
		return fmt.Errorf("strategy: %q not supported", c.config.Source.Strategy)
	}

	return nil
}

type CollectorOption func(*Collector)

func CollectorWithLogger(l *zap.Logger) CollectorOption {
	return func(c *Collector) {
		c.logger = l
	}
}

func NewCollector(config *Config, opts ...CollectorOption) (*Collector, error) {
	c := &Collector{
		config: config,
		now: func() time.Time {
			return time.Now().UTC()
		},
	}
	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}
