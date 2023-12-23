package collector

import (
	"github.com/stretchr/testify/assert"
	"github.com/turbolytics/collector/internal"
	"github.com/turbolytics/collector/internal/metrics"
	"testing"
)

func TestCollector_Transform_AddTagsFromConfig(t *testing.T) {

	coll, err := New(&internal.Config{
		Metric: internal.Metric{
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
	err = coll.Transform(ms)
	assert.NoError(t, err)
	assert.Equal(t, []*metrics.Metric{{
		Tags: map[string]string{
			"key1": "val1",
			"key2": "val2",
		},
	}}, ms)
}
