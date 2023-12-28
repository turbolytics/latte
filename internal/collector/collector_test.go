package collector

import (
	"github.com/stretchr/testify/assert"
	"github.com/turbolytics/collector/internal"
	"github.com/turbolytics/collector/internal/metrics"
	"testing"
	"time"
)

func TestCollector_Transform_AddTagsFromConfig(t *testing.T) {
	d := time.Duration(0)
	coll, err := New(&internal.Config{
		Metric: internal.Metric{
			Grain: &d,
			Tags: []internal.Tag{
				{"key1", "val1"},
				{"key2", "val2"},
			},
		},
	})
	assert.NoError(t, err)

	ms := []*metrics.Metric{{
		Tags: make(map[string]string),
	}}

	gt := time.Now().UTC()
	err = coll.Transform(gt, ms)

	assert.NoError(t, err)
	assert.Equal(t, []*metrics.Metric{{
		GrainDatetime: gt,
		Tags: map[string]string{
			"key1": "val1",
			"key2": "val2",
		},
	}}, ms)
}
