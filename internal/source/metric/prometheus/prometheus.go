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
	"github.com/turbolytics/latte/internal/metric"
	scsql "github.com/turbolytics/latte/internal/sql"
	"github.com/turbolytics/latte/internal/timeseries"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type config struct {
	SQL   string
	Query string
	URI   string
	Time  *timeseries.WindowConfig

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
		zap.String("url", uri),
	)

	// make request to prometheus
	resp, err := http.Get(uri)
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

func (p *Prometheus) applySQL(ctx context.Context, resp *apiResponse) ([]map[string]any, error) {

	// initialize db and write all results to db, to expose sql interface
	// over the data....
	connector, err := duckdb.NewConnector("", func(execer driver.ExecerContext) error {
		bootQueries := []string{
			"INSTALL 'json'",
			"LOAD 'json'",
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

	for _, result := range resp.Data.Result {
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

	return scsql.RowsToMaps(rows)
}

func (p *Prometheus) Window() *time.Duration {
	return p.config.Time.Duration()
}

func (p *Prometheus) Source(ctx context.Context) ([]*metric.Metric, error) {
	// This is already starting to devolve :sweat:
	windowStart := ctx.Value("window.start").(time.Time)
	windowEnd := ctx.Value("window.end").(time.Time)

	u, _ := url.Parse(p.config.url.String())
	q := u.Query()
	q.Add("query", p.config.Query)
	q.Add("time", strconv.FormatInt(windowEnd.Unix(), 10))

	u.RawQuery = q.Encode()

	promMetrics, err := p.promMetrics(ctx, u.String())
	if err != nil {
		return nil, err
	}
	results, err := p.applySQL(ctx, promMetrics)
	if err != nil {
		return nil, err
	}

	ms, err := metric.MapsToMetrics(results)
	if err == nil {
		for _, m := range ms {
			m.Window = &windowStart
		}
	}
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

	if conf.Time != nil {
		if err := conf.Time.Init(); err != nil {
			return nil, err
		}
	}

	p := &Prometheus{
		config: conf,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p, nil
}
