package timeseries

import (
	"fmt"
	"time"
)

type Window struct {
	Start time.Time
	End   time.Time
}

type HistoricTumblingWindowerOption func(*HistoricTumblingWindower)

func WithHistoricTumblingWindowerNow(now func() time.Time) HistoricTumblingWindowerOption {
	return func(w *HistoricTumblingWindower) {
		w.now = now
	}
}

func NewHistoricTumblingWindower(opts ...HistoricTumblingWindowerOption) HistoricTumblingWindower {
	w := HistoricTumblingWindower{
		now: func() time.Time {
			return time.Now().UTC()
		},
	}
	for _, opt := range opts {
		opt(&w)
	}
	return w
}

type HistoricTumblingWindower struct {
	now func() time.Time
}

// FullWindowsSince represents all missing full windows since the time
func (hw HistoricTumblingWindower) FullWindowsSince(t *time.Time, d time.Duration) ([]Window, error) {
	// no time is provided, just get the last complete window
	if t == nil {
		window := LastCompleteWindow(hw.now(), d)
		return []Window{window}, nil
	}

	// time is provided, get all complete windows from the time provided
	windows, err := TimeWindows(*t, hw.now(), d)
	return windows, err
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
