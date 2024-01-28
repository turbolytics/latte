package sink

import (
	"github.com/turbolytics/collector/internal/sinks"
	"github.com/turbolytics/collector/internal/sinks/console"
	"github.com/turbolytics/collector/internal/sinks/file"
	"github.com/turbolytics/collector/internal/sinks/http"
	"github.com/turbolytics/collector/internal/sinks/kafka"
	"go.uber.org/zap"
)

type Config struct {
	Type   sinks.Type
	Sinker sinks.Sinker
	Config map[string]any
}

func (c *Config) Init(validate bool, logger *zap.Logger) error {
	var sink sinks.Sinker
	var err error

	switch c.Type {
	case sinks.TypeConsole:
		sink, err = console.NewFromGenericConfig(c.Config)
	case sinks.TypeKafka:
		sink, err = kafka.NewFromGenericConfig(c.Config)
	case sinks.TypeHTTP:
		sink, err = http.NewFromGenericConfig(
			c.Config,
			http.WithLogger(logger),
		)
	case sinks.TypeFile:
		sink, err = file.NewFromGenericConfig(
			c.Config,
			validate,
		)
	}

	if err != nil {
		return err
	}

	c.Sinker = sink
	return nil

}
