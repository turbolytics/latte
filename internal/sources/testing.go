package sources

import (
	"context"
	"github.com/turbolytics/collector/internal/metrics"
	"time"
)

type TestSourcer struct {
	Ms  []*metrics.Metric
	Err error

	WindowDuration time.Duration
	SourceCalls    []context.Context
}

func (ts *TestSourcer) Window() *time.Duration {
	return &ts.WindowDuration
}

func (ts *TestSourcer) Source(ctx context.Context) ([]*metrics.Metric, error) {
	ts.SourceCalls = append(ts.SourceCalls, ctx)
	return ts.Ms, ts.Err
}
