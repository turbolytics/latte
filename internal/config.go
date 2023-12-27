package internal

import (
	"bytes"
	"context"
	"github.com/turbolytics/collector/internal/metrics"
	"github.com/turbolytics/collector/internal/sinks"
	"github.com/turbolytics/collector/internal/sinks/console"
	"github.com/turbolytics/collector/internal/sinks/file"
	"github.com/turbolytics/collector/internal/sinks/http"
	"github.com/turbolytics/collector/internal/sinks/kafka"
	"github.com/turbolytics/collector/internal/sources"
	"github.com/turbolytics/collector/internal/sources/mongodb"
	"github.com/turbolytics/collector/internal/sources/postgres"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"text/template"
	"time"
)

type Tag struct {
	Key   string
	Value string
}

type Metric struct {
	Name string
	Type metrics.Type
	Tags []Tag
}

type Schedule struct {
	Interval time.Duration
	Cron     string
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

type Source struct {
	Type    sources.Type
	Sourcer sources.Sourcer
	Config  map[string]any
}

type ConfigOption func(*Config)

func WithJustValidation(validate bool) ConfigOption {
	return func(c *Config) {
		c.validate = validate
	}
}

type Config struct {
	Name     string
	Metric   Metric
	Schedule Schedule
	Source   Source
	Sinks    map[string]Sink

	// validate will skip initializing network dependencies
	validate bool
}

// initSource initializes the correct source.
func initSource(c *Config) error {
	var s sources.Sourcer
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

func NewConfigsFromDir(dirname string) ([]*Config, error) {
	files, err := os.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	var configs []*Config

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		n := path.Join(dirname, f.Name())
		c, err := NewConfigFromFile(n)
		if err != nil {
			return nil, err
		}
		configs = append(configs, c)
	}
	return configs, nil
}

func NewConfigFromFile(name string, opts ...ConfigOption) (*Config, error) {
	bs, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return NewConfig(bs, opts...)
}

// NewConfig initializes a config from yaml bytes.
// NewConfig initializes all subtypes as well.
func NewConfig(raw []byte, opts ...ConfigOption) (*Config, error) {
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

	if err := initSource(&conf); err != nil {
		return nil, err
	}

	if err := initSinks(&conf); err != nil {
		return nil, err
	}
	return &conf, nil
}
