package invoker

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/turbolytics/latte/internal/record"
	"go.uber.org/zap"
	"testing"
)

func TestInvoker_Invoke_UnknownStrategy(t *testing.T) {
	i := &Invoker{
		logger:    zap.NewNop(),
		Collector: TestConfig{},
	}
	err := i.Invoke(context.Background())
	assert.EqualError(t, err, "strategy: \"\" not supported")
}

func TestCollector_Close(t *testing.T) {
	sink := &TestSink{}
	i := &Invoker{
		logger: zap.NewNop(),
		Collector: TestConfig{
			sinks: []*TestSink{sink},
		},
	}
	err := i.Close()
	assert.NoError(t, err)
	assert.Equal(t, 1, sink.closes)
}

func TestInvoker_Invoke_Tick_Success(t *testing.T) {
	sink := &TestSink{}
	i := &Invoker{
		logger: zap.NewNop(),
		Collector: TestConfig{
			invocationStrategy: TypeStrategyTick,
			sinks:              []*TestSink{sink},
			sourcer: TestSourcer{
				tr: TestResult{
					records: []*TestRecord{
						{
							m: map[string]any{
								"key": "value",
							},
						},
					},
				},
			},
		},
	}
	err := i.Invoke(context.Background())
	assert.NoError(t, err)
	assert.Equal(t,
		[]record.Record{
			&TestRecord{
				m: map[string]any{
					"key_transformed": "value_transformed",
				},
			},
		},
		sink.writes,
	)
}

/*
func TestCollector_Transform_AddTagsFromConfig(t *testing.T) {
	coll, err := NewCollector(&Collector{
		Metric: Metric{
			Tags: []Tag{
				{"key1", "val1"},
				{"key2", "val2"},
			},
		},
	})
	assert.NoError(t, err)

	ms := []*latteMetric.Metric{{
		Tags: make(map[string]string),
	}}

	err = coll.Transform(ms)

	assert.NoError(t, err)
	assert.Equal(t, []*latteMetric.Metric{{
		Tags: map[string]string{
			"key1": "val1",
			"key2": "val2",
		},
	}}, ms)
}

func TestCollector_Source_ValidMetrics(t *testing.T) {
	expectedMetrics := []*latteMetric.Metric{
		{
			Name: "test.metric",
		},
	}

	ts := &metric.TestSourcer{
		Ms: expectedMetrics,
	}

	coll := &Collector{
		Collector: &Collector{
			Source: source.Collector{
				MetricSourcer: ts,
				Strategy:      source.TypeStrategyTick,
			},
		},
	}
	ms, err := coll.Source(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedMetrics, ms)
}

func TestCollector_invokeWindow_NoPreviousInvocations(t *testing.T) {
	expectedMetrics := []*latteMetric.Metric{
		{
			Name: "test.metric",
		},
	}

	now := time.Date(2024, 1, 1, 1, 1, 0, 0, time.UTC)

	ts := &metric.TestSourcer{
		Ms:             expectedMetrics,
		WindowDuration: time.Minute,
	}
	ss, _ := state.NewMemoryStoreFromGenericConfig(map[string]any{})

	coll := &Collector{
		logger: zap.NewNop(),
		now: func() time.Time {
			return now
		},
		Collector: &Collector{
			Name: "test_collector",
			StateStore: state.Collector{
				Storer: ss,
			},
			Source: source.Collector{
				MetricSourcer: ts,
				Strategy:      source.TypeStrategyHistoricTumblingWindow,
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
		coll.Collector.Name,
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

	ts := &metric.TestSourcer{
		WindowDuration: time.Hour,
	}
	ss, _ := state.NewMemoryStoreFromGenericConfig(map[string]any{})
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
		Collector: &Collector{
			Name: "test_collector",
			StateStore: state.Collector{
				Storer: ss,
			},
			Source: source.Collector{
				MetricSourcer: ts,
				Strategy:      source.TypeStrategyHistoricTumblingWindow,
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
	expectedMetrics := []*latteMetric.Metric{
		{
			Name: "test.metric",
		},
	}

	now := time.Date(2024, 1, 1, 3, 1, 0, 0, time.UTC)

	ts := &metric.TestSourcer{
		Ms:             expectedMetrics,
		WindowDuration: time.Hour,
	}
	ss, _ := state.NewMemoryStoreFromGenericConfig(map[string]any{})
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
		Collector: &Collector{
			Name: "test_collector",
			StateStore: state.Collector{
				Storer: ss,
			},
			Source: source.Collector{
				MetricSourcer: ts,
				Strategy:      source.TypeStrategyHistoricTumblingWindow,
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
*/
