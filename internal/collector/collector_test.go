package collector

import (
	"github.com/stretchr/testify/assert"
	"github.com/turbolytics/latte/internal/collector/metric"
	"github.com/turbolytics/latte/internal/collector/sink"
	"testing"
)

func TestCollector_Close(t *testing.T) {
	ts := &metric.TestSink{}
	coll := &Invoker{
		Config: metric.Config{
			Sinks: map[string]sink.Config{
				"sink1": {
					Sinker: ts,
				},
				"sink2": {
					Sinker: ts,
				},
			},
		},
	}
	err := coll.Close()
	assert.NoError(t, err)
	assert.Equal(t, 2, ts.Closes)
}
