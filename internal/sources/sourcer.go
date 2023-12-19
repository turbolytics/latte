package sources

import (
	"context"
	"github.com/turbolytics/collector/internal/metrics"
)

type Type string

const (
	TypePostgres Type = "postgres"
	TypeMongoDB  Type = "mongodb"
)

type Sourcer interface {
	Source(ctx context.Context) ([]*metrics.Metric, error)
}
