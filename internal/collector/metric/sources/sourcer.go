package sources

import (
	"context"
	"github.com/turbolytics/collector/internal/metrics"
	"time"
)

type MetricSourcer interface {
	Source(ctx context.Context) ([]*metrics.Metric, error)
	Window() *time.Duration
}
