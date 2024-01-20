package timeseries

import (
	"fmt"
	"time"
)

type Window struct {
	Start time.Time
	End   time.Time
}

// TimeWindows calculates all complete buckets (of duration d)
// from t until now.
func TimeWindows(start time.Time, end time.Time, d time.Duration) ([]Window, error) {
	if start.After(end) {
		return nil, fmt.Errorf("start datetime (%s) must be before end datetime (%s)", start, end)
	}

	if end.Sub(start) < d {
		return nil, nil
	}

	var buckets []Window
	currStart := start
	currEnd := start.Add(d)
	for {
		if currEnd.After(end) {
			break
		}

		b := Window{
			Start: currStart,
			End:   currEnd,
		}

		currStart = currEnd
		currEnd = currEnd.Add(d)

		buckets = append(buckets, b)
	}
	return buckets, nil
}

func LastCompleteWindow(ct time.Time, d time.Duration) Window {
	b := Window{}
	b.End = ct.Truncate(d)
	b.Start = b.End.Add(-d)
	return b
}
