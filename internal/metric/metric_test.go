package metric

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_MapsToMetrics_MissingValue(t *testing.T) {
	rs := []map[string]any{
		{
			"key1": "value1",
		},
	}

	_, err := MapsToMetrics(rs)
	assert.EqualError(t, err, "each row must contain a \"value\" key")
}

func Test_MapsToMetrics_SingleMetric(t *testing.T) {
	rs := []map[string]any{
		{
			"key1":  "value1",
			"value": 2,
		},
	}

	ms, err := MapsToMetrics(rs)
	assert.NoError(t, err)
	// not sure how to cleanly do this
	// right now just remove the dynamic field from
	// output...
	for _, m := range ms {
		m.UUID = ""
		m.Timestamp = time.Time{}
	}
	assert.Equal(t, []*Metric{
		{
			Value: 2,
			Tags: map[string]string{
				"key1": "value1",
			},
		},
	}, ms)
}
