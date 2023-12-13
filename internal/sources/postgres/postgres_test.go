package postgres

import (
	"github.com/stretchr/testify/assert"
	"github.com/turbolytics/collector/internal/metrics"
	"testing"
)

func Test_resultsToMetrics_missingValue(t *testing.T) {
	rs := []map[string]any{
		{
			"key1": "value1",
		},
	}

	_, err := resultsToMetrics(rs)
	assert.EqualError(t, err, "each row must contain a \"value\" key")
}

func Test_resultsToMetrics_singleMetric(t *testing.T) {
	rs := []map[string]any{
		{
			"key1":  "value1",
			"value": 2,
		},
	}

	ms, err := resultsToMetrics(rs)
	assert.NoError(t, err)
	assert.Equal(t, []*metrics.Metric{
		{
			Value: 2,
			Tags: map[string]string{
				"key1": "value1",
			},
		},
	}, ms)
}
