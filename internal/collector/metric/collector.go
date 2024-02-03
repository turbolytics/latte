package metric

import (
	latteMetric "github.com/turbolytics/latte/internal/metric"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"time"
)

var meter = otel.Meter("latte-collector")

type Collector struct {
	Config *Config

	logger *zap.Logger
	now    func() time.Time
}

func (c *Collector) Transform(ms []*latteMetric.Metric) error {
	for _, m := range ms {
		m.Name = c.Config.Metric.Name
		m.Type = c.Config.Metric.Type

		// enrich with tags
		// should these be copied?
		for _, t := range c.Config.Metric.Tags {
			m.Tags[t.Key] = t.Value
		}

	}
	return nil
}

/*
func (c *Collector) Source(ctx context.Context) (ms []*latteMetric.Metric, err error) {
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
				attribute.String("source.type", string(c.Config.Source.Type)),
			),
		))

		meter.Int64ObservableGauge(
			"collector.source.metrics.total",
			metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
				o.Observe(int64(len(ms)), metric.WithAttributeSet(
					attribute.NewSet(
						attribute.String("source.type", string(c.Config.Source.Type)),
					),
				))
				return nil
			}),
		)

	}()

	ms, err = c.Config.Source.MetricSourcer.Source(ctx)
	if err != nil {
		return nil, err
	}

	return ms, nil
}

func (c *Collector) Sink(ctx context.Context, metrics []*latteMetric.Metric) error {

	histogram, _ := meter.Float64Histogram(
		"collector.sink.duration",
		metric.WithUnit("s"),
	)
	// add tags from config

	// need to add a serializer
	for _, m := range metrics {
		bs, err := json.Marshal(m)
		if err != nil {
			return err
		}
		for _, s := range c.Config.Sinks {
			start := time.Now().UTC()

			_, err := s.Sinker.Write(bs)

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
	}
	return nil
}

func (c *Collector) invokeTick(ctx context.Context, id uuid.UUID) ([]*latteMetric.Metric, error) {
	ms, err := c.Source(ctx)
	if err != nil {
		return nil, err
	}

	// only sink if metrics are present:
	if len(ms) == 0 {
		c.logger.Warn(
			"collector.Invoke",
			zap.String("msg", "no metrics found"),
			zap.String("id", id.String()),
			zap.String("name", c.Config.Name),
		)
		return ms, err
	}

	if err = c.Transform(ms); err != nil {
		return ms, err
	}

	if err = c.Sink(ctx, ms); err != nil {
		return ms, err
	}

	return ms, err
}

func (c *Collector) invokeWindowSourceAndSave(ctx context.Context, id uuid.UUID, window timeseries.Window) ([]*latteMetric.Metric, error) {
	c.logger.Info(
		"collector.invokeWindow",
		zap.String("msg", "invoking for window"),
		zap.String("window.start", window.Start.String()),
		zap.String("window.end", window.End.String()),
		zap.String("id", id.String()),
		zap.String("name", c.Config.CollectorName()),
	)

	// it is passed the window, collect data for the window
	// get start of window and end of window
	ctx = context.WithValue(ctx, "window.start", window.Start)
	ctx = context.WithValue(ctx, "window.end", window.End)
	ms, err := c.Source(ctx)
	if err != nil {
		return ms, err
	}

	err = c.Config.StateStore.Storer.SaveInvocation(&state.Invocation{
		CollectorName: c.Config.Name,
		Time:          c.now(),
		Window:        &window,
	})
	return ms, err
}

// invokeWindow uses the state store to check if a full window has elapsed.
// invokeWindow will only source data when a full window has elapsed.
// TODO - What happens when a window is changed in the config?
func (c *Collector) invokeWindow(ctx context.Context, id uuid.UUID) ([]*latteMetric.Metric, error) {
	var ms []*latteMetric.Metric
	var err error
	// TODO invokeWindow should handle gaps in windows.
	i, err := c.Config.StateStore.Storer.MostRecentInvocation(
		ctx,
		c.Config.Name,
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
		*(c.Config.Source.MetricSourcer.Window()),
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
		ms, err = c.invokeWindowSourceAndSave(ctx, id, windows[0])
	default:
		// multiple windows have been found, currently
		// not supported
		c.logger.Error(
			"collector.invokeWindow",
			zap.String("msg", "multiple windows detected"),
			zap.Int("windows", len(windows)),
			zap.String("id", id.String()),
			zap.String("name", c.Config.Name),
		)
		return nil, fmt.Errorf("backfilling multiple windows not yet supported: %v", windows)
	}

	if err = c.Transform(ms); err != nil {
		return ms, err
	}

	if err = c.Sink(ctx, ms); err != nil {
		return ms, err
	}

	return ms, err
}

func (c *Collector) Invoke(ctx context.Context) (err error) {
	id := ctx.Value("id").(uuid.UUID)
	// Collector supports multiple sourcing strategies.
	// The simplest is "tick" strategy which just invokes
	// the sourcer without any additional state necessary
	// The windowing strategy may result in multiple source
	// invocations for each window that needs to be executed.
	switch c.Config.Source.Strategy {
	case source.TypeStrategyTick:
		_, err = c.invokeTick(ctx, id)
	case source.TypeStrategyHistoricTumblingWindow:
		_, err = c.invokeWindow(ctx, id)
	default:
		return fmt.Errorf("strategy: %q not supported", c.Config.Source.Strategy)
	}

	return err
}

*/

type CollectorOption func(*Collector)

func CollectorWithLogger(l *zap.Logger) CollectorOption {
	return func(c *Collector) {
		c.logger = l
	}
}

func NewCollector(config *Config, opts ...CollectorOption) (*Collector, error) {
	c := &Collector{
		Config: config,
		now: func() time.Time {
			return time.Now().UTC()
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}
