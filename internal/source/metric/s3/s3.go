package s3

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"github.com/marcboeker/go-duckdb"
	"github.com/mitchellh/mapstructure"
	"github.com/turbolytics/latte/internal/metric"
	"github.com/turbolytics/latte/internal/partition"
	"github.com/turbolytics/latte/internal/record"
	"github.com/turbolytics/latte/internal/source"
	scsql "github.com/turbolytics/latte/internal/sql"
	"go.uber.org/zap"
	"text/template"
	"time"
)

type config struct {
	Bucket    string
	Endpoint  *string
	Partition string
	Prefix    string
	Region    string
	SQL       string
}

type Option func(*S3)

func WithLogger(l *zap.Logger) Option {
	return func(s *S3) {
		s.logger = l
	}
}

type S3 struct {
	config      config
	logger      *zap.Logger
	partitioner *partition.Partitioner
}

func (s *S3) Type() source.Type {
	return source.TypeMetricS3
}

func (s *S3) WindowDuration() *time.Duration {
	return nil
}

func (s *S3) createSecreteDDL() (string, error) {
	t := `
CREATE SECRET (
	TYPE S3,
	{{ if .Endpoint -}}
	ENDPOINT '{{.Endpoint}}',
	{{- end }}
	REGION '{{.Region}}',
	PROVIDER CREDENTIAL_CHAIN
);
`
	tmpl, err := template.New("source.metric.s3.secretDDL").Parse(t)
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	err = tmpl.Execute(&b, s.config)
	return b.String(), err
}

func (s *S3) Source(ctx context.Context) (record.Result, error) {
	// if invocation.start is set than this is a point in time lookup
	// we need a more explicit way to configure this for the source
	// instead of passing through the context...
	// This is leaky at this point since only tick will have invocation.start
	// and window-based collectors will contain a window...
	t := ctx.Value("invocation.start").(time.Time)
	p, err := s.partitioner.Render(t)
	if err != nil {
		return nil, err
	}

	// build the select path

	// initialize duckdb
	connector, err := duckdb.NewConnector("", func(execer driver.ExecerContext) error {
		bootQueries := []string{
			"INSTALL 'httpfs'",
			"LOAD 'httpfs'",
			"INSTALL 'json'",
			"LOAD 'json'",
			"INSTALL 'aws'",
			"LOAD 'aws'",
		}

		for _, qry := range bootQueries {
			_, err := execer.ExecContext(context.TODO(), qry, nil)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	db := sql.OpenDB(connector)
	defer db.Close()

	// add duckdb S3 configuration based on configs
	secretDDL, err := s.createSecreteDDL()
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(secretDDL)
	if err != nil {
		return nil, err
	}
	s.logger.Debug("source.metric.s3.query",
		zap.String("secret", secretDDL),
	)

	// query
	tmpl, err := template.New("source.metric.s3.query").Parse(s.config.SQL)
	if err != nil {
		return nil, err
	}
	var b bytes.Buffer
	err = tmpl.Execute(&b, struct {
		Partition string
	}{
		Partition: p,
	})
	if err != nil {
		return nil, err
	}

	s.logger.Debug("source.metric.s3.query",
		zap.String("query", b.String()),
	)

	rows, err := db.Query(b.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results, err := scsql.RowsToMaps(rows)
	if err != nil {
		return nil, err
	}

	ms, err := metric.MapsToMetrics(results)
	metricsResult := metric.NewMetricsResult(ms)

	return metricsResult, err
}

func NewFromGenericConfig(m map[string]any, opts ...Option) (*S3, error) {
	var conf config

	if err := mapstructure.Decode(m, &conf); err != nil {
		return nil, err
	}

	p, err := partition.New(conf.Partition)
	if err != nil {
		return nil, err
	}

	s := &S3{
		config:      conf,
		partitioner: p,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s, nil
}
