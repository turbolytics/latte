package initializer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMetricCollectorFromConfig_Success(t *testing.T) {
	bs := []byte(`
name: postgres_users_total_30s

collector:
  type: metric

metric:
  name: core.users.total
  type: COUNT
  tags:
    - key: env
      value: prod

schedule:
  interval: 30s

source:
  type: metric.postgres
  config:
    uri: 'postgresql://test:test@{{ getEnvOrDefault "SC_POSTGRES_HOST" "127.0.0.1" }}:5432/test?sslmode=disable'
    sql: |
      SELECT 
        account as customer,
        COUNT(*) as value
      FROM 
        users
      GROUP BY
        account
      
sinks:
  audit:
    type: file
    config:
      path: /tmp/log/latte.audit.log
`)

	c, err := NewCollector(
		bs,
		WithJustValidation(true),
	)
	assert.NoError(t, err)
	assert.Equal(t, "postgres_users_total_30s", c.Name())
	fmt.Printf("%+v", c)
}
