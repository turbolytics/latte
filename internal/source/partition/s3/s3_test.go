package s3

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSource_Source_PartitionSuccess(t *testing.T) {
	s := &Source{
		config: config{
			Bucket:    "test-bucket",
			Region:    "us-east-2",
			Prefix:    "prefix-1/prefix-2",
			Partition: "year={{.Year}}/month={{.Month}}/day={{.Day}}/hour={{.Hour}}/min={{.Minute}}",
		},
	}

	ctx := context.Background()
	start := time.Date(2024, 1, 1, 1, 1, 1, 0, time.UTC)
	ctx = context.WithValue(ctx, "window.start", start)

	p, err := s.Source(ctx)
	assert.NoError(t, err)
	assert.Equal(t,
		"https://test-bucket.s3.us-east-2.amazonaws.com/prefix-1/prefix-2/year=2024/month=1/day=1/hour=1/min=1",
		p.URI.String(),
	)
}
