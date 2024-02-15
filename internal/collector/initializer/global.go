package initializer

import (
	"context"
	"fmt"
	"github.com/turbolytics/latte/internal/invoker"
	"github.com/turbolytics/latte/internal/sink"
	"github.com/turbolytics/latte/internal/sink/console"
	"github.com/turbolytics/latte/internal/sink/file"
	"github.com/turbolytics/latte/internal/sink/http"
	"github.com/turbolytics/latte/internal/sink/kafka"
	"github.com/turbolytics/latte/internal/source"
	"github.com/turbolytics/latte/internal/source/metric/mongodb"
	"github.com/turbolytics/latte/internal/state"
	"go.uber.org/zap"
)

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
		s, err = mongodb.NewFromGenericConfig(
			context.TODO(),
			sc.Config,
			validate,
		)
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
