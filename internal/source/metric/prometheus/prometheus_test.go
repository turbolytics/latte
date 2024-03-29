package prometheus

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/turbolytics/latte/internal/metric"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestPrometheus_Source_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{
  "status": "success",
  "data": {
    "resultType": "vector",
    "result": [
      {
        "metric": {
          "collector_name": "mongo.users.10s",
          "instance": "host.docker.internal:12223",
          "job": "latte",
          "otel_scope_name": "latte",
          "result_status_code": "OK"
        },
        "value": [
          1704156155.891,
          "8555.054635645161"
        ]
      },
      {
        "metric": {
          "collector_name": "postgres.users.total.1m",
          "instance": "host.docker.internal:12223",
          "job": "latte",
          "otel_scope_name": "latte",
          "result_status_code": "OK"
        },
        "value": [
          1704156155.891,
          "1472.9032501210243"
        ]
      },
      {
        "metric": {
          "collector_name": "postgres.users.total.24h",
          "instance": "host.docker.internal:12223",
          "job": "latte",
          "otel_scope_name": "latte",
          "result_status_code": "ERROR"
        },
        "value": [
          1704156155.891,
          "0"
        ]
      },
      {
        "metric": {
          "collector_name": "postgres.users.total.24h",
          "instance": "host.docker.internal:12223",
          "job": "latte",
          "otel_scope_name": "latte",
          "result_status_code": "OK"
        },
        "value": [
          1704156155.891,
          "1.0012938478049112"
        ]
      },
      {
        "metric": {
          "collector_name": "postgres.users.total.30s",
          "instance": "host.docker.internal:12223",
          "job": "latte",
          "otel_scope_name": "latte",
          "result_status_code": "ERROR"
        },
        "value": [
          1704156155.891,
          "0"
        ]
      },
      {
        "metric": {
          "collector_name": "postgres.users.total.30s",
          "instance": "host.docker.internal:12223",
          "job": "latte",
          "otel_scope_name": "latte",
          "result_status_code": "OK"
        },
        "value": [
          1704156155.891,
          "2940.800031003024"
        ]
      }
    ]
  }
}`)
	}))
	defer ts.Close()

	u, _ := url.Parse(ts.URL)

	p := &Prometheus{
		logger: zap.NewNop(),
		config: config{
			url: u,
			SQL: `
SELECT
	metric->>'collector_name' as collector_name,
	CASE 
		WHEN metric->>'result_status_code' = 'OK'
		THEN false 
		ELSE true
	END as error,
	round(value::DOUBLE, 0) as value
FROM 
	prom_metrics
`,
		},
	}

	start := time.Date(2024, 1, 1, 1, 1, 1, 0, time.UTC)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "window.start", start)
	ctx = context.WithValue(ctx, "window.end", time.Now().UTC())

	rs, err := p.Source(ctx)
	assert.NoError(t, err)

	ms := rs.(*metric.Metrics)

	for _, m := range ms.Metrics {
		m.UUID = ""
		m.Timestamp = time.Time{}
	}

	assert.Equal(t, []*metric.Metric{
		{
			Window: &start,
			Value:  8555,
			Tags: map[string]string{
				"collector_name": "mongo.users.10s",
				"error":          "false",
			},
		},
		{
			Window: &start,
			Value:  1473,
			Tags: map[string]string{
				"collector_name": "postgres.users.total.1m",
				"error":          "false",
			},
		},
		{
			Window: &start,
			Value:  0,
			Tags: map[string]string{
				"collector_name": "postgres.users.total.24h",
				"error":          "true",
			},
		},
		{
			Window: &start,
			Value:  1,
			Tags: map[string]string{
				"collector_name": "postgres.users.total.24h",
				"error":          "false",
			},
		},
		{
			Window: &start,
			Value:  0,
			Tags: map[string]string{
				"collector_name": "postgres.users.total.30s",
				"error":          "true",
			},
		},
		{
			Window: &start,
			Value:  2941,
			Tags: map[string]string{
				"collector_name": "postgres.users.total.30s",
				"error":          "false",
			},
		},
	}, ms.Metrics)
}
