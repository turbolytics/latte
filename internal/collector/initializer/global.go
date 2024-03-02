package initializer

import (
	"context"
	"fmt"
	"github.com/turbolytics/latte/internal/collector/template"
	"github.com/turbolytics/latte/internal/invoker"
	"github.com/turbolytics/latte/internal/sink"
	"github.com/turbolytics/latte/internal/sink/console"
	"github.com/turbolytics/latte/internal/sink/file"
	"github.com/turbolytics/latte/internal/sink/http"
	"github.com/turbolytics/latte/internal/sink/kafka"
	"github.com/turbolytics/latte/internal/sink/s3"
	"github.com/turbolytics/latte/internal/source"
	"github.com/turbolytics/latte/internal/source/metric/mongodb"
	"github.com/turbolytics/latte/internal/source/metric/postgres"
	"github.com/turbolytics/latte/internal/source/metric/prometheus"
	"github.com/turbolytics/latte/internal/state"
	"go.uber.org/zap"
)

func NewSink(c sink.Config, l *zap.Logger, validate bool) (invoker.Sinker, error) {
	var s invoker.Sinker
	var err error
	switch c.Type {
	case sink.TypeConsole:
		s, err = console.NewFromGenericConfig(c.Config)
	case sink.TypeFile:
		s, err = file.NewFromGenericConfig(
			c.Config,
			validate,
		)
	case sink.TypeHTTP:
		s, err = http.NewFromGenericConfig(
			c.Config,
			http.WithLogger(l),
		)
	case sink.TypeKafka:
		s, err = kafka.NewFromGenericConfig(c.Config)
	case sink.TypeS3:
		s, err = s3.NewFromGenericConfig(
			c.Config,
			s3.WithLogger(l),
		)
	default:
		return nil, fmt.Errorf("sink type: %q not supported", c.Type)
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

	// enabling templating across a couple of fixed, known configuration fields
	bs, err := template.Parse([]byte(sc.Config["uri"].(string)))
	if err != nil {
		return nil, err
	}
	sc.Config["uri"] = string(bs)

	switch sc.Type {
	case source.TypeMetricPostgres:
		s, err = postgres.NewFromGenericConfig(
			sc.Config,
			validate,
		)
	case source.TypeMetricMongoDB:
		s, err = mongodb.NewFromGenericConfig(
			context.TODO(),
			sc.Config,
			validate,
		)
	case source.TypePrometheus:
		s, err = prometheus.NewFromGenericConfig(
			sc.Config,
			prometheus.WithLogger(l),
		)
	default:
		return nil, fmt.Errorf("source type: %q unknown", sc.Type)
	}

	return s, err
}

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
