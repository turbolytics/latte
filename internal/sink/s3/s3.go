package s3

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/turbolytics/latte/internal/encoding"
	"github.com/turbolytics/latte/internal/partition"
	"github.com/turbolytics/latte/internal/record"
	"github.com/turbolytics/latte/internal/sink"
	"go.uber.org/zap"
	"path"
	"time"
)

type config struct {
	BatchSize        int `mapstructure:"batch_size"`
	Bucket           string
	Encoding         encoding.Config
	Prefix           string
	Endpoint         *string
	Region           string
	S3ForcePathStyle bool `mapstructure:"force_path_style"`
	Partition        string
}

type Option func(*S3)

func WithLogger(l *zap.Logger) Option {
	return func(s *S3) {
		s.logger = l
	}
}

type S3 struct {
	buf     *bytes.Buffer
	config  config
	encoder encoding.Encoder

	logger      *zap.Logger
	partitioner *partition.Partitioner
	uploader    *s3manager.Uploader
}

func (s *S3) Close() error {
	return nil
}

func (s *S3) Flush(ctx context.Context) error {
	start := ctx.Value("window.start").(time.Time)
	fname := fmt.Sprintf("%s.json", uuid.New().String())
	p, err := s.partitioner.Render(start)
	if err != nil {
		return err
	}

	k := path.Join(
		s.config.Prefix,
		p,
		fname,
	)

	s.logger.Debug("sinks.S3.uploading",
		zap.String("bucket", s.config.Bucket),
		zap.String("key", k),
	)

	_, err = s.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.config.Bucket),

		// Can also use the `filepath` standard library package to modify the
		// filename as need for an S3 object key. Such as turning absolute path
		// to a relative path.
		Key: aws.String(k),

		// The file to be uploaded. io.ReadSeeker is preferred as the Uploader
		// will be able to optimize memory when uploading large content. io.Reader
		// is supported, but will require buffering of the reader's bytes for
		// each part.
		Body: bufio.NewReader(s.buf),
	})

	return err
}

func (s *S3) Type() sink.Type {
	return sink.TypeS3
}

func (s *S3) Write(ctx context.Context, r record.Record) (int, error) {
	if s.buf == nil {
		s.buf = &bytes.Buffer{}
	}
	if err := s.encoder.Init(s.buf); err != nil {
		return 0, err
	}

	if err := s.encoder.Write(r.Map()); err != nil {
		return 0, err
	}

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

	p, err := partition.New(conf.Partition)
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
	uploader := s3manager.NewUploader(sess)

	s := &S3{
		config:      conf,
		encoder:     e,
		partitioner: p,
		uploader:    uploader,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s, nil
}
