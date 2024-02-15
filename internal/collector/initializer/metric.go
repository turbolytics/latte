package initializer

import (
	"github.com/turbolytics/latte/internal/collector/metric"
	"github.com/turbolytics/latte/internal/schedule"
	"github.com/turbolytics/latte/internal/transform"
	"go.uber.org/zap"
)

func NewMetricCollectorFromConfig(bs []byte, validate bool, l *zap.Logger) (*metric.Collector, error) {
	conf, err := metric.NewConfig(bs)
	if err != nil {
		return nil, err
	}

	stateStore, err := NewStorer(
		conf.StateStore,
		l,
	)

	if err != nil {
		return nil, err
	}

	sinks, err := NewSinks(
		conf.Sinks,
		l,
		validate,
	)

	if err != nil {
		return nil, err
	}

	sourcer, err := NewSourcer(
		conf.Source,
		l,
		validate,
	)
	if err != nil {
		return nil, err
	}

	sch := schedule.New(conf.Schedule)

	coll, err := metric.NewCollector(
		conf,
		metric.WithLogger(l),
		metric.WithValidation(validate),
		metric.WithSchedule(sch),
		metric.WithSinks(sinks),
		metric.WithSourcer(sourcer),
		metric.WithStateStore(stateStore),
		metric.WithTransformer(transform.Noop{}),
	)

	return coll, nil
}
