package metrics

type Type string

const (
	TypeCount Type = "COUNT"
	TypeGauge Type = "GAUGE"
)

type Metric struct {
	Name  string            `json:"name"`
	Value float64           `json:"value"`
	Type  Type              `json:"type"`
	Tags  map[string]string `json:"tags"`
}

func New() Metric {
	return Metric{
		Tags: make(map[string]string),
	}
}
