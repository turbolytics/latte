package metrics

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/turbolytics/collector/internal/collector"
	"github.com/turbolytics/collector/internal/collector/metrics/config"
	"github.com/turbolytics/collector/internal/collector/metrics/sources"
	"github.com/turbolytics/collector/internal/collector/state"
	"github.com/turbolytics/collector/internal/collector/state/memory"
	"github.com/turbolytics/collector/internal/metrics"
	"github.com/turbolytics/collector/internal/timeseries"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestCollector_Transform_AddTagsFromConfig(t *testing.T) {
	coll, err := New(&config.Config{
		Metric: config.Metric{
			Tags: []config.Tag{
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

func TestCollector_Close(t *testing.T) {
	ts := &collector.TestSink{}
	coll := &Collector{
		Config: &config.Config{
			Sinks: map[string]config.Sink{
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

func TestCollector_Source_ValidMetrics(t *testing.T) {
	expectedMetrics := []*metrics.Metric{
		{
			Name: "test.metric",
		},
	}

	ts := &sources.TestSourcer{
		Ms: expectedMetrics,
	}

	coll := &Collector{
		Config: &config.Config{
			Source: config.Source{
				Sourcer:  ts,
				Strategy: config.TypeSourceStrategyTick,
			},
		},
	}
	ms, err := coll.Source(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedMetrics, ms)
}

func TestCollector_invokeWindow_NoPreviousInvocations(t *testing.T) {
	expectedMetrics := []*metrics.Metric{
		{
			Name: "test.metric",
		},
	}

	now := time.Date(2024, 1, 1, 1, 1, 0, 0, time.UTC)

	ts := &sources.TestSourcer{
		Ms:             expectedMetrics,
		WindowDuration: time.Minute,
	}
	ss, _ := memory.NewFromGenericConfig(map[string]any{})

	coll := &Collector{
		logger: zap.NewNop(),
		now: func() time.Time {
			return now
		},
		Config: &config.Config{
			Name: "test_collector",
			StateStore: config.StateStore{
				Storer: ss,
			},
			Source: config.Source{
				Sourcer:  ts,
				Strategy: config.TypeSourceStrategyWindow,
			},
		},
	}
	ms, err := coll.invokeWindow(
		context.Background(),
		uuid.New(),
	)

	assert.NoError(t, err)
	assert.Equal(t, expectedMetrics, ms)
	assert.Equal(t, 1, len(ts.SourceCalls))

	callCtx := ts.SourceCalls[0]
	assert.Equal(t,
		time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC),
		callCtx.Value("window.start").(time.Time),
	)
	assert.Equal(t,
		time.Date(2024, 1, 1, 1, 1, 0, 0, time.UTC),
		callCtx.Value("window.end").(time.Time),
	)

	i, err := ss.MostRecentInvocation(
		context.Background(),
		coll.Config.Name,
	)

	assert.NoError(t, err)
	assert.Equal(t, &state.Invocation{
		CollectorName: "test_collector",
		Time:          now,
		Window: &timeseries.Window{
			Start: time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC),
			End:   time.Date(2024, 1, 1, 1, 1, 0, 0, time.UTC),
		},
	}, i)
}

func TestCollector_invokeWindow_PreviousInvocations_MultipleWindowsPassed(t *testing.T) {
	now := time.Date(2024, 1, 1, 4, 1, 0, 0, time.UTC)

	ts := &sources.TestSourcer{
		WindowDuration: time.Hour,
	}
	ss, _ := memory.NewFromGenericConfig(map[string]any{})
	ss.SaveInvocation(&state.Invocation{
		CollectorName: "test_collector",
		Window: &timeseries.Window{
			Start: time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC),
			End:   time.Date(2024, 1, 1, 2, 0, 0, 0, time.UTC),
		},
	})

	coll := &Collector{
		logger: zap.NewNop(),
		now: func() time.Time {
			return now
		},
		Config: &config.Config{
			Name: "test_collector",
			StateStore: config.StateStore{
				Storer: ss,
			},
			Source: config.Source{
				Sourcer:  ts,
				Strategy: config.TypeSourceStrategyWindow,
			},
		},
	}
	ms, err := coll.invokeWindow(
		context.Background(),
		uuid.New(),
	)

	assert.EqualError(t, err, "backfilling multiple windows not yet supported: [{2024-01-01 02:00:00 +0000 UTC 2024-01-01 03:00:00 +0000 UTC} {2024-01-01 03:00:00 +0000 UTC 2024-01-01 04:00:00 +0000 UTC}]")
	assert.Nil(t, ms)
}

func TestCollector_invokeWindow_PreviousInvocations_SingleWindowPassed(t *testing.T) {
	expectedMetrics := []*metrics.Metric{
		{
			Name: "test.metric",
		},
	}

	now := time.Date(2024, 1, 1, 3, 1, 0, 0, time.UTC)

	ts := &sources.TestSourcer{
		Ms:             expectedMetrics,
		WindowDuration: time.Hour,
	}
	ss, _ := memory.NewFromGenericConfig(map[string]any{})
	err := ss.SaveInvocation(&state.Invocation{
		CollectorName: "test_collector",
		Window: &timeseries.Window{
			Start: time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC),
			End:   time.Date(2024, 1, 1, 2, 0, 0, 0, time.UTC),
		},
	})
	assert.NoError(t, err)

	coll := &Collector{
		logger: zap.NewNop(),
		now: func() time.Time {
			return now
		},
		Config: &config.Config{
			Name: "test_collector",
			StateStore: config.StateStore{
				Storer: ss,
			},
			Source: config.Source{
				Sourcer:  ts,
				Strategy: config.TypeSourceStrategyWindow,
			},
		},
	}
	ms, err := coll.invokeWindow(
		context.Background(),
		uuid.New(),
	)

	assert.NoError(t, err)
	assert.Equal(t, expectedMetrics, ms)
	assert.Equal(t, 1, len(ts.SourceCalls))

	callCtx := ts.SourceCalls[0]
	assert.Equal(t,
		time.Date(2024, 1, 1, 2, 0, 0, 0, time.UTC),
		callCtx.Value("window.start").(time.Time),
	)
	assert.Equal(t,
		time.Date(2024, 1, 1, 3, 0, 0, 0, time.UTC),
		callCtx.Value("window.end").(time.Time),
	)

	i, err := ss.MostRecentInvocation(
		context.Background(),
		"test_collector",
	)
	fmt.Println(i.Window, err)
	assert.NoError(t, err)
	assert.Equal(t, &state.Invocation{
		CollectorName: "test_collector",
		Time:          now,
		Window: &timeseries.Window{
			Start: time.Date(2024, 1, 1, 2, 0, 0, 0, time.UTC),
			End:   time.Date(2024, 1, 1, 3, 0, 0, 0, time.UTC),
		},
	}, i)
}
