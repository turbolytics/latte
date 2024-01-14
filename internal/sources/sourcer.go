package sources

import (
	"context"
	"github.com/turbolytics/collector/internal/metrics"
	"time"
)

type Type string

const (
	TypePostgres   Type = "postgres"
	TypeMongoDB    Type = "mongodb"
	TypePrometheus Type = "prometheus"
)

type Sourcer interface {
	Source(ctx context.Context) ([]*metrics.Metric, error)
	Window() *time.Duration
}
