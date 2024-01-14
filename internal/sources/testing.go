package sources

import (
	"context"
	"github.com/turbolytics/collector/internal/metrics"
	"time"
)

type TestSourcer struct {
	Ms  []*metrics.Metric
	Err error

	NumSourceCalls int
	NumWindowCalls int

	WindowDuration time.Duration
}

func (ts *TestSourcer) Window() *time.Duration {
	ts.NumWindowCalls++
	return &ts.WindowDuration
}

func (ts *TestSourcer) Source(ctx context.Context) ([]*metrics.Metric, error) {
	ts.NumSourceCalls++
	return ts.Ms, ts.Err
}
