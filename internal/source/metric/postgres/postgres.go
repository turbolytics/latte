package postgres

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/mitchellh/mapstructure"
	"github.com/turbolytics/latte/internal/metric"
	"github.com/turbolytics/latte/internal/record"
	"github.com/turbolytics/latte/internal/source"
	scsql "github.com/turbolytics/latte/internal/sql"
	"time"
)

type config struct {
	URI string
	SQL string
}

type Postgres struct {
	config config
	db     *sql.DB
}

func (p *Postgres) Window() *time.Duration {
	return nil
}

func (p *Postgres) Source(ctx context.Context) (record.Result, error) {
	rows, err := p.db.QueryContext(ctx, p.config.SQL)
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

func (p *Postgres) Type() source.Type {
	return source.TypeMetricPostgres
}

func NewFromGenericConfig(m map[string]any, validate bool) (*Postgres, error) {
	var conf config
	if err := mapstructure.Decode(m, &conf); err != nil {
		return nil, err
	}

	var db *sql.DB
	var err error
	if !validate {
		db, err = sql.Open("postgres", conf.URI)
		if err != nil {
			return nil, err
		}

		if err := db.Ping(); err != nil {
			return nil, err
		}
	}

	return &Postgres{
		config: conf,
		db:     db,
	}, nil
}
