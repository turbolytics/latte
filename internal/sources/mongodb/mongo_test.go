package mongodb

import (
	"github.com/stretchr/testify/assert"
	"github.com/turbolytics/collector/internal/metrics"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
	"time"
)

func Test_resultsToMetrics_MissingValue(t *testing.T) {
	rs := []bson.M{
		{
			"account": "google",
		},
	}

	_, err := resultsToMetrics(rs)
	assert.EqualError(t, err, "each row must contain a \"value\" key")

}

func Test_resultsToMetrics_SingleMetric(t *testing.T) {
	rs := []bson.M{
		{
			"key1":  "value1",
			"value": 2,
		},
	}

	ms, err := resultsToMetrics(rs)
	assert.NoError(t, err)
	// not sure how to cleanly do this
	// right now just remove the dynamic field from
	// output...
	for _, m := range ms {
		m.UUID = ""
		m.Time = time.Time{}
	}
	assert.Equal(t, []*metrics.Metric{
		{
			Value: 2,
			Tags: map[string]string{
				"key1": "value1",
			},
		},
	}, ms)
}
