package timeseries

import (
	"fmt"
	"time"
)

type Bucket struct {
	Start time.Time
	End   time.Time
}

// TimeBuckets calculates all complete buckets (of duration d)
// from t until now.
func TimeBuckets(start time.Time, end time.Time, d time.Duration) ([]Bucket, error) {
	if start.After(end) {
		return nil, fmt.Errorf("start datetime (%s) must be before end datetime (%s)", start, end)
	}

	if end.Sub(start) < d {
		return nil, nil
	}

	var buckets []Bucket
	currStart := start
	currEnd := start.Add(d)
	for {
		if currEnd.After(end) {
			break
		}

		b := Bucket{
			Start: currStart,
			End:   currEnd,
		}

		currStart = currEnd
		currEnd = currEnd.Add(d)

		buckets = append(buckets, b)
	}
	return buckets, nil
}

func LastCompleteBucket(ct time.Time, d time.Duration) Bucket {
	b := Bucket{}
	b.End = ct.Truncate(d)
	b.Start = b.End.Add(-d)
	return b
}
