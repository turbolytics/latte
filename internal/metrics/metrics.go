package metrics

type Type string

const (
	TypeCount Type = "COUNT"
	TypeGauge Type = "GAUGE"
)

type Metric struct {
	Name  string
	Value float64
	Type  Type
	Tags  map[string]string
}

func New() Metric {
	return Metric{
		Tags: make(map[string]string),
	}
}
