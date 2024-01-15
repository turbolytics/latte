package timeseries

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLastCompleteBucket(t *testing.T) {
	testCases := []struct {
		name           string
		ct             time.Time
		d              time.Duration
		expectedBucket Bucket
	}{
		{
			name: "start_before_end_bucket_aligned",
			ct:   time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			d:    time.Hour * 24,
			expectedBucket: Bucket{
				Start: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "start_before_end_bucket_over_end",
			ct:   time.Date(2024, 1, 2, 1, 15, 0, 0, time.UTC),
			d:    time.Hour * 24,
			expectedBucket: Bucket{
				Start: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bucket := LastCompleteBucket(
				tc.ct,
				tc.d,
			)
			assert.Equal(t, tc.expectedBucket, bucket)
		})
	}
}

func TestTimeBuckets(t *testing.T) {
	testCases := []struct {
		name    string
		start   time.Time
		end     time.Time
		d       time.Duration
		buckets []Bucket
		err     error
	}{
		{
			name:    "invalid_end_before_start",
			start:   time.Date(2024, 01, 01, 01, 20, 00, 00, time.UTC),
			end:     time.Date(2024, 01, 01, 01, 00, 00, 00, time.UTC),
			d:       time.Hour,
			err:     errors.New("start datetime (2024-01-01 01:20:00 +0000 UTC) must be before end datetime (2024-01-01 01:00:00 +0000 UTC)"),
			buckets: nil,
		},
		{
			name:    "no_full_window",
			start:   time.Date(2024, 01, 01, 01, 0, 00, 00, time.UTC),
			end:     time.Date(2024, 01, 01, 01, 10, 00, 00, time.UTC),
			d:       time.Hour,
			err:     nil,
			buckets: nil,
		},
		{
			name:  "multiple_valid_buckets_window_aligned",
			start: time.Date(2024, 01, 01, 01, 00, 00, 00, time.UTC),
			end:   time.Date(2024, 01, 01, 01, 10, 00, 00, time.UTC),
			d:     time.Minute * 5,
			err:   nil,
			buckets: []Bucket{
				{
					Start: time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC),
					End:   time.Date(2024, 1, 1, 1, 5, 0, 0, time.UTC),
				},
				{
					Start: time.Date(2024, 1, 1, 1, 5, 0, 0, time.UTC),
					End:   time.Date(2024, 1, 1, 1, 10, 0, 0, time.UTC),
				},
			},
		},
		{
			name:  "single_valid_bucket_window_aligned",
			start: time.Date(2024, 01, 01, 01, 00, 00, 00, time.UTC),
			end:   time.Date(2024, 01, 01, 01, 10, 00, 00, time.UTC),
			d:     time.Minute * 10,
			err:   nil,
			buckets: []Bucket{
				{
					Start: time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC),
					End:   time.Date(2024, 1, 1, 1, 10, 0, 0, time.UTC),
				},
			},
		},
		{
			name:  "single_valid_bucket_into_second_window",
			start: time.Date(2024, 01, 01, 01, 00, 00, 00, time.UTC),
			end:   time.Date(2024, 01, 01, 01, 12, 00, 00, time.UTC),
			d:     time.Minute * 10,
			err:   nil,
			buckets: []Bucket{
				{
					Start: time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC),
					End:   time.Date(2024, 1, 1, 1, 10, 0, 0, time.UTC),
				},
			},
		},
		{
			name:  "multiple_valid_buckets_into_next_window",
			start: time.Date(2024, 01, 01, 01, 00, 00, 00, time.UTC),
			end:   time.Date(2024, 01, 01, 01, 12, 00, 00, time.UTC),
			d:     time.Minute * 5,
			err:   nil,
			buckets: []Bucket{
				{
					Start: time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC),
					End:   time.Date(2024, 1, 1, 1, 5, 0, 0, time.UTC),
				},
				{
					Start: time.Date(2024, 1, 1, 1, 5, 0, 0, time.UTC),
					End:   time.Date(2024, 1, 1, 1, 10, 0, 0, time.UTC),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buckets, err := TimeBuckets(
				tc.start,
				tc.end,
				tc.d,
			)

			assert.Equal(t, tc.buckets, buckets)

			if tc.err == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, tc.err, err)
			}
		})
	}
}
