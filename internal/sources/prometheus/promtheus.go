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
	"io"
	"net/http"
	"net/url"
)

type config struct {
	URL *url.URL
	SQL string
}

type Prometheus struct {
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
	/*
	 param := make(url.Values)
	*/
	// make request to prometheus
	resp, err := http.Get(p.config.URL.String())
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
	promMetrics, err := p.promMetrics(ctx, p.config.URL.String())
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

func NewFromGenericConfig(m map[string]any, validate bool) (*Prometheus, error) {
	var conf config

	if err := mapstructure.Decode(m, &conf); err != nil {
		return nil, err
	}

	var err error
	return nil, err
}
