package postgres

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/mitchellh/mapstructure"
	"github.com/turbolytics/collector/internal/metrics"
)

type config struct {
	URI string
	SQL string
}

type Postgres struct {
	config config
	db     *sql.DB
}

func (p *Postgres) Source(ctx context.Context) ([]metrics.Metric, error) {
	fmt.Println("here")
	return nil, nil
}

func NewFromGenericConfig(m map[string]any) (*Postgres, error) {
	var conf config
	if err := mapstructure.Decode(m, &conf); err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", conf.URI)
	if err != nil {
		return nil, err
	}
	fmt.Println(conf.URI)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Postgres{
		config: conf,
		db:     db,
	}, nil
}
