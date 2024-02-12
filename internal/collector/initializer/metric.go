package initializer

import (
	"fmt"
	"github.com/turbolytics/latte/internal/collector/config"
	"github.com/turbolytics/latte/internal/collector/metric"
	"github.com/turbolytics/latte/internal/invoker"
	"github.com/turbolytics/latte/internal/schedule"
	"github.com/turbolytics/latte/internal/sink"
	"github.com/turbolytics/latte/internal/sink/console"
	"github.com/turbolytics/latte/internal/sink/file"
	"github.com/turbolytics/latte/internal/sink/http"
	"github.com/turbolytics/latte/internal/sink/kafka"
	"github.com/turbolytics/latte/internal/source"
	"github.com/turbolytics/latte/internal/state"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

func NewStorer(c state.Config, l *zap.Logger) (invoker.Storer, error) {
	var s invoker.Storer
	var err error
	switch c.Type {
	case state.StoreTypeMemory:
		s, err = state.NewMemoryStoreFromGenericConfig(
			c.Config,
			state.MemoryStoreWithLogger(l),
		)
	}
	if err != nil {
		return nil, err
	}
	return s, err
}

type RootOption func(*config.Root)

func WithJustValidation(validate bool) RootOption {
	return func(c *config.Root) {
		c.Validate = validate
	}
}

func RootWithLogger(l *zap.Logger) RootOption {
	return func(c *config.Root) {
		c.Logger = l
	}
}

func NewCollectorsFromGlob(glob string, opts ...RootOption) ([]invoker.Collector, error) {
	files, err := filepath.Glob(glob)
	if err != nil {
		return nil, err
	}
	var collectors []invoker.Collector

	for _, fName := range files {
		c, err := NewCollectorFromFile(fName, opts...)
		if err != nil {
			return nil, err
		}
		collectors = append(collectors, c)
	}
	return collectors, nil
}

func NewCollectorFromFile(fpath string, opts ...RootOption) (invoker.Collector, error) {
	fmt.Printf("loading config from file: %q\n", fpath)

	bs, err := os.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	return NewCollector(bs, opts...)
}

func NewCollector(bs []byte, opts ...RootOption) (invoker.Collector, error) {
	var conf config.Root

	for _, opt := range opts {
		opt(&conf)
	}

	if err := yaml.Unmarshal(bs, &conf); err != nil {
		return nil, err
	}

	var coll invoker.Collector
	var err error

	switch conf.Collector.Type {
	case config.TypeMetric:
		coll, err = NewMetricCollectorFromConfig(
			bs,
			conf.Validate,
			conf.Logger,
		)
	case config.TypePartition:
		/*
			collConfig, err := partition.NewConfig(
				bs,
				partition.ConfigWithJustValidation(conf.validate),
				partition.ConfigWithLogger(conf.logger),
			)
		*/
	default:
		return nil, fmt.Errorf("collector type: %v not supported", conf.Collector.Type)
	}

	if err != nil {
		return nil, err
	}

	return coll, err
}

func NewSink(c sink.Config, l *zap.Logger, validate bool) (invoker.Sinker, error) {
	var s invoker.Sinker
	var err error
	switch c.Type {
	case sink.TypeConsole:
		s, err = console.NewFromGenericConfig(c.Config)
	case sink.TypeKafka:
		s, err = kafka.NewFromGenericConfig(c.Config)
	case sink.TypeHTTP:
		s, err = http.NewFromGenericConfig(
			c.Config,
			http.WithLogger(l),
		)
	case sink.TypeFile:
		s, err = file.NewFromGenericConfig(
			c.Config,
			validate,
		)
	}
	return s, err
}

func NewSinks(cs map[string]sink.Config, l *zap.Logger, validate bool) (map[string]invoker.Sinker, error) {
	sinks := make(map[string]invoker.Sinker)

	for k, conf := range cs {
		i, err := NewSink(conf, l, validate)
		if err != nil {
			return nil, err
		}
		sinks[k] = i
	}
	return sinks, nil
}

func NewSourcer(sc source.Config, l *zap.Logger, validate bool) (invoker.Sourcer, error) {
	var err error
	var s invoker.Sourcer
	switch sc.Type {
	case source.TypePartitionS3:
		/*
			s, err = s3.NewFromGenericConfig(
				sc.Config,
			)
		*/
	case source.TypePostgres:
		/*
			s, err = postgres.NewFromGenericConfig(
				c.Collector,
				validate,
			)
		*/
	case source.TypeMongoDB:
		/*
			s, err = mongodb.NewFromGenericConfig(
				context.TODO(),
				sc.Config,
				validate,
			)
		*/
	case source.TypePrometheus:
		/*
			s, err = prometheus.NewFromGenericConfig(
				c.Collector,
				prometheus.WithLogger(l),
		*/
	default:
		return nil, fmt.Errorf("source type: %q unknown", sc.Type)
	}

	return s, err
}

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

	)

	return coll, nil
}
