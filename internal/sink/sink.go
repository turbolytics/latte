package sink

import (
	"github.com/turbolytics/latte/internal/sink/console"
	"github.com/turbolytics/latte/internal/sink/file"
	"github.com/turbolytics/latte/internal/sink/http"
	"github.com/turbolytics/latte/internal/sink/kafka"
	"go.uber.org/zap"
)

type Config struct {
	Type   Type
	Sinker Sinker
	Config map[string]any
}

func (c *Config) Init(validate bool, logger *zap.Logger) error {
	var sink Sinker
	var err error

	switch c.Type {
	case TypeConsole:
		sink, err = console.NewFromGenericConfig(c.Config)
	case TypeKafka:
		sink, err = kafka.NewFromGenericConfig(c.Config)
	case TypeHTTP:
		sink, err = http.NewFromGenericConfig(
			c.Config,
			http.WithLogger(logger),
		)
	case TypeFile:
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
