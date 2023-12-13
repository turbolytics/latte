package metrics

import "github.com/google/uuid"

type Type string

const (
	TypeCount Type = "COUNT"
	TypeGauge Type = "GAUGE"
)

type MetricOption func(metric *Metric)

type Metric struct {
	UUID  string            `json:"uuid"`
	Name  string            `json:"name"`
	Value float64           `json:"value"`
	Type  Type              `json:"type"`
	Tags  map[string]string `json:"tags"`
}

func WithUUID(id string) MetricOption {
	return func(m *Metric) {
		m.UUID = id
	}
}

func New(opts ...MetricOption) Metric {
	m := Metric{
		UUID: uuid.New().String(),
		Tags: make(map[string]string),
	}
	for _, opt := range opts {
		opt(&m)
	}
	return m
}
