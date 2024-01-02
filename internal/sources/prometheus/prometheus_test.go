package prometheus

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
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
          "job": "signals-collector",
          "otel_scope_name": "signals-collector",
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
          "job": "signals-collector",
          "otel_scope_name": "signals-collector",
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
          "job": "signals-collector",
          "otel_scope_name": "signals-collector",
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
          "job": "signals-collector",
          "otel_scope_name": "signals-collector",
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
          "job": "signals-collector",
          "otel_scope_name": "signals-collector",
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
          "job": "signals-collector",
          "otel_scope_name": "signals-collector",
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
		config: config{
			URL: u,
			SQL: `
SELECT
	COUNT(*) as value
FROM 
	prom_metrics
`,
		},
	}

	_, err := p.Source(context.Background())
	assert.NoError(t, err)
}
