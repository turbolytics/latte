package postgres

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/mitchellh/mapstructure"
	"github.com/turbolytics/collector/internal/metrics"
	scsql "github.com/turbolytics/collector/internal/sources/sql"
	"strconv"
)

type config struct {
	URI string
	SQL string
}

type Postgres struct {
	config config
	db     *sql.DB
}

func resultsToMetrics(results []map[string]any) ([]*metrics.Metric, error) {
	var ms []*metrics.Metric
	for _, r := range results {
		val, ok := r["value"]
		if !ok {
			return nil, fmt.Errorf("each row must contain a %q key", "value")
		}

		m := metrics.New()

		switch v := val.(type) {
		case int:
			m.Value = float64(v)
		case string:
			tv, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, fmt.Errorf("unable to parse string to float: %q", v)
			}
			m.Value = tv
		}
		delete(r, "value")
		for k, v := range r {
			m.Tags[k] = v.(string)
		}
		ms = append(ms, &m)
	}

	return ms, nil
}

func (p *Postgres) Source(ctx context.Context) ([]*metrics.Metric, error) {
	rows, err := p.db.QueryContext(ctx, p.config.SQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results, err := scsql.RowsToMaps(rows)
	if err != nil {
		return nil, err
	}

	ms, err := resultsToMetrics(results)
	return ms, err
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
