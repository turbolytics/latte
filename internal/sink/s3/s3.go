package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/mitchellh/mapstructure"
	"github.com/turbolytics/latte/internal/encoding"
	"github.com/turbolytics/latte/internal/record"
	"github.com/turbolytics/latte/internal/sink"
	"go.uber.org/zap"
)

type config struct {
	BatchSize        int `mapstructure:"batch_size"`
	Bucket           string
	Encoding         encoding.Config
	Prefix           string
	Endpoint         *string
	Region           string
	S3ForcePathStyle bool `mapstructure:"force_path_style"`
}

type Option func(*S3)

func WithLogger(l *zap.Logger) Option {
	return func(s *S3) {
		s.logger = l
	}
}

type S3 struct {
	client  *s3.S3
	config  config
	encoder encoding.Encoder

	logger *zap.Logger
}

func (s *S3) Close() error {
	return nil
}

func (s *S3) Flush() error {
	return nil
}

func (s *S3) Type() sink.Type {
	return sink.TypeS3
}

func (s *S3) Write(r record.Record) (int, error) {
	return 0, nil
}

func NewFromGenericConfig(m map[string]any, opts ...Option) (*S3, error) {
	var conf config

	if err := mapstructure.Decode(m, &conf); err != nil {
		return nil, err
	}

	e, err := encoding.NewEncoder(conf.Encoding)

	if err != nil {
		return nil, err
	}

	awsConfig := &aws.Config{
		Region: aws.String(conf.Region),
		// Credentials:      credentials.NewStaticCredentials("test", "test", ""),
		S3ForcePathStyle: aws.Bool(conf.S3ForcePathStyle),
	}

	if conf.Endpoint != nil {
		awsConfig.Endpoint = aws.String(*conf.Endpoint)
	}

	sess, _ := session.NewSession(awsConfig)
	client := s3.New(sess)

	s := &S3{
		client:  client,
		config:  conf,
		encoder: e,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s, nil
}
