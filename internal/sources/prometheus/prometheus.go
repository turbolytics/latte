package prometheus

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/marcboeker/go-duckdb"
	_ "github.com/marcboeker/go-duckdb"
	"github.com/mitchellh/mapstructure"
	"github.com/turbolytics/collector/internal/metrics"
	scsql "github.com/turbolytics/collector/internal/sources/sql"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
)

type Option func(*Prometheus)

func WithLogger(l *zap.Logger) Option {
	return func(p *Prometheus) {
		p.logger = l
	}
}

type config struct {
	SQL   string
	Query string
	URI   string

	url *url.URL
}

type Prometheus struct {
	logger *zap.Logger
	config config
}

type apiMetric struct {
	Metric interface{}
	Value  []any
}

type apiData struct {
	ResultType string `json:"resultType"`
	Result     []apiMetric
}

type apiResponse struct {
	Status string
	Data   apiData
}

func (p *Prometheus) promMetrics(ctx context.Context, uri string) (*apiResponse, error) {

	p.logger.Info(
		"prometheus.promMetrics",
		zap.String("url", p.config.url.String()),
	)

	// make request to prometheus
	resp, err := http.Get(p.config.url.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	apiResp := apiResponse{}
	if err := json.Unmarshal(bs, &apiResp); err != nil {
		return nil, err
	}

	return &apiResp, nil
}

func (p *Prometheus) Source(ctx context.Context) ([]*metrics.Metric, error) {
	promMetrics, err := p.promMetrics(ctx, p.config.url.String())
	if err != nil {
		return nil, err
	}

	// initialize db and write all results to db, to expose sql interface
	// over the data....
	connector, err := duckdb.NewConnector("", func(execer driver.ExecerContext) error {
		bootQueries := []string{
			"INSTALL 'json'",
			"LOAD 'json'",
		}

		for _, qry := range bootQueries {
			_, err = execer.ExecContext(context.TODO(), qry, nil)
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

	_, err = db.Exec(`
CREATE TABLE prom_metrics (
    metric JSON, 
    value VARCHAR 
)
`)
	if err != nil {
		return nil, err
	}

	conn, err := connector.Connect(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	appender, err := duckdb.NewAppenderFromConn(conn, "", "prom_metrics")
	if err != nil {
		return nil, err
	}
	defer appender.Close()

	for _, result := range promMetrics.Data.Result {
		bs, err := json.Marshal(result.Metric)
		if err != nil {
			return nil, err
		}
		// TODO - reevaluate duckdb here,
		// Check parameterize queries insert. Tried appender API and types
		// were not working...
		_, err = db.Exec(
			fmt.Sprintf(
				`INSERT INTO prom_metrics VALUES ('%s', '%s')`,
				string(bs),
				result.Value[1].(string),
			),
		)
		if err != nil {
			return nil, err
		}
	}

	rows, err := db.QueryContext(ctx, p.config.SQL)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	results, err := scsql.RowsToMaps(rows)
	if err != nil {
		return nil, err
	}

	ms, err := metrics.MapsToMetrics(results)
	return ms, err
}

func NewFromGenericConfig(m map[string]any, opts ...Option) (*Prometheus, error) {
	var conf config

	if err := mapstructure.Decode(m, &conf); err != nil {
		return nil, err
	}

	u, err := url.Parse(conf.URI)
	if err != nil {
		return nil, err
	}
	conf.url = u

	// initialize the query a single time
	q := conf.url.Query()
	q.Add("query", conf.Query)
	conf.url.RawQuery = q.Encode()

	p := &Prometheus{
		config: conf,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p, nil
}
