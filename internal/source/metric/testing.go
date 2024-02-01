package metric

import (
	"context"
	"github.com/turbolytics/latte/internal/metric"
	"time"
)

type TestSourcer struct {
	Ms  []*metric.Metric
	Err error

	WindowDuration time.Duration
	SourceCalls    []context.Context
}

func (ts *TestSourcer) Window() *time.Duration {
	return &ts.WindowDuration
}

func (ts *TestSourcer) Source(ctx context.Context) ([]*metric.Metric, error) {
	ts.SourceCalls = append(ts.SourceCalls, ctx)
	return ts.Ms, ts.Err
}
