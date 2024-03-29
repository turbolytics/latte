package metric

import (
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"time"
)

type Type string

const (
	TypeCount Type = "COUNT"
	TypeGauge Type = "GAUGE"
)

type MetricOption func(metric *Metric)

type Metric struct {
	UUID      string            `json:"uuid"`
	Name      string            `json:"name"`
	Value     float64           `json:"value"`
	Type      Type              `json:"type"`
	Tags      map[string]string `json:"tags"`
	Timestamp time.Time         `json:"timestamp"`

	// For 'tick' collectors, this time aligns with the `Timestamp`.
	// For `window` collectors, this time represents the windowing being collected.
	Window *time.Time `json:"window"`
}

func (m *Metric) Map() map[string]any {
	s := map[string]any{
		"uuid":      m.UUID,
		"name":      m.Name,
		"value":     m.Value,
		"type":      m.Type,
		"timestamp": m.Timestamp,
		"window":    m.Window,
	}

	for k, v := range m.Tags {
		tagK := fmt.Sprintf("tag.%s", k)
		fmt.Println(k, v)
		s[tagK] = v
	}
	return s
}

func New(opts ...MetricOption) Metric {
	m := Metric{
		UUID:      uuid.New().String(),
		Tags:      make(map[string]string),
		Timestamp: time.Now().UTC(),
	}
	for _, opt := range opts {
		opt(&m)
	}
	return m
}

func MapsToMetrics(results []map[string]any) ([]*Metric, error) {
	var ms []*Metric
	fmt.Println(results)
	for _, r := range results {
		val, ok := r["value"]
		if !ok {
			return nil, fmt.Errorf("each row must contain a %q key", "value")
		}

		m := New()

		switch v := val.(type) {
		case int:
			m.Value = float64(v)
		case int32:
			m.Value = float64(v)
		case int64:
			m.Value = float64(v)
		case string:
			tv, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, fmt.Errorf("unable to parse string to float: %q", v)
			}
			m.Value = tv
		}
		delete(r, "value")
		for k, v := range r {
			m.Tags[k] = v.(string)
		}
		ms = append(ms, &m)
	}
	return ms, nil
}
