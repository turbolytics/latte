package s3

import (
	"context"
	"github.com/mitchellh/mapstructure"
	"github.com/turbolytics/latte/internal/record"
	"github.com/turbolytics/latte/internal/source"
	"go.uber.org/zap"
	"time"
)

type config struct {
	Bucket    string
	Endpoint  *string
	Partition string
	Prefix    string
	Region    string
}

type Option func(*S3)

func WithLogger(l *zap.Logger) Option {
	return func(s *S3) {
		s.logger = l
	}
}

type S3 struct {
	config config
	logger *zap.Logger
}

func (s *S3) Type() source.Type {
	return source.TypeMetricS3
}

func (s *S3) WindowDuration() *time.Duration {
	return nil
}

func (s *S3) Source(ctx context.Context) (record.Result, error) {
	return nil, nil
}

func NewFromGenericConfig(m map[string]any, opts ...Option) (*S3, error) {
	var conf config

	if err := mapstructure.Decode(m, &conf); err != nil {
		return nil, err
	}

	s := &S3{
		config: conf,
	}

	for _, opt := range opts {
		opt(s)
	}

	return nil, nil
}
