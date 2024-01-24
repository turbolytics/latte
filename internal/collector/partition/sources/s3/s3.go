package s3

import (
	"context"
	"github.com/turbolytics/collector/internal/partition"
)

type Source struct{}

func (s *Source) Source(ctx context.Context) (*partition.Partition, error) {
	return nil, nil
}

func NewFromGenericConfig(m map[string]any) (*Source, error) {
	return nil, nil
}
