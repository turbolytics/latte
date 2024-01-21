package config

import (
	"bytes"
	"context"
	"fmt"
	"github.com/turbolytics/collector/internal/collector/metric/sources"
	"github.com/turbolytics/collector/internal/collector/metric/sources/mongodb"
	"github.com/turbolytics/collector/internal/collector/metric/sources/postgres"
	prometheus2 "github.com/turbolytics/collector/internal/collector/metric/sources/prometheus"
	"github.com/turbolytics/collector/internal/metrics"
	"github.com/turbolytics/collector/internal/sinks"
	"github.com/turbolytics/collector/internal/sinks/console"
	"github.com/turbolytics/collector/internal/sinks/file"
	"github.com/turbolytics/collector/internal/sinks/http"
	"github.com/turbolytics/collector/internal/sinks/kafka"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
	"text/template"
)

type validater interface {
	Validate() error
}

type Tag struct {
	Key   string
	Value string
}

type Metric struct {
	Name string
	Type metrics.Type
	Tags []Tag
}

type Collector struct {
	Type   string // enum
	Config map[string]any
}

type Sink struct {
	Type   sinks.Type
	Sinker sinks.Sinker
	Config map[string]any
}

type Config struct {
	Name       string
	Metric     Metric
	Schedule   Schedule
	Source     Source
	Sinks      map[string]Sink
	StateStore StateStore `yaml:"state_store"`

	logger *zap.Logger
	// validate will skip initializing network dependencies
	validate bool
}

type Option func(*Config)

func WithJustValidation(validate bool) Option {
	return func(c *Config) {
		c.validate = validate
	}
}

func WithLogger(l *zap.Logger) Option {
	return func(c *Config) {
		c.logger = l
	}
}

// initSource initializes the correct source.
func initSource(c *Config) error {
	var s sources.MetricSourcer
	var err error
	switch c.Source.Type {
	case sources.TypePostgres:
		s, err = postgres.NewFromGenericConfig(
			c.Source.Config,
			c.validate,
		)
	case sources.TypeMongoDB:
		s, err = mongodb.NewFromGenericConfig(
			context.TODO(),
			c.Source.Config,
			c.validate,
		)
	case sources.TypePrometheus:
		s, err = prometheus2.NewFromGenericConfig(
			c.Source.Config,
			prometheus2.WithLogger(c.logger),
		)
	default:
		return fmt.Errorf("source type: %q unknown", c.Source.Type)
	}

	if err != nil {
		return err
	}

	c.Source.Sourcer = s
	return nil
}

// initSinks initializes all the outputs
func initSinks(c *Config) error {
	for k, v := range c.Sinks {
		switch v.Type {
		case sinks.TypeConsole:
			sink, err := console.NewFromGenericConfig(v.Config)
			if err != nil {
				return err
			}
			v.Sinker = sink
			c.Sinks[k] = v
		case sinks.TypeKafka:
			sink, err := kafka.NewFromGenericConfig(v.Config)
			if err != nil {
				return err
			}
			v.Sinker = sink
			c.Sinks[k] = v
		case sinks.TypeHTTP:
			sink, err := http.NewFromGenericConfig(v.Config)
			if err != nil {
				return err
			}
			v.Sinker = sink
			c.Sinks[k] = v
		case sinks.TypeFile:
			sink, err := file.NewFromGenericConfig(
				v.Config,
				c.validate,
			)
			if err != nil {
				return err
			}
			v.Sinker = sink
			c.Sinks[k] = v
		}
	}
	return nil
}

func parseTemplate(bs []byte) ([]byte, error) {
	funcMap := template.FuncMap{
		"getEnv": func(key string) string {
			return os.Getenv(key)
		},
		"getEnvOrDefault": func(key string, d string) string {
			envVal := os.Getenv(key)
			if envVal == "" {
				return d
			}

			return envVal
		},
	}
	t, err := template.New("config").Funcs(funcMap).Parse(string(bs))
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	if err := t.Execute(&out, nil); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func defaults(c *Config) error {
	(&c.Source).SetDefaults()

	return nil
}

func validate(c Config) error {
	validators := []validater{
		c.Schedule,
		c.Source,
		c.StateStore,
	}

	for _, v := range validators {
		if err := v.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// New initializes a config from yaml bytes.
// New initializes all subtypes as well.
func New(raw []byte, opts ...Option) (*Config, error) {
	var conf Config

	for _, opt := range opts {
		opt(&conf)
	}

	bs, err := parseTemplate(raw)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(bs, &conf); err != nil {
		return nil, err
	}

	if err := defaults(&conf); err != nil {
		return nil, err
	}

	if err := validate(conf); err != nil {
		return nil, err
	}

	if err := initStateStore(&conf); err != nil {
		return nil, err
	}

	if err := initSource(&conf); err != nil {
		return nil, err
	}

	if err := initSinks(&conf); err != nil {
		return nil, err
	}
	return &conf, nil
}
