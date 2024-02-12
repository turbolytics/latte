package s3

import (
	"bytes"
	"context"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/turbolytics/latte/internal/partition"
	"github.com/turbolytics/latte/internal/source"
	"github.com/turbolytics/latte/internal/timeseries"
	"go.uber.org/zap"
	"net/url"
	"path"
	"text/template"
	"time"
)

type Option func(*Source)

func WithLogger(l *zap.Logger) Option {
	return func(s *Source) {
		s.logger = l
	}
}

type config struct {
	// ValidateExists bool `mapstructure:"validate_exists"`
	Bucket    string
	Region    string
	Prefix    string
	Partition string

	Time *timeseries.WindowConfig
}

func (c config) Host() string {
	return fmt.Sprintf(
		"%s.s3.%s.amazonaws.com",
		c.Bucket,
		c.Region,
	)
}

type tmplData struct {
	Year   int
	Month  int
	Day    int
	Hour   int
	Minute int
}

type Source struct {
	config config
	logger *zap.Logger
}

func (s *Source) Window() *time.Duration {
	return s.config.Time.Duration()
}

func (s *Source) Type() source.Type {
	return source.TypePartitionS3
}

func (s *Source) Source(ctx context.Context) (*partition.Partition, error) {
	windowStart := ctx.Value("window.start").(time.Time)

	// parse the template
	t, err := template.New("config").Parse(string(s.config.Partition))
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer

	td := tmplData{
		Year:   windowStart.Year(),
		Month:  int(windowStart.Month()),
		Day:    windowStart.Day(),
		Hour:   windowStart.Hour(),
		Minute: windowStart.Minute(),
	}

	if err := t.Execute(&out, td); err != nil {
		return nil, err
	}

	u := url.URL{
		Scheme: "https",
		Host:   s.config.Host(),
		Path:   path.Join(s.config.Prefix, string(out.Bytes())),
	}

	p := partition.Partition{
		URI: &u,
	}

	return &p, nil
}

func NewFromGenericConfig(m map[string]any, opts ...Option) (*Source, error) {
	var conf config

	if err := mapstructure.Decode(m, &conf); err != nil {
		return nil, err
	}

	if conf.Time != nil {
		if err := conf.Time.Init(); err != nil {
			return nil, err
		}
	}

	s := &Source{
		config: conf,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s, nil
}
