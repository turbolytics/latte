package sink

import (
	"github.com/turbolytics/latte/internal/sink/console"
	"github.com/turbolytics/latte/internal/sink/file"
	"github.com/turbolytics/latte/internal/sink/http"
	"github.com/turbolytics/latte/internal/sink/kafka"
	"github.com/turbolytics/latte/internal/sink/type"
	"go.uber.org/zap"
)

type Config struct {
	Type   _type.Type
	Sinker _type.Sinker
	Config map[string]any
}

func (c *Config) Init(validate bool, logger *zap.Logger) error {
	var sink _type.Sinker
	var err error

	switch c.Type {
	case _type.TypeConsole:
		sink, err = console.NewFromGenericConfig(c.Config)
	case _type.TypeKafka:
		sink, err = kafka.NewFromGenericConfig(c.Config)
	case _type.TypeHTTP:
		sink, err = http.NewFromGenericConfig(
			c.Config,
			http.WithLogger(logger),
		)
	case _type.TypeFile:
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
